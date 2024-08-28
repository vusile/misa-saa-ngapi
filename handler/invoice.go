package handler

import (
	"fmt"
	"net/http"
)

type Invoice struct {
}

func (i *Invoice) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Invoice")
}

func (i *Invoice) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List Invoice")
}

func (i *Invoice) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Invoice by ID")
}

func (i *Invoice) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Invoice by ID")
}

func (i *Invoice) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete Invoice by ID")
}
