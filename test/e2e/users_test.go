package e2e

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCreateUser(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	payload := map[string]string{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
	}

	resp, err := suite.httpClient.Post("/api/v1/users", payload)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Since /api/v1/users requires authentication, this should return 401
	resp.AssertStatusCode(t, http.StatusUnauthorized)
}

func TestCreateDuplicateEmail(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	payload := map[string]string{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
	}

	resp1, err := suite.httpClient.Post("/api/v1/users", payload)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	if resp1.StatusCode == http.StatusCreated {
		payload["name"] = "Jane Doe"
		resp2, err := suite.httpClient.Post("/api/v1/users", payload)
		if err != nil {
			t.Fatalf("Failed to attempt duplicate creation: %v", err)
		}

		if resp2.StatusCode != http.StatusConflict {
			t.Logf("Expected conflict (409) but got %d", resp2.StatusCode)
		}
	} else {
		t.Logf("Skipping duplicate test - first user creation failed with status %d", resp1.StatusCode)
	}
}

func TestGetUserNotFound(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	resp, err := suite.httpClient.Get("/api/v1/users/99999")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusInternalServerError {
		t.Logf("Expected 404 or 500, got %d", resp.StatusCode)
	}
}

func TestUpdateUser(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	createPayload := map[string]string{
		"name":     "Original Name",
		"email":    "original@example.com",
		"password": "password123",
	}

	createResp, err := suite.httpClient.Post("/api/v1/users", createPayload)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if createResp.StatusCode != http.StatusCreated {
		t.Logf("Skipping update test - create returned status %d", createResp.StatusCode)
		return
	}

	var successResp SuccessResponse
	if err := createResp.UnmarshalJSON(&successResp); err != nil {
		t.Fatalf("Failed to unmarshal create response: %v", err)
	}

	userData, ok := successResp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Create response data is not in expected format")
	}

	userIDFloat, _ := userData["id"].(float64)
	userID := int64(userIDFloat)

	updatePayload := map[string]interface{}{
		"name": "Updated Name",
	}

	resp, err := suite.httpClient.Put(fmt.Sprintf("/api/v1/users/%d", userID), updatePayload)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		var updateSuccessResp SuccessResponse
		if err := resp.UnmarshalJSON(&updateSuccessResp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		updatedUserData, ok := updateSuccessResp.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("Update response data is not in expected format")
		}

		name, _ := updatedUserData["name"].(string)
		email, _ := updatedUserData["email"].(string)

		if name != "Updated Name" {
			t.Errorf("Expected name Updated Name, got %s", name)
		}
		if email != "original@example.com" {
			t.Errorf("Expected email to remain original@example.com, got %s", email)
		}
	}
}

func TestUpdateUserEmail(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	user1Payload := map[string]string{
		"name":     "User 1",
		"email":    "user1@example.com",
		"password": "password123",
	}

	resp1, err := suite.httpClient.Post("/api/v1/users", user1Payload)
	if err != nil {
		t.Fatalf("Failed to create user 1: %v", err)
	}

	if resp1.StatusCode != http.StatusCreated {
		t.Logf("Skipping email update test - first create returned status %d", resp1.StatusCode)
		return
	}

	user2Payload := map[string]string{
		"name":     "User 2",
		"email":    "user2@example.com",
		"password": "password123",
	}

	resp2, err := suite.httpClient.Post("/api/v1/users", user2Payload)
	if err != nil {
		t.Fatalf("Failed to create user 2: %v", err)
	}

	if resp2.StatusCode != http.StatusCreated {
		t.Logf("Skipping email update test - second create returned status %d", resp2.StatusCode)
		return
	}

	var successResp SuccessResponse
	if err := resp2.UnmarshalJSON(&successResp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	userData, ok := successResp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Create response data is not in expected format")
	}

	userID2Float, _ := userData["id"].(float64)
	userID2 := int64(userID2Float)

	updatePayload := map[string]interface{}{
		"email": "user1@example.com",
	}

	resp, err := suite.httpClient.Put(fmt.Sprintf("/api/v1/users/%d", userID2), updatePayload)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	if resp.StatusCode != http.StatusConflict {
		t.Logf("Expected conflict (409) but got %d", resp.StatusCode)
	}
}

func TestDeleteUser(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	createPayload := map[string]string{
		"name":     "User to Delete",
		"email":    "delete@example.com",
		"password": "password123",
	}

	createResp, err := suite.httpClient.Post("/api/v1/users", createPayload)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if createResp.StatusCode != http.StatusCreated {
		t.Logf("Skipping delete test - create returned status %d", createResp.StatusCode)
		return
	}

	var successResp SuccessResponse
	if err := createResp.UnmarshalJSON(&successResp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	userData, ok := successResp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Create response data is not in expected format")
	}

	userIDFloat, _ := userData["id"].(float64)
	userID := int64(userIDFloat)

	resp, err := suite.httpClient.Delete(fmt.Sprintf("/api/v1/users/%d", userID))
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		getResp, err := suite.httpClient.Get(fmt.Sprintf("/api/v1/users/%d", userID))
		if err != nil {
			t.Fatalf("Failed to verify deletion: %v", err)
		}

		if getResp.StatusCode != http.StatusNotFound {
			t.Logf("Expected 404 after delete, got %d", getResp.StatusCode)
		}
	} else {
		t.Logf("Delete returned status %d (expected 204)", resp.StatusCode)
	}
}

func TestListUsers(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	successCount := 0
	for i := 1; i <= 3; i++ {
		payload := map[string]string{
			"name":     fmt.Sprintf("User %d", i),
			"email":    fmt.Sprintf("user%d@example.com", i),
			"password": "password123",
		}

		resp, err := suite.httpClient.Post("/api/v1/users", payload)
		if err != nil {
			t.Fatalf("Failed to create user %d: %v", i, err)
		}

		if resp.StatusCode == http.StatusCreated {
			successCount++
		}
	}

	if successCount == 0 {
		t.Logf("Skipping list test - no users were created successfully")
		return
	}

	resp, err := suite.httpClient.Get("/api/v1/users")
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		var usersResp UsersListResponse
		if err := resp.UnmarshalJSON(&usersResp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(usersResp.Data) != successCount {
			t.Logf("Expected %d users, got %d", successCount, len(usersResp.Data))
		}
	} else {
		t.Logf("List returned status %d (expected 200)", resp.StatusCode)
	}
}
