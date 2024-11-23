package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"order-api/internal/auth"
	"strings"
	"testing"

	"order-api/pkg/test"
)

func TestLoginSuccess(t *testing.T) {
	//Prepare DB
	db := test.InitDB()
	test.InitData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Phone: "+79993333333",
	})

	resp, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 201 {
		t.Fatalf("Expected status code 201, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Response body: %s", string(body))

	var resData auth.RegisterAndLoginResponse
	err = json.Unmarshal(body, &resData)

	if err != nil {
		t.Fatal(err)
	}

	if resData.SessionId == "" {
		t.Fatal("Empty SessionId")
	}
	test.RemoveData(db)
}

func TestLoginFail(t *testing.T) {
	db := test.InitDB()
	test.InitData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Phone: "+79993333330",
	})

	resp, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 401 {
		t.Fatalf("Expected status code 401, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	actualResponse := strings.TrimSpace(string(body))
	expectedResponse := "wrong phone"
	if actualResponse != expectedResponse {
		t.Fatalf("Expected response body '%s', got '%s'", expectedResponse, string(body))
	}
	test.RemoveData(db)
}
