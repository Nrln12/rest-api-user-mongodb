package repository

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"testing"
	"user-api/model"
)

func newMongoClient() *mongo.Client {
	mongoTestClient, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI("mongodb+srv://admin:admin@cluster0.9oypuag.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		log.Fatal("error while connecting to mongodb", err)
	}
	log.Println("mongodb connected successfully")
	if err = mongoTestClient.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatal("ping failed", err)
	}
	log.Println("ping success")
	return mongoTestClient
}

func TestMongoOperations(t *testing.T) {
	mongoTestClient := newMongoClient()
	defer mongoTestClient.Disconnect(context.Background())

	// mock data
	emp1 := uuid.New().String()
	//emp2 := uuid.New().String()

	// connect to collection
	collection := mongoTestClient.Database("companydb").Collection("employees_test")
	empRepo := EmployeeRepository{
		MongoCollection: collection,
	}

	// insert employee
	t.Run("Insert Employee 1", func(t *testing.T) {
		emp := model.Employee{
			EmployeeId: emp1,
			Name:       "John Doe",
			Department: "IT",
		}
		result, err := empRepo.InsertEmployee(&emp)
		if err != nil {
			t.Fatal("error while inserting employee", err)
		}
		t.Log("Inserted employee", result)
	})

	// Get Employee 1 Data
	t.Run("Get Employee 1", func(t *testing.T) {
		result, err := empRepo.FindEmployeeById(emp1)
		if err != nil {
			t.Fatal("error while finding employee", err)
		}
		t.Log("employee 1", result.Name)
	})

	// Get All Employees
	t.Run("Get All Employees", func(t *testing.T) {
		result, err := empRepo.FindAllEmployee()
		if err != nil {
			t.Fatal("error while finding all employee", err)
		}
		t.Log("All employees", result)
	})

	// Update Employee 1
	t.Run("Update Employee 1", func(t *testing.T) {
		emp := model.Employee{
			EmployeeId: emp1,
			Name:       "John Smith",
			Department: "IT",
		}
		result, err := empRepo.UpdateEmployeeById(emp1, &emp)
		if err != nil {
			t.Fatal("error while updating employee", err)
		}
		t.Log("Updated employee count", result)
	})

	// Delete Employee 1
	t.Run("Delete Employee 1", func(t *testing.T) {
		result, err := empRepo.DeleteEmployeeById(emp1)
		if err != nil {
			t.Fatal("error while deleting employee", err)
		}
		t.Log("Deleted employee count", result)
	})

	// Get All Employees After Delete
	t.Run("Get All Employees after Delete", func(t *testing.T) {
		result, err := empRepo.FindAllEmployee()
		if err != nil {
			t.Fatal("error while finding all employee", err)
		}
		t.Log("All employees", result)
	})

	// Delete All Employees
	t.Run("Delete All Employees", func(t *testing.T) {
		result, err := empRepo.FindAllEmployee()
		if err != nil {
			t.Fatal("error while finding all employee", err)
		}
		t.Log("All employees", result)
	})
}
