package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	mockUseCase "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/mockUseCase"
	utils "github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

// Signup test
func TestAdminSignUp(t *testing.T) {

	//NewController returns a new Controller.
	ctrl := gomock.NewController(t)
	//NewMockAdminUseCase creates a new mock instance.
	adminUseCase := mockUseCase.NewMockAdminUseCase(ctrl)
	adminHandler := NewAdminHandler(adminUseCase)

	// Create a new Gin router instance for testing
	router := gin.Default()
	router.POST("/admin/signup", adminHandler.AdminSignUp)

	tests := []struct {
		testName           string
		requestBody        gin.H
		mockExpectations   func()
		expectedStatusCode int
	}{
		{
			testName: "Successful admin signup",
			requestBody: gin.H{
				"firstName": "John",
				"lastName":  "Doe",
				"email":     "johndoe@example.com",
				"phoneNum":  "1234567890",
				"password":  "password123",
			},
			mockExpectations: func() {
				adminUseCase.EXPECT().SignUpAdmin(gomock.Any(), gomock.AssignableToTypeOf(utils.AdminSignUp{})).Times(1).Return(domain.Admin{}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			testName: "Failed admin signup - email format is incorrect",
			requestBody: gin.H{
				"firstName": "John",
				"lastName":  "Doe",
				"email":     "invalidemail",
				"phoneNum":  "1234567890",
				"password":  "password123",
			},
			mockExpectations: func() {
				// No expectations for the mock in this case
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			testName: "Failed admin signup - phone number format is incorrect",
			requestBody: gin.H{
				"firstName": "John",
				"lastName":  "Doe",
				"email":     "johndoe@example.com",
				"phoneNum":  "invalidphone",
				"password":  "password123",
			},
			mockExpectations: func() {
				// No expectations for the mock in this case
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			testName: "Failed admin signup - repository error",
			requestBody: gin.H{
				"firstName": "John",
				"lastName":  "Doe",
				"email":     "johndoe@example.com",
				"phoneNum":  "1234567890",
				"password":  "password123",
			},
			mockExpectations: func() {
				adminUseCase.EXPECT().SignUpAdmin(gomock.Any(), gomock.AssignableToTypeOf(utils.AdminSignUp{})).Times(1).Return(domain.Admin{}, errors.New("failed to add admin"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			// Set up the mock expectations
			tt.mockExpectations()

			// Create a new HTTP request with the request body
			//Marshal returns the JSON encoding of v.
			reqBody, _ := json.Marshal(tt.requestBody)
			//NewRequest wraps NewRequestWithContext using context.Background.
			req, _ := http.NewRequest(http.MethodPost, "/admin/signup", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Create a new HTTP recorder to capture the response
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Assert the response status code
			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
