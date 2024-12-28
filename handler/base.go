package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"github.com/vusile/misa-saa-ngapi/model"

	"gorm.io/gorm"
)

var ErrorAuth = errors.New("Unauthorized")

var LoopFuncMap = template.FuncMap{
	"loop": func(from, to int) []int {
		var Items []int
		for i := from; i <= to; i++ {
			Items = append(Items, i)
		}
		return Items
	},
}

var ModulusFuncMap = template.FuncMap{
	"mod": func(i, j int) bool {
		{
			return i%j == 0
		}
	},
}

type Repo struct {
	Client *gorm.DB
}

var DaysOfTheWeek = map[string]string{
	"1": "Jumapili",
	"2": "Jumatatu",
	"3": "Jumanne",
	"4": "Jumatano",
	"5": "Alhamisi",
	"6": "Ijumaa",
	"7": "Jumamosi",
}

// type Base struct {
// 	Repo *repository.JimboRepo
// }

func GetValidationMessage(err error) []string {
	var errorMessages []string
	for _, err := range err.(validator.ValidationErrors) {
		message := GetFieldErrorTranslations(err.Tag(), err.Field(), err.Param())
		errorMessages = append(errorMessages, message)
	}

	return errorMessages
}

func GetDbValidationMessages(err error) {

}

func GetFieldNameTranslations(fieldName string) string {
	switch fieldName {
	case "Name":
		return "Jina lako"

	case "Phone":
		return "Namba yako ya simu"

	case "Password":
		return "Neno siri"

	case "ConfirmPassword":
		return "Rudia Neno siri"

	case "Code":
		return "Namba ya Uhakiki"

	case "HudumaID":
		return "Huduma"

	case "LanguageID":
		return "Lugha"

	case "WeekDayID":
		return "Siku ya Juma"

	case "StartTime":
		return "Muda wa Kuanza"

	case "EndTime":
		return "Muda wa Kumaliza"
	}

	return fieldName
}

func GetFieldErrorTranslations(tag string, fieldName string, value string) string {
	switch tag {
	case "required":
		return "Ni lazima uweke " + GetFieldNameTranslations(fieldName)

	case "eqfield":
		return "Hakikisha " + GetFieldNameTranslations(fieldName) + " na " + GetFieldNameTranslations(value) + " zifanane"

	case "gte":
		return "Ni lazima " + GetFieldNameTranslations(fieldName) + " iwe na herufi " + value + " au zaidi"

	case "et_lt_st":
		return "Ni lazima muda wa kumaliza ibada uwe baada ya kuanza ibada"

	case "e164":
		return "Tafadhali hakikisha " + GetFieldNameTranslations(fieldName) + " ni sahihi"
	}

	return tag
}

func StringToInt(value string) uint64 {
	const base = 10
	const bitSize = 64

	intVal, _ := strconv.ParseUint(value, base, bitSize)

	return intVal
}

func Authorize(r *http.Request, client *gorm.DB) error {

	st, err := r.Cookie("session_token")
	user := model.User{}

	if err != nil || st.Value == "" {
		return ErrorAuth
	} else {
		csrf, err := r.Cookie("csrf_token")

		if csrf.Value != "" && err == nil {
			client.Where("session_token = ?", st.Value).First(&user)

			if csrf.Value != user.CsrfToken {
				return ErrorAuth
			}
		} else {
			return ErrorAuth
		}
	}

	return nil
}

func GetLoggedInUser(r *http.Request, client *gorm.DB) model.User {

	st, err := r.Cookie("session_token")
	user := model.User{}

	if err != nil || st.Value == "" {
		return user
	} else {
		client.Where("session_token = ?", st.Value).First(&user)
		return user
	}
}

func CreateSlug(input string) string {
	// Remove special characters
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	processedString := reg.ReplaceAllString(input, " ")

	// Remove leading and trailing spaces
	processedString = strings.TrimSpace(processedString)

	// Replace spaces with dashes
	slug := strings.ReplaceAll(processedString, " ", "-")

	// Convert to lowercase
	slug = strings.ToLower(slug)

	return slug
}

func SetupToken(r *http.Request) map[string]interface{} {
	token := make(map[string]interface{})
	token[csrf.TemplateTag] = csrf.TemplateField(r)
	return token
}
