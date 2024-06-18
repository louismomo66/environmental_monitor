package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/louismomo66/logger/controllers"
	"github.com/louismomo66/logger/database"
	"github.com/louismomo66/logger/helpers"
	"github.com/louismomo66/logger/midelware"
	"github.com/louismomo66/logger/models"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	repo := models.NewGormRepository(db)
	userRepo := models.NewUserRepo(db)
	controller := controllers.NewDeviceController(repo)
	otpManager := helpers.NewOTPManager()
	userController := controllers.NewUserController( userRepo, otpManager)
	r := mux.NewRouter()

	// r.Use(midelware.EnableCORS)

	r.HandleFunc("/signup", userController.SignUp).Methods("POST")
	r.HandleFunc("/signin", userController.Login).Methods("POST")
	r.HandleFunc("/forgot-password", userController.SentOTP).Methods("POST")
	r.HandleFunc("/verify-otp", userController.VeryfyOTP).Methods("POST")
	r.HandleFunc("/reset-password", userController.ResetPassword).Methods("POST")
	r.HandleFunc("/admin", midelware.IsAuthorized(controllers.AdminIndex)).Methods("GET")
	r.HandleFunc("/user", midelware.IsAuthorized(controllers.UserIndex)).Methods("GET")
	r.HandleFunc("/getreadings", controller.GetReadings).Methods("GET")
	r.HandleFunc("/update", controller.UpdateDevice).Methods("GET")
	r.HandleFunc("/list", controller.ListDevices).Methods("GET")
	r.HandleFunc("/users", userController.GetAllUsers).Methods("GET")
	r.HandleFunc("/device-imeis", controller.ListDeviceIMEIs).Methods("GET")
	r.HandleFunc("/devices", controller.CreateDevice).Methods("POST")
	r.HandleFunc("/devices/{imei}", controller.DeleteDevice).Methods("DELETE")

	if err := http.ListenAndServe(":9000",
	handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)(r)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
