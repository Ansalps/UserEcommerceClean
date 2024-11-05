package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
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
	//gomock.Controller object, which is responsible for tracking and managing mock expectations and assertions during the test
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
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   fmt.Sprintf(`{"message":"%v"}`, models.UserAlreadyExists),
			returnError:        errors.New(models.UserAlreadyExists),
			expectSignupCall:   true,
		},
		{
			name: "invalid input - empty first name",
			requestBody: models.User{
				FirstName: "",
				LastName:  "Doe",
				Email:     "johndoe@gmail.com",
				Password:  "password",
				Phone:     "1234567890",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"message":"%v"}`, "FirstName is required"),
			returnError:        nil,
			expectSignupCall:   false,
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
func TestLogin(t *testing.T) {
	router := gin.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockIUserService(ctrl)
	UserController := &UserController{UserService: mockUserService}
	router.POST("user-login", UserController.UserLogin)
	tests := []struct {
		name               string
		requestBody        models.UserLogin
		expectedStatusCode int
		mockError          error
		validateResponse   func(t *testing.T, response map[string]interface{})
	}{
		{
			name: "successful login",
			requestBody: models.UserLogin{
				Email:    "test@example.com",
				Password: "password",
			},
			expectedStatusCode: http.StatusOK,
			mockError:          nil,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, models.LoginSuccesful, response["message"])
				assert.NotEmpty(t, response["token"])
				user, ok := response["user"].(map[string]interface{})
				assert.True(t, ok)
				assert.NotNil(t, user)
			},
		},
		{
			name: "wrong password",
			requestBody: models.UserLogin{
				Email:    "test@example.com",
				Password: "WrongPass@123",
			},
			expectedStatusCode: http.StatusUnauthorized,
			mockError:          errors.New(models.InvalidInput),
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, models.InvalidInput, response["error"])
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.mockError != nil {
				mockUserService.EXPECT().UserLogin(&test.requestBody).Return(nil, test.mockError)
			} else if test.requestBody.Email != "" && test.requestBody.Password != "" {
				user := &models.User{
					Email: test.requestBody.Email,
				}
				mockUserService.EXPECT().UserLogin(&test.requestBody).Return(user, nil)
				mockUserService.EXPECT().ComparePassword(test.requestBody, *user).Return(true)
			}
			reqBody, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/user-login", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Use httptest.NewRecorder to record the response
			resp := httptest.NewRecorder()

			// Serve the request with the router
			router.ServeHTTP(resp, req)

			assert.Equal(t, test.expectedStatusCode, resp.Code)

			var response map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)

			test.validateResponse(t, response)
		})

	}
}
func TestGetProfile(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID != "" {
			c.Set("ID", userID) // Store userID in the Gin context
		}
		c.Next() // Continue to the next handler
	})
	// Example profile handler
	router.GET("/user/profile", func(c *gin.Context) {
		userID, exists := c.Get("ID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"userID": userID})
	})
	tests := []struct {
		name               string
		userID             string
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name:               "Successful profile retrieval",
			userID:             "123",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   map[string]interface{}{"userID": "123"},
		},
		{
			name:               "Unauthorized access - missing user ID",
			userID:             "",
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   map[string]interface{}{"error": "Unauthorized"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/user/profile", nil)
			req.Header.Set("Content-Type", "application/json")
			if test.userID != "" {
				req.Header.Set("X-User-ID", test.userID)
			}

			// Use httptest.NewRecorder to record the response
			resp := httptest.NewRecorder()

			// Serve the request with the router
			router.ServeHTTP(resp, req)
			//assert.NoError(t, err)
			assert.Equal(t, test.expectedStatusCode, resp.Code)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)
			assert.Equal(t, test.expectedResponse, response)
		})
	}
}
