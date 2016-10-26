package controller

import (
	"net/http"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc UserController) Index(w http.ResponseWriter, r *http.Request) {}
