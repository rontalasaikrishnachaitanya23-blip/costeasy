// backend/settings/internal/repository/shafafiya_repository.go
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

// ShafafiyaRepository implements ShafafiyaSettingsRepositoryInterface
type ShafafiyaRepository struct {
	pool *pgxpool.Pool
}

// NewShafafiyaRepository creates a new Shafafiya repository
func NewShafafiyaRepository(pool *pgxpool.Pool) *ShafafiyaRepository {
	return &ShafafiyaRepository{pool: pool}
}

// CreateShafafiyaSettings creates new Shafafiya settings for an organization
func (r *ShafafiyaRepository) CreateShafafiyaSettings(ctx context.Context, settings domain.ShafafiyaOrgSettings) (domain.ShafafiyaOrgSettings, error) {
	query := `
        INSERT INTO shafafiya_org_settings 
        (organization_id, username, password_encrypted, provider_code, default_currency_code, 
         default_language, include_sensitive_data, costing_method, allocation_method)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at, updated_at
    `

	err := r.pool.QueryRow(ctx, query,
		settings.OrganizationID,
		settings.Username,
		settings.Password, // Already encrypted by service
		settings.ProviderCode,
		settings.DefaultCurrencyCode,
		settings.DefaultLanguage,
		settings.IncludeSensitiveData,
		settings.CostingMethod,
		settings.AllocationMethod,
	).Scan(&settings.ID, &settings.CreatedAt, &settings.UpdatedAt)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.ShafafiyaOrgSettings{}, fmt.Errorf("failed to create Shafafiya settings: %w", err)
		}
		return domain.ShafafiyaOrgSettings{}, fmt.Errorf("failed to create Shafafiya settings: %w", err)
	}

	// Don't return password
	settings.Password = ""

	return settings, nil
}

// GetShafafiyaSettings retrieves settings WITHOUT password
func (r *ShafafiyaRepository) GetShafafiyaSettings(ctx context.Context, orgID uuid.UUID) (domain.ShafafiyaOrgSettings, error) {
	query := `
        SELECT id, organization_id, username, provider_code, default_currency_code, default_language,
               include_sensitive_data, costing_method, allocation_method, last_submission_at,
               last_submission_status, last_submission_error, created_at, updated_at
        FROM shafafiya_org_settings
        WHERE organization_id = $1
    `

	var settings domain.ShafafiyaOrgSettings
	err := r.pool.QueryRow(ctx, query, orgID).Scan(
		&settings.ID,
		&settings.OrganizationID,
		&settings.Username,
		&settings.ProviderCode,
		&settings.DefaultCurrencyCode,
		&settings.DefaultLanguage,
		&settings.IncludeSensitiveData,
		&settings.CostingMethod,
		&settings.AllocationMethod,
		&settings.LastSubmissionAt,
		&settings.LastSubmissionStatus,
		&settings.LastSubmissionError,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ShafafiyaOrgSettings{}, fmt.Errorf("shafafiya settings not found for organization %s", orgID)
		}
		return domain.ShafafiyaOrgSettings{}, fmt.Errorf("failed to get Shafafiya settings: %w", err)
	}

	// Don't return password
	settings.Password = ""

	return settings, nil
}

// GetShafafiyaSettingsWithPassword retrieves settings WITH encrypted password
// For internal use only - NEVER expose to API
func (r *ShafafiyaRepository) GetShafafiyaSettingsWithPassword(ctx context.Context, orgID uuid.UUID) (domain.ShafafiyaOrgSettings, error) {
	query := `
        SELECT id, organization_id, username, password_encrypted, provider_code, default_currency_code,
               default_language, include_sensitive_data, costing_method, allocation_method,
               last_submission_at, last_submission_status, last_submission_error, created_at, updated_at
        FROM shafafiya_org_settings
        WHERE organization_id = $1
    `

	var settings domain.ShafafiyaOrgSettings
	err := r.pool.QueryRow(ctx, query, orgID).Scan(
		&settings.ID,
		&settings.OrganizationID,
		&settings.Username,
		&settings.Password, // Include encrypted password
		&settings.ProviderCode,
		&settings.DefaultCurrencyCode,
		&settings.DefaultLanguage,
		&settings.IncludeSensitiveData,
		&settings.CostingMethod,
		&settings.AllocationMethod,
		&settings.LastSubmissionAt,
		&settings.LastSubmissionStatus,
		&settings.LastSubmissionError,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ShafafiyaOrgSettings{}, fmt.Errorf("shafafiya settings not found for organization %s", orgID)
		}
		return domain.ShafafiyaOrgSettings{}, fmt.Errorf("failed to get Shafafiya settings with password: %w", err)
	}

	return settings, nil
}

// UpdateShafafiyaCredentials updates only the credentials
func (r *ShafafiyaRepository) UpdateShafafiyaCredentials(ctx context.Context, orgID uuid.UUID, username, encryptedPassword, providerCode string) error {
	query := `
        UPDATE shafafiya_org_settings
        SET username = $1, password_encrypted = $2, provider_code = $3, updated_at = NOW()
        WHERE organization_id = $4
    `

	result, err := r.pool.Exec(ctx, query, username, encryptedPassword, providerCode, orgID)
	if err != nil {
		return fmt.Errorf("failed to update Shafafiya credentials: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("shafafiya settings not found for organization %s", orgID)
	}

	return nil
}

// UpdateShafafiyaCosting updates costing configuration
func (r *ShafafiyaRepository) UpdateShafafiyaCosting(ctx context.Context, orgID uuid.UUID, costingMethod, allocationMethod string) error {
	query := `
        UPDATE shafafiya_org_settings
        SET costing_method = $1, allocation_method = $2, updated_at = NOW()
        WHERE organization_id = $3
    `

	result, err := r.pool.Exec(ctx, query, costingMethod, allocationMethod, orgID)
	if err != nil {
		return fmt.Errorf("failed to update Shafafiya costing: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("shafafiya settings not found for organization %s", orgID)
	}

	return nil
}

// UpdateShafafiyaSubmission updates submission configuration
func (r *ShafafiyaRepository) UpdateShafafiyaSubmission(ctx context.Context, orgID uuid.UUID, language, currency string, includeSensitive bool) error {
	query := `
        UPDATE shafafiya_org_settings
        SET default_language = $1, default_currency_code = $2, include_sensitive_data = $3, updated_at = NOW()
        WHERE organization_id = $4
    `

	result, err := r.pool.Exec(ctx, query, language, currency, includeSensitive, orgID)
	if err != nil {
		return fmt.Errorf("failed to update Shafafiya submission settings: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("shafafiya settings not found for organization %s", orgID)
	}

	return nil
}

// UpdateSubmissionStatus updates submission tracking
func (r *ShafafiyaRepository) UpdateSubmissionStatus(ctx context.Context, orgID uuid.UUID, status, errorMsg string) error {
	query := `
        UPDATE shafafiya_org_settings
        SET last_submission_at = $1, last_submission_status = $2, last_submission_error = $3, updated_at = NOW()
        WHERE organization_id = $4
    `

	now := time.Now()
	result, err := r.pool.Exec(ctx, query, now, status, errorMsg, orgID)
	if err != nil {
		return fmt.Errorf("failed to update submission status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("shafafiya settings not found for organization %s", orgID)
	}

	return nil
}

// DeleteShafafiyaSettings deletes settings for an organization
func (r *ShafafiyaRepository) DeleteShafafiyaSettings(ctx context.Context, orgID uuid.UUID) error {
	query := `
        DELETE FROM shafafiya_org_settings
        WHERE organization_id = $1
    `

	result, err := r.pool.Exec(ctx, query, orgID)
	if err != nil {
		return fmt.Errorf("failed to delete Shafafiya settings: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("shafafiya settings not found for organization %s", orgID)
	}

	return nil
}

// ExistsForOrganization checks if settings exist
func (r *ShafafiyaRepository) ExistsForOrganization(ctx context.Context, orgID uuid.UUID) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM shafafiya_org_settings
            WHERE organization_id = $1
        )
    `

	var exists bool
	err := r.pool.QueryRow(ctx, query, orgID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check Shafafiya settings existence: %w", err)
	}

	return exists, nil
}

// ListBySubmissionStatus retrieves settings with specific submission status
func (r *ShafafiyaRepository) ListBySubmissionStatus(ctx context.Context, status string, limit int) ([]domain.ShafafiyaOrgSettings, error) {
	query := `
        SELECT id, organization_id, username, provider_code, default_currency_code,
               default_language, include_sensitive_data, costing_method, allocation_method,
               last_submission_at, last_submission_status, last_submission_error, created_at, updated_at
        FROM shafafiya_org_settings
        WHERE last_submission_status = $1
        ORDER BY last_submission_at DESC
        LIMIT $2
    `

	rows, err := r.pool.Query(ctx, query, status, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list Shafafiya settings by submission status: %w", err)
	}
	defer rows.Close()

	var settings []domain.ShafafiyaOrgSettings
	for rows.Next() {
		var s domain.ShafafiyaOrgSettings
		err := rows.Scan(
			&s.ID,
			&s.OrganizationID,
			&s.Username,
			&s.ProviderCode,
			&s.DefaultCurrencyCode,
			&s.DefaultLanguage,
			&s.IncludeSensitiveData,
			&s.CostingMethod,
			&s.AllocationMethod,
			&s.LastSubmissionAt,
			&s.LastSubmissionStatus,
			&s.LastSubmissionError,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan Shafafiya settings: %w", err)
		}

		// Don't return password
		s.Password = ""
		settings = append(settings, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating Shafafiya settings: %w", err)
	}

	return settings, nil
}
