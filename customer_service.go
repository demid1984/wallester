package main

import "time"

const (
	BirthdayFormat = "02.01.2006"
)

func handleError(err interface{}) {
	if err != nil {
		panic(err)
	}
}

type CustomerDto struct {
	Id        uint64
	FirstName string
	LastName  string
	Birthday  string
	Gender    string
	Email     string
	Address   string
	Version   uint64
	Error     string
}

func convertEntityToDto(entity Customer) CustomerDto {
	return CustomerDto{
		entity.id,
		entity.firstName,
		entity.lastName,
		entity.birthday.Format(BirthdayFormat),
		entity.gender,
		entity.email,
		entity.address,
		entity.version,
		"",
	}
}

type CustomersData struct {
	errMessage string
	Customers  []CustomerDto
}

type CustomerService struct {
}

func (s CustomerService) List() CustomersData {
	var dao CustomerDao
	entities := dao.list()
	customers := make([]CustomerDto, len(entities))
	for i := range entities {
		customers[i] = convertEntityToDto(entities[i])
	}
	return CustomersData{"", customers}
}

func (s CustomerService) Search(firstName, lastName string) CustomersData {
	var dao CustomerDao
	entities := dao.search(firstName, lastName)
	customers := make([]CustomerDto, len(entities))
	for i := range entities {
		customers[i] = convertEntityToDto(entities[i])
	}
	return CustomersData{"", customers}
}

func (s CustomerService) Add(customer CustomerDto) {
	var dao CustomerDao

	birthday, err := time.Parse(BirthdayFormat, customer.Birthday)
	handleError(err)
	var entity Customer = Customer{
		0,
		customer.FirstName,
		customer.LastName,
		birthday,
		customer.Gender,
		customer.Email,
		customer.Address,
		0}
	dao.create(&entity)
}

func (s CustomerService) Update(customer CustomerDto) error {
	var dao CustomerDao

	birthday, err := time.Parse(BirthdayFormat, customer.Birthday)
	handleError(err)
	var entity Customer = Customer{
		customer.Id,
		customer.FirstName,
		customer.LastName,
		birthday,
		customer.Gender,
		customer.Email,
		customer.Address,
		customer.Version}
	return dao.update(entity)
}

func (s CustomerService) Get(id uint64) CustomerDto {
	var dao CustomerDao
	entity := dao.get(id)
	return convertEntityToDto(entity)
}
