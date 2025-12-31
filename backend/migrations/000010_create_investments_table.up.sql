-- Create investments table
CREATE TABLE IF NOT EXISTS investments (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    account_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    ticker VARCHAR(20),
    purchase_date DATE NOT NULL,
    purchase_amount BIGINT NOT NULL,
    purchase_currency VARCHAR(3) NOT NULL DEFAULT 'BRL',
    current_value BIGINT NOT NULL,
    current_currency VARCHAR(3) NOT NULL DEFAULT 'BRL',
    quantity DECIMAL(15,4),
    context VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_investments_user_id ON investments(user_id);
CREATE INDEX idx_investments_account_id ON investments(account_id);
CREATE INDEX idx_investments_type ON investments(type);
CREATE INDEX idx_investments_user_context ON investments(user_id, context);
CREATE INDEX idx_investments_deleted_at ON investments(deleted_at);

-- Add foreign key constraints
ALTER TABLE investments
    ADD CONSTRAINT fk_investments_user_id
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE investments
    ADD CONSTRAINT fk_investments_account_id
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE;

