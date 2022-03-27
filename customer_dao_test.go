package main

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	_ "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	_ "github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
	"time"
)

func CreateTestContainer(ctx context.Context, dbname string) (testcontainers.Container, error) {
	user := "user"
	password := "password"
	var env = map[string]string{
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_USER":     "user",
		"POSTGRES_DB":       dbname,
	}
	seedDataPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	mountPath := seedDataPath + "/env/sql"
	var port = "5432/tcp"
	dbURL := func(port nat.Port) string {
		return fmt.Sprintf("postgres://user:password@localhost:%s/%s?sslmode=disable", port.Port(), dbname)
	}
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13.4-alpine",
		ExposedPorts: []string{port},
		Cmd:          []string{"postgres", "-c", "fsync=off"},
		Env:          env,
		BindMounts: map[string]string{
			"/docker-entrypoint-initdb.d/": mountPath,
		},
		WaitingFor: wait.ForSQL(nat.Port(port), "postgres", dbURL).Timeout(time.Second * 5),
		AutoRemove: true,
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return container, fmt.Errorf("failed to start container: %s", err)
	}
	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return container, fmt.Errorf("failed to get container external port: %s", err)
	}
	os.Setenv("DB_USER", user)
	os.Setenv("DB_PASSWORD", password)
	os.Setenv("DB_PORT", mappedPort.Port())
	os.Setenv("DB_NAME", dbname)
	return container, nil
}

func checkCustomer(t *testing.T, customer Customer) {
	if customer.id <= 0 {
		t.Error("Customer id is not read")
	}
	if customer.version <= 0 {
		t.Error("Customer version is not read")
	}
	if len(customer.firstName) == 0 {
		t.Error("Customer first name is not read")
	}
	if len(customer.lastName) == 0 {
		t.Error("Customer last name is not read")
	}
	if len(customer.address) == 0 {
		t.Error("Customer address is not read")
	}
	if len(customer.email) == 0 {
		t.Error("Customer email is not read")
	}
	if customer.birthday.IsZero() {
		t.Error("Customer birthday is not read")
	}
}

func TestWithPostgreSql(t *testing.T) {
	ctx := context.Background()
	container, err := CreateTestContainer(ctx, "testdb")
	if err != nil {
		panic(err)
	}

	var daoService CustomerDao
	customers := daoService.list()
	if len(customers) != 3 {
		t.Error("Cannot load predefined customers from database")
	}
	checkCustomer(t, customers[0])
	checkCustomer(t, customers[1])
	checkCustomer(t, customers[2])

	customer := customers[2]
	searchResult := daoService.search(customer.firstName, customer.lastName)
	if len(searchResult) != 1 {
		t.Error("Cannot search customer by first name and last name")
	}
	checkCustomer(t, searchResult[0])

	searchResult = daoService.search("111", "222")
	if len(searchResult) != 0 {
		t.Error("Search customer does not work correctly")
	}

	customerById, _ := daoService.get(customer.id)
	checkCustomer(t, customerById)

	var notExistsCustomer = Customer{}
	updateError := daoService.update(notExistsCustomer)
	if updateError == nil {
		t.Error("No error after updating non-exists customer")
	}

	updateError = daoService.update(customer)
	if updateError != nil {
		t.Error("Cannot update customer")
	}

	updateError = daoService.update(customer)
	if updateError == nil {
		t.Error("Update customer with incorrect version")
	}

	var newCustomer = Customer{
		0,
		"FNAME",
		"LNAME",
		time.Now(),
		"Male",
		"test111@test.ee",
		"qwerty asdfghh zxcvbn",
		0,
	}
	newId, createErr := daoService.create(newCustomer)
	if createErr != nil {
		t.Error("Cannot create customer")
	}
	if newId <= 0 {
		t.Error("Incorrect new customer id")
	}

	defer container.Terminate(ctx)
}
