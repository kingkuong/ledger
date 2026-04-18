-- Seed data
INSERT INTO accounts (account_number, nickname, balance)
VALUES
    ('ACC-001', 'Main Checking', 500000), -- $5000.00 (assuming cents/smallest unit)
    ('ACC-002', 'Savings Account', 1250050), -- $12500.50
    ('ACC-003', 'Emergency Fund', 1000000), -- $10000.00
    ('ACC-004', 'Travel Fund', 25075)       -- $250.75
ON CONFLICT (account_number) DO NOTHING;
