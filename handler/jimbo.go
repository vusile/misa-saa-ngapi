package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vusile/misa-saa-ngapi/model"
	"github.com/vusile/misa-saa-ngapi/repository"
)

type Jimbo struct {
	Repo *repository.JimboRepo
}

func (d *Jimbo) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name       string `json:"name"`
		IsJimboKuu bool   `json:"is_jimbo_kuu"`
		CountryID  uint   `json:"country_id"`
		ChurchID   uint   `json:"church_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("failed to decode:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jimbo := model.Jimbo{
		Name:       body.Name,
		IsJimboKuu: body.IsJimboKuu,
		ChurchID:   body.ChurchID,
		CountryID:  body.CountryID,
	}

	err := d.Repo.Insert(r.Context(), jimbo)

	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(jimbo)

	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (d *Jimbo) List(w http.ResponseWriter, r *http.Request) {

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
		Items []model.Jimbo `json:"items"`
		Page  int           `json:"page,omitempty"`
	}

	response.Items = res.Majimbo
	response.Page = res.Page

	data, err := json.Marshal(response)

	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (d *Jimbo) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	jimboID, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jimbo, err := d.Repo.FindByID(r.Context(), jimboID)

	if errors.Is(err, repository.ErrorNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(jimbo); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *Jimbo) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name       string `json:"name"`
		IsJimboKuu bool   `json:"is_jimbo_kuu"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("failed to decode:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	jimboID, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jimbo, err := d.Repo.FindByID(r.Context(), jimboID)

	if errors.Is(err, repository.ErrorNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jimbo.Name = body.Name
	jimbo.IsJimboKuu = body.IsJimboKuu

	err = d.Repo.Update(r.Context(), jimbo)

	if err != nil {
		fmt.Println("failed to update:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(jimbo); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *Jimbo) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	jimboID, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = d.Repo.DeleteByID(r.Context(), jimboID)

	if errors.Is(err, repository.ErrorNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
