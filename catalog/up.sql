CREATE TABLE IF NOT EXISTS catalog (
  id CHAR(27) PRIMARY KEY,
  name VARCHAR(24) NOT NULL,
  description VARCHAR(255), 
  price DECIMAL(10, 2) NOT NULL CHECK (price >= 0) -- Added check for non-negative price
);
