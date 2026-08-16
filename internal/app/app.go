package app

import (
	"ab_system/internal/clients/telegram_notifications"
	"ab_system/internal/configs"
	"ab_system/internal/domain/service"
	"ab_system/internal/http/handlers"
	"ab_system/internal/http/middlewares"
	"ab_system/internal/jobs"
	"ab_system/internal/lib/metrics"
	"ab_system/internal/notifications/slack"
	"ab_system/internal/repository/postgres"
	"ab_system/internal/repository/redis"
	"time"

	"context"

	"github.com/gin-gonic/gin"
)

func NewServer(
	r *gin.Engine,
	ctx context.Context,
	cfg configs.Configs,

	usersRepository *postgres.UsersRepository,
	featureFlagRepository *postgres.FeatureFlagRepository,
	experimentsRepository *postgres.ExperimentsRepository,
	experimentVersionsRepository *postgres.ExperimentVersionsRepository,
	variantRepository *postgres.VariantRepository,
	decisionRepository *postgres.DecisionRepository,
	teamRepository *postgres.TeamRepository,
	approverGroupRepository *postgres.ApproverGroupRepository,
	approvalRepository *postgres.ApprovalRepository,
	eventRepository *postgres.EventRepository,
	eventTypeRepository *postgres.EventTypeRepository,
	metricRepository *postgres.MetricRepository,
	expMetricRepository *postgres.ExperimentMetricRepository,
	guardrailTriggerRepository *postgres.GuardrailTriggerRepository,
	notificationSettingsRepository *postgres.NotificationSettingsRepository,

	redisRepository *redis.RedisRepository,

	telegramClient *telegram_notifications.NotifyClient,
) (router *gin.Engine) {
	pingHandler := handlers.NewPingHandler(telegramClient)

	userCreatorService := service.NewUserCreator(usersRepository, usersRepository)
	userService := service.NewUserService(
		usersRepository,
		usersRepository,
		usersRepository,
		usersRepository,
		usersRepository,
	)
	userHandlers := handlers.NewUserHandler(*userService, *userCreatorService)

	authService := service.NewAuthService(usersRepository, cfg)
	authHandlers := handlers.NewAuthHandler(*authService, *userCreatorService, cfg)

	teamService := service.NewTeamService(
		teamRepository,
		teamRepository,
		teamRepository,
		usersRepository,
		usersRepository,
	)
	teamHandlers := handlers.NewTeamHandler(teamService)

	featureFlagService := service.NewFeatureFlagService(
		featureFlagRepository,
		featureFlagRepository,
		featureFlagRepository,
		usersRepository,
	)
	featureFlagHandlers := handlers.NewFeatureFlagHandler(*featureFlagService)

	userActiveChecker := middlewares.UserActiveMiddleware(usersRepository)

	slackNotifier := slack.NewNotifier(
		redisRepository,
		time.Duration(cfg.NotificationParams.TtlNotificationsMinutes)*time.Minute,
	)

	experimentOwnerChecker := middlewares.CheckExperimenterPermission(experimentsRepository)

	experimentsService := service.NewExperimentService(
		experimentsRepository,
		experimentsRepository,
		experimentsRepository,
		experimentVersionsRepository,
		experimentVersionsRepository,
		notificationSettingsRepository,
		notificationSettingsRepository,
		usersRepository,
		telegramClient,
		slackNotifier,
	)
	experimentsHandlers := handlers.NewExperimentHandler(*experimentsService)

	decisionService := service.NewDecisionService(
		featureFlagRepository,
		experimentsRepository,
		decisionRepository)
	decisionHandlers := handlers.NewDecisionHandler(*decisionService)

	approverGroupService := service.NewApproverGroupService(
		approverGroupRepository,
		approverGroupRepository,
		usersRepository,
	)
	approverGroupHandlers := handlers.NewApproverGroupHandler(approverGroupService)

	approvalService := service.NewApprovalService(
		approvalRepository,
		approvalRepository,
		approvalRepository,
		experimentsRepository,
		experimentsRepository,
		approverGroupRepository,
		approverGroupRepository,
		usersRepository,
		telegramClient,
	)
	approvalHandlers := handlers.NewApprovalHandlers(*approvalService)

	eventService := service.NewEventService(
		eventTypeRepository,
		eventRepository,
		eventRepository,
		decisionRepository,
	)
	eventHandlers := handlers.NewEventHandler(eventService)

	eventTypeService := service.NewEventTypeService(
		eventTypeRepository,
		eventTypeRepository,
		eventTypeRepository,
	)
	eventTypeHandlers := handlers.NewEventTypeHandler(eventTypeService)

	metricsService := service.NewMetricService(
		metricRepository,
		metricRepository,
		metricRepository,
	)
	metricsHandlers := handlers.NewMetricHandler(metricsService)

	expMetricService := service.NewExperimentMetricService(
		expMetricRepository,
		expMetricRepository,
		expMetricRepository,
		experimentsRepository,
		metricRepository,
	)
	expMetricHandler := handlers.NewExperimentMetricHandler(expMetricService)

	metricLib := metrics.NewMetricLib(eventRepository)

	reportService := service.NewReportService(
		eventRepository,
		decisionRepository,
		metricRepository,
		expMetricRepository,
		variantRepository,
		metricLib,
	)

	reportHandlers := handlers.NewReportHandler(*reportService)

	notifSettingsService := service.NewNotificationSettingsService(
		notificationSettingsRepository,
		notificationSettingsRepository,
		notificationSettingsRepository,
		experimentsRepository,
	)
	notifSettingsHandler := handlers.NewNotificationSettingsHandler(notifSettingsService)

	guardrailChecker := jobs.NewGuardrailChecker(
		cfg.AppParams,

		experimentsRepository,
		expMetricRepository,
		metricRepository,
		eventRepository,
		decisionRepository,
		experimentsRepository,
		guardrailTriggerRepository,
		metricLib,
		usersRepository,

		notificationSettingsRepository,

		telegramClient,
		slackNotifier,

		time.Duration(cfg.GuardrailParams.CheckerTicketMinutes)*time.Minute,
	)

	guardrailChecker.Start(ctx)

	autopilotChecker := jobs.NewAutopilotRampUpChecker(
		experimentsRepository,
		experimentsRepository,

		time.Duration(cfg.AutopilotParams.CheckerTicketMinutes)*time.Minute,
	)

	autopilotChecker.Start(ctx)

	r.Use(
		gin.Recovery(),
		middlewares.TraceMiddleware,
		middlewares.ValidateUUID,
		userActiveChecker,
	)
	//r.Use(middlewares.DebugBody())

	api := r.Group("/api/v1")

	api.GET("health", pingHandler.Ping)
	api.GET("ready", pingHandler.Ready)

	pingGroup := api.Group("/ping")
	{
		pingGroup.GET("", pingHandler.Ping)
	}

	authGroup := api.Group("/auth")
	{
		//authGroup.POST("/register", authHandlers.RegisterUser)
		authGroup.POST("/login", authHandlers.LoginUser)
	}

	usersGroup := api.Group(
		"/users",
		middlewares.CheckUserAuthentication,
	)
	{
		usersGroup.GET(
			"",
			middlewares.CheckUserAdmin,
			userHandlers.GetAllUsers,
		)

		usersGroup.GET("/me", userHandlers.GetUserMe)

		usersGroup.GET(
			"/:id",
			middlewares.UserProtection,
			userHandlers.GetUserByID,
		)

		usersGroup.GET(
			"/team/:id",
			middlewares.UserProtection,
			userHandlers.GetUsersByTeamID,
		)

		usersGroup.POST(
			"",
			middlewares.CheckUserAdmin,
			userHandlers.CreateUser,
		)

		usersGroup.PUT("/me", userHandlers.UpdateUserMe)

		usersGroup.PUT(
			"/:id",
			middlewares.UserProtection,
			userHandlers.UpdateUser,
		)

		usersGroup.DELETE(
			"/:id",
			middlewares.CheckUserAdmin,
			userHandlers.DeactivateUser,
		)
	}

	teamGroup := api.Group("/teams", middlewares.CheckUserAuthentication)
	{
		teamGroup.GET("", teamHandlers.GetAllTeams)
		teamGroup.GET("/:id", teamHandlers.GetTeamByID)
		teamGroup.POST("", middlewares.CheckUserAdmin, teamHandlers.CreateTeam)
		teamGroup.POST("/member", middlewares.CheckUserAdmin, teamHandlers.AddMemberToTeam)
		teamGroup.PATCH("/:id", middlewares.CheckUserAdmin, teamHandlers.UpdateTeam)
		teamGroup.DELETE("/:id", middlewares.CheckUserAdmin, teamHandlers.DeleteTeam)
	}

	featureFlag := api.Group("/feature-flag", middlewares.CheckUserAuthentication, middlewares.CheckUserAdmin)
	{
		featureFlag.GET("", featureFlagHandlers.GetAllFeatureFlags)
		featureFlag.GET("/:id", featureFlagHandlers.GetFeatureFlagById)
		featureFlag.GET("/key/:key", featureFlagHandlers.GetFeatureFlagsByKey)
		featureFlag.GET("/owner/:id", featureFlagHandlers.GetFeatureFlagsByOwner)
		featureFlag.POST("", featureFlagHandlers.CreateFeatureFlag)
		featureFlag.PATCH("/:id", featureFlagHandlers.UpdateFeatureFlag)
		featureFlag.DELETE("/:id", featureFlagHandlers.DeleteFeatureFlagById)
	}

	experimentsGroup := api.Group("/experiments", middlewares.CheckUserAuthentication)
	{
		experimentsGroup.GET("", experimentsHandlers.GetAllExperiments)
		experimentsGroup.GET("/:id", experimentsHandlers.GetExperimentByID)
		experimentsGroup.GET("/flag/:flag/:status", experimentsHandlers.GetRunningExperimentByFlag)
		experimentsGroup.GET("/status/:status", experimentsHandlers.GetAllExperimentsByStatus)
		experimentsGroup.POST("", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, experimentsHandlers.CreateExperiment)
		experimentsGroup.PATCH("review/:id", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, experimentOwnerChecker, experimentsHandlers.SendExperimentToReview)
		experimentsGroup.PATCH("run/:id", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, experimentOwnerChecker, experimentsHandlers.RunExperiment)
		experimentsGroup.PATCH("complete/:id", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, experimentOwnerChecker, experimentsHandlers.CompleteExperiment)
		experimentsGroup.PUT("/:id", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, experimentOwnerChecker, experimentsHandlers.UpdateExperiment)
		experimentsGroup.DELETE("/archive/:id", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, experimentOwnerChecker, experimentsHandlers.ArchiveExperiment)
	}

	expGroup := experimentsGroup.Group("/:id")
	{
		expGroup.GET("/metrics", expMetricHandler.GetExperimentMetrics)
		expGroup.GET("/guardrails", expMetricHandler.GetGuardrailsForExperiment)
		expGroup.POST("/metrics", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, experimentOwnerChecker, expMetricHandler.AddMetricToExperiment)
		expGroup.PUT("/metrics/:metric_id", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, experimentOwnerChecker, expMetricHandler.UpdateExperimentMetric)
		expGroup.DELETE("/metrics/:metric_id", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, experimentOwnerChecker, expMetricHandler.RemoveMetricFromExperiment)
	}

	expNotificationsGroup := experimentsGroup.Group("/:id")
	{
		expNotificationsGroup.GET("/notification-settings", experimentOwnerChecker, notifSettingsHandler.GetNotificationSettingsByExperimentID)
		expNotificationsGroup.POST("/notification-settings", experimentOwnerChecker, notifSettingsHandler.CreateNotificationSettings)
		expNotificationsGroup.PATCH("/notification-settings", experimentOwnerChecker, notifSettingsHandler.UpdateNotificationSettings)
		//expNotificationsGroup.DELETE("/notification-settings", experimentOwnerChecker, notifSettingsHandler.DeleteNotificationSettings)
	}

	experimentVersionsGroup := experimentsGroup.Group("/versions")
	{
		experimentVersionsGroup.GET("/:id", experimentsHandlers.GetExperimentVersionsByID)
	}

	decisionGroup := api.Group("/decisions")
	{
		decisionGroup.POST("/decide", decisionHandlers.Decide)
	}

	approverGroup := api.Group("/approvers", middlewares.CheckUserAuthentication)
	{
		approverGroup.GET("/:id", approverGroupHandlers.GetApproverGroupByExperimenterID)
		approverGroup.GET("/experiment/:id", approverGroupHandlers.GetApproverGroupByExperimentID)
		approverGroup.POST("", middlewares.CheckUserAdmin, approverGroupHandlers.CreateApproverGroup)
		approverGroup.PUT("/:id", middlewares.CheckUserAdmin, approverGroupHandlers.UpdateApproverGroup)
		approverGroup.DELETE("/:id", middlewares.CheckUserAdmin, approverGroupHandlers.DeleteApproverGroup)
	}

	approvalGroup := api.Group("/approvals", middlewares.CheckUserAuthentication)
	{
		approvalGroup.GET("/:id", approvalHandlers.GetApprovalByID)
		approvalGroup.GET("/experiment/:id", approvalHandlers.GetApprovalsByExperimentID)
		approvalGroup.POST("", middlewares.CheckUserAdminSoft, middlewares.CheckUserApprover, approvalHandlers.CreateApproval)
		//approvalGroup.DELETE("/:id", middlewares.ValidateUUID, approvalHandlers.DeleteApproval)
	}

	eventsGroup := api.Group("/events")
	{
		eventsGroup.POST("", eventHandlers.PostEvent)
		eventsGroup.POST("/batch", eventHandlers.PostEvents)
	}

	eventTypesGroup := api.Group("/event-types", middlewares.CheckUserAuthentication)
	{
		eventTypesGroup.GET("", eventTypeHandlers.GetAllEventTypes)
		eventTypesGroup.GET("/:id", eventTypeHandlers.GetEventTypeByID)
		eventTypesGroup.POST("", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, eventTypeHandlers.CreateEventType)
		eventTypesGroup.PUT("/:id", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, eventTypeHandlers.UpdateEventType)
		//eventTypesGroup.DELETE("/:id", middlewares.CheckUserAdmin, middlewares.CheckUserExperimenter, middlewares.ValidateUUID, eventTypeHandlers.DeleteEventType)
	}

	metricsGroup := api.Group("/metrics", middlewares.CheckUserAuthentication)
	{
		metricsGroup.GET("", metricsHandlers.GetAllMetrics)
		metricsGroup.GET("/:id", metricsHandlers.GetMetricByID)
		metricsGroup.GET("title/:title", metricsHandlers.GetMetricByTitle)
		metricsGroup.POST("", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, metricsHandlers.CreateMetric)
		metricsGroup.PUT("/:id", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, metricsHandlers.UpdateMetric)
		metricsGroup.DELETE("/:id", middlewares.CheckUserAdminSoft, middlewares.CheckUserExperimenter, metricsHandlers.DeleteMetric)
	}

	reportsGroup := api.Group("/reports", middlewares.CheckUserAuthentication)
	{
		reportsGroup.GET("/:id", reportHandlers.GetExperimentReport)
	}

	return r
}
