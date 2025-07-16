INSERT INTO accounts (
    id, pan, number, merchant_account, bank_identifier_code, merchant_password,
    merchant_id, expiration_date, card_holder_name, security_code, balance
) VALUES
      (gen_random_uuid(), '5111222233334444', '8701123399', false, 'ISSUER001', 'pass991', 'issuer-001', NOW() + INTERVAL '2 years', 'Ana Anić', '774', 12500.00),
      (gen_random_uuid(), '4111333344445555', '8701123300', false, 'ISSUER002', 'pass992', 'issuer-002', NOW() + INTERVAL '3 years', 'Bojan Bajić', '223', 8800.00),
      (gen_random_uuid(), '6011555566667777', '8701123311', false, 'ISSUER003', 'pass993', 'issuer-003', NOW() + INTERVAL '1 year', 'Ceca Cvetković', '901', 9700.50),
      (gen_random_uuid(), '4539123412341234', '8701123322', false, 'ISSUER004', 'pass994', 'issuer-004', NOW() + INTERVAL '2 years', 'Dragan Dimić', '119', 6000.00),
      (gen_random_uuid(), '379144444444444',  '8701123333', false, 'ISSUER005', 'pass995', 'issuer-005', NOW() + INTERVAL '4 years', 'Emilija Erić', '886', 11000.25),
      (gen_random_uuid(), '6011999999999999', '8701123344', false, 'ISSUER006', 'pass996', 'issuer-006', NOW() + INTERVAL '3 years', 'Filip Filipović', '455', 7200.00),
      (gen_random_uuid(), '3530111333300000', '8701123355', false, 'ISSUER007', 'pass997', 'issuer-007', NOW() + INTERVAL '2 years', 'Gorana Gajić', '644', 13550.00),
      (gen_random_uuid(), '5555555555554444', '8701123366', false, 'ISSUER008', 'pass998', 'issuer-008', NOW() + INTERVAL '5 years', 'Hana Hrnić', '312', 9100.99),
      (gen_random_uuid(), '5105105105105100', '8701123377', false, 'ISSUER009', 'pass999', 'issuer-009', NOW() + INTERVAL '2 years', 'Igor Ilić', '808', 4300.00),
      (gen_random_uuid(), '6011000990139424', '8701123388', false, 'ISSUER010', 'pass910', 'issuer-010', NOW() + INTERVAL '1 year', 'Jovana Jović', '510', 14800.00),
      (gen_random_uuid(), '4007000000027',     '8701123390', false, 'ISSUER011', 'pass911', 'issuer-011', NOW() + INTERVAL '3 years', 'Kosta Krstić', '305', 5600.00),
      (gen_random_uuid(), '6221260000000000', '8701123391', false, 'ISSUER012', 'pass912', 'issuer-012', NOW() + INTERVAL '4 years', 'Lena Lukić', '722', 10000.00),
      (gen_random_uuid(), '3528000000000000', '8701123392', false, 'ISSUER013', 'pass913', 'issuer-013', NOW() + INTERVAL '2 years', 'Milan Matić', '066', 6700.50),
      (gen_random_uuid(), '4012888888881881', '8701123393', false, 'ISSUER014', 'pass914', 'issuer-014', NOW() + INTERVAL '1 year', 'Nina Novak', '448', 3900.00),
      (gen_random_uuid(), '6011111111111117', '8701123394', false, 'ISSUER015', 'pass915', 'issuer-015', NOW() + INTERVAL '5 years', 'Ognjen Ostojić', '917', 8700.00),
      (gen_random_uuid(), '3530111333300001', '8701123395', false, 'ISSUER016', 'pass916', 'issuer-016', NOW() + INTERVAL '2 years', 'Petra Popović', '129', 7200.20),
      (gen_random_uuid(), '5105105105105101', '8701123396', false, 'ISSUER017', 'pass917', 'issuer-017', NOW() + INTERVAL '3 years', 'Rade Radić', '631', 8900.00),
      (gen_random_uuid(), '4111111111111111', '8701123397', false, 'ISSUER018', 'pass918', 'issuer-018', NOW() + INTERVAL '4 years', 'Sofija Simić', '299', 11100.00),
      (gen_random_uuid(), '4012888888881882', '8701123398', false, 'ISSUER019', 'pass919', 'issuer-019', NOW() + INTERVAL '2 years', 'Teodora Tasić', '555', 5050.00),
      (gen_random_uuid(), '4222222222222',    '8701123399', false, 'ISSUER020', 'pass920', 'issuer-020', NOW() + INTERVAL '3 years', 'Uroš Uzelac', '777', 3000.00);
