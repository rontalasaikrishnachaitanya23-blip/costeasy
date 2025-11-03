-- backend/database/migrations/000004_create_organizations.up.sql

CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    country VARCHAR(50) NOT NULL DEFAULT 'AE',
    emirate VARCHAR(50),
    area VARCHAR(100),
    currency VARCHAR(3) NOT NULL DEFAULT 'AED',
    tax_id VARCHAR(100),
    license_number VARCHAR(100),
    establishment_id VARCHAR(100),
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT organizations_establishment_id_unique UNIQUE(establishment_id)
);

CREATE INDEX idx_organizations_type ON organizations(type);
CREATE INDEX idx_organizations_emirate ON organizations(emirate);
CREATE INDEX idx_organizations_is_active ON organizations(is_active);
CREATE INDEX idx_organizations_establishment_id ON organizations(establishment_id);
