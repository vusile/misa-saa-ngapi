package handler

import (
	"fmt"
	"net/http"
)

type User struct {
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Promise")
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

func (p *User) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete Promise by ID")
}
