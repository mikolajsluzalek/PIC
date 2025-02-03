package storage

import (
	"context"

	"api/internal/models"
	"github.com/pkg/errors"
)

func (s *Service) Projects(ctx context.Context) ([]models.Project, error) {
	sql := "SELECT p.Id_Project AS ProjectId, p.Name AS ProjectName, p.Office_Address AS ProjectAddress, p.Project_NIP AS ProjectNIP, COALESCE(emp_data.EmployeeCount, 0) AS EmployeeCount, COALESCE(acc_data.FreeAccommodationPlaces, 0) AS FreeAccommodationPlaces, COALESCE(car_data.CarCount, 0) AS CarCount FROM Project p LEFT JOIN (SELECT ep.Id_Project, COUNT(DISTINCT ep.Id_Employee) AS EmployeeCount FROM Employee_Project ep GROUP BY ep.Id_Project) emp_data ON p.Id_Project = emp_data.Id_Project LEFT JOIN (SELECT a.Id_Project, SUM(a.Number_Of_Places - COALESCE(assigned.CountAssignedEmployees, 0)) AS FreeAccommodationPlaces FROM Accommodation a LEFT JOIN (SELECT ea.Id_Accommodation, COUNT(ea.Id_Employee) AS CountAssignedEmployees FROM Employee_Accommodation ea GROUP BY ea.Id_Accommodation) assigned ON a.Id_Accommodation = assigned.Id_Accommodation GROUP BY a.Id_Project) acc_data ON p.Id_Project = acc_data.Id_Project LEFT JOIN (SELECT c.Id_Project, COUNT(DISTINCT c.Id_Car) AS CarCount FROM Car c GROUP BY c.Id_Project) car_data ON p.Id_Project = car_data.Id_Project;"

	rows, err := s.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for projects")
	}
	defer rows.Close()

	results := make([]models.Project, 0)

	for rows.Next() {
		var project models.Project
		err = rows.Scan(&project.ID, &project.Name, &project.OfficeAddress, &project.ProjectNIP, &project.EmployeeAmount, &project.FreePlaces, &project.AmountCars)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		results = append(results, project)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate rows")
	}

	return results, nil
}

func (s *Service) GetProject(ctx context.Context, id int) (models.Project, error) {
	sql := "select p.Id_Project, p.name, p.office_address, p.project_NIP, c.First_name, c.Last_name, c.Phone, c.Position from project p left join Contact_Person c on p.Id_Project = c.Id_Project where p.id_project = @p1"

	var project models.Project

	err := s.DB.QueryRowContext(ctx, sql, id).Scan(&project.ID, &project.Name, &project.OfficeAddress, &project.ProjectNIP, &project.FirstName, &project.LastName, &project.Phone, &project.Position)

	return project, errors.Wrap(err, "failed to retrieve project")
}

func (s *Service) AddProject(ctx context.Context, newProject models.NewProject) (id int, err error) {
	sql := "INSERT INTO Project (Name, Office_Address, Project_NIP) VALUES (@p1,@p2,@p3); SELECT SCOPE_IDENTITY() AS Id_Project;"

	err = s.DB.QueryRowContext(ctx, sql, newProject.Name, newProject.OfficeAddress, newProject.ProjectNIP).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add project")
	}

	sql = "INSERT INTO Contact_Person (Id_Project, First_Name, Last_Name, Phone, Position) VALUES (@p1,@p2,@p3,@p4,@p5);"
	_, err = s.DB.ExecContext(ctx, sql, id, newProject.FirstName, newProject.LastName, newProject.Phone, newProject.Position)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add contact person")
	}

	return id, nil
}

func (s *Service) RemoveProject(ctx context.Context, id int) error {

	sql := "DELETE FROM Payments WHERE Id_Accommodation IN (SELECT Id_Accommodation FROM Accommodation WHERE Id_Project = @p1);" +
		"DELETE FROM Contact WHERE Id_Accommodation IN (SELECT Id_Accommodation FROM Accommodation WHERE Id_Project = @p1);" +
		"DELETE FROM Employee_Accommodation WHERE Id_Accommodation IN (SELECT Id_Accommodation FROM Accommodation WHERE Id_Project = @p1)" +
		"DELETE FROM Accommodation WHERE Id_Project = @p1;" +
		"DELETE FROM Leasing WHERE Id_Car IN (SELECT Id_Car FROM Car WHERE Id_Project = @p1);" +
		"DELETE FROM Service WHERE Id_Car IN (SELECT Id_Car FROM Car WHERE Id_Project = @p1);" +
		"DELETE FROM Employee_Car WHERE Id_Car IN (SELECT Id_Car FROM Car WHERE Id_Project = @p1);" +
		"DELETE FROM Car WHERE Id_Project = @p1;" +
		"DELETE FROM Employee_Project WHERE Id_Project = @p1;" +
		"DELETE FROM Contact_Person WHERE Id_Project = @p1;" +
		"DELETE FROM Project WHERE Id_Project = @p1;"

	_, err := s.DB.ExecContext(ctx, sql, id)

	return errors.Wrap(err, "failed to remove project")
}

func (s *Service) UpdateProject(ctx context.Context, id int, updateProject models.UpdateProject) error {
	sql := "UPDATE Project SET [Name] = @p1, Office_Address = @p2, Project_NIP = @p3 WHERE Id_Project = @p8; UPDATE Contact_Person SET First_Name = @p4, Last_Name = @p5, Phone = @p6, Position = @p7 WHERE Id_Project = @p8;"

	_, err := s.DB.ExecContext(ctx, sql, updateProject.Name, updateProject.OfficeAddress, updateProject.ProjectNIP, updateProject.FirstName, updateProject.LastName, updateProject.Phone, updateProject.Position, id)

	return errors.Wrap(err, "failed to update project")
}

func (s *Service) GetProjectNames(ctx context.Context) ([]models.ProjectNames, error) {
	sql := "SELECT p.Id_project, p.Name FROM Project p ORDER BY p.Id_Project;"

	rows, err := s.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for project names")
	}
	defer rows.Close()

	results := make([]models.ProjectNames, 0)

	for rows.Next() {
		var project models.ProjectNames
		err = rows.Scan(&project.ID, &project.Name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		results = append(results, project)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate rows")
	}

	return results, errors.Wrap(err, "failed to iterate rows")
}
