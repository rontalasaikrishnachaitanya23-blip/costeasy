package imports

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

const employeeTemplateSheet = "Employees"
const referenceSheet = "Reference"

// Small helper (excelize does NOT have StringPtr)
func strPtr(s string) *string {
	return &s
}

func GenerateEmployeeTemplate() (*excelize.File, error) {
	f := excelize.NewFile()

	// Create sheets
	f.NewSheet(employeeTemplateSheet)
	f.NewSheet(referenceSheet)
	f.DeleteSheet("Sheet1")

	// Header row
	headers := []string{
		"employee_code", "first_name", "last_name", "email", "phone",
		"country_code", "department_name", "designation_name",
		"joined_at(YYYY-MM-DD)", "base_salary", "salary_currency",
		"employment_status", "termination_date(YYYY-MM-DD)", "termination_reason",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellStr(employeeTemplateSheet, cell, h)
	}

	// Style the header
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#EEEEEE"},
			Pattern: 1,
		},
	})

	f.SetCellStyle(employeeTemplateSheet, "A1", "N1", headerStyle)

	// ------- Reference Sheet -------
	f.SetCellStr(referenceSheet, "A1", "employment_status")
	statuses := []string{"active", "resigned", "terminated"}

	for i, s := range statuses {
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		f.SetCellStr(referenceSheet, cell, s)
	}

	// Define named range
	_ = f.SetDefinedName(&excelize.DefinedName{
		Name:     "EmploymentStatuses",
		RefersTo: fmt.Sprintf("%s!$A$2:$A$%d", referenceSheet, len(statuses)+1),
	})

	// -------- Data Validation (fixed API) --------
	dv := &excelize.DataValidation{
		Sqref:            "L2:L10000", // ✅ lowercase q
		Type:             "list",
		Formula1:         "=EmploymentStatuses",
		AllowBlank:       false,
		ShowErrorMessage: true,
		ErrorTitle:       strPtr("Invalid Value"),
		Error:            strPtr("Choose a value from the dropdown"),
	}

	if err := f.AddDataValidation(employeeTemplateSheet, dv); err != nil {
		return nil, err
	}

	// Freeze top row
	_ = f.SetPanes(employeeTemplateSheet, &excelize.Panes{
		Freeze:      true,
		Split:       true,
		YSplit:      1,
		TopLeftCell: "A2",
	})

	return f, nil
}
