package repository

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/chaitu35/costeasy/backend/internal/settings/domain"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationRepository struct {
    pool *pgxpool.Pool
}

func NewOrganizationRepository(pool *pgxpool.Pool) *OrganizationRepository {
    return &OrganizationRepository{pool: pool}
}

// Create creates a new organization
func (r *OrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
    query := `
        INSERT INTO organizations (
            id, name, code, display_name, type, country, emirate, area,
            address, city, state, postal_code, phone, email, website,
            currency, tax_id, license_number, license_expiry, establishment_id,
            description, is_active,
            ip_whitelist_enabled, allowed_ips, allowed_ip_ranges,
            mfa_enabled, mfa_enforced, mfa_method, allowed_mfa_methods,
            created_at, updated_at, created_by
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
            $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29,
            $30, $31, $32
        )
        RETURNING created_at, updated_at
    `

    org.ID = uuid.New()
    now := time.Now()

    err := r.pool.QueryRow(ctx, query,
        org.ID, org.Name, org.Code, org.DisplayName, org.Type,
        org.Country, org.Emirate, org.Area,
        org.Address, org.City, org.State, org.PostalCode,
        org.Phone, org.Email, org.Website,
        org.Currency, org.TaxID, org.LicenseNumber, org.LicenseExpiry, org.EstablishmentID,
        org.Description, org.IsActive,
        org.IPWhitelistEnabled, org.AllowedIPs, org.AllowedIPRanges,
        org.MFAEnabled, org.MFAEnforced, org.MFAMethod, org.AllowedMFAMethods,
        now, now, org.CreatedBy,
    ).Scan(&org.CreatedAt, &org.UpdatedAt)

    if err != nil {
        return fmt.Errorf("failed to create organization: %w", err)
    }

    return nil
}

// GetByID retrieves an organization by ID
func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
    query := `
        SELECT 
            id, name, code, display_name, type, country, emirate, area,
            address, city, state, postal_code, phone, email, website,
            currency, tax_id, license_number, license_expiry, establishment_id,
            description, is_active,
            ip_whitelist_enabled, allowed_ips, allowed_ip_ranges,
            mfa_enabled, mfa_enforced, mfa_method, allowed_mfa_methods,
            created_at, updated_at, created_by, updated_by
        FROM organizations
        WHERE id = $1
    `

    var org domain.Organization
    err := r.pool.QueryRow(ctx, query, id).Scan(
        &org.ID, &org.Name, &org.Code, &org.DisplayName, &org.Type,
        &org.Country, &org.Emirate, &org.Area,
        &org.Address, &org.City, &org.State, &org.PostalCode,
        &org.Phone, &org.Email, &org.Website,
        &org.Currency, &org.TaxID, &org.LicenseNumber, &org.LicenseExpiry, &org.EstablishmentID,
        &org.Description, &org.IsActive,
        &org.IPWhitelistEnabled, &org.AllowedIPs, &org.AllowedIPRanges,
        &org.MFAEnabled, &org.MFAEnforced, &org.MFAMethod, &org.AllowedMFAMethods,
        &org.CreatedAt, &org.UpdatedAt, &org.CreatedBy, &org.UpdatedBy,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("organization not found")
        }
        return nil, fmt.Errorf("failed to get organization: %w", err)
    }

    return &org, nil
}

// GetByCode retrieves organization by code
func (r *OrganizationRepository) GetByCode(ctx context.Context, code string) (*domain.Organization, error) {
    query := `
        SELECT 
            id, name, code, display_name, type, country, emirate, area,
            address, city, state, postal_code, phone, email, website,
            currency, tax_id, license_number, license_expiry, establishment_id,
            description, is_active,
            ip_whitelist_enabled, allowed_ips, allowed_ip_ranges,
            mfa_enabled, mfa_enforced, mfa_method, allowed_mfa_methods,
            created_at, updated_at, created_by, updated_by
        FROM organizations
        WHERE code = $1
    `

    var org domain.Organization
    err := r.pool.QueryRow(ctx, query, code).Scan(
        &org.ID, &org.Name, &org.Code, &org.DisplayName, &org.Type,
        &org.Country, &org.Emirate, &org.Area,
        &org.Address, &org.City, &org.State, &org.PostalCode,
        &org.Phone, &org.Email, &org.Website,
        &org.Currency, &org.TaxID, &org.LicenseNumber, &org.LicenseExpiry, &org.EstablishmentID,
        &org.Description, &org.IsActive,
        &org.IPWhitelistEnabled, &org.AllowedIPs, &org.AllowedIPRanges,
        &org.MFAEnabled, &org.MFAEnforced, &org.MFAMethod, &org.AllowedMFAMethods,
        &org.CreatedAt, &org.UpdatedAt, &org.CreatedBy, &org.UpdatedBy,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("organization not found")
        }
        return nil, fmt.Errorf("failed to get organization: %w", err)
    }

    return &org, nil
}

// GetByEstablishmentID retrieves organization by establishment ID
func (r *OrganizationRepository) GetByEstablishmentID(ctx context.Context, establishmentID string) (*domain.Organization, error) {
    query := `
        SELECT 
            id, name, code, display_name, type, country, emirate, area,
            address, city, state, postal_code, phone, email, website,
            currency, tax_id, license_number, license_expiry, establishment_id,
            description, is_active,
            ip_whitelist_enabled, allowed_ips, allowed_ip_ranges,
            mfa_enabled, mfa_enforced, mfa_method, allowed_mfa_methods,
            created_at, updated_at, created_by, updated_by
        FROM organizations
        WHERE establishment_id = $1
    `

    var org domain.Organization
    err := r.pool.QueryRow(ctx, query, establishmentID).Scan(
        &org.ID, &org.Name, &org.Code, &org.DisplayName, &org.Type,
        &org.Country, &org.Emirate, &org.Area,
        &org.Address, &org.City, &org.State, &org.PostalCode,
        &org.Phone, &org.Email, &org.Website,
        &org.Currency, &org.TaxID, &org.LicenseNumber, &org.LicenseExpiry, &org.EstablishmentID,
        &org.Description, &org.IsActive,
        &org.IPWhitelistEnabled, &org.AllowedIPs, &org.AllowedIPRanges,
        &org.MFAEnabled, &org.MFAEnforced, &org.MFAMethod, &org.AllowedMFAMethods,
        &org.CreatedAt, &org.UpdatedAt, &org.CreatedBy, &org.UpdatedBy,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("organization not found")
        }
        return nil, fmt.Errorf("failed to get organization: %w", err)
    }

    return &org, nil
}

// Update updates an existing organization
func (r *OrganizationRepository) Update(ctx context.Context, org *domain.Organization) error {
    query := `
        UPDATE organizations
        SET 
            name = $2, code = $3, display_name = $4, type = $5,
            country = $6, emirate = $7, area = $8,
            address = $9, city = $10, state = $11, postal_code = $12,
            phone = $13, email = $14, website = $15,
            currency = $16, tax_id = $17, license_number = $18, license_expiry = $19, establishment_id = $20,
            description = $21, is_active = $22,
            ip_whitelist_enabled = $23, allowed_ips = $24, allowed_ip_ranges = $25,
            mfa_enabled = $26, mfa_enforced = $27, mfa_method = $28, allowed_mfa_methods = $29,
            updated_at = $30, updated_by = $31
        WHERE id = $1
        RETURNING updated_at
    `

    now := time.Now()

    err := r.pool.QueryRow(ctx, query,
        org.ID, org.Name, org.Code, org.DisplayName, org.Type,
        org.Country, org.Emirate, org.Area,
        org.Address, org.City, org.State, org.PostalCode,
        org.Phone, org.Email, org.Website,
        org.Currency, org.TaxID, org.LicenseNumber, org.LicenseExpiry, org.EstablishmentID,
        org.Description, org.IsActive,
        org.IPWhitelistEnabled, org.AllowedIPs, org.AllowedIPRanges,
        org.MFAEnabled, org.MFAEnforced, org.MFAMethod, org.AllowedMFAMethods,
        now, org.UpdatedBy,
    ).Scan(&org.UpdatedAt)

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return fmt.Errorf("organization not found")
        }
        return fmt.Errorf("failed to update organization: %w", err)
    }

    return nil
}

// Delete soft-deletes an organization
func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
    query := `UPDATE organizations SET is_active = false, updated_at = $1 WHERE id = $2`

    result, err := r.pool.Exec(ctx, query, time.Now(), id)
    if err != nil {
        return fmt.Errorf("failed to delete organization: %w", err)
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("organization not found")
    }

    return nil
}

// List retrieves all active organizations with pagination
func (r *OrganizationRepository) List(ctx context.Context, limit, offset int) ([]*domain.Organization, error) {
    query := `
        SELECT 
            id, name, code, display_name, type, country, emirate, area,
            address, city, state, postal_code, phone, email, website,
            currency, tax_id, license_number, license_expiry, establishment_id,
            description, is_active,
            ip_whitelist_enabled, allowed_ips, allowed_ip_ranges,
            mfa_enabled, mfa_enforced, mfa_method, allowed_mfa_methods,
            created_at, updated_at, created_by, updated_by
        FROM organizations
        WHERE is_active = true
        ORDER BY name
        LIMIT $1 OFFSET $2
    `

    rows, err := r.pool.Query(ctx, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to list organizations: %w", err)
    }
    defer rows.Close()

    var organizations []*domain.Organization
    for rows.Next() {
        org := &domain.Organization{}
        err := rows.Scan(
            &org.ID, &org.Name, &org.Code, &org.DisplayName, &org.Type,
            &org.Country, &org.Emirate, &org.Area,
            &org.Address, &org.City, &org.State, &org.PostalCode,
            &org.Phone, &org.Email, &org.Website,
            &org.Currency, &org.TaxID, &org.LicenseNumber, &org.LicenseExpiry, &org.EstablishmentID,
            &org.Description, &org.IsActive,
            &org.IPWhitelistEnabled, &org.AllowedIPs, &org.AllowedIPRanges,
            &org.MFAEnabled, &org.MFAEnforced, &org.MFAMethod, &org.AllowedMFAMethods,
            &org.CreatedAt, &org.UpdatedAt, &org.CreatedBy, &org.UpdatedBy,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan organization: %w", err)
        }
        organizations = append(organizations, org)
    }

    return organizations, rows.Err()
}

// Count returns total number of organizations
func (r *OrganizationRepository) Count(ctx context.Context) (int, error) {
    query := `SELECT COUNT(*) FROM organizations WHERE is_active = true`
    
    var count int
    err := r.pool.QueryRow(ctx, query).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count organizations: %w", err)
    }
    
    return count, nil
}

// CountByType returns count of organizations by type
func (r *OrganizationRepository) CountByType(ctx context.Context, orgType domain.OrganizationType) (int, error) {
    query := `SELECT COUNT(*) FROM organizations WHERE type = $1 AND is_active = true`
    
    var count int
    err := r.pool.QueryRow(ctx, query, orgType).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count organizations by type: %w", err)
    }
    
    return count, nil
}

// CountByEmirate returns count of organizations by emirate
func (r *OrganizationRepository) CountByEmirate(ctx context.Context, emirate domain.UAEmirate) (int, error) {
    query := `SELECT COUNT(*) FROM organizations WHERE emirate = $1 AND is_active = true`
    
    var count int
    err := r.pool.QueryRow(ctx, query, emirate).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count organizations by emirate: %w", err)
    }
    
    return count, nil
}

// ActivateOrganization activates an organization
func (r *OrganizationRepository) ActivateOrganization(ctx context.Context, id uuid.UUID) error {
    query := `UPDATE organizations SET is_active = true, updated_at = $1 WHERE id = $2`

    result, err := r.pool.Exec(ctx, query, time.Now(), id)
    if err != nil {
        return fmt.Errorf("failed to activate organization: %w", err)
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("organization not found")
    }

    return nil
}

// DeactivateOrganization deactivates an organization
func (r *OrganizationRepository) DeactivateOrganization(ctx context.Context, id uuid.UUID) error {
    query := `UPDATE organizations SET is_active = false, updated_at = $1 WHERE id = $2`

    result, err := r.pool.Exec(ctx, query, time.Now(), id)
    if err != nil {
        return fmt.Errorf("failed to deactivate organization: %w", err)
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("organization not found")
    }

    return nil
}

// UpdateMFASettings updates only MFA-related settings for an organization
func (r *OrganizationRepository) UpdateMFASettings(ctx context.Context, orgID uuid.UUID, mfaEnabled, mfaEnforced bool, mfaMethod domain.MFAMethod, allowedMethods []domain.MFAMethod) error {
    query := `
        UPDATE organizations SET
            mfa_enabled = $2,
            mfa_enforced = $3,
            mfa_method = $4,
            allowed_mfa_methods = $5,
            updated_at = $6
        WHERE id = $1
    `

    result, err := r.pool.Exec(ctx, query,
        orgID, mfaEnabled, mfaEnforced, mfaMethod, allowedMethods, time.Now(),
    )
    if err != nil {
        return fmt.Errorf("failed to update MFA settings: %w", err)
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("organization not found")
    }

    return nil
}
