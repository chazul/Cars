package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	carsPrefixURLRe   = regexp.MustCompile(`^\/cars[\/]*$`)
	carsSpecificURLRe = regexp.MustCompile(`^\/cars\/(\d+)$`)
)

type car struct {
	Id           string `json:"id"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	PackageLevel string `json:"package_level"`
	Color        string `json:"color"`
	Year         string `json:"year"`
	Category     string `json:"category"`
	Mileage      string `json:"mileage"`
	Price        string `json:"price"`
}

type carDatastore struct {
	m map[string]car
}

type carHandler struct {
	store *carDatastore
}

func newCarHandler() *carHandler {

	carH := &carHandler{
		store: &carDatastore{
			m: make(map[string]car),
		},
	}
	return carH
}

func (carhandler *carHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	start := time.Now()
	traceId := req.Header["Trace_id"]
	fmt.Println("Handling request: ", traceId)

	writer.Header().Set("content-type", "application/json")
	switch {
	case req.Method == http.MethodGet && carsPrefixURLRe.MatchString(req.URL.Path):
		carhandler.GetAll(writer, req)
		logInfo(fmt.Sprintln("Served request ", traceId, "in ", time.Now().Sub(start)))
		return
	case req.Method == http.MethodGet && carsSpecificURLRe.MatchString(req.URL.Path):
		carhandler.Get(writer, req)
		logInfo(fmt.Sprintln("Served request ", traceId, "in ", time.Now().Sub(start)))
		return
	case req.Method == http.MethodPost && carsPrefixURLRe.MatchString(req.URL.Path):
		carhandler.Create(writer, req)
		logInfo(fmt.Sprintln("Served request ", traceId, "in ", time.Now().Sub(start)))
		return
	case req.Method == http.MethodPut && carsSpecificURLRe.MatchString(req.URL.Path):
		carhandler.Update(writer, req)
		logInfo(fmt.Sprintln("Served request ", traceId, "in ", time.Now().Sub(start)))
		return

	//TODO
	// implement Delete

	default:
		notFound(writer, req)
		return
	}
}

// GET: all cars
func (carHandler *carHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	cars := make([]car, 0, len(carHandler.store.m))
	for _, v := range carHandler.store.m {
		cars = append(cars, v)
	}
	jsonBytes, err := json.Marshal(cars)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// GET: one car by id
func (carHandler *carHandler) Get(writer http.ResponseWriter, req *http.Request) {
	matches := carsSpecificURLRe.FindStringSubmatch(req.URL.Path)
	if len(matches) < 2 {
		notFound(writer, req)
		return
	}

	u, ok := carHandler.store.m[matches[1]]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("not found"))
		return
	}
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		internalServerError(writer, req)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonBytes)
}

// POST: Create a car
func (carHandler *carHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var c car
	if err := json.NewDecoder(req.Body).Decode(&c); err != nil {
		fmt.Println("Req Parsing Error: ", err)
		internalServerError(writer, req)
		return
	}

	tmpLen := len(carHandler.store.m) + 1
	c.Id = strconv.Itoa(tmpLen)
	carHandler.store.m[c.Id] = c

	jsonBytes2, err2 := json.Marshal(c)
	if err2 != nil {
		internalServerError(writer, req)
		return
	}
	fmt.Println(req.RequestURI)
	writer.Header().Add("Location", req.URL.String()+"/"+c.Id) //Let the client know about the new resource
	writer.WriteHeader(http.StatusCreated)                     //Created instead Success
	writer.Write(jsonBytes2)
}

// PUT: update a car by Id
func (carHandler *carHandler) Update(writer http.ResponseWriter, req *http.Request) {

	//Get updated car info
	var updatedCar car
	if err := json.NewDecoder(req.Body).Decode(&updatedCar); err != nil {
		fmt.Println("Req Parsing Error: ", err)
		internalServerError(writer, req)
		return
	}

	matches := carsSpecificURLRe.FindStringSubmatch(req.URL.Path)
	if len(matches) < 2 {
		notFound(writer, req)
		return
	}

	//Get existing car
	carToUpdate, found := carHandler.store.m[matches[1]]
	if !found {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("not found"))
		return
	}

	//update the car in dataStore
	updatedCar.Id = carToUpdate.Id
	carHandler.store.m[carToUpdate.Id] = updatedCar

	jsonBytes, err := json.Marshal(updatedCar)
	if err != nil {
		internalServerError(writer, req)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonBytes)
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("internal server error"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

func logInfo(m string) {
	fmt.Println("info: ", m)
}

func logError(m string, err error) {
	fmt.Println("error: ", m, ", ", err)
}
