package main

type SortType struct {
	Code  string
	field string
}

var Sort = struct {
	Unordered SortType
	FirstName SortType
	LastName  SortType
	Birthday  SortType
	Gender    SortType
	Email     SortType
	Address   SortType
}{
	Unordered: SortType{"UNORDERED", "id"},
	FirstName: SortType{"FIRST_NAME", "first_name"},
	LastName:  SortType{"LAST_NAME", "last_name"},
	Birthday:  SortType{"BIRTHDAY", "birthdate"},
	Gender:    SortType{"GENDER", "gender"},
	Email:     SortType{"EMAIL", "email"},
	Address:   SortType{"ADDRESS", "address"},
}

func FindSortByCode(code string) SortType {
	switch code {
	case Sort.FirstName.Code:
		return Sort.FirstName
	case Sort.LastName.Code:
		return Sort.LastName
	case Sort.Birthday.Code:
		return Sort.Birthday
	case Sort.Gender.Code:
		return Sort.Gender
	case Sort.Email.Code:
		return Sort.Email
	case Sort.Address.Code:
		return Sort.Address
	default:
		return Sort.Unordered
	}
}
