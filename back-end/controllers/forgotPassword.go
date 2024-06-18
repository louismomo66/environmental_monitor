package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/louismomo66/logger/midelware"
	"github.com/louismomo66/logger/utils"
	"gorm.io/gorm"
)


func (u UserController) SentOTP(w http.ResponseWriter, r *http.Request){
	var request struct{
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil{
		Err := utils.SetError(err, "INvalid request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Err)
		return
	}

	  _, err = u.Repo.GetUserEmail(request.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				Err := utils.SetError(nil, "Email not found")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(Err)
				return
			}
			Err := utils.SetError(err, "Error checking email")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Err)
			return
		}
	token, err := u.OTPManager.GenerateOTP(request.Email, 5*time.Minute)
	if err != nil {
		Err := utils.SetError(err,"Error generating OTP")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Err)
		return	
	}

	if err := utils.SendEmail(request.Email,"Your OTP Code", "Your OTP is: "+token); err != nil{
		Err := utils.SetError(err,"Failed to send OTP")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (u *UserController) VeryfyOTP(w http.ResponseWriter, r *http.Request){
	var request struct {
		Email string `json:"email"`
		OTP string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil{
	http.Error(w, "Invalid request", http.StatusBadRequest)
	return
	}
if u.OTPManager. VeryfyOTP(request.Email,request.OTP){
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message":"OTP verified"})
}else {
	http.Error(w,"Invalid or expired OTP", http.StatusUnauthorized)
}
}

func (u UserController) ResetPassword(w http.ResponseWriter, r *http.Request){
	var request struct {
		Email string `json:"email"`
		NewPassword string `json:"newPassword"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		Err := utils.SetError(err, "Error in reading body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Err)
		return
	}
	hashedPassword, err := midelware.GeneratehashPassword(request.NewPassword)
	if err != nil {
		Err := utils.SetError(err, "Failed to hash password")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Err)
		return
	}

	err = u.Repo.UpdatePasswordByEmail(request.Email, hashedPassword)
	if err != nil {
		Err := utils.SetError(err, "Failed to update password")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"success": "Password reset"})
}