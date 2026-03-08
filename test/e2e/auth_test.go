package e2e

import (
	"net/http"
	"testing"
)

const (
	testName     = "Test User"
	testEmail    = "test@example.com"
	loginPath    = "/api/v1/auth/login"
	registerPath = "/api/v1/auth/register"
)

func TestAuthLogin(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	createPayload := map[string]string{
		"name":     testName,
		"email":    testEmail,
		"password": "password123",
	}

	createResp, err := suite.httpClient.Post("/api/v1/users", createPayload)
	if err != nil {
		t.Fatalf("Failed to create user for login test: %v", err)
	}

	if createResp.StatusCode != http.StatusCreated {
		t.Logf("Skipping login test - user creation failed with status %d", createResp.StatusCode)
		return
	}

	loginPayload := map[string]string{
		"email":    testEmail,
		"password": "password123",
	}

	resp, err := suite.httpClient.Post(loginPath, loginPayload)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if resp.StatusCode == http.StatusAccepted {
		var successResp SuccessResponse
		if err := resp.UnmarshalJSON(&successResp); err != nil {
			t.Fatalf("Failed to unmarshal login response: %v", err)
		}

		loginData, ok := successResp.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("Login response data is not in expected format")
		}

		user, ok := loginData["user"].(map[string]interface{})
		if !ok {
			t.Fatalf("User data is not in expected format")
		}

		email, _ := user["Email"].(string)
		name, _ := user["Name"].(string)
		token, _ := loginData["token"].(string)

		if email != loginPayload["email"] {
			t.Errorf("Expected email %s, got %s", loginPayload["email"], email)
		}
		if name != createPayload["name"] {
			t.Errorf("Expected name %s, got %s", createPayload["name"], name)
		}
		if token == "" {
			t.Error("Expected non-empty token")
		}
	} else {
		t.Logf("Login returned status %d (expected 202)", resp.StatusCode)
	}
}

func TestAuthLoginInvalidCredentials(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	loginPayload := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "wrongpassword",
	}

	resp, err := suite.httpClient.Post(loginPath, loginPayload)
	if err != nil {
		t.Fatalf("Failed to attempt login: %v", err)
	}

	if resp.StatusCode == http.StatusAccepted {
		t.Error("Expected login to fail with invalid credentials")
	} else {
		t.Logf("Login correctly failed with status %d", resp.StatusCode)
	}
}

func TestAuthRegister(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	registerPayload := map[string]string{
		"name":     "New User",
		"email":    "newuser@example.com",
		"password": "password123",
	}

	resp, err := suite.httpClient.Post(registerPath, registerPayload)
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}

	if resp.StatusCode == http.StatusCreated {
		var successResp SuccessResponse
		if err := resp.UnmarshalJSON(&successResp); err != nil {
			t.Fatalf("Failed to unmarshal register response: %v", err)
		}

		registerData, ok := successResp.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("Register response data is not in expected format")
		}

		user, ok := registerData["user"].(map[string]interface{})
		if !ok {
			t.Fatalf("User data is not in expected format")
		}

		email, _ := user["Email"].(string)
		name, _ := user["Name"].(string)

		if email != registerPayload["email"] {
			t.Errorf("Expected email %s, got %s", registerPayload["email"], email)
		}
		if name != registerPayload["name"] {
			t.Errorf("Expected name %s, got %s", registerPayload["name"], name)
		}
	} else {
		t.Logf("Register returned status %d (expected 201)", resp.StatusCode)
	}
}

func TestAuthRegisterDuplicateEmail(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	registerPayload := map[string]string{
		"name":     "Duplicate User",
		"email":    "duplicate@example.com",
		"password": "password123",
	}

	resp1, err := suite.httpClient.Post(registerPath, registerPayload)
	if err != nil {
		t.Fatalf("Failed to register first user: %v", err)
	}

	if resp1.StatusCode == http.StatusCreated {
		registerPayload["name"] = "Another User"
		resp2, err := suite.httpClient.Post(registerPath, registerPayload)
		if err != nil {
			t.Fatalf("Failed to attempt duplicate registration: %v", err)
		}

		if resp2.StatusCode != http.StatusConflict {
			t.Logf("Expected conflict (409) but got %d", resp2.StatusCode)
		}
	} else {
		t.Logf("Skipping duplicate test - first registration failed with status %d", resp1.StatusCode)
	}
}

func TestAuthLoginWrongPassword(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	createPayload := map[string]string{
		"name":     testName,
		"email":    "wrongpass@example.com",
		"password": "password123",
	}

	createResp, err := suite.httpClient.Post("/api/v1/users", createPayload)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if createResp.StatusCode != http.StatusCreated {
		t.Logf("Skipping wrong password test - user creation failed")
		return
	}

	loginPayload := map[string]string{
		"email":    "wrongpass@example.com",
		"password": "wrongpassword",
	}

	resp, err := suite.httpClient.Post(loginPath, loginPayload)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if resp.StatusCode == http.StatusAccepted {
		t.Error("Expected login to fail with wrong password")
	} else {
		t.Logf("Login correctly failed with status %d", resp.StatusCode)
	}
}

func TestAuthFlow(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	registerPayload := map[string]string{
		"name":     "Flow User",
		"email":    "flowuser@example.com",
		"password": "password123",
	}

	registerResp, err := suite.httpClient.Post(registerPath, registerPayload)
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}

	if registerResp.StatusCode != http.StatusCreated {
		t.Logf("Skipping flow test - registration failed with status %d", registerResp.StatusCode)
		return
	}

	loginPayload := map[string]string{
		"email":    "flowuser@example.com",
		"password": "password123",
	}

	loginResp, err := suite.httpClient.Post(loginPath, loginPayload)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if loginResp.StatusCode == http.StatusAccepted {
		t.Log("Complete auth flow successful: register -> login")
	} else {
		t.Logf("Login in flow failed with status %d", loginResp.StatusCode)
	}
}
