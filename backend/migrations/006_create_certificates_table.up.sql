-- Create certificates table
CREATE TABLE IF NOT EXISTS certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    template_url TEXT,
    event_date TIMESTAMP WITH TIME ZONE,
    issued_by VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'archived')),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_certificates_status ON certificates(status);
CREATE INDEX idx_certificates_created_by ON certificates(created_by);
CREATE INDEX idx_certificates_deleted_at ON certificates(deleted_at) WHERE deleted_at IS NULL;

-- Create certificate_recipients table
CREATE TABLE IF NOT EXISTS certificate_recipients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    certificate_id UUID NOT NULL REFERENCES certificates(id) ON DELETE CASCADE,
    recipient_name VARCHAR(255) NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    delivery_status VARCHAR(50) DEFAULT 'pending' CHECK (delivery_status IN ('pending', 'queued', 'processing', 'sent', 'failed')),
    queued_at TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cert_recipients_certificate_id ON certificate_recipients(certificate_id);
CREATE INDEX idx_cert_recipients_delivery_status ON certificate_recipients(delivery_status);
CREATE INDEX idx_cert_recipients_email ON certificate_recipients(recipient_email);
