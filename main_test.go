package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"dav_docker/routes"
	"dav_docker/utils"

	"github.com/gofiber/fiber/v2"
)

func TestAPIWorkflow(t *testing.T) {
	os.Remove("./db/database.db")
	os.Remove("./db/.env")

	if err := utils.InitDB(); err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	if err := utils.InitEnv(); err != nil {
		t.Fatalf("InitEnv failed: %v", err)
	}

	app := fiber.New()

	app.Get("/api/init", routes.HandleInit)
	app.Post("/api/register", routes.HandleRegister)
	app.Post("/api/login", routes.HandleLogin)
	app.Get("/api/refresh", routes.HandleRefresh)
	app.Get("/api/auth", routes.AuthMiddleware, routes.HandleAuth)

	// Test init before registration
	req := httptest.NewRequest("GET", "/api/init", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Test register
	body, _ := json.Marshal(routes.UserRequest{Username: "admin", Password: "password123"})
	req = httptest.NewRequest("POST", "/api/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Register status 200 expected, got %d", resp.StatusCode)
	}

	// Test login
	body, _ = json.Marshal(routes.UserRequest{Username: "admin", Password: "password123"})
	req = httptest.NewRequest("POST", "/api/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Login status 200 expected, got %d", resp.StatusCode)
	}

	var loginResp struct {
		Ok   bool   `json:"ok"`
		Data string `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&loginResp)
	if !loginResp.Ok || loginResp.Data == "" {
		t.Errorf("Login failed: %+v", loginResp)
	}

	accessToken := loginResp.Data

	// Test auth with token header
	req = httptest.NewRequest("GET", "/api/auth", nil)
	req.Header.Set("token", accessToken)
	resp, _ = app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Auth status 200 expected, got %d", resp.StatusCode)
	}
}
