package controllers

import (
	"github.com/laterius/service_architecture_hw3/app/models"
	"github.com/laterius/service_architecture_hw3/app/views"
)

// Users Struct for holding Users variables
type Users struct {
	NewView   *views.View
	us        models.UserService
	LoginView *views.View
}

// SignupForm for handling the metdata of the user
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// LoginForm Struct data for /login endpoint
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
