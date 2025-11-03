package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrganizationRepository implements OrganizationRepositoryInterface
type OrganizationRepository struct {
	pool *pgxpool.Pool
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(pool *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{pool: pool}
}

// CreateOrganization creates a new organization in the database
func (r *OrganizationRepository) CreateOrganization(ctx context.Context, org domain.Organization) (domain.Organization, error) {
	query := `
        INSERT INTO organizations 
        (name, type, country, emirate, area, currency, tax_id, license_number, establishment_id, description, is_active)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, created_at, updated_at
    `

	err := r.pool.QueryRow(ctx, query,
		org.Name,
		org.Type,
		org.Country,
		org.Emirate,
		org.Area,
		org.Currency,
		org.TaxID,
		org.LicenseNumber,
		org.EstablishmentID,
		org.Description,
		true, // is_active always true on creation
	).Scan(&org.ID, &org.CreatedAt, &org.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "establishment_id") {
				return domain.Organization{}, fmt.Errorf("organization with establishment ID %s already exists: %w", org.EstablishmentID, err)
			}
		}
		return domain.Organization{}, fmt.Errorf("failed to create organization: %w", err)
	}

	org.IsActive = true
	return org, nil
}

// GetOrganizationByID retrieves an organization by its ID
func (r *OrganizationRepository) GetOrganizationByID(ctx context.Context, id uuid.UUID) (domain.Organization, error) {
	query := `
        SELECT id, name, type, country, emirate, area, currency, tax_id, license_number,
               establishment_id, description, is_active, created_at, updated_at
        FROM organizations
        WHERE id = $1
    `

	var org domain.Organization
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&org.ID,
		&org.Name,
		&org.Type,
		&org.Country,
		&org.Emirate,
		&org.Area,
		&org.Currency,
		&org.TaxID,
		&org.LicenseNumber,
		&org.EstablishmentID,
		&org.Description,
		&org.IsActive,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Organization{}, fmt.Errorf("organization with ID %s not found", id)
		}
		return domain.Organization{}, fmt.Errorf("failed to get organization: %w", err)
	}

	return org, nil
}

// GetOrganizationByEstablishmentID retrieves an organization by establishment ID
func (r *OrganizationRepository) GetOrganizationByEstablishmentID(ctx context.Context, establishmentID string) (domain.Organization, error) {
	query := `
        SELECT id, name, type, country, emirate, area, currency, tax_id, license_number,
               establishment_id, description, is_active, created_at, updated_at
        FROM organizations
        WHERE establishment_id = $1
    `

	var org domain.Organization
	err := r.pool.QueryRow(ctx, query, establishmentID).Scan(
		&org.ID,
		&org.Name,
		&org.Type,
		&org.Country,
		&org.Emirate,
		&org.Area,
		&org.Currency,
		&org.TaxID,
		&org.LicenseNumber,
		&org.EstablishmentID,
		&org.Description,
		&org.IsActive,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Organization{}, fmt.Errorf("organization with establishment ID %s not found", establishmentID)
		}
		return domain.Organization{}, fmt.Errorf("failed to get organization by establishment ID: %w", err)
	}

	return org, nil
}

// ListOrganizations retrieves organizations with optional filters
func (r *OrganizationRepository) ListOrganizations(ctx context.Context, filters *OrganizationFilters) ([]domain.Organization, error) {
	if filters == nil {
		filters = DefaultFilters()
	}

	query := `
        SELECT id, name, type, country, emirate, area, currency, tax_id, license_number,
               establishment_id, description, is_active, created_at, updated_at
        FROM organizations
        WHERE 1=1
    `

	args := []interface{}{}
	argNum := 1

	// Filter by type
	if filters.Type != nil {
		query += fmt.Sprintf(" AND type = $%d", argNum)
		args = append(args, *filters.Type)
		argNum++
	}

	// Filter by emirate
	if filters.Emirate != nil {
		query += fmt.Sprintf(" AND emirate = $%d", argNum)
		args = append(args, *filters.Emirate)
		argNum++
	}

	// Filter by country
	if filters.Country != nil {
		query += fmt.Sprintf(" AND country = $%d", argNum)
		args = append(args, *filters.Country)
		argNum++
	}

	// Filter by active status
	if filters.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argNum)
		args = append(args, *filters.IsActive)
		argNum++
	}

	query += " ORDER BY name ASC"

	// Add pagination
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argNum)
		args = append(args, filters.Limit)
		argNum++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argNum)
		args = append(args, filters.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	defer rows.Close()

	var orgs []domain.Organization
	for rows.Next() {
		var org domain.Organization
		err := rows.Scan(
			&org.ID,
			&org.Name,
			&org.Type,
			&org.Country,
			&org.Emirate,
			&org.Area,
			&org.Currency,
			&org.TaxID,
			&org.LicenseNumber,
			&org.EstablishmentID,
			&org.Description,
			&org.IsActive,
			&org.CreatedAt,
			&org.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organization: %w", err)
		}
		orgs = append(orgs, org)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating organizations: %w", err)
	}

	return orgs, nil
}

// UpdateOrganization updates an existing organization
func (r *OrganizationRepository) UpdateOrganization(ctx context.Context, org domain.Organization) (domain.Organization, error) {
	if org.ID == uuid.Nil {
		return domain.Organization{}, fmt.Errorf("organization ID is required for update")
	}

	query := `
        UPDATE organizations
        SET name = $1, emirate = $2, area = $3, currency = $4, description = $5, updated_at = NOW()
        WHERE id = $6
        RETURNING created_at, updated_at
    `

	err := r.pool.QueryRow(ctx, query,
		org.Name,
		org.Emirate,
		org.Area,
		org.Currency,
		org.Description,
		org.ID,
	).Scan(&org.CreatedAt, &org.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Organization{}, fmt.Errorf("organization with ID %s not found", org.ID)
		}
		return domain.Organization{}, fmt.Errorf("failed to update organization: %w", err)
	}

	return org, nil
}

// DeactivateOrganization deactivates an organization
func (r *OrganizationRepository) DeactivateOrganization(ctx context.Context, id uuid.UUID) error {
	query := `
        UPDATE organizations
        SET is_active = FALSE, updated_at = NOW()
        WHERE id = $1
    `

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate organization: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("organization with ID %s not found", id)
	}

	return nil
}

// ActivateOrganization reactivates a deactivated organization
func (r *OrganizationRepository) ActivateOrganization(ctx context.Context, id uuid.UUID) error {
	query := `
        UPDATE organizations
        SET is_active = TRUE, updated_at = NOW()
        WHERE id = $1
    `

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to activate organization: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("organization with ID %s not found", id)
	}

	return nil
}

// ExistsWithName checks if organization with given name exists
func (r *OrganizationRepository) ExistsWithName(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM organizations
            WHERE LOWER(name) = LOWER($1)
    `

	args := []interface{}{name}

	if excludeID != nil {
		query += " AND id != $2"
		args = append(args, *excludeID)
	}

	query += ")"

	var exists bool
	err := r.pool.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check organization name existence: %w", err)
	}

	return exists, nil
}

// ExistsWithEstablishmentID checks if organization with given establishment ID exists
func (r *OrganizationRepository) ExistsWithEstablishmentID(ctx context.Context, establishmentID string, excludeID *uuid.UUID) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM organizations
            WHERE establishment_id = $1
    `

	args := []interface{}{establishmentID}

	if excludeID != nil {
		query += " AND id != $2"
		args = append(args, *excludeID)
	}

	query += ")"

	var exists bool
	err := r.pool.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check establishment ID existence: %w", err)
	}

	return exists, nil
}

// CountByType counts organizations by type
func (r *OrganizationRepository) CountByType(ctx context.Context, orgType domain.OrganizationType) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM organizations
        WHERE type = $1 AND is_active = TRUE
    `

	var count int
	err := r.pool.QueryRow(ctx, query, orgType).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count organizations by type: %w", err)
	}

	return count, nil
}

// CountByEmirate counts healthcare organizations by emirate
func (r *OrganizationRepository) CountByEmirate(ctx context.Context, emirate domain.UAEmirate) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM organizations
        WHERE emirate = $1 AND type = $2 AND is_active = TRUE
    `

	var count int
	err := r.pool.QueryRow(ctx, query, emirate, domain.OrganizationTypeHealthcare).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count organizations by emirate: %w", err)
	}

	return count, nil
}
