package main

import (
	"html/template"
	"net/http"
	"os"
)

func handleError(err interface{}) {
	if err != nil {
		panic(err)
	}
}

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
		requestParam, data, sortable, err := service.Search(r)
		var errorMessage string
		if err != nil {
			errorMessage = err.Error()
		}
		var result = map[string]interface{}{
			"ErrMessage": errorMessage,
			"Customers":  data,
			"Sortable":   sortable,
			"Sort":       r.URL.Query().Get("sort"),
			"Request":    requestParam,
			"SortTypes":  Sort,
		}
		customerError := customerTemplate.Execute(w, result)
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
		id, err := service.Add(r)
		var data = map[string]uint64{"Id": 0}
		if err == nil {
			data["Id"] = id
		}
		customerError := customerTemplate.Execute(w, data)
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
