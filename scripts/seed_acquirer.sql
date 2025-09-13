-- Seed Acquirer Accounts (UUIDs are hard-coded, no extensions needed)

INSERT INTO accounts (
    id, pan, number, merchant_account, bank_identifier_code, merchant_password,
    merchant_id, expiration_date, card_holder_name, security_code, balance
) VALUES
-- Merchant accounts
('1e2f3a4b-5c6d-7e8f-9012-3456789abc01', '4111111111111111', 'ACC-1001', true, 'BICRSBG10', 'm3rchant!', 'MERCH-1001', NOW() + INTERVAL '2 years', 'Nikola Nikolic', '321', 23450.75),
('2b3c4d5e-6f7a-8b9c-0123-456789abcdef', '5500000000000004', 'ACC-1002', true, 'BICRSBG11', 's3cur3Pa$$', 'MERCH-1002', NOW() + INTERVAL '3 years', 'Milena Milenkovic', '453', 15200.00),
('3c4d5e6f-7a8b-9c0d-1234-56789abcdef0', '4007000000027', 'ACC-1003', true, 'BICRSBG12', 'P@ssw0rd', 'MERCH-1003', NOW() + INTERVAL '1 year', 'Petar Petrovic', '987', 19999.99),
('4d5e6f7a-8b9c-0d1e-2345-6789abcdef01', '6011000000000004', 'ACC-1004', true, 'BICRSBG13', 'safe1234', 'MERCH-1004', NOW() + INTERVAL '4 years', 'Jelena Jankovic', '129', 18250.10),
('5e6f7a8b-9c0d-1e2f-3456-789abcdef012', '3530111333300000', 'ACC-1005', true, 'BICRSBG14', 'merchant5!', 'MERCH-1005', NOW() + INTERVAL '2 years', 'Marko Markovic', '213', 22500.00),

-- Regular personal accounts
('6f7a8b9c-0d1e-2f3a-4567-89abcdef0123', '4012888888881881', 'ACC-2001', false, 'BICRSBG21', NULL, NULL, NOW() + INTERVAL '2 years', 'Ana Anic', '621', 6200.50),
('7a8b9c0d-1e2f-3a4b-5678-9abcdef01234', '6011111111111117', 'ACC-2002', false, 'BICRSBG22', NULL, NULL, NOW() + INTERVAL '1 year', 'Milan Milic', '874', 4300.00),
('8b9c0d1e-2f3a-4b5c-6789-abcdef012345', '5105105105105100', 'ACC-2003', false, 'BICRSBG23', NULL, NULL, NOW() + INTERVAL '3 years', 'Ivana Ivanovic', '753', 550.00),
('9c0d1e2f-3a4b-5c6d-789a-bcdef0123456', '4111111111111234', 'ACC-2004', false, 'BICRSBG24', NULL, NULL, NOW() + INTERVAL '2 years', 'Stefan Stevic', '214', 8800.90),
('0d1e2f3a-4b5c-6d7e-89ab-cdef01234567', '6011000990139424', 'ACC-2005', false, 'BICRSBG25', NULL, NULL, NOW() + INTERVAL '4 years', 'Teodora Todorovic', '634', 12750.25),
('1f2e3d4c-5b6a-7f8e-9012-3456789abc11', '378282246310005', 'ACC-2006', false, 'BICRSBG26', NULL, NULL, NOW() + INTERVAL '2 years', 'Vuk Vukovic', '147', 400.00),
('2f3e4d5c-6b7a-8f9e-0123-456789abcd12', '6011000000000004', 'ACC-2007', false, 'BICRSBG27', NULL, NULL, NOW() + INTERVAL '1 year', 'Tamara Tamic', '369', 320.75),
('3f4e5d6c-7b8a-9f0e-1234-56789abcde23', '4222222222222', 'ACC-2008', false, 'BICRSBG28', NULL, NULL, NOW() + INTERVAL '5 years', 'Aleksa Aleksic', '852', 2300.00),
('4f5e6d7c-8b9a-0f1e-2345-6789abcdef34', '4000000000000002', 'ACC-2009', false, 'BICRSBG29', NULL, NULL, NOW() + INTERVAL '2 years', 'Nina Nikolic', '951', 10320.00),
('5f6e7d8c-9b0a-1f2e-3456-789abcdef045', '4000001234567890', 'ACC-2010', false, 'BICRSBG30', NULL, NULL, NOW() + INTERVAL '2 years', 'Dusan Dusanic', '142', 100.00),

-- More merchant accounts
('6f7e8d9c-0b1a-2f3e-4567-89abcdef0567', '5555555555554444', 'ACC-1006', true, 'BICRSBG15', 'merchant6!', 'MERCH-1006', NOW() + INTERVAL '2 years', 'Kristina Krstic', '777', 26500.00),
('7f8e9d0c-1b2a-3f4e-5678-9abcdef06789', '2223000048400011', 'ACC-1007', true, 'BICRSBG16', 'merchant7!', 'MERCH-1007', NOW() + INTERVAL '3 years', 'Filip Filipovic', '543', 17900.00),
('8f9e0d1c-2b3a-4f5e-6789-abcdef07890a', '4111111133334444', 'ACC-1008', true, 'BICRSBG17', 'merchant8!', 'MERCH-1008', NOW() + INTERVAL '1 year', 'Natalija Novak', '456', 20100.00),
('9f0e1d2c-3b4a-5f6e-7890-bcdef08901ab', '60115564485789458','ACC-1009', true, 'BICRSBG18', 'merchant9!', 'MERCH-1009', NOW() + INTERVAL '2 years', 'Sava Savic', '234', 30000.00),
('0f1e2d3c-4b5a-6f7e-8901-cdef0123456b', '343434343434343',  'ACC-1010', true, 'BICRSBG19', 'merchant10!','MERCH-1010', NOW() + INTERVAL '4 years', 'Marija Maric', '999', 18800.00);
