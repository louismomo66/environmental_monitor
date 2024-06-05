package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/louismomo66/logger/models"
	"github.com/louismomo66/logger/utils"
)

type DeviceController struct {
	Repo models.DeviceRepository
}

func NewDeviceController(repo models.DeviceRepository) *DeviceController {
	return &DeviceController{Repo: repo}
}
func (d *DeviceController) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	readings, err := d.ParseReadings(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Validate the readings
	if err := readings.Validate(); err != nil {
		log.Printf("Validation failed: %v", err) // Log for debugging
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify the device exists
	if _, err := d.Repo.GetDeviceByIMEI(readings.IMEI); err != nil {
		http.Error(w, "Device with this IMEI doesn't exist: "+err.Error(), http.StatusNotFound)
		return
	}

	// Attempt to create the readings in the database
	if err := d.Repo.CreateReadings(&readings); err != nil {
		http.Error(w, "Failed to create readings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, readings)
}

func (d *DeviceController) ParseReadings(r *http.Request) (models.Readings, error) {
	// Extract query parameters
	imei := r.URL.Query().Get("imei")
	temperatureV := r.URL.Query().Get("t")
	humidityV := r.URL.Query().Get("h")
	soilMoistureV := r.URL.Query().Get("s")
	countV := r.URL.Query().Get("count") // Getting the count parameter

	temperature, err := strconv.ParseFloat(temperatureV, 64)
	if err != nil {

		return models.Readings{}, fmt.Errorf("invalid temperature value")
	}
	humidity, err := strconv.ParseFloat(humidityV, 64)
	if err != nil {

		return models.Readings{}, fmt.Errorf("invalid humidity value")
	}
	soilMoisture, err := strconv.ParseFloat(soilMoistureV, 64)
	if err != nil {

		return models.Readings{}, fmt.Errorf("invalid moisture value")
	}
	count, err := strconv.Atoi(countV) // Converting count from string to integer
	if err != nil {

		return models.Readings{}, fmt.Errorf("invalid count")
	}

	// Create readings struct
	return models.Readings{
		IMEI:         imei,
		Temperature:  temperature,
		Humidity:     humidity,
		SoilMoisture: soilMoisture,
		Count:        count, // Including count in the readings
	}, nil
}

func (d *DeviceController) GetReadings(w http.ResponseWriter, r *http.Request) {
	imei := r.URL.Query().Get("imei")
	if imei == "" {
		http.Error(w, "IMEI is required", http.StatusBadRequest)
		return
	}

	// Fetch readings from the repository using the provided IMEI
	data, err := d.Repo.GetReadingsByIMEI(imei)
	if err != nil {
		log.Printf("Error retrieving readings for IMEI %s: %v", imei, err)
		http.Error(w, "Failed to retrieve device readings", http.StatusInternalServerError)
		return
	}

	// If no readings are found, you might want to return a different status or message
	if len(data) == 0 {
		http.Error(w, "No readings found for this IMEI", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding readings to JSON: %v", err)
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}
