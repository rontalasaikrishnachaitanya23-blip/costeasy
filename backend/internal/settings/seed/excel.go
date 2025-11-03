// backend/internal/settings/seed/excel.go
package seed

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain" // Update with your actual module path
)

type Seeder struct {
	pool *pgxpool.Pool
}

// NewSeeder creates a new Seeder instance
func NewSeeder(pool *pgxpool.Pool) *Seeder {
	return &Seeder{pool: pool}
}

// SeedOrganizationsFromExcel loads organizations from Excel file with domain validation
func (s *Seeder) SeedOrganizationsFromExcel(ctx context.Context, filePath string) error {
	log.Printf("Loading organizations from Excel: %s", filePath)

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Organizations")
	if err != nil {
		return fmt.Errorf("failed to read Organizations sheet: %w", err)
	}

	if len(rows) < 2 {
		return fmt.Errorf("excel file must have header row and at least one data row")
	}

	count := 0
	skipped := 0

	for i, row := range rows[1:] {
		rowNum := i + 2

		if len(row) < 8 {
			log.Printf("⚠️  Skipping row %d: insufficient columns", rowNum)
			skipped++
			continue
		}

		// Parse isActive
		isActive := true
		if len(row) > 9 && strings.ToLower(strings.TrimSpace(row[9])) == "false" {
			isActive = false
		}

		// Create domain object for validation
		org := &domain.Organization{
			ID:              uuid.New(), // Temporary ID for validation
			Name:            strings.TrimSpace(row[0]),
			Type:            domain.OrganizationType(strings.ToUpper(strings.TrimSpace(row[1]))),
			Country:         strings.ToUpper(strings.TrimSpace(row[2])),
			Emirate:         domain.UAEmirate(strings.ToUpper(strings.TrimSpace(row[3]))),
			Area:            strings.TrimSpace(row[4]),
			Currency:        strings.ToUpper(strings.TrimSpace(row[5])),
			LicenseNumber:   strings.TrimSpace(row[6]),
			EstablishmentID: strings.TrimSpace(row[7]),
			Description:     getColumnOrDefault(row, 8, ""),
			IsActive:        isActive,
		}

		// Validate using domain rules
		if domainErr := org.Validate(); domainErr != nil {
			log.Printf("⚠️  Skipping row %d: validation failed - %v", rowNum, domainErr)
			skipped++
			continue
		}

		query := `
            INSERT INTO organizations (name, type, country, emirate, area, currency, license_number, establishment_id, description, is_active)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
            ON CONFLICT (establishment_id) DO UPDATE SET
                name = EXCLUDED.name,
                type = EXCLUDED.type,
                emirate = EXCLUDED.emirate,
                description = EXCLUDED.description,
                currency = EXCLUDED.currency,
                updated_at = NOW()
        `

		_, err := s.pool.Exec(ctx, query,
			org.Name,
			org.Type,
			org.Country,
			org.Emirate,
			org.Area,
			org.Currency,
			org.LicenseNumber,
			org.EstablishmentID,
			org.Description,
			org.IsActive,
		)

		if err != nil {
			log.Printf("⚠️  Failed to insert organization row %d (%s): %v", rowNum, org.Name, err)
			skipped++
			continue
		}

		log.Printf("✓ Inserted organization: %s (Establishment ID: %s)", org.Name, org.EstablishmentID)
		count++
	}

	log.Printf("✓ Loaded %d organizations from Excel (%d skipped)", count, skipped)
	return nil
}

// SeedAccountsFromExcel loads accounts with account type validation
func (s *Seeder) SeedAccountsFromExcel(ctx context.Context, filePath, orgID string) error {
	log.Printf("Loading accounts from Excel: %s for organization: %s", filePath, orgID)

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Accounts")
	if err != nil {
		return fmt.Errorf("failed to read Accounts sheet: %w", err)
	}

	if len(rows) < 2 {
		return fmt.Errorf("excel file must have header row and at least one data row")
	}

	count := 0
	skipped := 0

	for i, row := range rows[1:] {
		rowNum := i + 2

		if len(row) < 3 {
			log.Printf("⚠️  Skipping row %d: insufficient columns", rowNum)
			skipped++
			continue
		}

		// Parse isActive
		isActive := true
		if len(row) > 4 && strings.ToLower(strings.TrimSpace(row[4])) == "false" {
			isActive = false
		}

		code := strings.TrimSpace(row[0])
		name := strings.TrimSpace(row[1])
		accountType := strings.ToUpper(strings.TrimSpace(row[2]))

		// Basic validation
		if code == "" || name == "" {
			log.Printf("⚠️  Skipping row %d: Code or Name missing", rowNum)
			skipped++
			continue
		}

		// Validate account type
		if !isValidAccountType(accountType) {
			log.Printf("⚠️  Skipping row %d: invalid account type '%s'. Valid types: ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE", rowNum, accountType)
			skipped++
			continue
		}

		query := `
            INSERT INTO accounts (organization_id, code, name, account_type, description, is_active)
            VALUES ($1, $2, $3, $4, $5, $6)
            ON CONFLICT (organization_id, code) DO UPDATE SET
                name = EXCLUDED.name,
                account_type = EXCLUDED.account_type,
                description = EXCLUDED.description,
                is_active = EXCLUDED.is_active,
                updated_at = NOW()
        `

		_, err := s.pool.Exec(ctx, query,
			orgID,
			code,
			name,
			accountType,
			getColumnOrDefault(row, 3, ""),
			isActive,
		)

		if err != nil {
			log.Printf("⚠️  Failed to insert account row %d (%s): %v", rowNum, name, err)
			skipped++
			continue
		}

		log.Printf("✓ Inserted account: %s - %s", code, name)
		count++
	}

	log.Printf("✓ Loaded %d accounts from Excel (%d skipped)", count, skipped)
	return nil
}

// SeedShafafiyaFromExcel loads Shafafiya settings with domain validation and password encryption
func (s *Seeder) SeedShafafiyaFromExcel(ctx context.Context, filePath string) error {
	log.Printf("Loading Shafafiya settings from Excel: %s", filePath)

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Shafafiya")
	if err != nil {
		return fmt.Errorf("failed to read Shafafiya sheet: %w", err)
	}

	if len(rows) < 2 {
		return fmt.Errorf("excel file must have header row and at least one data row")
	}

	count := 0
	skipped := 0

	for i, row := range rows[1:] {
		rowNum := i + 2

		if len(row) < 4 {
			log.Printf("⚠️  Skipping row %d: insufficient columns", rowNum)
			skipped++
			continue
		}

		establishmentID := strings.TrimSpace(row[0])
		if establishmentID == "" {
			log.Printf("⚠️  Skipping row %d: Establishment ID missing", rowNum)
			skipped++
			continue
		}

		// Get organization ID
		var orgID string
		err := s.pool.QueryRow(ctx, "SELECT id FROM organizations WHERE establishment_id = $1", establishmentID).Scan(&orgID)
		if err != nil {
			log.Printf("⚠️  Organization not found for establishment_id %s at row %d: %v", establishmentID, rowNum, err)
			skipped++
			continue
		}

		orgUUID, err := uuid.Parse(orgID)
		if err != nil {
			log.Printf("⚠️  Invalid organization UUID at row %d: %v", rowNum, err)
			skipped++
			continue
		}

		// Parse includeSensitive
		includeSensitive := false
		if len(row) > 6 && strings.ToLower(strings.TrimSpace(row[6])) == "true" {
			includeSensitive = true
		}

		password := strings.TrimSpace(row[2])

		// Create domain object for validation (before encryption)
		settings := &domain.ShafafiyaOrgSettings{
			ID:                   uuid.New(), // Temporary ID for validation
			OrganizationID:       orgUUID,
			Username:             strings.TrimSpace(row[1]),
			Password:             password,
			ProviderCode:         strings.TrimSpace(row[3]),
			DefaultCurrencyCode:  strings.ToUpper(strings.TrimSpace(getColumnOrDefault(row, 4, "AED"))),
			DefaultLanguage:      strings.ToLower(strings.TrimSpace(getColumnOrDefault(row, 5, "en"))),
			IncludeSensitiveData: includeSensitive,
			CostingMethod:        strings.ToUpper(strings.TrimSpace(getColumnOrDefault(row, 7, "DEPARTMENTAL"))),
			AllocationMethod:     strings.ToUpper(strings.TrimSpace(getColumnOrDefault(row, 8, "WEIGHTED"))),
		}

		// Validate using domain rules
		if domainErr := settings.Validate(); domainErr != nil {
			log.Printf("⚠️  Skipping row %d: validation failed - %v", rowNum, domainErr)
			skipped++
			continue
		}

		// Additional business logic validation
		if !settings.CanSubmit() {
			log.Printf("⚠️  Skipping row %d: settings incomplete for submission", rowNum)
			skipped++
			continue
		}

		// Encrypt password
		encryptedPassword, err := encryptPassword(password)
		if err != nil {
			log.Printf("⚠️  Failed to encrypt password at row %d: %v", rowNum, err)
			skipped++
			continue
		}

		query := `
            INSERT INTO shafafiya_org_settings (
                organization_id, username, password_encrypted, provider_code,
                default_currency_code, default_language, include_sensitive_data,
                costing_method, allocation_method
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            ON CONFLICT (organization_id) DO UPDATE SET
                username = EXCLUDED.username,
                password_encrypted = EXCLUDED.password_encrypted,
                provider_code = EXCLUDED.provider_code,
                default_currency_code = EXCLUDED.default_currency_code,
                default_language = EXCLUDED.default_language,
                include_sensitive_data = EXCLUDED.include_sensitive_data,
                costing_method = EXCLUDED.costing_method,
                allocation_method = EXCLUDED.allocation_method,
                updated_at = NOW()
        `

		_, err = s.pool.Exec(ctx, query,
			settings.OrganizationID,
			settings.Username,
			encryptedPassword,
			settings.ProviderCode,
			settings.DefaultCurrencyCode,
			settings.DefaultLanguage,
			settings.IncludeSensitiveData,
			settings.CostingMethod,
			settings.AllocationMethod,
		)

		if err != nil {
			log.Printf("⚠️  Failed to insert Shafafiya settings for org %s at row %d: %v", orgID, rowNum, err)
			skipped++
			continue
		}

		log.Printf("✓ Inserted Shafafiya settings for Establishment ID: %s", establishmentID)
		count++
	}

	log.Printf("✓ Loaded %d Shafafiya settings from Excel (%d skipped)", count, skipped)
	return nil
}

// SeedAllFromExcel performs complete seeding from a single Excel file
func (s *Seeder) SeedAllFromExcel(ctx context.Context, filePath string) error {
	log.Printf("Starting complete seed from Excel: %s", filePath)

	// Step 1: Seed organizations
	if err := s.SeedOrganizationsFromExcel(ctx, filePath); err != nil {
		return fmt.Errorf("failed to seed organizations: %w", err)
	}

	// Step 2: Get all organization IDs for seeding accounts
	rows, err := s.pool.Query(ctx, "SELECT id, establishment_id FROM organizations ORDER BY created_at DESC")
	if err != nil {
		return fmt.Errorf("failed to query organizations: %w", err)
	}
	defer rows.Close()

	var orgIDs []string
	orgCount := 0
	for rows.Next() {
		var orgID, establishmentID string
		if err := rows.Scan(&orgID, &establishmentID); err != nil {
			log.Printf("⚠️  Failed to scan organization: %v", err)
			continue
		}
		orgIDs = append(orgIDs, orgID)
		orgCount++
	}

	log.Printf("Found %d organizations for account seeding", orgCount)

	// Step 3: Seed accounts for each organization
	// Note: This assumes all accounts in the sheet should go to all organizations
	// Modify this logic if you need org-specific account assignment
	if len(orgIDs) > 0 {
		// For simplicity, seed accounts to the first organization
		// In production, you might want org-specific account sheets
		if err := s.SeedAccountsFromExcel(ctx, filePath, orgIDs[0]); err != nil {
			return fmt.Errorf("failed to seed accounts: %w", err)
		}
	}

	// Step 4: Seed Shafafiya settings
	if err := s.SeedShafafiyaFromExcel(ctx, filePath); err != nil {
		return fmt.Errorf("failed to seed Shafafiya settings: %w", err)
	}

	log.Printf("✓ Complete seed finished successfully")
	return nil
}

// ExportTemplatesToExcel creates Excel template files with proper allowed values
func (s *Seeder) ExportTemplatesToExcel(filePath string) error {
	log.Printf("Creating Excel template: %s", filePath)

	f := excelize.NewFile()

	// ==================== Organizations Sheet ====================
	orgSheet := "Organizations"
	index, err := f.NewSheet(orgSheet)
	if err != nil {
		return fmt.Errorf("failed to create Organizations sheet: %w", err)
	}
	f.SetActiveSheet(index)

	// Set headers with bold style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 11},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})

	headers := []string{"Name", "Type", "Country", "Emirate", "Area", "Currency", "License Number", "Establishment ID", "Description", "Is Active"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(orgSheet, cell, header)
		f.SetCellStyle(orgSheet, cell, cell, headerStyle)
	}

	// Set column widths
	f.SetColWidth(orgSheet, "A", "A", 30)
	f.SetColWidth(orgSheet, "B", "B", 20)
	f.SetColWidth(orgSheet, "C", "C", 10)
	f.SetColWidth(orgSheet, "D", "D", 20)
	f.SetColWidth(orgSheet, "E", "E", 20)
	f.SetColWidth(orgSheet, "F", "F", 10)
	f.SetColWidth(orgSheet, "G", "G", 20)
	f.SetColWidth(orgSheet, "H", "H", 20)
	f.SetColWidth(orgSheet, "I", "I", 30)
	f.SetColWidth(orgSheet, "J", "J", 12)

	// Add example row
	f.SetCellValue(orgSheet, "A2", "Al Shifa Medical Center")
	f.SetCellValue(orgSheet, "B2", "HEALTHCARE")
	f.SetCellValue(orgSheet, "C2", "AE")
	f.SetCellValue(orgSheet, "D2", "ABU_DHABI")
	f.SetCellValue(orgSheet, "E2", "Healthcare City")
	f.SetCellValue(orgSheet, "F2", "AED")
	f.SetCellValue(orgSheet, "G2", "DHA-12345")
	f.SetCellValue(orgSheet, "H2", "EST-001")
	f.SetCellValue(orgSheet, "I2", "Multi-specialty clinic")
	f.SetCellValue(orgSheet, "J2", "true")

	// Add notes/instructions row with info style
	noteStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Italic: true, Size: 9, Color: "666666"},
		Alignment: &excelize.Alignment{WrapText: true, Vertical: "top"},
	})

	f.SetCellValue(orgSheet, "B3", "HEALTHCARE, RETAIL, MANUFACTURING, FINANCE, EDUCATION, HOSPITALITY, LOGISTICS, REAL_ESTATE, SERVICE, OTHER")
	f.SetCellStyle(orgSheet, "B3", "B3", noteStyle)

	f.SetCellValue(orgSheet, "D3", "ABU_DHABI, DUBAI, SHARJAH, RAS_AL_KHAIMAH, UMM_AL_QUWAIN, FUJAIRAH, AJMAN")
	f.SetCellStyle(orgSheet, "D3", "D3", noteStyle)

	f.SetCellValue(orgSheet, "F3", "AED, USD, EUR, GBP, INR, SAR, KWD, QAR")
	f.SetCellStyle(orgSheet, "F3", "F3", noteStyle)

	f.SetCellValue(orgSheet, "J3", "true or false")
	f.SetCellStyle(orgSheet, "J3", "J3", noteStyle)

	// ==================== Accounts Sheet ====================
	accSheet := "Accounts"
	f.NewSheet(accSheet)

	accHeaders := []string{"Code", "Name", "Account Type", "Description", "Is Active"}
	for i, header := range accHeaders {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(accSheet, cell, header)
		f.SetCellStyle(accSheet, cell, cell, headerStyle)
	}

	f.SetColWidth(accSheet, "A", "A", 15)
	f.SetColWidth(accSheet, "B", "B", 35)
	f.SetColWidth(accSheet, "C", "C", 20)
	f.SetColWidth(accSheet, "D", "D", 40)
	f.SetColWidth(accSheet, "E", "E", 12)

	// Add example rows
	f.SetCellValue(accSheet, "A2", "1000")
	f.SetCellValue(accSheet, "B2", "Cash")
	f.SetCellValue(accSheet, "C2", "ASSET")
	f.SetCellValue(accSheet, "D2", "Cash on hand")
	f.SetCellValue(accSheet, "E2", "true")

	f.SetCellValue(accSheet, "A3", "1100")
	f.SetCellValue(accSheet, "B3", "Accounts Receivable")
	f.SetCellValue(accSheet, "C3", "ASSET")
	f.SetCellValue(accSheet, "D3", "Amount owed by customers")
	f.SetCellValue(accSheet, "E3", "true")

	f.SetCellValue(accSheet, "A4", "2000")
	f.SetCellValue(accSheet, "B4", "Accounts Payable")
	f.SetCellValue(accSheet, "C4", "LIABILITY")
	f.SetCellValue(accSheet, "D4", "Amount owed to suppliers")
	f.SetCellValue(accSheet, "E4", "true")

	f.SetCellValue(accSheet, "A5", "4000")
	f.SetCellValue(accSheet, "B5", "Medical Services Revenue")
	f.SetCellValue(accSheet, "C5", "REVENUE")
	f.SetCellValue(accSheet, "D5", "Income from patient services")
	f.SetCellValue(accSheet, "E5", "true")

	f.SetCellValue(accSheet, "A6", "5000")
	f.SetCellValue(accSheet, "B6", "Salaries Expense")
	f.SetCellValue(accSheet, "C6", "EXPENSE")
	f.SetCellValue(accSheet, "D6", "Employee salaries and wages")
	f.SetCellValue(accSheet, "E6", "true")

	// Add notes
	f.SetCellValue(accSheet, "C7", "ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE")
	f.SetCellStyle(accSheet, "C7", "C7", noteStyle)

	f.SetCellValue(accSheet, "E7", "true or false")
	f.SetCellStyle(accSheet, "E7", "E7", noteStyle)

	// ==================== Shafafiya Sheet ====================
	shafSheet := "Shafafiya"
	f.NewSheet(shafSheet)

	shafHeaders := []string{"Establishment ID", "Username", "Password (Encrypted)", "Provider Code", "Currency Code", "Language", "Include Sensitive Data", "Costing Method", "Allocation Method"}
	for i, header := range shafHeaders {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(shafSheet, cell, header)
		f.SetCellStyle(shafSheet, cell, cell, headerStyle)
	}

	f.SetColWidth(shafSheet, "A", "A", 20)
	f.SetColWidth(shafSheet, "B", "B", 20)
	f.SetColWidth(shafSheet, "C", "C", 25)
	f.SetColWidth(shafSheet, "D", "D", 20)
	f.SetColWidth(shafSheet, "E", "E", 15)
	f.SetColWidth(shafSheet, "F", "F", 12)
	f.SetColWidth(shafSheet, "G", "G", 20)
	f.SetColWidth(shafSheet, "H", "H", 20)
	f.SetColWidth(shafSheet, "I", "I", 20)

	// Add example row
	f.SetCellValue(shafSheet, "A2", "EST-001")
	f.SetCellValue(shafSheet, "B2", "clinic_user")
	f.SetCellValue(shafSheet, "C2", "MySecurePassword123")
	f.SetCellValue(shafSheet, "D2", "DHA")
	f.SetCellValue(shafSheet, "E2", "AED")
	f.SetCellValue(shafSheet, "F2", "en")
	f.SetCellValue(shafSheet, "G2", "false")
	f.SetCellValue(shafSheet, "H2", "DEPARTMENTAL")
	f.SetCellValue(shafSheet, "I2", "WEIGHTED")

	// Add notes
	f.SetCellValue(shafSheet, "C3", "Password will be encrypted automatically on import")
	f.SetCellStyle(shafSheet, "C3", "C3", noteStyle)

	f.SetCellValue(shafSheet, "E3", "AED, USD, EUR, GBP, INR, SAR, KWD, QAR")
	f.SetCellStyle(shafSheet, "E3", "E3", noteStyle)

	f.SetCellValue(shafSheet, "F3", "en, ar")
	f.SetCellStyle(shafSheet, "F3", "F3", noteStyle)

	f.SetCellValue(shafSheet, "G3", "true or false")
	f.SetCellStyle(shafSheet, "G3", "G3", noteStyle)

	f.SetCellValue(shafSheet, "H3", "DEPARTMENTAL, ACTIVITY_BASED, SERVICE_BASED")
	f.SetCellStyle(shafSheet, "H3", "H3", noteStyle)

	f.SetCellValue(shafSheet, "I3", "WEIGHTED, PERCENTAGE, FIXED")
	f.SetCellStyle(shafSheet, "I3", "I3", noteStyle)

	// Delete default Sheet1
	f.DeleteSheet("Sheet1")

	// Save file
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	log.Printf("✓ Template created: %s", filePath)
	log.Printf("  - Organizations sheet with 10 columns")
	log.Printf("  - Accounts sheet with 5 example accounts")
	log.Printf("  - Shafafiya sheet with complete configuration")
	return nil
}

// ==================== Helper Functions ====================

// getColumnOrDefault returns column value or default if not present
func getColumnOrDefault(row []string, index int, defaultValue string) string {
	if index < len(row) && strings.TrimSpace(row[index]) != "" {
		return row[index]
	}
	return defaultValue
}

// encryptPassword encrypts password using bcrypt
func encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt password: %w", err)
	}
	return string(hash), nil
}

// isValidAccountType validates account type
func isValidAccountType(accountType string) bool {
	validTypes := map[string]bool{
		"ASSET":     true,
		"LIABILITY": true,
		"EQUITY":    true,
		"REVENUE":   true,
		"EXPENSE":   true,
	}
	return validTypes[accountType]
}

// ValidateExcelStructure validates the Excel file structure before seeding
func (s *Seeder) ValidateExcelStructure(filePath string) error {
	log.Printf("Validating Excel structure: %s", filePath)

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	// Check for required sheets
	requiredSheets := []string{"Organizations", "Accounts", "Shafafiya"}
	sheets := f.GetSheetList()
	sheetMap := make(map[string]bool)
	for _, sheet := range sheets {
		sheetMap[sheet] = true
	}

	for _, required := range requiredSheets {
		if !sheetMap[required] {
			return fmt.Errorf("missing required sheet: %s", required)
		}
	}

	// Validate Organizations sheet headers
	orgRows, err := f.GetRows("Organizations")
	if err != nil || len(orgRows) < 1 {
		return fmt.Errorf("organizations sheet is empty or invalid")
	}

	expectedOrgHeaders := []string{"Name", "Type", "Country", "Emirate", "Area", "Currency", "License Number", "Establishment ID"}
	if len(orgRows[0]) < len(expectedOrgHeaders) {
		return fmt.Errorf("organizations sheet has insufficient columns")
	}

	// Validate Accounts sheet headers
	accRows, err := f.GetRows("Accounts")
	if err != nil || len(accRows) < 1 {
		return fmt.Errorf("accounts sheet is empty or invalid")
	}

	expectedAccHeaders := []string{"Code", "Name", "Account Type"}
	if len(accRows[0]) < len(expectedAccHeaders) {
		return fmt.Errorf("accounts sheet has insufficient columns")
	}

	// Validate Shafafiya sheet headers
	shafRows, err := f.GetRows("Shafafiya")
	if err != nil || len(shafRows) < 1 {
		return fmt.Errorf("shafafiya sheet is empty or invalid")
	}

	expectedShafHeaders := []string{"Establishment ID", "Username", "Password (Encrypted)", "Provider Code"}
	if len(shafRows[0]) < len(expectedShafHeaders) {
		return fmt.Errorf("shafafiya sheet has insufficient columns")
	}

	log.Printf("✓ Excel structure validation passed")
	return nil
}
