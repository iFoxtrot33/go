package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"order-api/internal/auth"
	"order-api/internal/order"
	"order-api/pkg/test"
	"testing"
)

func TestLoginAndMakeOrder(t *testing.T) {
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

	defer resp.Body.Close()

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

	dataForJWT, err := json.Marshal(&auth.AuthRequest{
		SessionId: resData.SessionId,
		Code:      "1234",
	})

	if err != nil {
		t.Fatal(err)
	}

	resp2, err := http.Post(ts.URL+"/auth/session", "application/json", bytes.NewReader(dataForJWT))

	if err != nil {
		t.Fatal(err)
	}

	defer resp2.Body.Close()

	if resp2.StatusCode != 201 {
		t.Fatalf("Expected status code 201, got %d", resp.StatusCode)
	}

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		t.Fatal(err)
	}

	var authResponse auth.AuthResponse
	err = json.Unmarshal(body2, &authResponse)

	if err != nil {
		t.Fatal(err)
	}

	if authResponse.Token == "" {
		t.Fatal("Empty JWT token received")
	}

	dataForOrder, err := json.Marshal(&order.OrderCreateRequest{
		Products: []order.OrderProductItem{
			{
				ProductID: 1,
				Quantity:  1,
			},
			{
				ProductID: 2,
				Quantity:  2,
			},
		},
		Description: "Pizza with cheese",
	})

	req3, err := http.NewRequest("POST", ts.URL+"/order", bytes.NewReader(dataForOrder))

	req3.Header.Set("Authorization", "Bearer "+authResponse.Token)
	req3.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp3, err := client.Do(req3)
	if err != nil {
		t.Fatal(err)
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != 201 {
		t.Fatalf("Expected status code 201 for order creation, got %d", resp3.StatusCode)
	}

	body3, err := io.ReadAll(resp3.Body)
	if err != nil {
		t.Fatal(err)
	}

	var createdOrder order.Order

	err = json.Unmarshal(body3, &createdOrder)
	if err != nil {
		t.Fatal(err)
	}

	req4, err := http.NewRequest("GET", ts.URL+"/my-orders", nil)
	if err != nil {
		t.Fatal(err)
	}

	req4.Header.Set("Authorization", "Bearer "+authResponse.Token)
	req4.Header.Set("Content-Type", "application/json")

	resp4, err := client.Do(req4)
	if err != nil {
		t.Fatal(err)
	}
	defer resp4.Body.Close()

	if resp4.StatusCode != 201 {
		t.Fatalf("Expected status code 200 for get orders, got %d", resp4.StatusCode)
	}

	body4, err := io.ReadAll(resp4.Body)
	if err != nil {
		t.Fatal(err)
	}

	var orders []order.Order
	err = json.Unmarshal(body4, &orders)
	if err != nil {
		t.Fatal(err)
	}

	if len(orders) == 0 {
		t.Fatal("No orders found")
	}

	var found bool

	for _, order := range orders {
		if order.ID == createdOrder.ID {
			found = true
			if order.Description != "Pizza with cheese" {
				t.Fatalf("Wrong order description: expected 'Pizza with cheese', got '%s'",
					order.Description)
			}
			break
		}
	}

	if !found {
		t.Fatal("Created order not found in orders list")
	}

	test.RemoveData(db)
}
