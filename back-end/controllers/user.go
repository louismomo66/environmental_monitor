package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/louismomo66/logger/midelware"
	"github.com/louismomo66/logger/models"
	"github.com/louismomo66/logger/utils"
)

type UserController struct {
	Repo models.UserRepository
}


func NewUserController(repo models.UserRepository) * UserController{
	return &UserController{Repo:repo}
}



func (u UserController) SignUp(w http.ResponseWriter,r *http.Request){
var user models.User 
err := json.NewDecoder(r.Body).Decode(&user)
if err != nil{
		Err := utils.SetError(err, "Error in reading body")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Err)
		return
}

 if email, err := u.Repo.GetUserEmail(user.Email);email != ""{
	Err := utils.SetError(err, "Email already in use")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Err)
	return
}
user.Password, err = midelware.GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}

	if err = u.Repo.CreateUser(&user);err != nil{
		Err := utils.SetError(err, "Failed to create User:")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Err)
		return	
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}




func (u UserController) Login(w http.ResponseWriter, r *http.Request){
	var authDetails models.Authentication
	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		Err :=utils.SetError(err, "Error reading payload")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Err)
		return
	}

	var authUser models.User
	if email, err := u.Repo.GetUserEmail(authUser.Email);email == ""{
		Err := utils.SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Err)
		return
	}

	check := midelware.CheckPasswordHash(authDetails.Password,authUser.Password)
	if !check {
		
		Err := utils.SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Err)
		return
	}


	validToken, err := midelware.GenerateJWT(authUser.Email, authUser.Role)
	if err != nil {
		
		Err := utils.SetError(err, "Failed to generate token")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Err)
		return
	}

	var token models.Token
	token.Email = authUser.Email
	token.Role = authUser.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
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