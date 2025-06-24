package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"user-api/repository"
	"user-api/service"
)

var client *mongo.Client

func init() {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("env file loaded")

	// set creadentials
	credential := options.Credential{
		Username:   os.Getenv("MONGO_USERNAME"),
		Password:   os.Getenv("MONGO_PASSWORD"),
		AuthSource: os.Getenv("MONGO_AUTH_SOURCE"),
	}
	// create mongo client
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")).SetAuth(credential))
	if err != nil {
		log.Fatal("connection error", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("ping failed", err)
	}
	log.Println("Connected to MongoDB")
}

func main() {
	// close mongo connection
	defer client.Disconnect(context.Background())

	collection := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))

	// create employee service
	employeeRepository := repository.EmployeeRepository{MongoCollection: collection}
	empService := service.EmployeeService{EmployeeRepository: &employeeRepository}

	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	subRouter.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	employeeRouter := subRouter.PathPrefix("/employees").Subrouter()
	employeeRouter.HandleFunc("", empService.CreateEmployee).Methods(http.MethodPost)
	employeeRouter.HandleFunc("/{id}", empService.GetEmployeeById).Methods(http.MethodGet)
	employeeRouter.HandleFunc("", empService.GetAllEmployees).Methods(http.MethodGet)
	employeeRouter.HandleFunc("/{id}", empService.UpdateEmployee).Methods(http.MethodPut)
	employeeRouter.HandleFunc("/{id}", empService.DeleteEmployee).Methods(http.MethodDelete)
	employeeRouter.HandleFunc("", empService.DeleteAllEmployees).Methods(http.MethodDelete)
	log.Println("Starting server on 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("running..."))
}
