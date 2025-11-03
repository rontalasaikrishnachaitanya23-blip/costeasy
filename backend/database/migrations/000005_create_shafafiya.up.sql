-- backend/database/migrations/000005_create_shafafiya.up.sql

CREATE TABLE IF NOT EXISTS shafafiya_org_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL UNIQUE REFERENCES organizations(id) ON DELETE CASCADE,
    
    username VARCHAR(255) NOT NULL,
    password_encrypted TEXT NOT NULL,
    provider_code VARCHAR(100) NOT NULL,
    
    default_currency_code VARCHAR(3) DEFAULT 'AED',
    default_language VARCHAR(2) DEFAULT 'en',
    include_sensitive_data BOOLEAN DEFAULT FALSE,
    
    costing_method VARCHAR(50) DEFAULT 'DEPARTMENTAL',
    allocation_method VARCHAR(50) DEFAULT 'WEIGHTED',
    
    last_submission_at TIMESTAMP,
    last_submission_status VARCHAR(50),
    last_submission_error TEXT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_shafafiya_organization ON shafafiya_org_settings(organization_id);
CREATE INDEX idx_shafafiya_provider_code ON shafafiya_org_settings(provider_code);
CREATE INDEX idx_shafafiya_submission_status ON shafafiya_org_settings(last_submission_status);
