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

type Parokia struct {
	Repo *repository.ParokiaRepo
}

func (d *Parokia) Create(w http.ResponseWriter, r *http.Request) {
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

	parokia := model.Parokia{
		Name:      body.Name,
		JimboID:   body.JimboID,
		IsKigango: body.IsKigango,
	}

	err := d.Repo.Insert(r.Context(), parokia)

	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(parokia)

	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (d *Parokia) List(w http.ResponseWriter, r *http.Request) {

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
		Items []model.Parokia `json:"items"`
		Page  int             `json:"page,omitempty"`
	}

	response.Items = res.Parokia
	response.Page = res.Page

	data, err := json.Marshal(response)

	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
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
