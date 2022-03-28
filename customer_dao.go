package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type ICustomerDao interface {
	create(customer Customer) (uint64, error)
	update(customer Customer) error
	get(id uint64) (Customer, error)
	search(firstName, lastName string) ([]Customer, error)
	list() ([]Customer, error)
}

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
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "user"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "test"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
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

func (d CustomerDao) create(customer Customer) (uint64, error) {
	connection := open()
	var id uint64
	sqlErr := connection.QueryRow("INSERT INTO customers(first_name, last_name, birthdate, gender, email, address) "+
		"VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		customer.firstName,
		customer.lastName,
		customer.birthday,
		customer.gender,
		customer.email,
		customer.address).Scan(&id)
	if sqlErr != nil {
		return 0, sqlErr
	}
	defer connection.Close()
	return id, nil
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
		if versionErr.Error() == "sql: no rows in result set" {
			return errors.New("Cannot find customer by id")
		} else {
			return versionErr
		}
	}

	if version == customer.version {
		_, err := tx.Exec("UPDATE customers SET first_name=$1, last_name=$2, birthdate=$3, gender=$4, email=$5, address=$6, version=$7 WHERE id=$8",
			customer.firstName, customer.lastName, customer.birthday, customer.gender, customer.email, customer.address,
			customer.version+1, customer.id)
		if err != nil {
			return err
		}
		commitErr := tx.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return errors.New("Customer info is not the last version. Please refresh page.")
	}
}

func (d CustomerDao) get(id uint64) (Customer, error) {
	connection := open()
	row := connection.QueryRow("SELECT id, first_name, last_name, birthdate, gender, email, address, version FROM customers WHERE id=$1", id)
	var customer Customer
	err := row.Scan(&customer.id, &customer.firstName, &customer.lastName, &customer.birthday, &customer.gender,
		&customer.email, &customer.address, &customer.version)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return Customer{}, err
	}
	return customer, nil
}

func (d CustomerDao) search(firstName, lastName string) ([]Customer, error) {
	connection := open()
	rows, err := connection.Query("SELECT id, first_name, last_name, birthdate, gender, email, address, version "+
		"FROM customers WHERE first_name=$1 AND last_name=$2 ORDER BY id", firstName, lastName)
	if err != nil {
		return []Customer{}, err
	}
	defer connection.Close()
	defer rows.Close()
	var customers []Customer
	for rows.Next() {
		var customer Customer
		err := rows.Scan(
			&customer.id, &customer.firstName, &customer.lastName, &customer.birthday, &customer.gender,
			&customer.email, &customer.address, &customer.version)
		if err != nil && err.Error() != "sql: no rows in result set" {
			return []Customer{}, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (d CustomerDao) list() ([]Customer, error) {
	connection := open()
	rows, err := connection.Query("SELECT id, first_name, last_name, birthdate, gender, email, address, version FROM customers ORDER BY id")
	if err != nil {
		return []Customer{}, err
	}
	defer connection.Close()
	defer rows.Close()
	var customers []Customer
	for rows.Next() {
		var customer Customer
		err := rows.Scan(
			&customer.id, &customer.firstName, &customer.lastName, &customer.birthday, &customer.gender,
			&customer.email, &customer.address, &customer.version)
		if err != nil && err.Error() != "sql: no rows in result set" {
			return []Customer{}, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}
