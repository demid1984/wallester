package main

import (
	"html/template"
	"net/http"
	"os"
)

func mainHandler() {
	mainTemplate := template.Must(template.ParseFiles("static/main.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mainError := mainTemplate.Execute(w, nil)
		handleError(mainError)
	})
}

func listHandler() {
	customerTemplate := template.Must(template.ParseFiles("static/customers.html"))
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		var service HttpCustomerService
		data := service.Search(r)
		customerError := customerTemplate.Execute(w, data)
		handleError(customerError)
	})
}

func editHandler() {
	customerTemplate := template.Must(template.ParseFiles("static/edit.html"))
	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		var service HttpCustomerService
		customer := service.Update(r)
		customerError := customerTemplate.Execute(w, customer)
		handleError(customerError)
	})
}

func addHandler() {
	customerTemplate := template.Must(template.ParseFiles("static/add.html"))
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		var service HttpCustomerService
		service.Add(r)
		customerError := customerTemplate.Execute(w, nil)
		handleError(customerError)
	})
}

func main() {

	listHandler()
	mainHandler()
	editHandler()
	addHandler()

	port := os.Getenv("HTTP_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	serverError := http.ListenAndServe(":"+port, nil)
	handleError(serverError)
}
