package models_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/louismomo66/logger/mocks"
	"github.com/louismomo66/logger/models"
	"gorm.io/gorm"
)

func TestCreateDevicek(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockDeviceRepository(ctrl)
	testDevice := &models.Device{IMEI: "123456789012345", PhoneNumber: "099929234", DeviceType: "logger"}
	tests := []struct {
		name      string
		device    *models.Device
		setupMock func()
		wantErr   bool
	}{
		{
			name:   "Success - Create Device",
			device: testDevice,
			setupMock: func() {
				mockRepo.EXPECT().CreateDevicek(testDevice).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Failure - DB Error",
			device: testDevice,
			setupMock: func() {
				mockRepo.EXPECT().CreateDevicek(testDevice).Return(gorm.ErrInvalidDB)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mockRepo.CreateDevicek(tt.device)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateReadings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeviceRepository(ctrl)
	testReadings := &models.Readings{
		IMEI:         "123456789012345",
		Temperature:  25.0,
		Humidity:     50.0,
		SoilMoisture: 10.0,
		Count:        1,
	}

	tests := []struct {
		name      string
		readings  *models.Readings
		setupMock func()
		wantErr   bool
	}{
		{
			name:     "Success - Create Readings",
			readings: testReadings,
			setupMock: func() {
				mockRepo.EXPECT().CreateReadings(testReadings).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Failure - DB Error",
			readings: testReadings,
			setupMock: func() {
				mockRepo.EXPECT().CreateReadings(testReadings).Return(gorm.ErrInvalidDB)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mockRepo.CreateReadings(tt.readings)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateReadings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDeviceByIMEI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeviceRepository(ctrl)
	imei := "123456789012345"
	expectedDevice := &models.Device{
		IMEI: imei,
		// Other fields as needed...
	}

	// Test scenarios
	tests := []struct {
		name           string
		imei           string
		setupMock      func()
		wantErr        bool
		expectedDevice *models.Device
	}{
		{
			name:           "Success - Device Found",
			imei:           imei,
			wantErr:        false,
			expectedDevice: expectedDevice,
			setupMock: func() {
				mockRepo.EXPECT().GetDeviceByIMEI(imei).Return(expectedDevice, nil).Times(1)
			},
		},
		{
			name:    "Failure - Device Not Found",
			imei:    "nonexistent",
			wantErr: true,
			setupMock: func() {
				mockRepo.EXPECT().GetDeviceByIMEI("nonexistent").Return(nil, gorm.ErrRecordNotFound).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			device, err := mockRepo.GetDeviceByIMEI(tt.imei)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDeviceByIMEI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && device.IMEI != tt.expectedDevice.IMEI {
				t.Errorf("GetDeviceByIMEI() got = %v, want %v", device, tt.expectedDevice)
			}
		})
	}
}
