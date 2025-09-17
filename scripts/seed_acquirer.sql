INSERT INTO accounts (
    id, primary_account_number, number, merchant_account, bank_identifier_code, merchant_password,
    merchant_id, expiration_date, card_holder_name, security_code, balance
) VALUES
-- Merchant accounts
('1e2f3a4b-5c6d-7e8f-9012-3456789abc01', '3879520000012345', '160-13646185005023-13', true, '387952', 'm3rchant!', 'MERCH-1001', NOW() + INTERVAL '2 years', 'Nikola Nikolic', '321', 234500.75),
('2b3c4d5e-6f7a-8b9c-0123-456789abcdef', '3879520000023456', '160-44811483878583-30', true, '387952', 's3cur3Pa$$', 'MERCH-1002', NOW() + INTERVAL '3 years', 'Milena Milenkovic', '453', 152000.00),
('3c4d5e6f-7a8b-9c0d-1234-56789abcdef0', '3879520000034567', '160-33752021114200-38', true, '387952', 'P@ssw0rd', 'MERCH-1003', NOW() + INTERVAL '1 year', 'Petar Petrovic', '987', 199990.99),
('4d5e6f7a-8b9c-0d1e-2345-6789abcdef01', '3879520000045678', '160-12188161996163-87', true, '387952', 'safe1234', 'MERCH-1004', NOW() + INTERVAL '4 years', 'Jelena Jankovic', '129', 182500.10),
('5e6f7a8b-9c0d-1e2f-3456-789abcdef012', '3879520000056789', '160-32450648275404-70', true, '387952', 'merchant5!', 'MERCH-1005', NOW() + INTERVAL '2 years', 'Marko Markovic', '213', 225000.00),

-- Regular personal accounts
('6f7a8b9c-0d1e-2f3a-4567-89abcdef0123', '3879520000067890', '160-65935980837512-98', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'Ana Anic', '621', 62000.50),
('7a8b9c0d-1e2f-3a4b-5678-9abcdef01234', '3879520000078901', '160-78897075011199-57', false, '387952', NULL, NULL, NOW() + INTERVAL '1 year', 'Milan Milic', '874', 430000.00),
('8b9c0d1e-2f3a-4b5c-6789-abcdef012345', '3879520000089012', '160-50846027517296-38', false, '387952', NULL, NULL, NOW() + INTERVAL '3 years', 'Ivana Ivanovic', '753', 5500.00),
('9c0d1e2f-3a4b-5c6d-789a-bcdef0123456', '3879520000090123', '160-40040111647409-81', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'Stefan Stevic', '214', 88000.90),
('0d1e2f3a-4b5c-6d7e-89ab-cdef01234567', '3879520000101234', '160-19574843957189-95', false, '387952', NULL, NULL, NOW() + INTERVAL '4 years', 'Teodora Todorovic', '634', 127500.25),
('1f2e3d4c-5b6a-7f8e-9012-3456789abc11', '3879520000112345', '160-42910618998733-02', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'Vuk Vukovic', '147', 4000.00),
('2f3e4d5c-6b7a-8f9e-0123-456789abcd12', '3879520000123456', '160-41898313465951-95', false, '387952', NULL, NULL, NOW() + INTERVAL '1 year', 'Tamara Tamic', '369', 3200.75),
('3f4e5d6c-7b8a-9f0e-1234-56789abcde23', '3879520000134567', '160-99216092923556-56', false, '387952', NULL, NULL, NOW() + INTERVAL '5 years', 'Aleksa Aleksic', '852', 23000.00),
('4f5e6d7c-8b9a-0f1e-2345-6789abcdef34', '3879520000145678', '160-19861636440986-12', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'Nina Nikolic', '951', 103200.00),
('5f6e7d8c-9b0a-1f2e-3456-789abcdef045', '3879520000156789', '160-78054280842909-24', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'Dusan Dusanic', '142', 1000.00),

-- More merchant accounts
('6f7e8d9c-0b1a-2f3e-4567-89abcdef0567', '3879520000167890', '160-17747416133912-09', true, '387952', 'merchant6!', 'MERCH-1006', NOW() + INTERVAL '2 years', 'Kristina Krstic', '777', 265000.00),
('7f8e9d0c-1b2a-3f4e-5678-9abcdef06789', '3879520000178901', '160-65094541400868-04', true, '387952', 'merchant7!', 'MERCH-1007', NOW() + INTERVAL '3 years', 'Filip Filipovic', '543', 170900.00),
('8f9e0d1c-2b3a-4f5e-6789-abcdef07890a', '3879520000189012', '160-18488188011365-57', true, '387952', 'merchant8!', 'MERCH-1008', NOW() + INTERVAL '1 year', 'Natalija Novak', '456', 201000.00),
('9f0e1d2c-3b4a-5f6e-7890-bcdef08901ab', '3879520000190123', '160-85473844459373-60', true, '387952', 'merchant9!', 'MERCH-1009', NOW() + INTERVAL '2 years', 'Sava Savic', '234', 300000.00),
('0f1e2d3c-4b5a-6f7e-8901-cdef0123456b', '3879520000201234', '160-32957845108126-78', true, '387952', 'merchant10!','MERCH-1010', NOW() + INTERVAL '4 years', 'Marija Maric', '999', 188000.00);
