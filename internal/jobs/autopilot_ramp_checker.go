package jobs

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/pkg/logger"
	"context"
	"encoding/json"
	"time"
)

type AutopilotRampUpChecker struct {
	experimentReader repository.ExperimentReader
	experimentWriter repository.ExperimentWriter
	interval         time.Duration
}

func NewAutopilotRampUpChecker(
	experimentReader repository.ExperimentReader,
	experimentWriter repository.ExperimentWriter,
	interval time.Duration,
) *AutopilotRampUpChecker {
	return &AutopilotRampUpChecker{
		experimentReader: experimentReader,
		experimentWriter: experimentWriter,
		interval:         interval,
	}
}

func (c *AutopilotRampUpChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(c.interval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info.Println("Quitting Autopilot Ramp Up Checker")

				ticker.Stop()
				return
			case <-ticker.C:
				logger.Info.Println("Checking Autopilot Ramp Up Checker")
				c.checkExperiments(ctx)
			}
		}
	}()
}

func (c *AutopilotRampUpChecker) checkExperiments(ctx context.Context) {
	experiments, err := c.experimentReader.GetExperimentsWithRampEnabled(ctx)
	if err != nil {
		logger.Error.Printf("[AutopilotRampUpChecker] Failed to get experiments: %v", err)
		return
	}

	for _, exp := range *experiments {
		c.processExperiment(ctx, &exp)
	}
}

func (c *AutopilotRampUpChecker) processExperiment(ctx context.Context, exp *models.Experiment) {
	if exp.Status != models.StatusRunning {
		return
	}

	if time.Since(exp.RampUp.RampLastIncrease) < time.Duration(exp.RampUp.RampIntervalMinutes)*time.Minute {
		return
	}

	var steps []int
	if err := json.Unmarshal(exp.RampUp.RampSteps, &steps); err != nil || len(steps) == 0 {
		logger.Error.Printf("[service.AutopilotRampUpChecker] Invalid ramp steps for experiment %s: %v", exp.ID, err)
		return
	}

	if exp.RampUp.RampCurrentStep >= len(steps)-1 {
		return
	}

	if exp.GuardrailTriggered {
		return
	}

	nextStep := exp.RampUp.RampCurrentStep + 1
	newTraffic := steps[nextStep]

	if err := c.updateVariantWeights(exp, newTraffic); err != nil {
		logger.Error.Printf("[AutopilotRampUpChecker] Failed to update variant weights: %v", err)
		return
	}

	exp.TrafficPercent = newTraffic
	exp.RampUp.RampCurrentStep = nextStep
	exp.RampUp.RampLastIncrease = time.Now()

	if _, err := c.experimentWriter.UpdateExperiment(ctx, exp); err != nil {
		logger.Error.Printf("[AutopilotRampUpChecker] Failed to update experiment: %v", err)
		return
	}

	logger.Info.Printf("[AutopilotRampUpChecker] Experiment %s ramped up to %d%%", exp.ID, newTraffic)
}

func (c *AutopilotRampUpChecker) updateVariantWeights(exp *models.Experiment, newTraffic int) error {
	if exp.TrafficPercent == 0 {
		return nil
	}

	ratio := float64(newTraffic) / float64(exp.TrafficPercent)
	total := 0

	for i := range exp.Variants {
		w := int(float64(exp.Variants[i].Weight) * ratio)
		if w < 1 {
			w = 1
		}
		exp.Variants[i].Weight = w
		total += w
	}

	if diff := newTraffic - total; diff != 0 {
		for i := range exp.Variants {
			if !exp.Variants[i].IsControl {
				exp.Variants[i].Weight += diff
				break
			}
		}
	}

	return nil
}
