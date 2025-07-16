CREATE EXTENSION IF NOT EXISTS "pgcrypto";

INSERT INTO accounts (
    id, pan, number, merchant_account, bank_identifier_code, merchant_password,
    merchant_id, expiration_date, card_holder_name, security_code, balance
) VALUES
-- Merchant accounts
(gen_random_uuid(), '4111111111111111', 'ACC-1001', true, 'BICRSBG10', 'm3rchant!', 'MERCH-1001', NOW() + INTERVAL '2 years', 'Nikola Nikolic', '321', 23450.75),
(gen_random_uuid(), '5500000000000004', 'ACC-1002', true, 'BICRSBG11', 's3cur3Pa$$', 'MERCH-1002', NOW() + INTERVAL '3 years', 'Milena Milenkovic', '453', 15200.00),
(gen_random_uuid(), '4007000000027',    'ACC-1003', true, 'BICRSBG12', 'P@ssw0rd', 'MERCH-1003', NOW() + INTERVAL '1 year', 'Petar Petrovic', '987', 19999.99),
(gen_random_uuid(), '6011000000000004', 'ACC-1004', true, 'BICRSBG13', 'safe1234', 'MERCH-1004', NOW() + INTERVAL '4 years', 'Jelena Jankovic', '129', 18250.10),
(gen_random_uuid(), '3530111333300000', 'ACC-1005', true, 'BICRSBG14', 'merchant5!', 'MERCH-1005', NOW() + INTERVAL '2 years', 'Marko Markovic', '213', 22500.00),

-- Regular personal accounts
(gen_random_uuid(), '4012888888881881', 'ACC-2001', false, 'BICRSBG21', NULL, NULL, NOW() + INTERVAL '2 years', 'Ana Anic', '621', 6200.50),
(gen_random_uuid(), '6011111111111117', 'ACC-2002', false, 'BICRSBG22', NULL, NULL, NOW() + INTERVAL '1 year', 'Milan Milic', '874', 4300.00),
(gen_random_uuid(), '5105105105105100', 'ACC-2003', false, 'BICRSBG23', NULL, NULL, NOW() + INTERVAL '3 years', 'Ivana Ivanovic', '753', 550.00),
(gen_random_uuid(), '4111111111111234', 'ACC-2004', false, 'BICRSBG24', NULL, NULL, NOW() + INTERVAL '2 years', 'Stefan Stevic', '214', 8800.90),
(gen_random_uuid(), '6011000990139424', 'ACC-2005', false, 'BICRSBG25', NULL, NULL, NOW() + INTERVAL '4 years', 'Teodora Todorovic', '634', 12750.25),
(gen_random_uuid(), '378282246310005',  'ACC-2006', false, 'BICRSBG26', NULL, NULL, NOW() + INTERVAL '2 years', 'Vuk Vukovic', '147', 400.00),
(gen_random_uuid(), '6011000000000004', 'ACC-2007', false, 'BICRSBG27', NULL, NULL, NOW() + INTERVAL '1 year', 'Tamara Tamic', '369', 320.75),
(gen_random_uuid(), '4222222222222',    'ACC-2008', false, 'BICRSBG28', NULL, NULL, NOW() + INTERVAL '5 years', 'Aleksa Aleksic', '852', 2300.00),
(gen_random_uuid(), '4000000000000002', 'ACC-2009', false, 'BICRSBG29', NULL, NULL, NOW() + INTERVAL '2 years', 'Nina Nikolic', '951', 10320.00),
(gen_random_uuid(), '4000001234567890', 'ACC-2010', false, 'BICRSBG30', NULL, NULL, NOW() + INTERVAL '2 years', 'Dusan Dusanic', '142', 100.00),

-- More merchant accounts
(gen_random_uuid(), '5555555555554444', 'ACC-1006', true, 'BICRSBG15', 'merchant6!', 'MERCH-1006', NOW() + INTERVAL '2 years', 'Kristina Krstic', '777', 26500.00),
(gen_random_uuid(), '2223000048400011', 'ACC-1007', true, 'BICRSBG16', 'merchant7!', 'MERCH-1007', NOW() + INTERVAL '3 years', 'Filip Filipovic', '543', 17900.00),
(gen_random_uuid(), '4111111133334444', 'ACC-1008', true, 'BICRSBG17', 'merchant8!', 'MERCH-1008', NOW() + INTERVAL '1 year', 'Natalija Novak', '456', 20100.00),
(gen_random_uuid(), '60115564485789458','ACC-1009', true, 'BICRSBG18', 'merchant9!', 'MERCH-1009', NOW() + INTERVAL '2 years', 'Sava Savic', '234', 30000.00),
(gen_random_uuid(), '343434343434343',  'ACC-1010', true, 'BICRSBG19', 'merchant10!','MERCH-1010', NOW() + INTERVAL '4 years', 'Marija Maric', '999', 18800.00);
