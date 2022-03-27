package main

import (
	"errors"
	"testing"
	"time"
)

const (
	Id        = 5
	Version   = 1
	FirstName = "FNAME"
	LastName  = "LNAME"
	Birthday  = "15.05.2000"
	Gender    = "Male"
	Email     = "email"
	Address   = "asdjaskdjlkjalkjsdk"
)

func createDto(id, version uint64) CustomerDto {
	return CustomerDto{
		id,
		FirstName,
		LastName,
		Birthday,
		Gender,
		Email,
		Address,
		version,
		"",
	}
}

type StubDao struct {
	err error
	t   *testing.T
}

func (d *StubDao) create(customer Customer) (uint64, error) {
	if customer.firstName != FirstName {
		d.t.Error("Incorrect first name")
	}
	if customer.lastName != LastName {
		d.t.Error("Incorrect last name")
	}
	birthday, _ := time.Parse(BirthdayFormat, Birthday)
	if customer.birthday != birthday {
		d.t.Error("Incorrect birthday")
	}
	if customer.gender != Gender {
		d.t.Error("Incorrect gender")
	}
	if customer.email != Email {
		d.t.Error("Incorrect email")
	}
	if customer.address != Address {
		d.t.Error("Incorrect address")
	}
	return Id, d.err
}

func (d *StubDao) update(customer Customer) error {
	if customer.id != Id {
		d.t.Error("Incorrect first name")
	}
	if customer.version != Version {
		d.t.Error("Incorrect first name")
	}
	if customer.firstName != FirstName {
		d.t.Error("Incorrect first name")
	}
	if customer.lastName != LastName {
		d.t.Error("Incorrect last name")
	}
	birthday, _ := time.Parse(BirthdayFormat, Birthday)
	if customer.birthday != birthday {
		d.t.Error("Incorrect birthday")
	}
	if customer.gender != Gender {
		d.t.Error("Incorrect gender")
	}
	if customer.email != Email {
		d.t.Error("Incorrect email")
	}
	if customer.address != Address {
		d.t.Error("Incorrect address")
	}
	return d.err
}

func (d *StubDao) get(id uint64) (Customer, error) {
	birthday, _ := time.Parse(BirthdayFormat, Birthday)
	return Customer{
		Id,
		FirstName,
		LastName,
		birthday,
		Gender,
		Email,
		Address,
		Version,
	}, d.err
}

func (d *StubDao) search(firstName, lastName string) []Customer {
	birthday, _ := time.Parse(BirthdayFormat, Birthday)
	if firstName != FirstName {
		d.t.Error("Incorrect first name")
	}
	if lastName != LastName {
		d.t.Error("Incorrect last name")
	}
	return []Customer{Customer{
		Id,
		firstName,
		lastName,
		birthday,
		Gender,
		Email,
		Address,
		Version,
	}}
}

func (d *StubDao) list() []Customer {
	birthday, _ := time.Parse(BirthdayFormat, Birthday)
	return []Customer{Customer{
		Id,
		FirstName,
		LastName,
		birthday,
		Gender,
		Email,
		Address,
		Version,
	}}
}

var victim CustomerService
var stubDao = StubDao{}

func init() {
	victim.dao = &stubDao
}

func checkCustomerDto(customerDto CustomerDto, t *testing.T) {
	if customerDto.Id != Id {
		t.Error("Incorrect first name")
	}
	if customerDto.Version != Version {
		t.Error("Incorrect first name")
	}
	if customerDto.FirstName != FirstName {
		t.Error("Incorrect first name")
	}
	if customerDto.LastName != LastName {
		t.Error("Incorrect last name")
	}
	if customerDto.Birthday != Birthday {
		t.Error("Incorrect birthday")
	}
	if customerDto.Gender != Gender {
		t.Error("Incorrect gender")
	}
	if customerDto.Email != Email {
		t.Error("Incorrect email")
	}
	if customerDto.Address != Address {
		t.Error("Incorrect address")
	}
}

func TestCustomerService_Add(t *testing.T) {
	stubDao.t = t
	stubDao.err = nil

	_, err := victim.Add(CustomerDto{})
	if err == nil {
		t.Error("Cannot add empty customer dto")
	}

	dto := createDto(0, 0)
	id, addErr := victim.Add(dto)
	if addErr != nil {
		t.Error("Cannot add customer: " + addErr.Error())
	}
	if id != Id {
		t.Error("Invalid customer ID")
	}

	stubDao.err = errors.New("test error")
	_, addErr2 := victim.Add(dto)
	if addErr2 != stubDao.err {
		t.Error("Incorrect returned error: " + addErr2.Error())
	}
}

func TestCustomerService_Update(t *testing.T) {
	stubDao.t = t
	stubDao.err = nil

	err := victim.Update(CustomerDto{})
	if err == nil {
		t.Error("Cannot update empty customer")
	}

	dto := createDto(Id, Version)
	updateErr := victim.Update(dto)
	if updateErr != nil {
		t.Error("Cannot update customer: " + updateErr.Error())
	}

	stubDao.err = errors.New("test error")
	updateErr2 := victim.Update(dto)
	if updateErr2 != stubDao.err {
		t.Error("Incorrect returned error: " + updateErr2.Error())
	}
}

func TestCustomerService_Get(t *testing.T) {
	stubDao.t = t

	stubDao.err = errors.New("test error")
	_, err := victim.Get(Id)
	if err != stubDao.err {
		t.Error("Incorrect error from dao")
	}

	stubDao.err = nil
	customerDto, err := victim.Get(Id)
	if err != nil {
		t.Error("Incorrect error from dao")
	}

	checkCustomerDto(customerDto, t)
}

func TestCustomerService_Search(t *testing.T) {
	stubDao.t = t
	customerData := victim.Search(FirstName, LastName)
	if customerData.errMessage != "" {
		t.Error("Incorrect search request")
	}
	if len(customerData.Customers) != 1 {
		t.Error("Incorrect customers list")
	}

	checkCustomerDto(customerData.Customers[0], t)
}

func TestCustomerService_List(t *testing.T) {
	stubDao.t = t
	customerData := victim.List()
	if customerData.errMessage != "" {
		t.Error("Incorrect search request")
	}
	if len(customerData.Customers) != 1 {
		t.Error("Incorrect customers list")
	}

	checkCustomerDto(customerData.Customers[0], t)
}
