package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str != "" && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// parse string
	t, err := time.Parse(time.DateOnly, str)
	if err == nil {
		*d = Date(t)
		return nil
	}
	return fmt.Errorf("invalid duration type %T, value: '%s'", b, b)
}

func (d Date) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(d).Format(time.DateOnly))
	return []byte(stamp), nil
}

func (d Date) ConvertToTime() time.Time {
	return time.Time(d)
}

type NullableDate struct {
	isSet bool
	date  *time.Time
}

func (d *NullableDate) UnmarshalJSON(b []byte) error {
	str := string(b)

	if str == "" || str == "null" {
		return nil
	}

	if str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	t, err := time.Parse(time.DateOnly, str)
	if err == nil {
		d.isSet = true
		d.date = &t
		return nil
	}

	return fmt.Errorf("invalid duration type %T, value: '%s'", b, b)
}

func (d NullableDate) MarshalJSON() ([]byte, error) {
	if d.isSet && d.date != nil {
		stamp := fmt.Sprintf("\"%s\"", d.date.Format(time.DateOnly))
		return []byte(stamp), nil
	}

	return []byte{}, nil
}

func (d NullableDate) ConvertToTime() *time.Time {
	if d.isSet && d.date != nil {
		return d.date
	}
	return nil
}

type JWTCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type DashboardEmployeesProject struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type DashboardAccommodation struct {
	Name  string `json:"name"`
	Taken int    `json:"taken"`
	Free  int    `json:"free"`
}

type DashboardCarInspection struct {
	Date               string `json:"date"`
	RegistrationNumber string `json:"registration_number"`
}

type DashboardEmployeePermits struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Document  string `json:"document"`
	Date      string `json:"date"`
}

type Dashboard struct {
	EmployeesProject []DashboardEmployeesProject `json:"employees_project"`
	Accommodations   []DashboardAccommodation    `json:"accommodations"`
	CarInspections   []DashboardCarInspection    `json:"car_inspections"`
	EmployeePermits  []DashboardEmployeePermits  `json:"employee_permits"`
}

type Car struct {
	ID                 int     `json:"id"`
	Model              string  `json:"model"`
	Color              string  `json:"color"`
	RegistrationNumber string  `json:"registration_number"`
	VIN                string  `json:"vin"`
	InspectionFrom     Date    `json:"inspection_from"`
	InspectionTo       Date    `json:"inspection_to"`
	InsuranceFrom      Date    `json:"insurance_from"`
	InsuranceTo        Date    `json:"insurance_to"`
	FleetCardNumber    string  `json:"fleet_card_number"`
	ProjectID          int     `json:"project_id"`
	ProjectName        string  `json:"project_name"`
	Service            Service `json:"service"`
	Leasing            Leasing `json:"leasing"`
}

type CarNumbers struct {
	ID                 int    `json:"id"`
	RegistrationNumber string `json:"registration_number"`
}

type NewCar struct {
	Model              string  `json:"model"`
	Color              string  `json:"color"`
	RegistrationNumber string  `json:"registrationNumber"`
	VIN                string  `json:"vin"`
	InspectionFrom     Date    `json:"inspectionFrom"`
	InspectionTo       Date    `json:"inspectionTo"`
	InsuranceFrom      Date    `json:"insuranceFrom"`
	InsuranceTo        Date    `json:"insuranceTo"`
	FleetCardNumber    string  `json:"fleetCardNumber"`
	IdProject          int     `json:"idProject"`
	Service            Service `json:"service"`
	Leasing            Leasing `json:"leasing"`
}
type Service struct {
	ID          int    `json:"id_service"`
	ServiceName string `json:"serviceName"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
}

type Leasing struct {
	Amount         float64 `json:"amount"`
	MonthlyPayment float64 `json:"monthlyPayment"`
	PaymentDay     int     `json:"paymentDay"`
}

type UpdateCar struct {
	Model              string  `json:"model"`
	Color              string  `json:"color"`
	RegistrationNumber string  `json:"registrationNumber"`
	VIN                string  `json:"vin"`
	InspectionFrom     Date    `json:"inspectionFrom"`
	InspectionTo       Date    `json:"inspectionTo"`
	InsuranceFrom      Date    `json:"insuranceFrom"`
	InsuranceTo        Date    `json:"insuranceTo"`
	FleetCardNumber    string  `json:"fleetCardNumber"`
	IdProject          int     `json:"idProject"`
	Service            Service `json:"service"`
	Leasing            Leasing `json:"leasing"`
}

type Project struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	OfficeAddress  string `json:"office_address"`
	ProjectNIP     string `json:"project_nip"`
	EmployeeAmount int    `json:"employee_amount"`
	FreePlaces     int    `json:"free_places"`
	AmountCars     int    `json:"amount_cars"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Phone          string `json:"phone"`
	Position       string `json:"position"`
}

type ProjectNames struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type NewProject struct {
	Name          string `json:"name"`
	OfficeAddress string `json:"office_address"`
	ProjectNIP    string `json:"project_nip"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	Position      string `json:"position"`
}

type UpdateProject struct {
	Name          string `json:"name"`
	OfficeAddress string `json:"office_address"`
	ProjectNIP    string `json:"project_nip"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	Position      string `json:"position"`
}

type Accommodation struct {
	ID                   int            `json:"id"`
	ProjectID            int            `json:"project_id"`
	ProjectName          string         `json:"project_name"`
	City                 string         `json:"city"`
	AccommodationAddress string         `json:"accommodation_address"`
	NumberOfPlaces       int            `json:"number_of_places"`
	Contact              ContactDetails `json:"contact"`
	Payment              PaymentDetails `json:"payment"`
}

type AccommodationAddresses struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
}

type NewAccommodation struct {
	ProjectID      int                  `json:"idProject"`
	City           string               `json:"city"`
	Address        string               `json:"accommodationAddress"`
	NumberOfPlaces int                  `json:"numberOfPlaces"`
	Contact        UpdateContactDetails `json:"contact"`
	Payment        UpdatePaymentDetails `json:"payment"`
}

type UpdateAccommodation struct {
	ProjectID            int                  `json:"projectId"`
	City                 string               `json:"city"`
	AccommodationAddress string               `json:"accommodationAddress"`
	NumberOfPlaces       int                  `json:"numberOfPlaces"`
	Contact              UpdateContactDetails `json:"contact"`
	Payment              UpdatePaymentDetails `json:"payment"`
}

type ContactDetails struct {
	ID          *int    `json:"id"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	PhoneNumber *string `json:"phone_number"`
}

type UpdateContactDetails struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
}

type PaymentDetails struct {
	ID            *int     `json:"id"`
	Cost          *float64 `json:"cost"`
	Deposit       *float64 `json:"deposit"`
	Contract      *string  `json:"contract"`
	AccountNumber *string  `json:"account_number"`
	PaymentDay    *int     `json:"payment_day"`
}

type UpdatePaymentDetails struct {
	Cost          float64 `json:"cost"`
	Deposit       float64 `json:"deposit"`
	Contract      string  `json:"contract"`
	AccountNumber string  `json:"accountNumber"`
	PaymentDay    int     `json:"paymentDay"`
}

type Employee struct {
	ID               int                  `json:"id"`
	LastName         string               `json:"last_name"`
	FirstName        string               `json:"first_name"`
	PassportNumber   string               `json:"passport_number"`
	Pesel            string               `json:"pesel"`
	Email            string               `json:"email"`
	DateOfBirth      Date                 `json:"date_of_birth"`
	FatherName       string               `json:"father_name"`
	MotherName       string               `json:"mother_name"`
	MaidenName       string               `json:"maiden_name"`
	MotherMaidenName string               `json:"mother_maiden_name"`
	BankAccount      string               `json:"bank_account"`
	AddressPoland    string               `json:"address_poland"`
	HomeAddress      *string              `json:"home_address"`
	Login            *string              `json:"login,omitempty"`
	Password         *string              `json:"password,omitempty"`
	ResidenceCard    ResidenceCardDetails `json:"residence_card"`
	Employment       EmploymentDetails    `json:"employment"`
	Medicals         MedicalDetails       `json:"medicals"`
	ProjectId        int                  `json:"project_id"`
	AccommodationId  *int                 `json:"accommodation_id"`
	CarId            *int                 `json:"car_id"`
}

type ResidenceCardDetails struct {
	Bio   *Date `json:"bio,omitempty"`
	Visa  *Date `json:"visa,omitempty"`
	TCard *Date `json:"tcard,omitempty"`
}

type EmploymentDetails struct {
	ContractType   string `json:"contract_type"`
	StartDate      Date   `json:"start_date"`
	EndDate        *Date  `json:"end_date,omitempty"`
	Authorizations string `json:"authorizations,omitempty"`
}

type MedicalDetails struct {
	OSHValidUntil         Date  `json:"osh_valid_until,omitempty"`
	PsychotestsValidUntil *Date `json:"psychotests_valid_until,omitempty"`
	MedicalValidUntil     Date  `json:"medical_valid_until,omitempty"`
	SanitaryValidUntil    *Date `json:"sanitary_valid_until,omitempty"`
}

type NewEmployee struct {
	LastName         string                  `json:"lastName"`
	FirstName        string                  `json:"firstName"`
	PassportNumber   string                  `json:"passportNumber"`
	Pesel            string                  `json:"pesel"`
	Email            string                  `json:"email"`
	DateOfBirth      Date                    `json:"dateOfBirth"`
	FatherName       string                  `json:"fatherName"`
	MotherName       string                  `json:"motherName"`
	MaidenName       string                  `json:"maidenName"`
	MotherMaidenName string                  `json:"motherMaidenName"`
	BankAccount      string                  `json:"bankAccount"`
	AddressPoland    string                  `json:"addressPoland"`
	HomeAddress      string                  `json:"homeAddress"`
	ResidenceCard    NewResidenceCardDetails `json:"residenceCard"`
	Employment       NewEmploymentDetails    `json:"employment"`
	Medicals         NewMedicalDetails       `json:"medicals"`
	ProjectId        int                     `json:"projectId"`
	AccommodationId  int                     `json:"accommodationId"`
	CarId            int                     `json:"carId"`
}

type UpdateEmployee struct {
	LastName         string                  `json:"lastName"`
	FirstName        string                  `json:"firstName"`
	PassportNumber   string                  `json:"passportNumber"`
	Pesel            string                  `json:"pesel"`
	Email            string                  `json:"email"`
	DateOfBirth      Date                    `json:"dateOfBirth"`
	FatherName       string                  `json:"fatherName"`
	MotherName       string                  `json:"motherName"`
	MaidenName       string                  `json:"maidenName"`
	MotherMaidenName string                  `json:"motherMaidenName"`
	BankAccount      string                  `json:"bankAccount"`
	AddressPoland    string                  `json:"addressPoland"`
	HomeAddress      string                  `json:"homeAddress"`
	ResidenceCard    NewResidenceCardDetails `json:"residenceCard"`
	Employment       NewEmploymentDetails    `json:"employment"`
	Medicals         NewMedicalDetails       `json:"medicals"`
	ProjectId        int                     `json:"projectId"`
	AccommodationId  int                     `json:"accommodationId"`
	CarId            int                     `json:"carId"`
}

type NewResidenceCardDetails struct {
	Bio   NullableDate `json:"bio,omitempty"`
	Visa  NullableDate `json:"visa,omitempty"`
	TCard NullableDate `json:"tcard,omitempty"`
}

type NewEmploymentDetails struct {
	ContractType   string       `json:"contractType"`
	StartDate      Date         `json:"startDate"`
	EndDate        NullableDate `json:"endDate,omitempty"`
	Authorizations string       `json:"authorizations,omitempty"`
}

type NewMedicalDetails struct {
	OSHValidUntil         Date         `json:"oshValidUntil,omitempty"`
	PsychotestsValidUntil NullableDate `json:"psychotestsValidUntil,omitempty"`
	MedicalValidUntil     Date         `json:"medicalValidUntil,omitempty"`
	SanitaryValidUntil    NullableDate `json:"sanitaryValidUntil,omitempty"`
}
