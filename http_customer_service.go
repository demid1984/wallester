package main

import (
	"net/http"
	"strconv"
)

func convertRequestToDto(r *http.Request) CustomerDto {
	var dto CustomerDto
	dto.FirstName = r.PostFormValue("firstName")
	dto.LastName = r.PostFormValue("lastName")
	dto.Birthday = r.PostFormValue("birthday")
	dto.Gender = r.PostFormValue("gender")
	dto.Email = r.PostFormValue("email")
	dto.Address = r.PostFormValue("address")
	versionStr := r.PostFormValue("version")
	version, err := strconv.Atoi(versionStr)
	if err == nil {
		dto.Version = uint64(version)
	}
	return dto
}

type HttpCustomerService struct {
	customerService CustomerService
}

func (s *HttpCustomerService) Add(r *http.Request) (uint64, error) {
	if r.Method == "POST" {
		err := r.ParseForm()
		handleError(err)
		return s.customerService.Add(convertRequestToDto(r))
	}
	return 0, nil
}

func (s *HttpCustomerService) Update(r *http.Request) CustomerDto {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	var customer CustomerDto
	if err == nil {
		var serviceErr error
		if r.Method == "POST" {
			parseFormErr := r.ParseForm()
			if parseFormErr != nil {
				customer.Error = parseFormErr.Error()
			}
			customerDto := convertRequestToDto(r)
			customerDto.Id = uint64(id)
			serviceErr = s.customerService.Update(customerDto)
		}
		customer, serviceErr = s.customerService.Get(uint64(id))
		if serviceErr != nil {
			customer.Error = serviceErr.Error()
		}
	}
	return customer
}

func (s *HttpCustomerService) Search(r *http.Request) (map[string]string, []CustomerDto, bool, error) {
	firstName := r.URL.Query().Get("firstName")
	lastName := r.URL.Query().Get("lastName")
	if len(firstName) > 0 && len(lastName) > 0 {
		request := map[string]string{"FirstName": firstName, "LastName": lastName}
		sortType := FindSortByCode(r.URL.Query().Get("sort"))
		data, err := s.customerService.Search(firstName, lastName, sortType)
		return request, data, true, err
	} else {
		data, err := s.customerService.List()
		return map[string]string{}, data, false, err
	}
}
