package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/louismomo66/logger/utils"
	"github.com/louismomo66/logger/midelware"
	"gorm.io/gorm"
)

var OTPStore = make(map[string]string)

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (u UserController) SentOTP(w http.ResponseWriter, r *http.Request){
	var request struct{
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil{
		Err := utils.SetError(err, "Error in reading body")
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
	otp := generateOTP()
	OTPStore[request.Email] = otp

	if err := utils.SendEmail(request.Email,"Your OTP Code", "Your OTP is: "+otp); err != nil{
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

func (u UserController) ResetPassword(w http.ResponseWriter, r *http.Request){
	var request struct {
		Email string `json:"email"`
		OTP string `json:"opt"`
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

	storedOtp, ok := OTPStore[request.Email]
	if !ok || storedOtp != request.OTP{
		Err := utils.SetError(nil,"Invalid OTP")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
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

	delete(OTPStore, request.Email) // clear the used OTP

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}