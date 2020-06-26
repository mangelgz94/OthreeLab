package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Client
}

func (app *App) Initialize() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	app.DB = client
	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func (app *App) createCustomer(writer http.ResponseWriter, request *http.Request) {
	var customer Customer
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&customer); err != nil {
		respondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}
	defer request.Body.Close()

	if err := customer.createCustomer(app.DB); err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusCreated, customer)

}

func (app *App) getCustomers(writer http.ResponseWriter, request *http.Request) {

	customers, err := getCustomers(app.DB)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, customers)
}

func respondWithError(writer http.ResponseWriter, code int, message string) {
	respondWithJSON(writer, code, map[string]string{"error": message})
}

func respondWithJSON(writer http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(response)
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/customers", app.createCustomer).Methods("POST")
	app.Router.HandleFunc("/customers", app.getCustomers).Methods("GET")
}
