-- Migration 000003: Create full platform schema for Orders, Invitations, Domains, RSVP, Greetings, Transactions, and Settings
-- Sesuai SSOT §45 (Orders), §40 (Invitations), §24 (Domains), §52 (RSVP), §53 (Greetings), §47 (Transactions), §56 (Database)

-- Tabel Orders
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) NOT NULL UNIQUE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    package_id UUID REFERENCES packages(id) ON DELETE SET NULL,
    template_id UUID REFERENCES templates(id) ON DELETE SET NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(50) NOT NULL,
    groom_name VARCHAR(150),
    bride_name VARCHAR(150),
    event_date DATE,
    event_location TEXT,
    custom_domain VARCHAR(255),
    package_name VARCHAR(100) NOT NULL,
    amount BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'accepted', 'rejected', 'completed', 'cancelled')),
    payment_method VARCHAR(50) DEFAULT 'bank_transfer',
    payment_proof_url TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders (created_at DESC);

-- Tabel Invitations
CREATE TABLE IF NOT EXISTS invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    template_id UUID REFERENCES templates(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(150) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'scheduled', 'active', 'expired', 'disabled')),
    start_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    end_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '180 days'),
    rsvp_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    greeting_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    content JSONB NOT NULL DEFAULT '{}'::jsonb,
    theme_config JSONB NOT NULL DEFAULT '{}'::jsonb,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_invitations_user_id ON invitations (user_id);
CREATE INDEX IF NOT EXISTS idx_invitations_status ON invitations (status);

-- Tabel Invitation Domains
CREATE TABLE IF NOT EXISTS invitation_domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invitation_id UUID NOT NULL REFERENCES invitations(id) ON DELETE CASCADE,
    hostname VARCHAR(255) NOT NULL,
    normalized_hostname VARCHAR(255) NOT NULL UNIQUE,
    is_primary BOOLEAN NOT NULL DEFAULT TRUE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'dns_verified', 'ssl_pending', 'active', 'failed', 'disabled')),
    dns_verified_at TIMESTAMPTZ,
    ssl_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    ssl_verified_at TIMESTAMPTZ,
    last_checked_at TIMESTAMPTZ,
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_invitation_domains_status ON invitation_domains (status);

-- Tabel RSVP Responses
CREATE TABLE IF NOT EXISTS rsvp_responses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invitation_id UUID NOT NULL REFERENCES invitations(id) ON DELETE CASCADE,
    guest_code VARCHAR(50),
    guest_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(50),
    attendance_status VARCHAR(50) NOT NULL CHECK (attendance_status IN ('attending', 'not_attending', 'unsure')),
    guest_count INT NOT NULL DEFAULT 1,
    message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rsvp_invitation_id ON rsvp_responses (invitation_id);
CREATE INDEX IF NOT EXISTS idx_rsvp_attendance ON rsvp_responses (attendance_status);

-- Tabel Guest Messages (Ucapan Tamu)
CREATE TABLE IF NOT EXISTS guest_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invitation_id UUID NOT NULL REFERENCES invitations(id) ON DELETE CASCADE,
    guest_name VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'approved' CHECK (status IN ('pending', 'approved', 'hidden', 'rejected')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_guest_messages_invitation ON guest_messages (invitation_id);
CREATE INDEX IF NOT EXISTS idx_guest_messages_status ON guest_messages (status);

-- Tabel Transactions
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    amount BIGINT NOT NULL DEFAULT 0,
    payment_method VARCHAR(50) NOT NULL DEFAULT 'bank_transfer',
    payment_reference VARCHAR(100),
    status VARCHAR(50) NOT NULL DEFAULT 'verified' CHECK (status IN ('pending', 'verified', 'failed', 'refunded')),
    proof_url TEXT,
    verified_by UUID REFERENCES users(id) ON DELETE SET NULL,
    verified_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions (status);

-- Tabel Application Settings
CREATE TABLE IF NOT EXISTS application_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    setting_key VARCHAR(100) NOT NULL UNIQUE,
    setting_value JSONB NOT NULL DEFAULT '{}'::jsonb,
    description TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
