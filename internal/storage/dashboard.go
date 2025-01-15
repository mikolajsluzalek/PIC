package storage

import (
	"context"

	"api/internal/models"
	"github.com/pkg/errors"
)

func (s *Service) DashboardEmployeeProjects(ctx context.Context) ([]models.DashboardEmployeesProject, error) {
	sql := "WITH ProjectCounts AS (SELECT COUNT(*) AS Count, [Name] FROM Employee_Project pp JOIN Project pro ON pro.Id_Project = pp.Id_Project GROUP BY [Name]), RankedProjects AS (SELECT [Name], Count, ROW_NUMBER() OVER (ORDER BY Count DESC) AS RowNum FROM ProjectCounts), TopProjects AS (SELECT [Name], Count FROM RankedProjects WHERE RowNum <= 10 UNION ALL SELECT 'Pozostałe' AS [Name], SUM(Count) AS Count FROM RankedProjects WHERE RowNum > 10) SELECT [Name], Count FROM TopProjects ORDER BY Count DESC;"

	rows, err := s.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for dashboard employee projects")
	}
	defer rows.Close()

	results := make([]models.DashboardEmployeesProject, 0)

	for rows.Next() {
		var (
			name  string
			count int
		)

		err = rows.Scan(&name, &count)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		results = append(results, models.DashboardEmployeesProject{Name: name, Count: count})
	}

	err = rows.Err()

	return results, errors.Wrap(err, "failed to iterate over rows")
}

func (s *Service) Accommodation(ctx context.Context) ([]models.DashboardAccommodation, error) {
	sql := "SELECT TOP 10 p.Name AS ProjectName, SUM(a.Number_Of_Places - COALESCE(ea.OccupiedPlaces, 0)) AS free, SUM(COALESCE(ea.OccupiedPlaces, 0)) AS taken FROM Project p LEFT JOIN Accommodation a ON p.Id_Project = a.Id_Project LEFT JOIN (SELECT ea.Id_Accommodation, COUNT(ea.Id_Employee) AS OccupiedPlaces FROM Employee_Accommodation ea GROUP BY ea.Id_Accommodation) ea ON a.Id_Accommodation = ea.Id_Accommodation GROUP BY p.Name ORDER BY free DESC"

	rows, err := s.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for dashboard accommodation")
	}
	defer rows.Close()

	results := make([]models.DashboardAccommodation, 0)

	for rows.Next() {
		var (
			name  string
			free  int
			taken int
		)

		err = rows.Scan(&name, &free, &taken)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		results = append(results, models.DashboardAccommodation{
			Name:  name,
			Free:  free,
			Taken: taken,
		})
	}

	err = rows.Err()

	return results, errors.Wrap(err, "failed to iterate over rows")
}

func (s *Service) CarInspections(ctx context.Context) ([]models.DashboardCarInspection, error) {
	sql := "Select TOP 5 Inspection_to, Registration_number from car order by Inspection_To"

	rows, err := s.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for dashboard car inspections")
	}
	defer rows.Close()

	results := make([]models.DashboardCarInspection, 0)

	for rows.Next() {
		var (
			date string
			reg  string
		)

		err = rows.Scan(&date, &reg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		results = append(results, models.DashboardCarInspection{
			Date:               date,
			RegistrationNumber: reg,
		})
	}

	err = rows.Err()

	return results, errors.Wrap(err, "failed to iterate over rows")
}

func (s *Service) EmployeePermits(ctx context.Context) ([]models.DashboardEmployeePermits, error) {
	sql := "SELECT TOP 50 E.First_Name, E.Last_Name, 'OSH' AS Document, M.OSH_Valid_Until AS Expiry_Date FROM Employee E INNER JOIN Medicals M ON E.Id_Employee = M.Id_Employee WHERE M.OSH_Valid_Until IS NOT NULL UNION ALL SELECT E.First_Name, E.Last_Name, 'Psychotests' AS Document, M.Psychotests_Valid_Until AS Expiry_Date FROM Employee E INNER JOIN Medicals M ON E.Id_Employee = M.Id_Employee WHERE M.Psychotests_Valid_Until IS NOT NULL UNION ALL SELECT E.First_Name, E.Last_Name, 'Medical' AS Document, M.Medical_Valid_Until AS Expiry_Date FROM Employee E INNER JOIN Medicals M ON E.Id_Employee = M.Id_Employee WHERE M.Medical_Valid_Until IS NOT NULL UNION ALL SELECT E.First_Name, E.Last_Name, 'Sanitary' AS Document, M.Sanitary_Valid_Until AS Expiry_Date FROM Employee E INNER JOIN Medicals M ON E.Id_Employee = M.Id_Employee WHERE M.Sanitary_Valid_Until IS NOT NULL UNION ALL SELECT E.First_Name, E.Last_Name, 'Bio' AS Document, R.Bio AS Expiry_Date FROM Employee E INNER JOIN Residence_Card R ON E.Id_Employee = R.Employee_Id WHERE R.Bio IS NOT NULL UNION ALL SELECT E.First_Name, E.Last_Name, 'Visa' AS Document, R.Visa AS Expiry_Date FROM Employee E INNER JOIN Residence_Card R ON E.Id_Employee = R.Employee_Id WHERE R.Visa IS NOT NULL UNION ALL SELECT E.First_Name, E.Last_Name, 'TCard' AS Document, R.Tcard AS Expiry_Date FROM Employee E INNER JOIN Residence_Card R ON E.Id_Employee = R.Employee_Id WHERE R.Tcard IS NOT NULL ORDER BY Expiry_Date ASC;\n"

	rows, err := s.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for dashboard car inspections")
	}
	defer rows.Close()

	results := make([]models.DashboardEmployeePermits, 0)

	for rows.Next() {
		var (
			first string
			last  string
			doc   string
			date  string
		)

		err = rows.Scan(&first, &last, &doc, &date)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		results = append(results, models.DashboardEmployeePermits{
			FirstName: first,
			LastName:  last,
			Document:  doc,
			Date:      date,
		})
	}

	err = rows.Err()

	return results, errors.Wrap(err, "failed to iterate over rows")
}
