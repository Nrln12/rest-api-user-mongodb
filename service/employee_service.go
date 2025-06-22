package service

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"user-api/model"
	"user-api/repository"
)

type EmployeeService struct {
	MongoCollection *mongo.Collection
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (service *EmployeeService) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := &Response{}
	defer json.NewEncoder(w).Encode(response)
	var emp model.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid request body: ", err)
		response.Error = err.Error()
		return
	}

	// assign new employee id
	emp.EmployeeId = uuid.NewString()
	empRepository := repository.EmployeeRepository{MongoCollection: service.MongoCollection}

	//insert employee
	employeeId, err := empRepository.InsertEmployee(&emp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error inserting employee: ", err)
		response.Error = err.Error()
		return
	}
	response.Data = employeeId
	w.WriteHeader(http.StatusCreated)
	log.Println("employee inserted with id: ", employeeId, emp)

}

func (service *EmployeeService) GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := &Response{}
	defer json.NewEncoder(w).Encode(response)

	// get employee id
	employeeId := mux.Vars(r)["id"]
	log.Println("employee id: ", employeeId)

	empRepository := repository.EmployeeRepository{MongoCollection: service.MongoCollection}
	emp, err := empRepository.FindEmployeeById(employeeId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			log.Println("employee not found with id: ", employeeId)
			response.Error = err.Error()
		}
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error finding employee: ", err)
		response.Error = err.Error()
	}
	response.Data = emp
	w.WriteHeader(http.StatusOK)
}

func (service *EmployeeService) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := &Response{}
	defer json.NewEncoder(w).Encode(response)

	// get all employees
	empRepository := repository.EmployeeRepository{MongoCollection: service.MongoCollection}
	employees, err := empRepository.FindAllEmployee()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			log.Println("employee not found")
			response.Error = err.Error()
		}
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error finding employee: ", err)
		response.Error = err.Error()
	}
	response.Data = employees
	w.WriteHeader(http.StatusOK)
}

func (service *EmployeeService) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := &Response{}
	defer json.NewEncoder(w).Encode(response)

	// get employee id
	empId := mux.Vars(r)["id"]
	log.Println("employee id: ", empId)

	// validation
	if empId == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("employee id is empty: ", empId)
		response.Error = "employee id is empty"
		return
	}

	// decode body
	var emp model.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid request body: ", err)
		response.Error = err.Error()
		return
	}

	emp.EmployeeId = empId
	empRepository := repository.EmployeeRepository{MongoCollection: service.MongoCollection}
	count, err := empRepository.UpdateEmployeeById(empId, &emp)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			log.Println("employee not found with id: ", empId)
			response.Error = err.Error()
		}
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error updating employee: ", err)
		response.Error = err.Error()
	}
	response.Data = count
	w.WriteHeader(http.StatusOK)
}

func (service *EmployeeService) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := &Response{}
	defer json.NewEncoder(w).Encode(response)

	// decode path variable
	empId := mux.Vars(r)["id"]
	log.Println("employee id: ", empId)

	// validation
	if empId == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("employee id is empty: ", empId)
		response.Error = "employee id is empty"
		return
	}

	empRepository := repository.EmployeeRepository{MongoCollection: service.MongoCollection}
	count, err := empRepository.DeleteEmployeeById(empId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			log.Println("employee not found with id: ", empId)
			response.Error = err.Error()
		}
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error deleting employee: ", err)
		response.Error = err.Error()
	}
	response.Data = count
	w.WriteHeader(http.StatusOK)
}

func (service *EmployeeService) DeleteAllEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := &Response{}
	defer json.NewEncoder(w).Encode(response)

	empRepository := repository.EmployeeRepository{MongoCollection: service.MongoCollection}
	count, err := empRepository.DeleteAllEmployee()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			log.Println("employee not found")
			response.Error = err.Error()
		}
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error deleting employee: ", err)
		response.Error = err.Error()
	}
	response.Data = count
	w.WriteHeader(http.StatusOK)
}
