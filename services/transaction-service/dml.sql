INSERT INTO accounts (id, user_id, name, type, balance, created_at, updated_at) 
VALUES 
    (UUID(), '7cc06ed7-75ee-4ae5-b74a-8dc5cd382553', 'Rezky\'s Main Account', 'MAIN', 10000000.00, NOW(), NOW()),
    (UUID(), 'merchant-kopi-kenangan', 'Kopi Kenangan Main Account','MAIN', 50000.00, NOW(), NOW());