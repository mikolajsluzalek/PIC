package storage

import (
	"api/internal/models"
	"context"
	"github.com/pkg/errors"
)

func (s *Service) Employees(ctx context.Context) ([]models.Employee, error) {
	sql := `
	SELECT 
    e.Id_Employee, 
    e.Last_Name, 
    e.First_Name, 
    e.Passport_Number, 
    e.Pesel, 
    e.Email, 
    e.Date_Of_Birth, 
    e.Father_Name, 
    e.Mother_Name, 
    e.Maiden_Name, 
    e.Mother_Maiden_Name, 
    e.Bank_Account, 
    e.Address_Poland, 
    e.Home_Address, 
    e.Login, 
    e.Password, 
    MAX(rc.Id_Card) AS Id_Card, 
    MAX(rc.Bio) AS Bio, 
    MAX(rc.Visa) AS Visa, 
    MAX(rc.Tcard) AS TCard, 
    em.Id_Employment, 
    em.Contract_Type, 
    em.Start_Date AS Start_Date, 
    em.End_Date AS End_Date, 
    em.Authorizations, 
    MAX(m.Id_Medicals) AS Id_Medicals, 
    MAX(m.OSH_Valid_Until) AS OSH_Valid_Until, 
    MAX(m.Psychotests_Valid_Until) AS Psychotests_Valid_Until, 
    MAX(m.Medical_Valid_Until) AS Medical_Valid_Until, 
    MAX(m.Sanitary_Valid_Until) AS Sanitary_Valid_Until, 
    MAX(p.Id_Project) AS Id_Project, 
    MAX(p.Name) AS Project_Name, 
    MAX(a.Id_Accommodation) AS Id_Accommodation, 
    MAX(a.Accommodation_Address) AS Accommodation_Address, 
    MAX(c.Id_Car) AS Id_Car, 
    MAX(c.Registration_Number) AS Registration_Number
FROM Employee e
LEFT JOIN Residence_Card rc ON e.Id_Employee = rc.Employee_Id 
LEFT JOIN (
    SELECT em1.*
    FROM Employment em1
    INNER JOIN (
        SELECT Id_Employee, MAX(Id_Employment) AS MaxEmploymentId
        FROM Employment
        GROUP BY Id_Employee
    ) em2 ON em1.Id_Employee = em2.Id_Employee AND em1.Id_Employment = em2.MaxEmploymentId
) em ON e.Id_Employee = em.Id_Employee
LEFT JOIN Medicals m ON e.Id_Employee = m.Id_Employee 
LEFT JOIN Employee_Project ep ON e.Id_Employee = ep.Id_Employee 
LEFT JOIN Project p ON ep.Id_Project = p.Id_Project 
LEFT JOIN Employee_Accommodation ea ON e.Id_Employee = ea.Id_Employee 
LEFT JOIN Accommodation a ON ea.Id_Accommodation = a.Id_Accommodation 
LEFT JOIN Employee_Car ec ON e.Id_Employee = ec.Id_Employee 
LEFT JOIN Car c ON ec.Id_Car = c.Id_Car 
GROUP BY 
    e.Id_Employee, 
    e.Last_Name, 
    e.First_Name, 
    e.Passport_Number, 
    e.Pesel, 
    e.Email, 
    e.Date_Of_Birth, 
    e.Father_Name, 
    e.Mother_Name, 
    e.Maiden_Name, 
    e.Mother_Maiden_Name, 
    e.Bank_Account, 
    e.Address_Poland, 
    e.Home_Address, 
    e.Login, 
    e.Password, 
    em.Id_Employment, 
    em.Contract_Type, 
    em.Start_Date, 
    em.End_Date, 
    em.Authorizations;
`

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
			&employee.ResidenceCard.ID,
			&employee.ResidenceCard.Bio,
			&employee.ResidenceCard.Visa,
			&employee.ResidenceCard.TCard,
			&employee.Employment.ID,
			&employee.Employment.ContractType,
			&employee.Employment.StartDate,
			&employee.Employment.EndDate,
			&employee.Employment.Authorizations,
			&employee.Medicals.ID,
			&employee.Medicals.OSHValidUntil,
			&employee.Medicals.PsychotestsValidUntil,
			&employee.Medicals.MedicalValidUntil,
			&employee.Medicals.SanitaryValidUntil,
			&employee.ProjectId,
			&employee.ProjectName,
			&employee.AccommodationId,
			&employee.AccommodationAddress,
			&employee.CarId,
			&employee.CarRegistrationNumber,
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
	sql := "SELECT e.Id_Employee, e.Last_Name, e.First_Name, e.Passport_Number, e.Pesel, e.Email, e.Date_Of_Birth, e.Father_Name, e.Mother_Name, e.Maiden_Name, e.Mother_Maiden_Name, e.Bank_Account, e.Address_Poland, e.Home_Address, e.Login, e.Password, MAX(rc.Id_Card), MAX(rc.Bio), MAX(rc.Visa), MAX(rc.Tcard), em.Id_Employment, em.Contract_Type, em.Start_Date, em.End_Date, em.Authorizations, MAX(m.Id_Medicals), MAX(m.OSH_Valid_Until), MAX(m.Psychotests_Valid_Until), MAX(m.Medical_Valid_Until), MAX(m.Sanitary_Valid_Until), MAX(p.Id_Project), MAX(p.Name), MAX(a.Id_Accommodation), MAX(a.Accommodation_Address), MAX(c.Id_Car), MAX(c.Registration_Number) FROM Employee e LEFT JOIN Residence_Card rc ON e.Id_Employee = rc.Employee_Id LEFT JOIN Employment em ON e.Id_Employee = em.Id_Employee LEFT JOIN Medicals m ON e.Id_Employee = m.Id_Employee LEFT JOIN Employee_Project ep ON e.Id_Employee = ep.Id_Employee LEFT JOIN Project p ON ep.Id_Project = p.Id_Project LEFT JOIN Employee_Accommodation ea ON e.Id_Employee = ea.Id_Employee LEFT JOIN Accommodation a ON ea.Id_Accommodation = a.Id_Accommodation LEFT JOIN Employee_Car ec ON e.Id_Employee = ec.Id_Employee LEFT JOIN Car c ON ec.Id_Car = c.Id_Car GROUP BY e.Id_Employee, e.Last_Name, e.First_Name, e.Passport_Number, e.Pesel, e.Email, e.Date_Of_Birth, e.Father_Name, e.Mother_Name, e.Maiden_Name, e.Mother_Maiden_Name, e.Bank_Account, e.Address_Poland, e.Home_Address, e.Login, e.Password, em.Id_Employment, em.Contract_Type, em.Start_Date, em.End_Date, em.Authorizations HAVING e.Id_Employee = @p1;"

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
		&employee.ResidenceCard.ID,
		&employee.ResidenceCard.Bio,
		&employee.ResidenceCard.Visa,
		&employee.ResidenceCard.TCard,
		&employee.Employment.ID,
		&employee.Employment.ContractType,
		&employee.Employment.StartDate,
		&employee.Employment.EndDate,
		&employee.Employment.Authorizations,
		&employee.Medicals.ID,
		&employee.Medicals.OSHValidUntil,
		&employee.Medicals.PsychotestsValidUntil,
		&employee.Medicals.MedicalValidUntil,
		&employee.Medicals.SanitaryValidUntil,
		&employee.ProjectId,
		&employee.ProjectName,
		&employee.AccommodationId,
		&employee.AccommodationAddress,
		&employee.CarId,
		&employee.CarRegistrationNumber,
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
		newEmployee.Pesel, newEmployee.Email, newEmployee.DateOfBirth,
		newEmployee.FatherName, newEmployee.MotherName, newEmployee.MaidenName,
		newEmployee.MotherMaidenName, newEmployee.BankAccount, newEmployee.AddressPoland,
		newEmployee.HomeAddress,
	).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add employee")
	}

	sql = "INSERT INTO Residence_Card (Employee_Id, Bio, Visa, TCard) VALUES (@p1, @p2, @p3, @p4);"
	_, err = s.DB.ExecContext(ctx, sql, id, newEmployee.ResidenceCard.Bio, newEmployee.ResidenceCard.Visa, newEmployee.ResidenceCard.TCard)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add residence card")
	}

	sql = `
	INSERT INTO Medicals (
		Id_Employee, OSH_Valid_Until, Psychotests_Valid_Until, Medical_Valid_Until, Sanitary_Valid_Until
	) VALUES (@p1, @p2, @p3, @p4, @p5);`
	_, err = s.DB.ExecContext(ctx, sql, id,
		newEmployee.Medicals.OSHValidUntil, newEmployee.Medicals.PsychotestsValidUntil,
		newEmployee.Medicals.MedicalValidUntil, newEmployee.Medicals.SanitaryValidUntil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add medical details")
	}

	sql = `
	INSERT INTO Employment (
		Id_Employee, Contract_Type, Start_Date, End_Date, Authorizations
	) VALUES (@p1, @p2, @p3, @p4, @p5);`
	_, err = s.DB.ExecContext(ctx, sql, id, newEmployee.Employment.ContractType,
		newEmployee.Employment.StartDate, newEmployee.Employment.EndDate, newEmployee.Employment.Authorizations)
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

	sql := `
	UPDATE Employee 
	SET Last_Name = @p1, First_Name = @p2, Passport_Number = @p3, Pesel = @p4, 
		Email = @p5, Date_Of_Birth = @p6, Father_Name = @p7, Mother_Name = @p8, 
		Maiden_Name = @p9, Mother_Maiden_Name = @p10, Bank_Account = @p11, 
		Address_Poland = @p12, Home_Address = @p13
	WHERE Id_Employee = @p14;`
	_, err := s.DB.ExecContext(ctx, sql,
		updateEmployee.LastName, updateEmployee.FirstName, updateEmployee.PassportNumber,
		updateEmployee.Pesel, updateEmployee.Email, updateEmployee.DateOfBirth,
		updateEmployee.FatherName, updateEmployee.MotherName, updateEmployee.MaidenName,
		updateEmployee.MotherMaidenName, updateEmployee.BankAccount, updateEmployee.AddressPoland,
		updateEmployee.HomeAddress, id)
	if err != nil {
		return errors.Wrap(err, "failed to update employee details")
	}

	sql = `
	UPDATE Residence_Card 
	SET Bio = @p1, Visa = @p2, TCard = @p3 
	WHERE Employee_Id = @p4;`
	_, err = s.DB.ExecContext(ctx, sql,
		updateEmployee.ResidenceCard.Bio, updateEmployee.ResidenceCard.Visa,
		updateEmployee.ResidenceCard.TCard, id)
	if err != nil {
		return errors.Wrap(err, "failed to update residence card")
	}

	sql = `
	UPDATE Medicals 
	SET OSH_Valid_Until = @p1, Psychotests_Valid_Until = @p2, 
		Medical_Valid_Until = @p3, Sanitary_Valid_Until = @p4 
	WHERE Id_Employee = @p5;`
	_, err = s.DB.ExecContext(ctx, sql,
		updateEmployee.Medicals.OSHValidUntil, updateEmployee.Medicals.PsychotestsValidUntil,
		updateEmployee.Medicals.MedicalValidUntil, updateEmployee.Medicals.SanitaryValidUntil, id)
	if err != nil {
		return errors.Wrap(err, "failed to update medical details")
	}

	sql = `
		UPDATE Employment 
		SET Contract_Type = @p1, Start_Date = @p2, End_Date = @p3, Authorizations = @p4 
		WHERE Id_Employee = @p6;`
	_, err = s.DB.ExecContext(ctx, sql,
		updateEmployee.Employment.ContractType, updateEmployee.Employment.StartDate, updateEmployee.Employment.EndDate,
		updateEmployee.Employment.Authorizations, id)
	if err != nil {
		return errors.Wrap(err, "failed to update employment details")
	}

	sql = `
	UPDATE Employee_Project 
	SET Id_Project = @p1 
	WHERE Id_Employee = @p2;`
	_, err = s.DB.ExecContext(ctx, sql, updateEmployee.ProjectId, id)
	if err != nil {
		return errors.Wrap(err, "failed to update project details")
	}

	sql = `
	UPDATE Employee_Accommodation 
	SET Id_Accommodation = @p1 
	WHERE Id_Employee = @p2;`
	_, err = s.DB.ExecContext(ctx, sql, updateEmployee.AccommodationId, id)
	if err != nil {
		return errors.Wrap(err, "failed to update accommodation details")
	}

	sql = `
	UPDATE Employee_Car 
	SET Id_Car = @p1 
	WHERE Id_Employee = @p2;`
	_, err = s.DB.ExecContext(ctx, sql, updateEmployee.CarId, id)
	if err != nil {
		return errors.Wrap(err, "failed to update car details")
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
