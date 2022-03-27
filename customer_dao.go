package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

const (
	host     = "localhost"
	port     = 9444
	user     = "wallester"
	password = "password"
	dbname   = "test"
)

type CustomerDao struct {
}

type Customer struct {
	id        uint64
	firstName string
	lastName  string
	birthday  time.Time
	gender    string
	email     string
	address   string
	version   uint64
}

func open() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func (d CustomerDao) create(customer *Customer) {
	connection := open()
	_, err := connection.Exec("INSERT INTO customers(first_name, last_name, birthdate, gender, email, address) VALUES($1, $2, $3, $4, $5, $6)",
		customer.firstName,
		customer.lastName,
		customer.birthday,
		customer.gender,
		customer.email,
		customer.address)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
}

func (d CustomerDao) update(customer Customer) error {
	connection := open()
	tx, txErr := connection.Begin()
	if txErr != nil {
		panic(txErr)
	}
	var version uint64
	versionErr := tx.QueryRow("SELECT version FROM customers WHERE id=$1 FOR UPDATE", customer.id).Scan(&version)
	if versionErr != nil {
		panic(versionErr)
	}

	if version == customer.version {
		_, err := tx.Exec("UPDATE customers SET first_name=$1, last_name=$2, birthdate=$3, gender=$4, email=$5, address=$6, version=$7 WHERE id=$8",
			customer.firstName, customer.lastName, customer.birthday, customer.gender, customer.email, customer.address,
			customer.version+1, customer.id)
		if err != nil {
			panic(err)
		}
		commitErr := tx.Commit()
		if commitErr != nil {
			panic(commitErr)
		}
		return nil
	} else {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			panic(rollbackErr)
		}
		return errors.New("Customer info is not the last version. Please refresh page.")
	}
}

func (d CustomerDao) get(id uint64) Customer {
	connection := open()
	row := connection.QueryRow("SELECT id, first_name, last_name, birthdate, gender, email, address, version FROM customers WHERE id=$1", id)
	var customer Customer
	err := row.Scan(&customer.id, &customer.firstName, &customer.lastName, &customer.birthday, &customer.gender,
		&customer.email, &customer.address, &customer.version)
	if err != nil && err.Error() != "sql: no rows in result set" {
		panic(err)
	}
	return customer
}

func (d CustomerDao) search(firstName, lastName string) []Customer {
	connection := open()
	rows, err := connection.Query("SELECT id, first_name, last_name, birthdate, gender, email, address FROM customers "+
		"WHERE first_name=$1 AND last_name=$2", firstName, lastName)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	defer rows.Close()
	var customers []Customer
	for rows.Next() {
		var customer Customer
		err := rows.Scan(
			&customer.id, &customer.firstName, &customer.lastName, &customer.birthday, &customer.gender,
			&customer.email, &customer.address)
		if err != nil && err.Error() != "sql: no rows in result set" {
			panic(err)
		}
		customers = append(customers, customer)
	}
	return customers
}

func (d CustomerDao) list() []Customer {
	connection := open()
	rows, err := connection.Query("SELECT id, first_name, last_name, birthdate, gender, email, address FROM customers")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	defer rows.Close()
	var customers []Customer
	for rows.Next() {
		var customer Customer
		err := rows.Scan(
			&customer.id, &customer.firstName, &customer.lastName, &customer.birthday, &customer.gender,
			&customer.email, &customer.address)
		if err != nil && err.Error() != "sql: no rows in result set" {
			panic(err)
		}
		customers = append(customers, customer)
	}
	return customers
}
