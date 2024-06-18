package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/louismomo66/logger/helpers"
	"github.com/louismomo66/logger/midelware"
	"github.com/louismomo66/logger/models"
	"github.com/louismomo66/logger/utils"
	"gorm.io/gorm"
)

type UserController struct {
	Repo models.UserRepository
	OTPManager *helpers.OTPManager
}


func NewUserController(repo models.UserRepository, otpManager *helpers.OTPManager) *UserController {
    return &UserController{
        Repo:       repo,
        OTPManager: otpManager,
    }
}


func (u UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		Err := utils.SetError(err, "Error in reading body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Err)
		return
	}

	email, err := u.Repo.GetUserEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		Err := utils.SetError(err, "Error checking email")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Err)
		return
	}
	if email != "" {
		Err := utils.SetError(nil, "Email already in use")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(Err)
		return
	}

	user.Password, err = midelware.GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}

	if err = u.Repo.CreateUser(&user); err != nil {
		Err := utils.SetError(err, "Failed to create User")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}



func (u UserController) Login(w http.ResponseWriter, r *http.Request) {
	var authDetails models.Authentication
	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		Err := utils.SetError(err, "Error reading payload")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Err)
		return
	}

	// Fetch user by email
	authUser, err := u.Repo.GetUserByEmail(authDetails.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Err := utils.SetError(nil, "Please insert correct email")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Err)
			return
		}
		Err := utils.SetError(err, "Error checking email")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Err)
		return
	}

	// Check password
	check := midelware.CheckPasswordHash(authDetails.Password, authUser.Password)
	if !check {
		Err := utils.SetError(nil, "Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Err)
		return
	}

	// Generate JWT token
	validToken, err := midelware.GenerateJWT(authUser.Email, authUser.Role)
	if err != nil {
		Err := utils.SetError(err, "Failed to generate token")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Err)
		return
	}

	// Prepare and send the token
	var token models.Token
	token.Email = authUser.Email
	token.Role = authUser.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

// In your UserController
func (u UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.Repo.GetAllUsers()
	if err != nil {
		Err := utils.SetError(err, "Failed to fetch users")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}










func AdminIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}
	w.Write([]byte("Welcome, Admin."))
}

func UserIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "user" {
		w.Write([]byte("Not Authorized."))
		return
	}
	w.Write([]byte("Welcome, User."))
}