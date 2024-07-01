-- Customers Table
CREATE TABLE Customers (
    customer_id INT PRIMARY KEY, -- Or consider using a generated sequence/auto-increment
    name VARCHAR(255) NOT NULL, 
    address VARCHAR(255),
    phone_number VARCHAR(25), 
    email VARCHAR(255)
);

-- Accounts Table
CREATE TABLE Accounts (
    account_id INT PRIMARY KEY, -- Or consider using a generated sequence/auto-increment 
    customer_id INT NOT NULL,
    account_number VARCHAR(50) NOT NULL, -- Adjust length if needed
    account_type VARCHAR(20) NOT NULL, -- Consider using an ENUM for structure
    balance DECIMAL(19, 4) NOT NULL, -- Scale for your currency/precision needs
    interest_rate DECIMAL(5, 4), -- If optional, allow NULL

    CONSTRAINT FK_Accounts_Customers FOREIGN KEY (customer_id) REFERENCES Customers(customer_id)
);

-- Transactions Table
CREATE TABLE Transactions (
    transaction_id INT PRIMARY KEY, -- Or consider using a generated sequence/auto-increment 
    account_id INT NOT NULL, 
    transaction_type VARCHAR(20) NOT NULL, -- Consider using an ENUM 
    amount DECIMAL(19, 4) NOT NULL, 
    date DATE NOT NULL,
    description VARCHAR(255), 

    CONSTRAINT FK_Transactions_Accounts FOREIGN KEY (account_id) REFERENCES Accounts(account_id)
);

-- General Ledger Table
CREATE TABLE General_Ledger (
    transaction_id INT PRIMARY KEY, 
    account_debited INT, 
    account_credited INT, 
    amount DECIMAL(19, 4) NOT NULL,

    CONSTRAINT FK_GeneralLedger_Transactions FOREIGN KEY (transaction_id) REFERENCES Transactions(transaction_id),
    CONSTRAINT FK_GeneralLedger_Accounts_Debited FOREIGN KEY (account_debited) REFERENCES Accounts(account_id),
    CONSTRAINT FK_GeneralLedger_Accounts_Credited FOREIGN KEY (account_credited) REFERENCES Accounts(account_id)
); 
