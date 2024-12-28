package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"github.com/vusile/misa-saa-ngapi/model"
	"github.com/vusile/misa-saa-ngapi/repository"
)

type Timing struct {
	Repo *repository.TimingRepo
}

type TimingValidator struct {
	ParokiaID  uint       `json:"parokia_id" validate:"required"`
	StartTime  string     `json:"start_time" validate:"required"`
	TStartTime *time.Time `json:"t_start_time"`
	TEndTime   *time.Time `json:"t_end_time"`
	LanguageID uint       `json:"language_id" validate:"required"`
	WeekDayID  string     `json:"days_of_the_week" validate:"required"`
	HudumaID   uint       `json:"huduma_id" validate:"required"`
}

func (d *Timing) Create(w http.ResponseWriter, r *http.Request) {
	// client, err := elasticsearch.NewDefaultClient()
	// CreateESIndex(client)

	var data struct {
		ErrorMessages []string
		Title         string
		Token         map[string]interface{}
	}

	st, error := time.Parse(time.TimeOnly, r.PostFormValue("start_time"))
	st = st.AddDate(1, 0, 0)

	if error != nil {
		fmt.Println("failed to convert time:", error)
	} else {
		fmt.Println("start time:", st)
	}

	et, error := time.Parse(time.TimeOnly, r.PostFormValue("end_time"))
	et = et.AddDate(1, 0, 0)

	if error != nil {
		fmt.Println("failed to convert time:", error)
	} else {
		fmt.Println("end time:", et)
	}

	timingValidator := &TimingValidator{}

	timingValidator.ParokiaID = uint(StringToInt(r.PostFormValue("parokia_id")))
	timingValidator.StartTime = r.PostFormValue("start_time")
	timingValidator.TStartTime = &st
	timingValidator.TEndTime = &et
	timingValidator.LanguageID = uint(StringToInt(r.PostFormValue("language_id")))
	timingValidator.HudumaID = uint(StringToInt(r.PostFormValue("huduma_id")))
	timingValidator.WeekDayID = r.PostFormValue("days_of_the_week")

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterStructValidation(TimeValidation, TimingValidator{})
	err := validate.Struct(timingValidator)

	// if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
	// 	fmt.Println("failed to decode:", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	if err != nil {
		data.ErrorMessages = GetValidationMessage(err)
		// data.Token = make(map[string]interface{})
		// data.Token[csrf.TemplateTag] = csrf.TemplateField(r)

		w.WriteHeader(http.StatusBadRequest)

		tmpl := template.Must(template.New("index.html").
			Funcs(LoopFuncMap).Funcs(ModulusFuncMap).ParseFiles(
			"/go/src/app/views/backend/timings/index.html",
			"/go/src/app/views/backend/template.html"))

		tmpl.ExecuteTemplate(w, "validation-errors", data)
	} else {

		r.ParseForm()

		var timings []model.Timing

		isPublicHoliday := false

		if r.PostFormValue("is_public_holiday") == "true" {
			isPublicHoliday = true
		}

		timing := model.Timing{
			ParokiaID:       uint(StringToInt(r.PostFormValue("parokia_id"))),
			StartTime:       &st,
			EndTime:         &et,
			LanguageID:      uint(StringToInt(r.PostFormValue("language_id"))),
			HudumaID:        uint(StringToInt(r.PostFormValue("huduma_id"))),
			IsPublicHoliday: isPublicHoliday,
			UserID:          GetLoggedInUser(r, d.Repo.Client).ID,
			InsertGroupBy:   time.Now().Unix(),
		}

		for _, value := range r.Form["days_of_the_week"] {
			timing.WeekDayID = uint(StringToInt(value))
			timings = append(timings, timing)
		}

		err := d.Repo.Insert(r.Context(), timings)

		if err != nil {
			// fmt.Println("failed to insert:", err)
			// w.WriteHeader(http.StatusInternalServerError)
			// return
			w.WriteHeader(http.StatusBadRequest)
			data.ErrorMessages = append(data.ErrorMessages, err.Error())
			tmpl := template.Must(template.New("index.html").
				Funcs(LoopFuncMap).Funcs(ModulusFuncMap).ParseFiles(
				"/go/src/app/views/backend/timings/index.html",
				"/go/src/app/views/backend/template.html"))

			tmpl.ExecuteTemplate(w, "validation-errors", data)
		} else {
			w.Header().Add("Hx-Redirect", "/timings/timingform/"+r.PostFormValue("parokia_id"))
		}

		// res, err := json.Marshal(timing)

		// if err != nil {
		// 	fmt.Println("failed to marshal:", err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// }

		// w.Write(res)
		// w.WriteHeader(http.StatusCreated)
	}
}

func (d *Timing) List(w http.ResponseWriter, r *http.Request) {

	pageNum, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	idParam := chi.URLParam(r, "id")
	parokiaID := StringToInt(idParam)

	res, err := d.Repo.FindByParishId(r.Context(), parokiaID, repository.FindAllPage{
		PageNum: pageNum,
		Size:    size,
	})

	if err != nil {
		fmt.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var hudumas []model.Huduma
	var languages []model.Language
	var parokia model.Parokia
	d.Repo.Client.Find(&hudumas)
	d.Repo.Client.Find(&parokia, "id = ?", parokiaID)
	d.Repo.Client.Find(&languages)

	var response struct {
		Items         []model.Timing `json:"items"`
		Page          int            `json:"page,omitempty"`
		Title         string
		Parokia       model.Parokia
		Hudumas       []model.Huduma
		Languages     []model.Language
		DaysOfTheWeek map[string]string
		Token         map[string]interface{}
	}

	response.Items = res.Timings
	response.Page = res.Page
	response.Parokia = parokia
	response.Title = "Parokia ya " + parokia.Name
	response.Hudumas = hudumas
	response.Languages = languages
	response.DaysOfTheWeek = DaysOfTheWeek
	response.Token = make(map[string]interface{})
	response.Token[csrf.TemplateTag] = csrf.TemplateField(r)

	// data, err := json.Marshal(response)

	// if err != nil {
	// 	fmt.Println("failed to marshal:", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// w.Write(data)

	tmpl := template.Must(
		template.New("index.html").Funcs(LoopFuncMap).Funcs(ModulusFuncMap).
			ParseFiles(
				"/go/src/app/views/backend/timings/index.html",
				"/go/src/app/views/backend/template.html"))
	tmpl.Execute(w, response)
}

func (d *Timing) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	timingID, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timing, err := d.Repo.FindByID(r.Context(), timingID)

	if errors.Is(err, repository.ErrorTimingNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(timing); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *Timing) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ParokiaID  uint   `json:"parokia_id"`
		StartTime  string `json:"start_time"`
		EndTime    string `json:"end_time"`
		Details    string `json:"details"`
		LanguageID uint   `json:"language_id"`
		WeekDayID  uint   `json:"week_day_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("failed to decode:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	timingID, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timing, err := d.Repo.FindByID(r.Context(), timingID)

	if errors.Is(err, repository.ErrorTimingNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	st, error := time.Parse(time.TimeOnly, body.StartTime)
	st = st.AddDate(1, 0, 0)

	if error != nil {
		st = st.Add(time.Duration(time.Now().Year()))
	}

	et, error := time.Parse(time.TimeOnly, body.EndTime)
	et = et.AddDate(1, 0, 0)

	if error != nil {
		et = et.Add(time.Duration(time.Now().Year()))
	}

	timing.ParokiaID = body.ParokiaID
	timing.StartTime = &st
	timing.EndTime = &et
	timing.Details = body.Details
	timing.LanguageID = body.LanguageID
	timing.WeekDayID = body.WeekDayID

	err = d.Repo.Update(r.Context(), timing)

	if err != nil {
		fmt.Println("failed to update:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(timing); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *Timing) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	parokiaIdParam := chi.URLParam(r, "parokiaId")

	const base = 10
	const bitSize = 64

	timingID, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = d.Repo.DeleteByID(r.Context(), timingID)

	if errors.Is(err, repository.ErrorTimingNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Hx-Redirect", "/timings/timingform/"+parokiaIdParam)
}

func TimeValidation(sl validator.StructLevel) {

	timing := sl.Current().Interface().(TimingValidator)

	if timing.TEndTime.Before(*timing.TStartTime) {
		sl.ReportError(timing.TEndTime, "t_end_time", "EndTime", "et_lt_st", "")
	}

	// plus can do more, even with different tag than "fnameorlname"
}
