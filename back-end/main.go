package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/louismomo66/logger/controllers"
	"github.com/louismomo66/logger/database"
	"github.com/louismomo66/logger/midelware"
	"github.com/louismomo66/logger/models"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Printf("Failed %v", err)
	}
	repo := models.NewGormRepository(db)
	controller := controllers.NewDeviceController(repo)
	r := mux.NewRouter()
	r.Use(midelware.EnableCORS)
	r.Handle("/devices", midelware.EnableCORS(http.HandlerFunc(controller.CreateDevice)))
	r.Handle("/devices/{imei}", midelware.EnableCORS(http.HandlerFunc(controller.DeleteDevice)))
	r.Handle("/device-imeis", midelware.EnableCORS(http.HandlerFunc(controller.ListDeviceIMEIs)))
	r.HandleFunc("/getreadings", controller.GetReadings).Methods("GET")
	r.HandleFunc("/update", controller.UpdateDevice).Methods("GET")
	r.HandleFunc("/list", controller.ListDevices).Methods("GET")
	r.HandleFunc("/device-imeis", controller.ListDeviceIMEIs).Methods("GET")
	r.HandleFunc("/devices", controller.CreateDevice).Methods("POST")
	r.HandleFunc("/devices/{imei}", controller.DeleteDevice).Methods("DELETE")
	if err := http.ListenAndServe(":9000", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
