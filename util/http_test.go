package util_test

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/http/httptest"
	"testing"
	
)

type TestEntity struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func TestValidateRequest(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   any
		expectedError bool
	}{
		{
			name:          "Valid request",
			requestBody:   TestEntity{Name: "John Doe", Email: "john@example.com"},
			expectedError: false,
		},
		{
			name:          "Invalid JSON",
			requestBody:   "{invalid-json}",
			expectedError: true,
		},
		{
			name:          "Validation error: missing name",
			requestBody:   TestEntity{Email: "john@example.com"},
			expectedError: true,
		},
		{
			name:          "Validation error: invalid email",
			requestBody:   TestEntity{Name: "John Doe", Email: "invalid-email"},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			var reqBody []byte
			var err error

			switch body := tt.requestBody.(type) {
			case string:
				reqBody = []byte(body)
			default:
				reqBody, _ = json.Marshal(body)
			}

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			c := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(c)

			c.Request().Header.SetMethod(http.MethodPost)
			c.Request().SetBody(reqBody)

			var entity TestEntity
			err = util.ValidateRequest(c, &entity)
			if tt.expectedError && err == nil {
				t.Errorf("expected an error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("expected no error but got %v", err)
			}
		})
	}
}
