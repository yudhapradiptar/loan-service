-- +goose Up
-- +goose StatementBegin
CREATE TABLE loan_approval_validators (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(255) NOT NULL,
    loan_approval_id INT NOT NULL,
    employee_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_uuid (uuid),
    INDEX idx_employee_id (employee_id),
    FOREIGN KEY (loan_approval_id) REFERENCES loan_approvals(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS loan_approval_validators;
-- +goose StatementEnd
