package controller

type UserController struct{}

func NewUserContoller() *UserController {
	return &UserContoller{}
}

func (uc UserController) Index(w http.ResponseWriter, r *http.Request) {}
