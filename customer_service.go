package main

import "time"

const (
	BirthdayFormat = "02.01.2006"
)

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

type CustomerService struct {
	daoVal ICustomerDao
}

func (s *CustomerService) dao() ICustomerDao {
	if s.daoVal == nil {
		return CustomerDao{}
	} else {
		return s.daoVal
	}
}

func (s *CustomerService) List() ([]CustomerDto, error) {
	entities, err := s.dao().list()
	if err != nil {
		return []CustomerDto{}, err
	}
	customers := make([]CustomerDto, len(entities))
	for i := range entities {
		customers[i] = convertEntityToDto(entities[i])
	}
	return customers, nil
}

func (s *CustomerService) Search(firstName, lastName string, sort SortType) ([]CustomerDto, error) {
	entities, err := s.dao().search(firstName, lastName, sort)
	if err != nil {
		return []CustomerDto{}, err
	}
	customers := make([]CustomerDto, len(entities))
	for i := range entities {
		customers[i] = convertEntityToDto(entities[i])
	}
	return customers, nil
}

func (s *CustomerService) Add(customer CustomerDto) (uint64, error) {
	birthday, err := time.Parse(BirthdayFormat, customer.Birthday)
	if err != nil {
		return 0, err
	}
	var entity = Customer{
		0,
		customer.FirstName,
		customer.LastName,
		birthday,
		customer.Gender,
		customer.Email,
		customer.Address,
		0}
	id, err := s.dao().create(entity)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *CustomerService) Update(customer CustomerDto) error {
	birthday, err := time.Parse(BirthdayFormat, customer.Birthday)
	if err != nil {
		return err
	}
	var entity Customer = Customer{
		customer.Id,
		customer.FirstName,
		customer.LastName,
		birthday,
		customer.Gender,
		customer.Email,
		customer.Address,
		customer.Version}
	return s.dao().update(entity)
}

func (s *CustomerService) Get(id uint64) (CustomerDto, error) {
	entity, err := s.dao().get(id)
	if err != nil {
		return CustomerDto{}, err
	}
	return convertEntityToDto(entity), nil
}
