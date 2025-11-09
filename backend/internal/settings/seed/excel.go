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

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
)

// Seeder loads settings (organizations, accounts, shafafiya) from Excel files.
// This file contains safer helpers (nil handling, header validation) and clearer logging.
type Seeder struct {
	pool *pgxpool.Pool
}

func nilIfEmptyEmiratePtr(e *domain.UAEmirate) interface{} {
	if e == nil || strings.TrimSpace(string(*e)) == "" {
		return nil
	}
	return string(*e)
}

func NewSeeder(pool *pgxpool.Pool) *Seeder {
	return &Seeder{pool: pool}
}

// ==================== Organizations ====================

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
		if len(row) > 9 && strings.EqualFold(strings.TrimSpace(row[9]), "false") {
			isActive = false
		}

		// Create domain object for validation
		org := &domain.Organization{
			Name:            strings.TrimSpace(row[0]),
			Type:            domain.OrganizationType(strings.ToUpper(strings.TrimSpace(row[1]))),
			Country:         strings.TrimSpace(row[2]),
			Emirate:         domain.EmiratePtr(strings.TrimSpace(row[3])),
			Area:            domain.StringPtr(strings.TrimSpace(row[4])),
			Address:         domain.StringPtr(strings.TrimSpace(row[5])),
			Phone:           domain.StringPtr(strings.TrimSpace(row[6])),
			Email:           domain.StringPtr(strings.TrimSpace(row[7])),
			Website:         domain.StringPtr(getColumnOrDefault(row, 8, "")),
			Currency:        strings.TrimSpace(getColumnOrDefault(row, 9, "AED")),
			TaxID:           domain.StringPtr(getColumnOrDefault(row, 10, "")),
			LicenseNumber:   domain.StringPtr(getColumnOrDefault(row, 11, "")),
			EstablishmentID: domain.StringPtr(getColumnOrDefault(row, 12, "")),
			Description:     domain.StringPtr(getColumnOrDefault(row, 13, "")),
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
			nilIfEmptyEmiratePtr(org.Emirate),
			nilIfEmptyStrPtr(org.Area),
			org.Currency,
			nilIfEmptyStrPtr(org.LicenseNumber),
			nilIfEmptyStrPtr(org.EstablishmentID),
			nilIfEmptyStrPtr(org.Description),
			org.IsActive,
		)

		if err != nil {
			log.Printf("⚠️  Failed to insert organization row %d (%s): %v", rowNum, org.Name, err)
			skipped++
			continue
		}

		log.Printf("✓ Inserted organization: %s (Establishment ID: %s)", org.Name, safeDeref(org.EstablishmentID))
		count++
	}

	log.Printf("✓ Loaded %d organizations from Excel (%d skipped)", count, skipped)
	return nil
}

// ==================== Accounts ====================

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
		if len(row) > 4 && strings.EqualFold(strings.TrimSpace(row[4]), "false") {
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

		desc := getColumnOrDefault(row, 3, "")
		var descArg interface{}
		if strings.TrimSpace(desc) == "" {
			descArg = nil
		} else {
			descArg = desc
		}

		_, err := s.pool.Exec(ctx, query,
			orgID,
			code,
			name,
			accountType,
			descArg,
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

// ==================== Shafafiya ====================

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
		if len(row) > 6 && strings.EqualFold(strings.TrimSpace(row[6]), "true") {
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

		// Hash password with bcrypt
		hashed, err := encryptPassword(password)
		if err != nil {
			log.Printf("⚠️  Failed to hash password at row %d: %v", rowNum, err)
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
			hashed,
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

// ==================== Complete & Templates ====================

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

	// Step 3: Seed accounts for each organization (simple approach uses first org)
	if len(orgIDs) > 0 {
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

	// Example row
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

	// Example rows
	f.SetCellValue(accSheet, "A2", "1000")
	f.SetCellValue(accSheet, "B2", "Cash")
	f.SetCellValue(accSheet, "C2", "ASSET")
	f.SetCellValue(accSheet, "D2", "Cash on hand")
	f.SetCellValue(accSheet, "E2", "true")

	// ==================== Shafafiya Sheet ====================
	shafSheet := "Shafafiya"
	f.NewSheet(shafSheet)

	shafHeaders := []string{"Establishment ID", "Username", "Password (Encrypted)", "Provider Code", "Currency Code", "Language", "Include Sensitive Data", "Costing Method", "Allocation Method"}
	for i, header := range shafHeaders {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(shafSheet, cell, header)
		f.SetCellStyle(shafSheet, cell, cell, headerStyle)
	}

	f.SetCellValue(shafSheet, "A2", "EST-001")
	f.SetCellValue(shafSheet, "B2", "clinic_user")
	f.SetCellValue(shafSheet, "C2", "MySecurePassword123")
	f.SetCellValue(shafSheet, "D2", "DHA")
	f.SetCellValue(shafSheet, "E2", "AED")
	f.SetCellValue(shafSheet, "F2", "en")
	f.SetCellValue(shafSheet, "G2", "false")
	f.SetCellValue(shafSheet, "H2", "DEPARTMENTAL")
	f.SetCellValue(shafSheet, "I2", "WEIGHTED")

	// Delete default Sheet1
	f.DeleteSheet("Sheet1")

	// Save file
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	log.Printf("✓ Template created: %s", filePath)
	return nil
}

// ==================== Helpers ====================

func getColumnOrDefault(row []string, index int, defaultValue string) string {
	if index < len(row) && strings.TrimSpace(row[index]) != "" {
		return row[index]
	}
	return defaultValue
}

// nilIfEmptyStrPtr returns nil when pointer is nil or points to empty string, otherwise returns the dereferenced string.
func nilIfEmptyStrPtr(s *string) interface{} {
	if s == nil || strings.TrimSpace(*s) == "" {
		return nil
	}
	return *s
}

// safeDeref returns a readable string for a *string (returns empty string if nil)
func safeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// encryptPassword hashes the password using bcrypt (non-reversible)
func encryptPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

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
