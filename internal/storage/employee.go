package storage

import (
	"api/internal/models"
	"context"
	"database/sql"
	mssql "github.com/microsoft/go-mssqldb"
	"github.com/pkg/errors"
)

func (s *Service) Employees(ctx context.Context) ([]models.Employee, error) {
	sql := `SELECT e.Id_Employee, e.First_name, e.Last_name, e.Pesel, e.Passport_number, e.Date_of_birth from employee e`

	rows, err := s.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for employees")
	}
	defer rows.Close()

	results := make([]models.Employee, 0)

	for rows.Next() {
		var employee models.Employee

		err = rows.Scan(
			&employee.ID,
			&employee.LastName,
			&employee.FirstName,
			&employee.Pesel,
			&employee.PassportNumber,
			&employee.DateOfBirth,
		)

		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		results = append(results, employee)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate rows")
	}

	return results, nil
}

func (s *Service) GetEmployee(ctx context.Context, id int) (models.Employee, error) {
	sql := `SELECT TOP 1 e.Id_Employee, e.Last_Name, e.First_Name, e.Passport_Number, e.Pesel, e.Email, e.Date_Of_Birth, e.Father_Name, e.Mother_Name, e.Maiden_Name, e.Mother_Maiden_Name, e.Bank_Account, e.Address_Poland, e.Home_Address, e.Login, e.Password, m.OSH_Valid_Until, m.Psychotests_Valid_Until, m.Medical_Valid_Until, m.Sanitary_Valid_Until, em.Contract_Type, em.Start_Date, em.End_Date, em.Authorizations, rc.Bio, rc.Visa, rc.Tcard, ea.Id_Accommodation, ep.Id_Project, ec.Id_Car FROM employee e LEFT JOIN (SELECT OSH_Valid_Until, Psychotests_Valid_Until, Medical_Valid_Until, Sanitary_Valid_Until, Id_employee FROM Medicals WHERE Id_employee = @p1) m ON e.Id_Employee = m.Id_employee LEFT JOIN (SELECT Contract_Type, Start_Date, End_Date, Authorizations, Id_Employee FROM Employment WHERE Id_Employee = @p1) em ON e.Id_Employee = em.Id_Employee LEFT JOIN (SELECT Bio, Visa, Tcard, Employee_Id FROM Residence_Card WHERE Employee_Id = @p1) rc ON e.Id_Employee = rc.Employee_Id LEFT JOIN (SELECT Id_Accommodation, Id_Employee FROM Employee_Accommodation WHERE Id_Employee = @p1) ea ON e.Id_Employee = ea.Id_Employee LEFT JOIN (SELECT Id_Project, Id_Employee FROM Employee_Project WHERE Id_Employee = @p1) ep ON e.Id_Employee = ep.Id_Employee LEFT JOIN (SELECT Id_Car, Id_Employee FROM Employee_Car WHERE Id_Employee = @p1) ec ON e.Id_Employee = ec.Id_Employee WHERE e.Id_Employee = @p1;`
	var employee models.Employee

	err := s.DB.QueryRowContext(ctx, sql, id).Scan(
		&employee.ID,
		&employee.LastName,
		&employee.FirstName,
		&employee.PassportNumber,
		&employee.Pesel,
		&employee.Email,
		&employee.DateOfBirth,
		&employee.FatherName,
		&employee.MotherName,
		&employee.MaidenName,
		&employee.MotherMaidenName,
		&employee.BankAccount,
		&employee.AddressPoland,
		&employee.HomeAddress,
		&employee.Login,
		&employee.Password,
		&employee.Medicals.OSHValidUntil,
		&employee.Medicals.PsychotestsValidUntil,
		&employee.Medicals.MedicalValidUntil,
		&employee.Medicals.SanitaryValidUntil,
		&employee.Employment.ContractType,
		&employee.Employment.StartDate,
		&employee.Employment.EndDate,
		&employee.Employment.Authorizations,
		&employee.ResidenceCard.Bio,
		&employee.ResidenceCard.Visa,
		&employee.ResidenceCard.TCard,
		&employee.AccommodationId,
		&employee.ProjectId,
		&employee.CarId,
	)

	return employee, errors.Wrap(err, "failed to retrieve employee")
}

func (s *Service) AddEmployee(ctx context.Context, newEmployee models.NewEmployee) (id int, err error) {

	sql := `
	INSERT INTO Employee (
		Last_Name, First_Name, Passport_Number, Pesel, Email, Date_Of_Birth, 
		Father_Name, Mother_Name, Maiden_Name, Mother_Maiden_Name, Bank_Account, 
		Address_Poland, Home_Address
	) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13);
	SELECT SCOPE_IDENTITY() AS Id_Employee;`
	err = s.DB.QueryRowContext(ctx, sql,
		newEmployee.LastName, newEmployee.FirstName, newEmployee.PassportNumber,
		newEmployee.Pesel, newEmployee.Email, mssql.DateTime1(newEmployee.DateOfBirth),
		newEmployee.FatherName, newEmployee.MotherName, newEmployee.MaidenName,
		newEmployee.MotherMaidenName, newEmployee.BankAccount, newEmployee.AddressPoland,
		newEmployee.HomeAddress,
	).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add employee")
	}

	sql = "INSERT INTO Residence_Card (Employee_Id, Bio, Visa, TCard) VALUES (@p1, @p2, @p3, @p4);"

	_, err = s.DB.ExecContext(ctx, sql, id, newEmployee.ResidenceCard.Bio.ConvertToTime(), newEmployee.ResidenceCard.Visa.ConvertToTime(), newEmployee.ResidenceCard.TCard.ConvertToTime())
	if err != nil {
		return 0, errors.Wrap(err, "failed to add residence card")
	}

	sql = `
	INSERT INTO Medicals (
		Id_Employee, OSH_Valid_Until, Psychotests_Valid_Until, Medical_Valid_Until, Sanitary_Valid_Until
	) VALUES (@p1, @p2, @p3, @p4, @p5);`
	_, err = s.DB.ExecContext(ctx, sql, id,
		mssql.DateTime1(newEmployee.Medicals.OSHValidUntil), newEmployee.Medicals.PsychotestsValidUntil.ConvertToTime(),
		mssql.DateTime1(newEmployee.Medicals.MedicalValidUntil), newEmployee.Medicals.SanitaryValidUntil.ConvertToTime())
	if err != nil {
		return 0, errors.Wrap(err, "failed to add medical details")
	}

	sql = `
	INSERT INTO Employment (
		Id_Employee, Contract_Type, Start_Date, End_Date, Authorizations
	) VALUES (@p1, @p2, @p3, @p4, @p5);`
	_, err = s.DB.ExecContext(ctx, sql, id, newEmployee.Employment.ContractType,
		mssql.DateTime1(newEmployee.Employment.StartDate), newEmployee.Employment.EndDate.ConvertToTime(), newEmployee.Employment.Authorizations)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add employment details")
	}

	sql = "INSERT INTO Employee_Project (Id_Employee, Id_Project) VALUES (@p1, @p2);"
	_, err = s.DB.ExecContext(ctx, sql, id, newEmployee.ProjectId)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add project")
	}

	sql = "INSERT INTO Employee_Accommodation (Id_Employee, Id_Accommodation) VALUES (@p1, @p2);"
	_, err = s.DB.ExecContext(ctx, sql, id, newEmployee.AccommodationId)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add accommodation")
	}

	sql = "INSERT INTO Employee_Car (Id_Employee, Id_Car) VALUES (@p1, @p2);"
	_, err = s.DB.ExecContext(ctx, sql, id, newEmployee.CarId)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add car")
	}

	return id, nil
}

func (s *Service) UpdateEmployee(ctx context.Context, id int, updateEmployee models.UpdateEmployee) error {
	query := `
	UPDATE Employee 
	SET Last_Name = @p1, First_Name = @p2, Passport_Number = @p3, Pesel = @p4, 
		Email = @p5, Date_Of_Birth = @p6, Father_Name = @p7, Mother_Name = @p8, 
		Maiden_Name = @p9, Mother_Maiden_Name = @p10, Bank_Account = @p11, 
		Address_Poland = @p12, Home_Address = @p13
	WHERE Id_Employee = @p14;`
	_, err := s.DB.ExecContext(ctx, query,
		updateEmployee.LastName,
		updateEmployee.FirstName,
		updateEmployee.PassportNumber,
		updateEmployee.Pesel,
		updateEmployee.Email,
		mssql.DateTime1(updateEmployee.DateOfBirth),
		updateEmployee.FatherName,
		updateEmployee.MotherName,
		updateEmployee.MaidenName,
		updateEmployee.MotherMaidenName,
		updateEmployee.BankAccount,
		updateEmployee.AddressPoland,
		updateEmployee.HomeAddress,
		id)
	if err != nil {
		return errors.Wrap(err, "failed to update employee details")
	}

	query = `
	UPDATE Residence_Card 
	SET Bio = @p1, Visa = @p2, TCard = @p3 
	WHERE Employee_Id = @p4;`
	_, err = s.DB.ExecContext(ctx, query,
		updateEmployee.ResidenceCard.Bio.ConvertToTime(), updateEmployee.ResidenceCard.Visa.ConvertToTime(),
		updateEmployee.ResidenceCard.TCard.ConvertToTime(), id)
	if err != nil {
		return errors.Wrap(err, "failed to update residence card")
	}

	query = `
	UPDATE Medicals 
	SET OSH_Valid_Until = @p1, Psychotests_Valid_Until = @p2, 
		Medical_Valid_Until = @p3, Sanitary_Valid_Until = @p4 
	WHERE Id_Employee = @p5;`
	_, err = s.DB.ExecContext(ctx, query,
		mssql.DateTime1(updateEmployee.Medicals.OSHValidUntil), updateEmployee.Medicals.PsychotestsValidUntil.ConvertToTime(),
		mssql.DateTime1(updateEmployee.Medicals.MedicalValidUntil), updateEmployee.Medicals.SanitaryValidUntil.ConvertToTime(), id)
	if err != nil {
		return errors.Wrap(err, "failed to update medical details")
	}

	query = `
		UPDATE Employment 
		SET Contract_Type = @p1, Start_Date = @p2, End_Date = @p3, Authorizations = @p4 
		WHERE Id_Employee = @p5;`
	_, err = s.DB.ExecContext(ctx, query,
		updateEmployee.Employment.ContractType, mssql.DateTime1(updateEmployee.Employment.StartDate), updateEmployee.Employment.EndDate.ConvertToTime(),
		updateEmployee.Employment.Authorizations, id)
	if err != nil {
		return errors.Wrap(err, "failed to update employment details")
	}

	query = `
	UPDATE Employee_Project 
	SET Id_Project = @p1 
	WHERE Id_Employee = @p2;`
	_, err = s.DB.ExecContext(ctx, query, updateEmployee.ProjectId, id)
	if err != nil {
		return errors.Wrap(err, "failed to update project details")
	}

	var rowExists bool

	query = `SELECT CASE WHEN EXISTS (SELECT 1 FROM Employee_Accommodation WHERE Id_Employee = @p1) THEN 1 ELSE 0 END AS RowExists;`
	err = s.DB.QueryRowContext(ctx, query, id).Scan(
		&rowExists,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(err, "failed to get car details")
	}

	if rowExists {
		query = `UPDATE Employee_Accommodation SET Id_Accommodation = @p1 WHERE Id_Employee = @p2;`
		_, err = s.DB.ExecContext(ctx, query, updateEmployee.AccommodationId, id)
		if err != nil {
			return errors.Wrap(err, "failed to update accommodation details")
		}
	} else {
		query = "INSERT INTO Employee_Accommodation (Id_Employee, Id_Accommodation) VALUES (@p1, @p2);"
		_, err = s.DB.ExecContext(ctx, query, id, updateEmployee.AccommodationId)
		if err != nil {
			return errors.Wrap(err, "failed to insert car details")
		}
	}

	query = `SELECT CASE WHEN EXISTS (SELECT 1 FROM Employee_Car WHERE Id_Employee = @p1) THEN 1 ELSE 0 END AS RowExists;`
	err = s.DB.QueryRowContext(ctx, query, id).Scan(
		&rowExists,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(err, "failed to get car details")
	}

	if rowExists {
		query = `UPDATE Employee_Car SET Id_Car = @p1 WHERE Id_Employee = @p2;`
		_, err = s.DB.ExecContext(ctx, query, updateEmployee.CarId, id)
		if err != nil {
			return errors.Wrap(err, "failed to update car details")
		}
	} else {
		query = "INSERT INTO Employee_Car (Id_Employee, Id_Car) VALUES (@p1, @p2);"
		_, err = s.DB.ExecContext(ctx, query, id, updateEmployee.CarId)
		if err != nil {
			return errors.Wrap(err, "failed to insert car details")
		}
	}

	return nil
}

func (s *Service) RemoveEmployee(ctx context.Context, id int) error {
	sql := `
		DELETE FROM Employee_Car WHERE Id_Employee = @p1;
		DELETE FROM Employee_Accommodation WHERE Id_Employee = @p1;
		DELETE FROM Employee_Project WHERE Id_Employee = @p1;
		DELETE FROM Employment WHERE Id_Employee = @p1;
		DELETE FROM Medicals WHERE Id_Employee = @p1;
		DELETE FROM Residence_Card WHERE Employee_Id = @p1;
		DELETE FROM Employee WHERE Id_Employee = @p1;
	`

	_, err := s.DB.ExecContext(ctx, sql, id)

	return errors.Wrap(err, "failed to remove employee")
}
