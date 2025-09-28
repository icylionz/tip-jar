-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    avatar TEXT,
    google_id VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tip jars table
CREATE TABLE tip_jars (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    invite_code VARCHAR(50) UNIQUE NOT NULL,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Jar memberships table
CREATE TABLE jar_memberships (
    id SERIAL PRIMARY KEY,
    jar_id INTEGER NOT NULL REFERENCES tip_jars(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member' CHECK (role IN ('admin', 'member')),
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(jar_id, user_id)
);

-- Offense types table
CREATE TABLE offense_types (
    id SERIAL PRIMARY KEY,
    jar_id INTEGER NOT NULL REFERENCES tip_jars(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    cost_type VARCHAR(20) NOT NULL CHECK (cost_type IN ('monetary', 'action', 'item', 'service')),
    cost_amount DECIMAL(10,2),
    cost_action TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Offenses table
CREATE TABLE offenses (
    id SERIAL PRIMARY KEY,
    jar_id INTEGER NOT NULL REFERENCES tip_jars(id) ON DELETE CASCADE,
    offense_type_id INTEGER NOT NULL REFERENCES offense_types(id),
    reporter_id INTEGER NOT NULL REFERENCES users(id),
    offender_id INTEGER NOT NULL REFERENCES users(id),
    notes TEXT,
    cost_override DECIMAL(10,2),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'disputed', 'forgiven')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Payments table
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    offense_id INTEGER NOT NULL REFERENCES offenses(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    amount DECIMAL(10,2),
    proof_type VARCHAR(20) CHECK (proof_type IN ('image', 'receipt', 'video')),
    proof_url TEXT,
    verified BOOLEAN NOT NULL DEFAULT false,
    verified_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_tip_jars_invite_code ON tip_jars(invite_code);
CREATE INDEX idx_tip_jars_created_by ON tip_jars(created_by);
CREATE INDEX idx_jar_memberships_jar_id ON jar_memberships(jar_id);
CREATE INDEX idx_jar_memberships_user_id ON jar_memberships(user_id);
CREATE INDEX idx_offense_types_jar_id ON offense_types(jar_id);
CREATE INDEX idx_offenses_jar_id ON offenses(jar_id);
CREATE INDEX idx_offenses_reporter_id ON offenses(reporter_id);
CREATE INDEX idx_offenses_offender_id ON offenses(offender_id);
CREATE INDEX idx_offenses_status ON offenses(status);
CREATE INDEX idx_payments_offense_id ON payments(offense_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);

-- Updated at triggers
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tip_jars_updated_at BEFORE UPDATE ON tip_jars FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_offense_types_updated_at BEFORE UPDATE ON offense_types FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_offenses_updated_at BEFORE UPDATE ON offenses FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
