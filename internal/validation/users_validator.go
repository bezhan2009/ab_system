package validation

import (
	"ab_system/internal/http/dto"
	"ab_system/pkg/errs"
	"regexp"
	"strings"
)

func ValidateCreateUser(user *dto.User) []errs.FieldError {
	var errors []errs.FieldError

	if user.Email == "" {
		errors = append(errors, errs.NewFieldError("email", "обязательное поле"))
	} else {
		if len(user.Email) > 254 {
			errors = append(errors, errs.NewFieldError("email", "не более 254 символов", user.Email))
		}
		emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		if !regexp.MustCompile(emailRegex).MatchString(user.Email) {
			errors = append(errors, errs.NewFieldError("email", "неверный формат email", user.Email))
		}
	}

	if user.Password == "" {
		errors = append(errors, errs.NewFieldError("password", "обязательное поле"))
	} else {
		if len(user.Password) < 8 || len(user.Password) > 72 {
			errors = append(errors, errs.NewFieldError("password", "должен быть от 8 до 72 символов"))
		}
		hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(user.Password)
		hasDigit := regexp.MustCompile(`\d`).MatchString(user.Password)
		if !hasLetter || !hasDigit {
			errors = append(errors, errs.NewFieldError("password", "должен содержать минимум 1 букву и 1 цифру"))
		}
	}

	if user.FullName == "" {
		errors = append(errors, errs.NewFieldError("fullName", "обязательное поле"))
	} else if len(user.FullName) < 2 || len(user.FullName) > 200 {
		errors = append(errors, errs.NewFieldError("fullName", "должно быть от 2 до 200 символов", user.FullName))
	}

	if user.Region != nil && len(*user.Region) > 32 {
		errors = append(errors, errs.NewFieldError("region", "не более 32 символов", *user.Region))
	}

	if user.Age != nil && (*user.Age < 18 || *user.Age > 120) {
		errors = append(errors, errs.NewFieldError("age", "должен быть от 18 до 120", *user.Age))
	}

	if user.Gender != nil {
		gender := strings.ToUpper(string(*user.Gender))
		if gender != "MALE" && gender != "FEMALE" {
			errors = append(errors, errs.NewFieldError("gender", "должен быть MALE или FEMALE", *user.Gender))
		}
	}

	if user.Role == "" {
		errors = append(errors, errs.NewFieldError("role", "обязательное поле"))
	} else {
		role := strings.ToLower(string(user.Role))
		validRoles := map[string]bool{
			"admin":        true,
			"viewer":       true,
			"experimenter": true,
			"approver":     true,
		}

		if !validRoles[role] {
			errors = append(errors, errs.NewFieldError("role", "должен быть admin, viewer, experimenter или approver", user.Role))
		}
	}

	return errors
}

func ValidateUpdateUserFull(user *dto.User) []errs.FieldError {
	var errors []errs.FieldError

	if user.FullName == "" {
		errors = append(errors, errs.NewFieldError("fullName", "обязательное поле"))
	} else if len(user.FullName) < 2 || len(user.FullName) > 200 {
		errors = append(errors, errs.NewFieldError("fullName", "должно быть от 2 до 200 символов", user.FullName))
	}

	if user.Region != nil && len(*user.Region) > 32 {
		errors = append(errors, errs.NewFieldError("region", "не более 32 символов", *user.Region))
	}

	if user.Age != nil && (*user.Age < 18 || *user.Age > 120) {
		errors = append(errors, errs.NewFieldError("age", "должен быть от 18 до 120", *user.Age))
	}

	if user.Gender != nil {
		gender := strings.ToUpper(string(*user.Gender))
		if gender != "MALE" && gender != "FEMALE" {
			errors = append(errors, errs.NewFieldError("gender", "должен быть MALE или FEMALE", *user.Gender))
		}
	}

	if user.Role != "" {
		role := strings.ToLower(string(user.Role))
		validRoles := map[string]bool{
			"admin":        true,
			"viewer":       true,
			"experimenter": true,
			"approver":     true,
		}

		if !validRoles[role] {
			errors = append(errors, errs.NewFieldError("role", "должен быть admin, viewer, experimenter или approver", user.Role))
		}
	}

	return errors
}

func ValidateLogin(email, password string) []errs.FieldError {
	var errors []errs.FieldError

	if email == "" {
		errors = append(errors, errs.NewFieldError("email", "обязательное поле"))
	} else {
		if len(email) > 254 {
			errors = append(errors, errs.NewFieldError("email", "не более 254 символов", email))
		}
		emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		if !regexp.MustCompile(emailRegex).MatchString(email) {
			errors = append(errors, errs.NewFieldError("email", "неверный формат email", email))
		}
	}

	if password == "" {
		errors = append(errors, errs.NewFieldError("password", "обязательное поле"))
	}

	return errors
}
