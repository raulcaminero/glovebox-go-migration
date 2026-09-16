-- +goose Up
CREATE TABLE policyholders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policyholder_id UUID NOT NULL REFERENCES policyholders(id) ON DELETE CASCADE,
    carrier TEXT NOT NULL,
    policy_number TEXT NOT NULL,
    line_of_business TEXT NOT NULL, -- auto, home, umbrella, etc.
    effective_date DATE NOT NULL,
    expiration_date DATE NOT NULL,
    premium_cents BIGINT NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active', -- active, expired, cancelled
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (carrier, policy_number)
);

CREATE INDEX idx_policies_policyholder_id ON policies(policyholder_id);
CREATE INDEX idx_policies_expiration_date ON policies(expiration_date);

CREATE TABLE notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policyholder_id UUID NOT NULL REFERENCES policyholders(id) ON DELETE CASCADE,
    author TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_notes_policyholder_id ON notes(policyholder_id);

-- +goose Down
DROP TABLE notes;
DROP TABLE policies;
DROP TABLE policyholders;
