package handler

import (
	"fmt"
	"net/http"
)

type Promise struct {
	// ID,
	// UserID,
	// Date,
	// InvoiceID,
	// Amount
}

func (p *Promise) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Promise")
}

func (p *Promise) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List Promises")
}

func (p *Promise) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Promise by ID")
}

func (p *Promise) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Promise by ID")
}

func (p *Promise) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete Promise by ID")
}
