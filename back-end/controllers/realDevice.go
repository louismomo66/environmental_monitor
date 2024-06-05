package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/louismomo66/logger/models"
	"gorm.io/gorm"
)

func (d *DeviceController) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var device models.Device
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := d.Repo.GetDeviceByIMEI(device.IMEI); err == nil {
		http.Error(w, "Device with this Imei already exists", http.StatusBadRequest)
		return
	}

	if err := d.Repo.CreateDevicek(&device); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// var existingDevice models.Device
	// if err := d.DB.Where("imei = ?", device.IMEI).First(&existingDevice).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	http.Error(w, "Device with this Imei already exists", http.StatusBadRequest)
	// 	return
	// }
	// d.DB.Create(&device)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(device); err != nil {
		// Handle the error appropriately, perhaps logging it and sending an HTTP error response
		log.Printf("Error encoding data to JSON: %v", err)
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}

func (d *DeviceController) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imei := vars["imei"]
	log.Printf("Attempting to delete device with IMEI: %s", imei)

	err := d.Repo.DeleteDeviceByIMEI(imei)
	if err != nil {
		log.Printf("Error deleting device: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Device not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListDeviceIMEIs handles the HTTP request to list all device IMEIs.
func (d *DeviceController) ListDeviceIMEIs(w http.ResponseWriter, r *http.Request) {
	imeis, err := d.Repo.GetAllDeviceIMEIs()
	if err != nil {
		http.Error(w, "Failed to fetch device IMEIs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(imeis); err != nil {
		http.Error(w, "Failed to encode IMEIs: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// func (d *DeviceController) ListDevices(w http.ResponseWriter, r *http.Request) {
// 	devices, err := d.Repo.GetAllDevices()
// 	if err != nil {
// 		http.Error(w, "Failed to fetch devices: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

//		w.Header().Set("Content-Type", "application/json")
//		if err := json.NewEncoder(w).Encode(devices); err != nil {
//			http.Error(w, "Failed to encode devices: "+err.Error(), http.StatusInternalServerError)
//			return
//		}
//	}
func (c *DeviceController) ListDevices(w http.ResponseWriter, r *http.Request) {
	// Log that the request to list devices is received
	log.Println("Request received: ListDevices")

	// Fetch all devices from the repository
	devices, err := c.Repo.GetAllDevices()
	if err != nil {
		// Log the error
		log.Printf("Failed to fetch devices: %v", err)
		// Return internal server error with error message
		http.Error(w, "Failed to fetch devices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode devices into JSON format and write response
	if err := json.NewEncoder(w).Encode(devices); err != nil {
		// Log encoding error
		log.Printf("Failed to encode devices: %v", err)
		// Return internal server error with error message
		http.Error(w, "Failed to encode devices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log that the response is successfully sent
	log.Println("Response sent: ListDevices")
}
