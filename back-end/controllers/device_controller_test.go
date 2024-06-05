package controllers_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/gorilla/mux"
// 	"github.com/louismomo66/logger/controllers"
// 	"github.com/louismomo66/logger/mocks"
// 	"github.com/louismomo66/logger/models"
// 	"gorm.io/gorm"
// )

// func TestDeleteDevice(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()

// 	mockRepo := mocks.NewMockDeviceRepository(mockCtrl)
// 	controller := controllers.NewDeviceController(mockRepo)

// 	tests := []struct {
// 		name           string
// 		imei           string
// 		setupMock      func()
// 		expectedStatus int
// 	}{
// 		{
// 			name: "Successfull Delete",
// 			imei: "123456789012345",
// 			setupMock: func() {
// 				mockRepo.EXPECT().DeleteDeviceByIMEI("123456789012345").Return(nil)
// 			},
// 			expectedStatus: http.StatusNoContent,
// 		},
// 		{
// 			name: "Device not found",
// 			imei: "nonexistentimei",
// 			setupMock: func() {
// 				mockRepo.EXPECT().DeleteDeviceByIMEI("nonexistentimei").Return(gorm.ErrRecordNotFound)
// 			},
// 			expectedStatus: http.StatusNotFound,
// 		},
// 		{
// 			name: "Internal Server Error on Delete",
// 			imei: "123450123456789",
// 			setupMock: func() {
// 				mockRepo.EXPECT().DeleteDeviceByIMEI("123450123456789").Return(models.ErrInternal)
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setupMock()

// 			// Creating a request with a mux.Var
// 			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/device/%s", tt.imei), nil)

// 			req = mux.SetURLVars(req, map[string]string{
// 				"imei": tt.imei,
// 			})

// 			resp := httptest.NewRecorder()
// 			controller.DeleteDevice(resp, req)

// 			if resp.Code != tt.expectedStatus {
// 				t.Errorf("%s: expected status %v, got %v", tt.name, tt.expectedStatus, resp.Code)
// 			}
// 		})
// 	}

// }

// func TestCreateDevice(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()

// 	mockRepo := mocks.NewMockDeviceRepository(mockCtrl)
// 	controller := controllers.NewDeviceController(mockRepo)
// 	tests := []struct {
// 		name           string
// 		deviceJSON     string
// 		setupMock      func()
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{
// 			name:       "Successful Creation",
// 			deviceJSON: `{"imei":"123456789012345","phone_number":"123-457-7890","device_type":"TestType","serial_number":"SN123456"}`,
// 			setupMock: func() {
// 				mockRepo.EXPECT().GetDeviceByIMEI(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
// 				mockRepo.EXPECT().CreateDevicek(gomock.Any()).Return(nil)
// 			},
// 			expectedStatus: http.StatusCreated,
// 			expectedBody:   "",
// 		},
// 		{
// 			name:       "Device Exists",
// 			deviceJSON: `{"imei":"123456789012345","phone_number":"123-456-7890","device_type":"TestType","serial_number":"SN123456"}`,
// 			setupMock: func() {
// 				mockRepo.EXPECT().GetDeviceByIMEI("123456789012345").Return(&models.Device{}, nil)
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedBody:   "Device with this Imei already exists",
// 		},
// 		// {
// 		// 	name:           "Invalid json",
// 		// 	deviceJSON:     `{"imei":"123"}`,
// 		// 	setupMock:      func() {},
// 		// 	expectedStatus: http.StatusBadRequest,
// 		// 	expectedBody:   "",
// 		// },
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setupMock()
// 			req, _ := http.NewRequest("POST", "/devices", bytes.NewBufferString(tt.deviceJSON))
// 			resp := httptest.NewRecorder()

// 			controller.CreateDevice(resp, req)
// 			if resp.Code != tt.expectedStatus {
// 				t.Errorf("%s: expected status %v, got %v", tt.name, tt.expectedStatus, resp.Code)
// 			}
// 			if tt.expectedBody != "" {
// 				if strings.TrimSpace(resp.Body.String()) != strings.TrimSpace(tt.expectedBody) {
// 					t.Errorf("%s: expected body %v, got %v", tt.name, tt.expectedBody, resp.Body.String())
// 				}
// 			}
// 		})
// 	}
// }

// func TestUpdateDevice(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()

// 	mockRepo := mocks.NewMockDeviceRepository(mockCtrl)
// 	controller := controllers.NewDeviceController(mockRepo)

// 	// Define your test cases
// 	tests := []struct {
// 		name           string
// 		imei           string
// 		temperature    string
// 		humidity       string
// 		soilMoisture   string
// 		count          string
// 		setupMock      func()
// 		expectedStatus int
// 	}{{
// 		name:         "Valid Update",
// 		imei:         "123456789012345",
// 		temperature:  "22.2",
// 		humidity:     "45.5",
// 		soilMoisture: "33.7",
// 		count:        "1",
// 		setupMock: func() {
// 			mockRepo.EXPECT().GetDeviceByIMEI("123456789012345").Return(&models.Device{}, nil).Times(1)
// 			mockRepo.EXPECT().CreateReadings(gomock.Any()).Return(nil).Times(1)
// 		},
// 		expectedStatus: http.StatusOK,
// 	},
// 		{
// 			name:         "Device Not Found",
// 			imei:         "000000000000000",
// 			temperature:  "22.2",
// 			humidity:     "45.5",
// 			soilMoisture: "33.7",
// 			count:        "1",
// 			setupMock: func() {
// 				mockRepo.EXPECT().GetDeviceByIMEI("000000000000000").Return(nil, errors.New("device not found")).Times(1)
// 			},
// 			expectedStatus: http.StatusNotFound,
// 		},
// 		{
// 			name:           "Invalid Temperature",
// 			imei:           "123456789012345",
// 			temperature:    "invalid",
// 			humidity:       "45.5",
// 			soilMoisture:   "33.7",
// 			count:          "1",
// 			setupMock:      func() {},
// 			expectedStatus: http.StatusBadRequest,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setupMock()
// 			req, _ := http.NewRequest("GET", fmt.Sprintf("/update?imei=%s&h=%s&t=%s&s=%s&count=%s",
// 				tt.imei, tt.humidity, tt.temperature, tt.soilMoisture, tt.count), nil)
// 			resp := httptest.NewRecorder()

// 			controller.UpdateDevice(resp, req)

// 			if resp.Code != tt.expectedStatus {
// 				t.Errorf("%s: expected status %v, got %v", tt.name, tt.expectedStatus, resp.Code)
// 			}
// 		})
// 	}
// }

// func TestGetReadings(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()

// 	mockRepo := mocks.NewMockDeviceRepository(mockCtrl)
// 	controller := controllers.NewDeviceController(mockRepo)

// 	tests := []struct {
// 		name             string
// 		imei             string
// 		setupMock        func()
// 		expectedStatus   int
// 		expectedReadings []models.Readings // Adjust based on your model
// 	}{
// 		{
// 			name: "Valid IMEI",
// 			imei: "123456789012345",
// 			setupMock: func() {
// 				readings := []models.Readings{ /* Populate with expected data */ }
// 				mockRepo.EXPECT().GetReadingsByIMEI("123456789012345").Return(readings, nil)
// 			},
// 			expectedStatus:   http.StatusOK,
// 			expectedReadings: []models.Readings{ /* Expected readings data */ },
// 		},
// 		{
// 			name: "IMEI Not Found",
// 			imei: "nonexistentimei",
// 			setupMock: func() {
// 				mockRepo.EXPECT().GetReadingsByIMEI("nonexistentimei").Return(nil, gorm.ErrRecordNotFound)
// 			},
// 			expectedStatus: http.StatusInternalServerError, // Or http.StatusNotFound if your implementation returns this
// 		},
// 		{
// 			name:           "Missing IMEI",
// 			imei:           "",
// 			setupMock:      func() {},
// 			expectedStatus: http.StatusBadRequest,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setupMock()

// 			req, _ := http.NewRequest("GET", "/readings?imei="+tt.imei, nil)
// 			resp := httptest.NewRecorder()

// 			http.HandlerFunc(controller.GetReadings).ServeHTTP(resp, req)

// 			if resp.Code != tt.expectedStatus {
// 				t.Errorf("%s: expected status %v, got %v", tt.name, tt.expectedStatus, resp.Code)
// 			}

// 			if tt.expectedStatus == http.StatusOK {
// 				var actualReadings []models.Readings
// 				if err := json.NewDecoder(resp.Body).Decode(&actualReadings); err != nil {
// 					t.Fatal("Failed to decode response body")
// 				}

// 			}
// 		})
// 	}
// }

// // func testParseReadings(t *testing.T) {
// // 	tests := []struct {
// // 		name                string
// // 		querry              string
// // 		exxpecterro         bool
// // 		expectedimei        string
// // 		expectedtemperature string
// // 		expectedhumidity    string
// // 		expectedmoisture    string
// // 		expectedcount       string
// // 	}{
// // 		{
// // 			name : "Valid parsing"
// // 			query: "/?imei"
// // 		}
// // 	}
// // }
