CREATE TABLE IF NOT EXISTS exchange_rates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cur_id INT NOT NULL,
    date DATE NOT NULL,
    cur_abbreviation VARCHAR(10) NOT NULL,
    cur_scale INT NOT NULL,
    cur_name VARCHAR(255) NOT NULL,
    cur_official_rate DECIMAL(10,4) NOT NULL,
    UNIQUE KEY unique_currency_date (cur_abbreviation, date)
);
