package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_CarRouter(t *testing.T) {
	t.Run("Create Car Happy path", createCar_Should_Return_the_Car_with_Auto_Generated_ID)
	t.Run("Get Car Happy Path", getAllCars)
	t.Run("Get Car by Id Happy Path", getCarById)
	t.Run("Get Car by Id Saf Path", getCarById_NotFound)
	t.Run("Update Car By ID", updateCarByID)
	t.Run("Update Car that doesn't exist", updateCarByID_IdNotFound)
}

func createCar_Should_Return_the_Car_with_Auto_Generated_ID(t *testing.T) {
	carHandler := newCarHandler()

	w, _ := addCar(carHandler)

	expected := `{"id":"1","make":"Ford","model":"F10","package_level":"Base","color":"Silver","year":"2010","category":"Truck","mileage":"200000","price":"1999900"}`
	actual := w.Body.String()
	fmt.Println("createCar", actual)

	//No assertions in standard testing pkg, so sad :(
	if string(expected) != w.Body.String() {
		t.Fatalf("addCar :: expected response body to be: `%s`, got: `%s`", expected, actual)
	}
}

func getAllCars(t *testing.T) {
	carHandler := newCarHandler()

	addCar(carHandler)
	addCar(carHandler)

	w2 := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/cars", nil)
	r.Header.Set("Trace_id", "getAllCars")
	r.Header.Set("Content-Type", "application/json")
	carHandler.ServeHTTP(w2, r)
	// fmt.Println("Get result", w2.Body.String())
	actual := w2.Body.String()

	//TODO: make a json object comparison, than a literal string comparison
	expected := `[{"id":"1","make":"Ford","model":"F10","package_level":"Base","color":"Silver","year":"2010","category":"Truck","mileage":"200000","price":"1999900"},{"id":"2","make":"Ford","model":"F10","package_level":"Base","color":"Silver","year":"2010","category":"Truck","mileage":"200000","price":"1999900"}]`

	if string(expected) != w2.Body.String() {
		t.Fatalf("getCar :: expected response body to be: `%s`, got: `%s`", expected, actual)
	}
}

func getCarById(t *testing.T) {
	carHandler := newCarHandler()

	//Lets add 2 cars to the DB
	addCar(carHandler)
	addCar(carHandler)

	w2 := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/cars/2", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Trace_id", "getCarbyId")
	carHandler.ServeHTTP(w2, r)
	fmt.Println("Get result", w2.Body.String())
	actual := w2.Body.String()

	expected := `{"id":"2","make":"Ford","model":"F10","package_level":"Base","color":"Silver","year":"2010","category":"Truck","mileage":"200000","price":"1999900"}`

	if string(expected) != w2.Body.String() {
		t.Fatalf("getCar :: expected response body to be: `%s`, got: `%s`", expected, actual)
	}
}

// trying to get a car by id that doesn't exists
func getCarById_NotFound(t *testing.T) {
	carHandler := newCarHandler()

	//Lets add 2 cars to the DB
	addCar(carHandler)
	addCar(carHandler)

	w2 := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/cars/3", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Trace_id", "getCarbyId_NotFound")
	carHandler.ServeHTTP(w2, r)
	// fmt.Println("Get result", w2.Body.String())
	actual := w2.Body.String()

	expected := `not found`

	if string(expected) != w2.Body.String() {
		t.Fatalf("getCar :: expected response body to be: `%s`, got: `%s`", expected, actual)
	}
}

// An exisiting car is updated using the Id
func updateCarByID(t *testing.T) {
	carHandler := newCarHandler()

	//Lets add 2 cars to the DB
	addCar(carHandler)
	addCar(carHandler)

	//updated car info
	values := map[string]string{
		"make":          "Ford",
		"model":         "BRONCO",
		"package_level": "Base",
		"color":         "Green",
		"year":          "2020",
		"category":      "SUV",
		"mileage":       "0",
		"price":         "1999900"}

	jsonValue, _ := json.Marshal(values)
	payload := bytes.NewBuffer(jsonValue)

	w2 := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/cars/2", payload)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Trace_id", "updateCar")
	carHandler.ServeHTTP(w2, r)
	// fmt.Println("Get result", w2.Body.String())
	actual := w2.Body.String()

	expected := `{"id":"2","make":"Ford","model":"BRONCO","package_level":"Base","color":"Green","year":"2020","category":"SUV","mileage":"0","price":"1999900"}`

	if string(expected) != w2.Body.String() {
		t.Fatalf("getCar :: expected response body to be: `%s`, got: `%s`", expected, actual)
	}
}

// try to update a non-exisiting car using the Id
func updateCarByID_IdNotFound(t *testing.T) {
	carHandler := newCarHandler()

	//Lets add 2 cars to the DB
	addCar(carHandler)

	//updated car info
	values := map[string]string{
		"make":          "Ford",
		"model":         "BRONCO",
		"package_level": "Base",
		"color":         "Green",
		"year":          "2020",
		"category":      "SUV",
		"mileage":       "0",
		"price":         "1999900"}

	jsonValue, _ := json.Marshal(values)
	payload := bytes.NewBuffer(jsonValue)

	w2 := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/cars/2", payload)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Trace_id", "updateCarThatDoesNotExist")
	carHandler.ServeHTTP(w2, r)
	// fmt.Println("Get result", w2.Body.String())
	actual := w2.Body.String()

	expected := `not found`

	if string(expected) != w2.Body.String() {
		t.Fatalf("getCar :: expected response body to be: `%s`, got: `%s`", expected, actual)
	}
}

func addCar(carHandler *carHandler) (*httptest.ResponseRecorder, *http.Request) {
	values := map[string]string{
		"make":          "Ford",
		"model":         "F10",
		"package_level": "Base",
		"color":         "Silver",
		"year":          "2010",
		"category":      "Truck",
		"mileage":       "200000",
		"price":         "1999900"}

	jsonValue, _ := json.Marshal(values)
	// fmt.Println("jsonValue", string(jsonValue))
	payload := bytes.NewBuffer(jsonValue)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/cars", payload)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Trace_id", "addCar")

	carHandler.ServeHTTP(w, r)
	return w, r
}
