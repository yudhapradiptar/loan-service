-- +goose Up
-- +goose StatementBegin
CREATE TABLE loan_disbursements (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(255) NOT NULL,
    loan_id INT NOT NULL,
    field_officer_employee_id VARCHAR(255) NOT NULL,
    signed_agreement_letter_url VARCHAR(255) NOT NULL,
    disbursed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_uuid (uuid),
    INDEX idx_loan_id (loan_id),
    INDEX idx_field_officer_employee_id (field_officer_employee_id),
    INDEX idx_disbursed_at (disbursed_at),
    FOREIGN KEY (loan_id) REFERENCES loans(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS loan_disbursements;
-- +goose StatementEnd
