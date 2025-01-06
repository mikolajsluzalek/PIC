package storage

import (
	"api/internal/models"
	"context"
	"github.com/pkg/errors"
)

func (s *Service) Accommodations(ctx context.Context) ([]models.Accommodation, error) {
	sql := "SELECT a.Id_Accommodation, a.Id_Project, pro.Name, a.City, a.Accommodation_Address, a.Number_Of_Places, c.Id_Contact, c.First_Name, c.Last_Name, c.Phone_Number, p.Id_Payment, p.Cost, p.Deposit, p.Contract, p.Account_Number, p.Payment_Day FROM Accommodation a LEFT JOIN Contact c ON a.Id_Accommodation = c.Id_Accommodation LEFT JOIN Payments p ON a.Id_Accommodation = p.Id_Accommodation LEFT JOIN Project pro ON a.Id_Project = pro.Id_Project;"
	rows, err := s.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for accommodations")
	}
	defer rows.Close()

	results := make([]models.Accommodation, 0)
	for rows.Next() {
		var acc models.Accommodation
		err = rows.Scan(&acc.ID, &acc.ProjectID, &acc.ProjectName, &acc.City, &acc.AccommodationAddress, &acc.NumberOfPlaces, &acc.Contact.ID, &acc.Contact.FirstName, &acc.Contact.LastName, &acc.Contact.PhoneNumber, &acc.Payment.ID, &acc.Payment.Cost, &acc.Payment.Deposit, &acc.Payment.Contract, &acc.Payment.AccountNumber, &acc.Payment.PaymentDay)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		results = append(results, acc)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate rows")
	}

	return results, nil
}

func (s *Service) AddAccommodation(ctx context.Context, newAccommodation models.NewAccommodation) (id int, err error) {

	sql := "INSERT INTO Accommodation (Id_Project, City, Accommodation_Address, Number_Of_Places) VALUES (@p1,@p2,@p3,@p4); SELECT SCOPE_IDENTITY() AS Id_Accommodation;;"
	err = s.DB.QueryRowContext(ctx, sql, newAccommodation.ProjectID, newAccommodation.City, newAccommodation.Address, newAccommodation.NumberOfPlaces).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add accommodation")
	}

	sql = "INSERT INTO Contact (Id_Accommodation, First_Name, Last_Name, Phone_Number) VALUES (@p1,@p2,@p3,@p4);"
	_, err = s.DB.ExecContext(ctx, sql, id, newAccommodation.Contact.FirstName, newAccommodation.Contact.LastName, newAccommodation.Contact.PhoneNumber)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add contact")
	}

	sql = "INSERT INTO Payments (Id_Accommodation, Cost, Deposit, Contract, Account_Number, Payment_Day) VALUES (@p1,@p2,@p3,@p4,@p5,@p6);"
	_, err = s.DB.ExecContext(ctx, sql, id, newAccommodation.Payment.Cost, newAccommodation.Payment.Deposit, newAccommodation.Payment.Contract, newAccommodation.Payment.AccountNumber, newAccommodation.Payment.PaymentDay)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add payment")
	}

	return id, nil
}

func (s *Service) GetAccommodation(ctx context.Context, id int) (models.Accommodation, error) {
	sql := "SELECT a.Id_Accommodation, a.Id_Project, pro.Name, a.City, a.Accommodation_Address, a.Number_Of_Places, c.Id_Contact, c.First_Name, c.Last_Name, c.Phone_Number, p.Id_Payment, p.Cost, p.Deposit, p.Contract, p.Account_Number, p.Payment_Day FROM Accommodation a LEFT JOIN Contact c ON a.Id_Accommodation = c.Id_Accommodation LEFT JOIN Payments p ON a.Id_Accommodation = p.Id_Accommodation LEFT JOIN Project pro ON a.Id_Project = pro.Id_Project WHERE a.Id_Accommodation = @p1;"

	var acc models.Accommodation

	err := s.DB.QueryRowContext(ctx, sql, id).Scan(&acc.ID, &acc.ProjectID, &acc.ProjectName, &acc.City, &acc.AccommodationAddress, &acc.NumberOfPlaces, &acc.Contact.ID, &acc.Contact.FirstName, &acc.Contact.LastName, &acc.Contact.PhoneNumber, &acc.Payment.ID, &acc.Payment.Cost, &acc.Payment.Deposit, &acc.Payment.Contract, &acc.Payment.AccountNumber, &acc.Payment.PaymentDay)

	return acc, errors.Wrap(err, "failed to retrieve accommodation")
}

func (s *Service) RemoveAccommodation(ctx context.Context, id int) error {
	sql := "DELETE FROM Contact WHERE Id_Accommodation = @p1; DELETE FROM Payments WHERE Id_Accommodation = @p1; DELETE FROM Accommodation WHERE Id_Accommodation = @p1; "

	_, err := s.DB.ExecContext(ctx, sql, id)

	return errors.Wrap(err, "failed to remove car")
}

func (s *Service) UpdateAccommodation(ctx context.Context, id int, updateAccommodation models.UpdateAccommodation) error {

	sql := "UPDATE Accommodation SET Id_Project = @p1, City = @p2, Accommodation_Address = @p3, Number_Of_Places = @p4 WHERE Id_Accommodation = @p5;"

	_, err := s.DB.ExecContext(ctx, sql, updateAccommodation.ProjectID, updateAccommodation.City, updateAccommodation.AccommodationAddress, updateAccommodation.NumberOfPlaces, id)

	sql = "UPDATE Contact SET First_Name = @p1, Last_Name = @p2, Phone_Number = @p3 WHERE Id_Accommodation = @p4;"

	_, err = s.DB.ExecContext(ctx, sql, updateAccommodation.Contact.FirstName, updateAccommodation.Contact.LastName, updateAccommodation.Contact.PhoneNumber, id)

	sql = "UPDATE Payments SET Cost = @p1, Deposit = @p2, Contract = @p3, Account_Number = @p4, Payment_Day = @p5 WHERE Id_Accommodation = @p6;"

	_, err = s.DB.ExecContext(ctx, sql, updateAccommodation.Payment.Cost, updateAccommodation.Payment.Deposit, updateAccommodation.Payment.Contract, updateAccommodation.Payment.AccountNumber, updateAccommodation.Payment.PaymentDay, id)

	return errors.Wrap(err, "failed to update accommodation")
}

func (s *Service) GetAccommodationAddresses(ctx context.Context) ([]models.AccommodationAddresses, error) {
	sql := "SELECT a.Id_Accommodation, CONCAT(a.City, ' ', a.Accommodation_Address) AS FullAddress FROM Accommodation a ORDER BY a.Id_Accommodation;"

	rows, err := s.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for accommodation addresses")
	}
	defer rows.Close()

	results := make([]models.AccommodationAddresses, 0)

	for rows.Next() {
		var accommodation models.AccommodationAddresses
		err = rows.Scan(&accommodation.ID, &accommodation.Address)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		results = append(results, accommodation)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate rows")
	}

	return results, errors.Wrap(err, "failed to iterate rows")
}
