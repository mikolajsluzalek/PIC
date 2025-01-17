package storage

import (
	"api/internal/models"
	"context"
	mssql "github.com/microsoft/go-mssqldb"
	"github.com/pkg/errors"
)

func (s *Service) Cars(ctx context.Context) ([]models.Car, error) {
	sql := "SELECT C.Id_Car, C.Model, C.Color, C.Registration_Number, C.VIN_Number, C.Inspection_From, C.Inspection_To, C.Insurance_From, C.Insurance_To, C.Fleet_Card_Number, C.Id_Project, Project.Name FROM Car C LEFT JOIN Project ON C.Id_Project = Project.Id_Project"

	rows, err := s.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for cars")
	}
	defer rows.Close()

	results := make([]models.Car, 0)

	for rows.Next() {
		var car models.Car
		err = rows.Scan(&car.ID, &car.Model, &car.Color, &car.RegistrationNumber, &car.VIN, &car.InspectionFrom, &car.InspectionTo, &car.InsuranceFrom, &car.InsuranceTo, &car.FleetCardNumber, &car.ProjectID, &car.ProjectName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		results = append(results, car)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate rows")
	}

	return results, nil
}

func (s *Service) GetCar(ctx context.Context, id int) (models.Car, error) {
	sql := "SELECT C.Id_Car,C.Model,C.Color,C.Registration_Number,C.VIN_Number,C.Inspection_From,C.Inspection_To,C.Insurance_From,C.Insurance_To,C.Fleet_Card_Number,C.Id_Project,P.Name,S.Id_Service,S.Service_Name,S.Address,S.Phone_Number,L.Amount,L.Monthly_Payment,L.Payment_Day FROM Car C LEFT JOIN Project P ON C.Id_Project=P.Id_Project LEFT JOIN Service S ON C.Id_Car=S.Id_Car LEFT JOIN Leasing L ON C.Id_Car=L.Id_Car WHERE C.Id_Car = @p1;"

	var car models.Car

	err := s.DB.QueryRowContext(ctx, sql, id).Scan(&car.ID, &car.Model, &car.Color, &car.RegistrationNumber, &car.VIN, &car.InspectionFrom, &car.InspectionTo, &car.InsuranceFrom, &car.InsuranceTo, &car.FleetCardNumber, &car.ProjectID, &car.ProjectName, &car.Service.ID, &car.Service.ServiceName, &car.Service.Address, &car.Service.PhoneNumber, &car.Leasing.Amount, &car.Leasing.MonthlyPayment, &car.Leasing.PaymentDay)

	return car, errors.Wrap(err, "failed to retrieve car")
}

func (s *Service) AddCar(ctx context.Context, newCar models.NewCar) (id int, err error) {
	sql := "INSERT INTO Car (Model, Color, Registration_Number, VIN_Number, Inspection_From, Inspection_To, Insurance_From, Insurance_To, Fleet_Card_Number, Id_Project) VALUES (@p1,@p2,@p3,@p4,@p5,@p6,@p7,@p8,@p9,@p10); SELECT SCOPE_IDENTITY() AS Id_Car;"

	err = s.DB.QueryRowContext(ctx, sql, newCar.Model, newCar.Color, newCar.RegistrationNumber, newCar.VIN, mssql.DateTime1(newCar.InspectionFrom), mssql.DateTime1(newCar.InspectionTo), mssql.DateTime1(newCar.InsuranceFrom), mssql.DateTime1(newCar.InsuranceTo), newCar.FleetCardNumber, newCar.IdProject).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add car")
	}

	sql = "INSERT INTO [Service] (Id_Car, Service_name, Address, Phone_Number) VALUES (@p1,@p2,@p3,@p4)"

	_, err = s.DB.ExecContext(ctx, sql, id, newCar.Service.ServiceName, newCar.Service.Address, newCar.Service.PhoneNumber)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add service")
	}

	sql = "INSERT INTO Leasing (Id_Car, Amount, Monthly_Payment, Payment_Day) VALUES (@p1,@p2,@p3,@p4)"

	_, err = s.DB.ExecContext(ctx, sql, id, newCar.Leasing.Amount, newCar.Leasing.MonthlyPayment, newCar.Leasing.PaymentDay)

	return id, errors.Wrap(err, "failed to add leasing")
}

func (s *Service) RemoveCar(ctx context.Context, id int) error {
	sql := "DELETE FROM Leasing WHERE Id_Car = @p1; DELETE FROM Service WHERE Id_Car = @p1;DELETE FROM Car WHERE Id_Car = @p1; "

	_, err := s.DB.ExecContext(ctx, sql, id)

	return errors.Wrap(err, "failed to remove car")
}

func (s *Service) UpdateCar(ctx context.Context, id int, updateCar models.UpdateCar) error {

	sql := "UPDATE Car SET Model = @p1, Color = @p2, Registration_Number = @p3, VIN_Number = @p4, Inspection_From = @p5, Inspection_To = @p6, Insurance_From = @p7, Insurance_To = @p8, Fleet_Card_Number = @p9, Id_Project = @p10 WHERE Id_Car = @p17;UPDATE Service SET Service_Name = @p11, Address = @p12, Phone_Number = @p13 WHERE Id_Car = @p17; UPDATE Leasing SET Amount = @p14, Monthly_Payment = @p15, Payment_Day = @p16 WHERE Id_Car = @p17;"

	_, err2 := s.DB.ExecContext(ctx, sql, updateCar.Model, updateCar.Color, updateCar.RegistrationNumber, updateCar.VIN, mssql.DateTime1(updateCar.InspectionFrom), mssql.DateTime1(updateCar.InspectionTo), mssql.DateTime1(updateCar.InsuranceFrom), mssql.DateTime1(updateCar.InsuranceTo), updateCar.FleetCardNumber, updateCar.IdProject, updateCar.Service.ServiceName, updateCar.Service.Address, updateCar.Service.PhoneNumber, updateCar.Leasing.Amount, updateCar.Leasing.MonthlyPayment, updateCar.Leasing.PaymentDay, id)

	return errors.Wrap(err2, "failed to update car")
}

func (s *Service) GetCarNumbers(ctx context.Context) ([]models.CarNumbers, error) {
	sql := "SELECT c.Id_Car, c.Registration_Number FROM Car c ORDER BY c.Registration_Number;"

	rows, err := s.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for car names")
	}
	defer rows.Close()

	results := make([]models.CarNumbers, 0)

	for rows.Next() {
		var car models.CarNumbers
		err = rows.Scan(&car.ID, &car.RegistrationNumber)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		results = append(results, car)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate rows")
	}

	return results, errors.Wrap(err, "failed to iterate rows")
}
