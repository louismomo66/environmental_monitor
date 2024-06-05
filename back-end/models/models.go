package models

import (
	"errors"
	"fmt"
	"log"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	IMEI        string `json:"imei" gorm:"uniqueIndex"`
	PhoneNumber string `json:"phone_number"`
	DeviceType  string `json:"device_type"`
	// SerialNumber string `gorm:"not null;unique"`
	// Readings []Readings
}

type Readings struct {
	gorm.Model
	IMEI         string  `json:"imei"`
	Temperature  float64 `json:"temperature"`
	Humidity     float64 `json:"humidity"`
	SoilMoisture float64 `json:"soilmoisture"`
	Count        int     `gorm:"not null"`
}

var ErrInternal = errors.New("internal server error")

func (data Readings) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.IMEI, validation.Required, is.Digit, validation.Length(15, 15)),
		validation.Field(&data.Temperature, validation.Required),
		validation.Field(&data.Humidity, validation.Required, validation.Min(0.0), validation.Max(100.0)),
		validation.Field(&data.SoilMoisture, validation.Required, validation.Min(0.0), validation.Max(100.0)),
		validation.Field(&data.Count, validation.Min(0.0)),
	)
}

// func (data Readings) Validate() error {
// 	return validation.ValidateStruct(&data,
// 		validation.Field(&data.IMEI, validation.Required, is.Digit, validation.Length(15, 15)),
// 		validation.Field(&data.Temperature, validation.Required, validation.Min(-50.0), validation.Max(100)),
// 		validation.Field(&data.Humidity, validation.Required, validation.Min(0.0), validation.Max(100.0)),
// 		validation.Field(&data.SoilMoisture, validation.Required, validation.Min(0.0), validation.Max(100.0)),
// 	)
// }

type GormDeviceRepo struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormDeviceRepo {
	return &GormDeviceRepo{DB: db}
}

type DeviceRepository interface {
	CreateDevicek(device *Device) error
	GetDeviceByIMEI(imei string) (*Device, error)
	CreateReadings(readings *Readings) error
	GetReadingsByIMEI(imei string) ([]Readings, error)
	DeleteDeviceByIMEI(imei string) error
	GetAllDevices() ([]Device, error)
	GetAllDeviceIMEIs() ([]string, error)
}

func (repo *GormDeviceRepo) CreateDevicek(device *Device) error {
	return repo.DB.Create(device).Error
}

func (repo *GormDeviceRepo) GetDeviceByIMEI(imei string) (*Device, error) {
	var device Device
	err := repo.DB.Where("imei = ?", imei).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (repo *GormDeviceRepo) CreateReadings(readings *Readings) error {
	return repo.DB.Create(readings).Error
}

func (repo *GormDeviceRepo) GetReadingsByIMEI(imei string) ([]Readings, error) {
	var readings []Readings
	err := repo.DB.Where("imei=?", imei).Find(&readings).Error
	if err != nil {
		return nil, err
	}
	return readings, nil
}

//	func (repo *GormDeviceRepo) GetAllDevices() ([]Device, error) {
//		var devices []Device
//		result := repo.DB.Find(&devices)
//		return devices, result.Error
//	}
func (repo *GormDeviceRepo) GetAllDevices() ([]Device, error) {
	var devices []Device
	result := repo.DB.Debug().Find(&devices)
	fmt.Println(result.Statement.SQL.String())
	return devices, result.Error
}

// GetAllDeviceIMEIs retrieves only the IMEI numbers of all devices from the database.
func (repo *GormDeviceRepo) GetAllDeviceIMEIs() ([]string, error) {
	var imeis []string
	result := repo.DB.Model(&Device{}).Pluck("imei", &imeis)
	return imeis, result.Error
}

func (repo *GormDeviceRepo) DeleteDeviceByIMEI(imei string) error {
	// Perform the delete operation
	result := repo.DB.Where("imei = ?", imei).Delete(&Device{})

	// Log the error if it exists
	if result.Error != nil {
		log.Printf("Failed to delete device with IMEI %s: %v", imei, result.Error)
		return result.Error
	}

	// Check if no rows were affected
	if result.RowsAffected == 0 {
		log.Printf("No device found with IMEI %s", imei)
		return errors.New("no device found") // You might want to define a custom error
	}

	log.Printf("Device with IMEI %s deleted successfully. Rows affected: %d", imei, result.RowsAffected)
	return nil
}
