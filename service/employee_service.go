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
	EmployeeRepository *repository.EmployeeRepository
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

	//insert employee
	employeeId, err := service.EmployeeRepository.InsertEmployee(&emp)
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

	emp, err := service.EmployeeRepository.FindEmployeeById(employeeId)
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
	employees, err := service.EmployeeRepository.FindAllEmployee()
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
	count, err := service.EmployeeRepository.UpdateEmployeeById(empId, &emp)
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

	count, err := service.EmployeeRepository.DeleteEmployeeById(empId)
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

	count, err := service.EmployeeRepository.DeleteAllEmployee()
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
