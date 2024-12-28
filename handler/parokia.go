package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"github.com/vusile/misa-saa-ngapi/model"
	"github.com/vusile/misa-saa-ngapi/repository"
)

type Parokia struct {
	Repo *repository.ParokiaRepo
}

func (d *Parokia) Create(w http.ResponseWriter, r *http.Request) {

	var response struct {
		ErrorMessages []string
		Title         string
	}
	// var body struct {
	// 	Name      string `json:"name"`
	// 	IsKigango bool   `json:"is_kigango"`
	// 	JimboID   uint   `json:"jimbo_id"`
	// }

	// if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
	// 	fmt.Println("failed to decode:", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	var parokiaValidation struct {
		Name     string `json:"name" validate:"required"`
		JimboID  string `json:"jimbo_id" validate:"required"`
		Location string `json:"location" validate:"required"`
	}

	parokiaValidation.JimboID = r.PostFormValue("jimbo_id")
	parokiaValidation.Name = r.PostFormValue("name")
	parokiaValidation.Location = r.PostFormValue("location")

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(parokiaValidation)

	if err != nil {
		response.ErrorMessages = GetValidationMessage(err)

		w.WriteHeader(http.StatusBadRequest)

		tmpl := template.Must(template.ParseFiles(
			"/go/src/app/views/backend/parokia/index.html",
			"/go/src/app/views/backend/template.html"))
		tmpl.Execute(w, response)

		tmpl.ExecuteTemplate(w, "validation-errors", response)
	} else {
		jimboID := StringToInt(r.PostFormValue("jimbo_id"))
		parokia := model.Parokia{
			Name:      r.PostFormValue("name"),
			JimboID:   uint(jimboID),
			IsKigango: false,
			UserID:    GetLoggedInUser(r, d.Repo.Client).ID,
			Location:  r.PostFormValue("location"),
			Slug:      CreateSlug(r.PostFormValue("name")),
		}

		err := d.Repo.Insert(r.Context(), &parokia)

		if err != nil {
			fmt.Println("failed to insert:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		parokia.AddToIndex(d.Repo.ESClient)

		tmpl := template.Must(template.ParseFiles(
			"/go/src/app/views/backend/parokia/index.html",
			"/go/src/app/views/backend/template.html"))

		tmpl.ExecuteTemplate(w, "parokia-list-element", parokia)

		// res, err := json.Marshal(parokia)

		// if err != nil {
		// 	fmt.Println("failed to marshal:", err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// }

		// w.Write(res)
		// w.WriteHeader(http.StatusCreated)
	}

}

func (d *Parokia) List(w http.ResponseWriter, r *http.Request) {

	pageNum, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	res, err := d.Repo.FindByUser(r.Context(), uint64(GetLoggedInUser(r, d.Repo.Client).ID), repository.FindAllPage{
		PageNum: pageNum,
		Size:    size,
	})

	if err != nil {
		fmt.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var majimbo []model.Jimbo
	d.Repo.Client.Order("name").Find(&majimbo)

	var response struct {
		Items   []model.Parokia `json:"items"`
		Majimbo []model.Jimbo   `json:"majimbo"`
		Page    int             `json:"page,omitempty"`
		Title   string
		Token   map[string]interface{}
	}

	response.Token = make(map[string]interface{})
	response.Token[csrf.TemplateTag] = csrf.TemplateField(r)
	response.Items = res.Parokia
	response.Page = res.Page
	response.Majimbo = majimbo
	response.Title = "Parokia"

	// data, err := json.Marshal(response)

	// if err != nil {
	// 	fmt.Println("failed to marshal:", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// w.Write(data)
	tmpl := template.Must(template.ParseFiles(
		"/go/src/app/views/backend/parokia/index.html",
		"/go/src/app/views/backend/template.html"))
	tmpl.Execute(w, response)
}

func (d *Parokia) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	parokiaID, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parokia, err := d.Repo.FindByID(r.Context(), parokiaID)

	if errors.Is(err, repository.ErrorParokiaNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(parokia); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *Parokia) Detail(w http.ResponseWriter, r *http.Request) {
	var response struct {
		Parokia  model.Parokia
		Timings  map[uint][]model.Timing
		Huduma   []model.Huduma
		WeekDays []model.WeekDay
		Title    string
		Token    map[string]interface{}
	}
	response.Token = SetupToken(r)
	response.Timings = make(map[uint][]model.Timing)

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	parokiaID, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		fmt.Println("failed to find:", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {

		parokia, err := d.Repo.FindByID(r.Context(), parokiaID)

		if errors.Is(err, repository.ErrorParokiaNotExist) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			fmt.Println("failed to find by id:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var huduma []model.Huduma
		d.Repo.Client.Order("priority").Find(&huduma)

		var weekdays []model.WeekDay
		d.Repo.Client.Order("priority").Find(&weekdays)

		response.Parokia = parokia
		response.Timings = parokia.GenerateTimings()
		response.WeekDays = weekdays
		response.Huduma = huduma
		response.Title = "Muda wa Ibada na Huduma Mbalimbali Parokia ya " + parokia.Name

		tmpl := template.Must(template.ParseFiles(
			"/go/src/app/views/frontend/parokia/detail.html",
			"/go/src/app/views/frontend/template.html"))
		tmpl.Execute(w, response)
	}
}

func (d *Parokia) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name      string `json:"name"`
		IsKigango bool   `json:"is_kigango"`
		JimboID   uint   `json:"jimbo_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("failed to decode:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	parokiaID, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parokia, err := d.Repo.FindByID(r.Context(), parokiaID)

	if errors.Is(err, repository.ErrorParokiaNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	parokia.Name = body.Name
	parokia.IsKigango = body.IsKigango
	parokia.JimboID = body.JimboID

	err = d.Repo.Update(r.Context(), parokia)

	if err != nil {
		fmt.Println("failed to update:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(parokia); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *Parokia) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	parokiaID, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = d.Repo.DeleteByID(r.Context(), parokiaID)

	if errors.Is(err, repository.ErrorParokiaNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
