package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var app App

func TestMain(m *testing.M) {
	app.Initialize()
	code := m.Run()
	os.Exit(code)
}

func TestGetCustomersEndpoint(test *testing.T) {
	request, _ := http.NewRequest("GET", "/customers", nil)
	response := executeRequest(request)
	checkResponseCode(test, http.StatusOK, response.Code)
}

func TestCreateCustomerEndpoint(test *testing.T) {
	var jsonStr = []byte(`{
        "name": "John",
        "last_name": "Doe",
        "birth_date": "2020-04-23T18:25:43.511Z",
        "email": "john@doe.com",
        "phone_number": "000-000-000"}`)
	request, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response := executeRequest(request)
	checkResponseCode(test, http.StatusCreated, response.Code)
	var customer map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &customer)
	if customer["name"] != "John" {
		test.Errorf("Expected  name to be 'John'. Got '%v'", customer["name"])
	}
	if customer["last_name"] != "Doe" {
		test.Errorf("Expected last name to be 'Doe'. Got '%v'", customer["last_name"])
	}
	if customer["birth_date"] != "2020-04-23T18:25:43.511Z" {
		test.Errorf("Expected birth_date to be '2020-04-23T18:25:43.511Z'. Got '%v'", customer["birth_date"])
	}
	if customer["email"] != "john@doe.com" {
		test.Errorf("Expected email to be 'john@doe.com'. Got '%v'", customer["email"])
	}
	if customer["phone_number"] != "000-000-000" {
		test.Errorf("Expected phone number to be '000-000-000'. Got '%v'", customer["phone_number"])
	}

}

func checkResponseCode(t *testing.T, expected, actual int) {
	log.Printf("check response code")
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}
