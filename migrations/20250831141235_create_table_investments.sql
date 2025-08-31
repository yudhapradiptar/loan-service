-- +goose Up
-- +goose StatementBegin
CREATE TABLE investments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(255) NOT NULL,
    loan_id INT NOT NULL,
    investor_id VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    invested_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_uuid (uuid),
    INDEX idx_loan_id (loan_id),
    INDEX idx_investor_id (investor_id),
    FOREIGN KEY (loan_id) REFERENCES loans(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS investments;
-- +goose StatementEnd
