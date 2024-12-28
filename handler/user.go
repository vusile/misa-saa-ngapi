package handler

import (
	crand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	// "strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"github.com/vusile/misa-saa-ngapi/model"
	"github.com/vusile/misa-saa-ngapi/repository"
)

type User struct {
	Repo *repository.UserRepo
}

func (u *User) LoginForm(w http.ResponseWriter, r *http.Request) {
	var response struct {
		Title         string
		Countries     []model.Country `json:"country"`
		Token         map[string]interface{}
		OtherMessages []string
	}

	response.Token = make(map[string]interface{})
	response.Token[csrf.TemplateTag] = csrf.TemplateField(r)
	response.Title = "Login"
	if r.URL.Query().Get("fromRegister") == "1" {
		response.OtherMessages = append(response.OtherMessages, "Tayari kuna akaunti yenye namba yako. Login sasa")
	}

	// data, err := json.Marshal(response)

	// if err != nil {
	// 	fmt.Println("failed to marshal:", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	tmpl := template.Must(template.ParseFiles(
		"/go/src/app/views/frontend/auth/login.html",
		"/go/src/app/views/frontend/template.html"))

	tmpl.Execute(w, response)
}

func (u *User) CodeForm(w http.ResponseWriter, r *http.Request) {
	var response struct {
		Title  string
		UserId string
		Token  map[string]interface{}
	}
	response.Token = make(map[string]interface{})
	response.Token[csrf.TemplateTag] = csrf.TemplateField(r)
	response.Title = "Hakiki Akaunti"

	response.UserId = chi.URLParam(r, "id")

	// data, err := json.Marshal(response)

	// if err != nil {
	// 	fmt.Println("failed to marshal:", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	tmpl := template.Must(template.ParseFiles(
		"/go/src/app/views/frontend/auth/code.html",
		"/go/src/app/views/frontend/template.html"))

	tmpl.Execute(w, response)
}

func (u *User) RegistrationForm(w http.ResponseWriter, r *http.Request) {
	var response struct {
		Title         string
		Countries     []model.Country `json:"country"`
		Token         map[string]interface{}
		OtherMessages []string
	}

	if r.URL.Query().Get("fromLogin") == "1" {
		response.OtherMessages = append(response.OtherMessages, "Hakuna akaunti yenye namba yako. Jisajili sasa")
	}

	response.Title = "Jisajili"
	response.Token = SetupToken(r)

	// data, err := json.Marshal(response)

	// if err != nil {
	// 	fmt.Println("failed to marshal:", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	tmpl := template.Must(template.ParseFiles(
		"/go/src/app/views/frontend/auth/register.html",
		"/go/src/app/views/frontend/template.html"))

	tmpl.Execute(w, response)
}

func (u *User) Logout(w http.ResponseWriter, r *http.Request) {
	if err := Authorize(r, u.Repo.Client); err != nil {
		http.Redirect(w, r, "/users/login", http.StatusTemporaryRedirect)
	}

	st, _ := r.Cookie("session_token")
	user := model.User{}
	u.Repo.Client.Where("session_token = ?", st.Value).First(&user)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
		Path:     "/",
	})

	user.SessionToken = ""
	user.CsrfToken = ""
	u.Repo.Update(r.Context(), user)
	http.Redirect(w, r, "/users/login", http.StatusTemporaryRedirect)
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {

	var data struct {
		ErrorMessages []string
		Title         string
		Token         map[string]interface{}
	}

	data.Token = SetupToken(r)

	var login struct {
		Phone    string `json:"phone" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	login.Password = r.PostFormValue("password")
	var country model.Country
	u.Repo.Client.First(&country, r.PostFormValue("country_id"))

	if len(r.PostFormValue("phone")) > 0 {
		login.Phone = country.CountryCode + r.PostFormValue("phone")[1:]
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(login)

	if err != nil {
		data.ErrorMessages = GetValidationMessage(err)

		w.WriteHeader(http.StatusBadRequest)

		tmpl := template.Must(template.ParseFiles(
			"/go/src/app/views/frontend/auth/login.html",
			"/go/src/app/views/frontend/template.html"))

		tmpl.ExecuteTemplate(w, "validation-errors", data)
	} else {
		var user model.User

		if err := u.Repo.Client.Where("phone = ?", login.Phone).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			w.Header().Add("Hx-Redirect", "/users/register?fromLogin=1")
		} else {
			if CheckPasswordHash(login.Password, user.Password) {
				sessionToken := GenerateToken(32)
				csrfToken := GenerateToken(32)
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    sessionToken,
					Expires:  time.Now().Add(2 * time.Hour),
					HttpOnly: true,
					Path:     "/",
				})

				http.SetCookie(w, &http.Cookie{
					Name:     "csrf_token",
					Value:    csrfToken,
					Expires:  time.Now().Add(2 * time.Hour),
					HttpOnly: false,
					Path:     "/",
				})

				user.SessionToken = sessionToken
				user.CsrfToken = csrfToken

				u.Repo.Update(r.Context(), user)

				//Todo: Redirect to admin page

			} else {
				w.WriteHeader(http.StatusBadRequest)
				data.ErrorMessages = append(data.ErrorMessages, "Tafadhali tumia namba na password sahihi")

				tmpl := template.Must(template.ParseFiles(
					"/go/src/app/views/frontend/auth/login.html",
					"/go/src/app/views/frontend/template.html"))

				tmpl.ExecuteTemplate(w, "validation-errors", data)
			}
		}
	}
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {

	var data struct {
		ErrorMessages []string
		UserId        uint
		Title         string
	}

	data.Title = "Jisajili"

	var register struct {
		Name            string `json:"name" validate:"required"`
		Phone           string `json:"phone" validate:"required,e164"`
		Password        string `json:"password" validate:"required,eqfield=ConfirmPassword,gte=10,alphanum"`
		ConfirmPassword string `json:"confirm_password" validate:"required"`
		CountryID       string `json:"country_id" validate:"required"`
		ChurchID        string `json:"church_id" validate:"required"`
	}

	// if err := json.NewDecoder(r.Body).Decode(&Register); err != nil {
	// 	fmt.Println("failed to decode:", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	register.Name = r.PostFormValue("name")
	register.Password = r.PostFormValue("password")
	register.ConfirmPassword = r.PostFormValue("confirm_password")
	register.CountryID = r.PostFormValue("country_id")
	register.ChurchID = r.PostFormValue("church_id")

	var country model.Country
	u.Repo.Client.First(&country, r.PostFormValue("country_id"))

	if len(r.PostFormValue("phone")) > 0 {
		register.Phone = country.CountryCode + r.PostFormValue("phone")[1:]
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(register)

	if err != nil {
		data.ErrorMessages = GetValidationMessage(err)

		w.WriteHeader(http.StatusBadRequest)

		tmpl := template.Must(template.ParseFiles(
			"/go/src/app/views/frontend/auth/register.html",
			"/go/src/app/views/frontend/template.html"))

		tmpl.ExecuteTemplate(w, "validation-errors", data)
	} else {
		var user model.User

		if err := u.Repo.Client.Where("phone = ?", register.Phone).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			churchId := StringToInt(r.PostFormValue("church_id"))
			countryId := StringToInt(r.PostFormValue("country_id"))
			user.Name = register.Name
			user.Phone = register.Phone
			user.ChurchID = uint(churchId)
			user.CountryID = uint(countryId)
			user.Code = GenerateConfirmationCode()
			password, _ := HashPassword(register.Password)
			user.Password = password

			userId, err := u.Repo.Insert(r.Context(), user)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				data.ErrorMessages = append(data.ErrorMessages, err.Error())
				tmpl := template.Must(template.ParseFiles(
					"/go/src/app/views/frontend/auth/register.html",
					"/go/src/app/views/frontend/template.html"))

				tmpl.ExecuteTemplate(w, "validation-errors", data)
				// fmt.Println("failed to insert:", err)
				// w.WriteHeader(http.StatusInternalServerError)
				// return
			} else {
				if SendConfirmationCode(user.Code) {
					w.Header().Add("Hx-Redirect", "/users/confirm-account/"+strconv.FormatUint(uint64(userId), 10))
				} else {
					w.Header().Add("Hx-Redirect", "/users/login")
				}
			}
		} else {
			w.Header().Add("Hx-Redirect", "/users/login?fromRegister=1")
		}

	}
}

func (u *User) ConfirmAccount(w http.ResponseWriter, r *http.Request) {

	var data struct {
		ErrorMessages []string
		Title         string
	}

	code := StringToInt(r.PostFormValue("code"))

	var confirm struct {
		Code uint64 `json:"code" validate:"required"`
	}

	confirm.Code = code

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(confirm)

	if err != nil {
		data.ErrorMessages = GetValidationMessage(err)

		w.WriteHeader(http.StatusBadRequest)

		tmpl := template.Must(template.ParseFiles(
			"/go/src/app/views/frontend/auth/code.html",
			"/go/src/app/views/frontend/template.html"))

		tmpl.ExecuteTemplate(w, "validation-errors", data)
	} else {
		if len(r.PostFormValue("user_id")) > 0 {
			var user model.User
			u.Repo.Client.First(&user, r.PostFormValue("user_id"))

			if user.Code == int(code) {
				activationTime := time.Now()
				user.ActivatedAt = &activationTime
				user.Code = 0

				u.Repo.Update(r.Context(), user)

				w.Header().Add("Hx-Redirect", "/users/login")
			} else {
				data.ErrorMessages = append(data.ErrorMessages, "Tafadhali hakikisha umeweka namba uliyotumiwa WhatsApp")
				tmpl := template.Must(template.ParseFiles(
					"/go/src/app/views/frontend/auth/code.html",
					"/go/src/app/views/frontend/template.html"))

				tmpl.ExecuteTemplate(w, "validation-errors", data)
			}
		} else {
			//User Id isn't there. But confirming accounts is not mission critical for us
			w.Header().Add("Hx-Redirect", "/users/login")
		}
	}
}

func (u *User) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List Promises")
}

func (u *User) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Promise by ID")
}

func (u *User) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Promise by ID")
}

func (u *User) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete Promise by ID")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)

	if _, err := crand.Read(bytes); err != nil {
		fmt.Println(err)
	}

	return base64.URLEncoding.EncodeToString(bytes)
}

func GenerateConfirmationCode() int {
	return rand.IntN(10000)

}

func SendConfirmationCode(code int) bool {
	messageSent := true
	//Send code for verification
	//Check if has WhatsApp
	//set messageSent to true if has whatsapp, false if not
	//If yes, make it confirmed.

	return messageSent
}
