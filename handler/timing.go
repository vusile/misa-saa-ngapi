package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vusile/misa-saa-ngapi/model"
	"github.com/vusile/misa-saa-ngapi/repository"
)

type Timing struct {
	Repo *repository.TimingRepo
}

func (d *Timing) Create(w http.ResponseWriter, r *http.Request) {
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

	st, error := time.Parse(time.TimeOnly, body.StartTime)
	st = st.AddDate(1, 0, 0)

	if error != nil {
		fmt.Println("failed to convert time:", error)
	} else {
		fmt.Println("start time:", st)
	}

	et, error := time.Parse(time.TimeOnly, body.EndTime)
	et = et.AddDate(1, 0, 0)

	if error != nil {
		fmt.Println("failed to convert time:", error)
	} else {
		fmt.Println("end time:", et)
	}

	timing := model.Timing{
		ParokiaID:  body.ParokiaID,
		StartTime:  &st,
		EndTime:    &et,
		Details:    body.Details,
		LanguageID: body.LanguageID,
		WeekDayID:  body.WeekDayID,
	}

	err := d.Repo.Insert(r.Context(), timing)

	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(timing)

	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (d *Timing) List(w http.ResponseWriter, r *http.Request) {

	pageNum, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	res, err := d.Repo.FindAll(r.Context(), repository.FindAllPage{
		PageNum: pageNum,
		Size:    size,
	})

	if err != nil {
		fmt.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var response struct {
		Items []model.Timing `json:"items"`
		Page  int            `json:"page,omitempty"`
	}

	response.Items = res.Timings
	response.Page = res.Page

	data, err := json.Marshal(response)

	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
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
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
