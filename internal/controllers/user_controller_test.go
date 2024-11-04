package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ansalps/UserEcommerceClean/internal/mocks"
	"github.com/Ansalps/UserEcommerceClean/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	router := gin.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockIUserService(ctrl)
	UserController := &UserController{UserService: mockUserService}
	router.POST("user-signup", UserController.UserSignUp)

	type TypeCase struct {
		name               string
		requestBody        models.User
		expectedStatusCode int
		expectedResponse   string
		returnError        error
		expectSignupCall   bool
	}
	tests := []TypeCase{
		{
			name: "successful signup",
			requestBody: models.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "johndoe@gmail.com",
				Password:  "password",
				Phone:     "1234567890",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   fmt.Sprintf(`{"message":"%v"}`, models.SignupSuccessful),
			returnError:        nil,
			expectSignupCall:   true,
		},
		{
			name: "user already exists",
			requestBody: models.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "johndoe@gmail.com",
				Password:  "password",
				Phone:     "1234567890",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   fmt.Sprintf(`{"message":"%v"}`, models.SignupSuccessful),
			returnError:        nil,
			expectSignupCall:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectSignupCall {
				if test.returnError != nil {
					mockUserService.EXPECT().UserSignUp(&test.requestBody).Return(test.returnError).Times(1)
				} else {
					mockUserService.EXPECT().UserSignUp(&test.requestBody).Return(nil).Times(1)
				}
			}

			reqBody, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/user-signup", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Use httptest.NewRecorder to record the response
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, test.expectedStatusCode, recorder.Code)

			assert.JSONEq(t, test.expectedResponse, recorder.Body.String())
		})
	}
}
