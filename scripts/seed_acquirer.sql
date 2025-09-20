-- INSERT INTO accounts (
--     id, primary_account_number, number, merchant_account, bank_identifier_code, merchant_password,
--     merchant_id, expiration_date, card_holder_name, security_code, balance
-- ) VALUES
-- -- Merchant accounts
-- ('1e2f3a4b-5c6d-7e8f-9012-3456789abc01', '3879520000012345', '160-13646185005023-13', true, '387952', 'm3rchant!', 'MERCH-1001', NOW() + INTERVAL '2 years', 'Nikola Nikolic', '321', 234500.75),
-- ('2b3c4d5e-6f7a-8b9c-0123-456789abcdef', '3879520000023456', '160-44811483878583-30', true, '387952', 's3cur3Pa$$', 'MERCH-1002', NOW() + INTERVAL '3 years', 'Milena Milenkovic', '453', 152000.00),
-- ('3c4d5e6f-7a8b-9c0d-1234-56789abcdef0', '3879520000034567', '160-33752021114200-38', true, '387952', 'P@ssw0rd', 'MERCH-1003', NOW() + INTERVAL '1 year', 'Petar Petrovic', '987', 199990.99),
-- ('4d5e6f7a-8b9c-0d1e-2345-6789abcdef01', '3879520000045678', '160-12188161996163-87', true, '387952', 'safe1234', 'MERCH-1004', NOW() + INTERVAL '4 years', 'Jelena Jankovic', '129', 182500.10),
-- ('5e6f7a8b-9c0d-1e2f-3456-789abcdef012', '3879520000056789', '160-32450648275404-70', true, '387952', 'merchant5!', 'MERCH-1005', NOW() + INTERVAL '2 years', 'Marko Markovic', '213', 225000.00),
--
-- -- Regular personal accounts
-- ('6f7a8b9c-0d1e-2f3a-4567-89abcdef0123', '3879520000067890', '160-65935980837512-98', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'Ana Anic', '621', 62000.50),
-- ('7a8b9c0d-1e2f-3a4b-5678-9abcdef01234', '3879520000078901', '160-78897075011199-57', false, '387952', NULL, NULL, NOW() + INTERVAL '1 year', 'Milan Milic', '874', 430000.00),
-- ('8b9c0d1e-2f3a-4b5c-6789-abcdef012345', '3879520000089012', '160-50846027517296-38', false, '387952', NULL, NULL, NOW() + INTERVAL '3 years', 'Ivana Ivanovic', '753', 5500.00),
-- ('9c0d1e2f-3a4b-5c6d-789a-bcdef0123456', '3879520000090123', '160-40040111647409-81', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'Stefan Stevic', '214', 88000.90),
-- ('0d1e2f3a-4b5c-6d7e-89ab-cdef01234567', '3879520000101234', '160-19574843957189-95', false, '387952', NULL, NULL, NOW() + INTERVAL '4 years', 'Teodora Todorovic', '634', 127500.25),
-- ('1f2e3d4c-5b6a-7f8e-9012-3456789abc11', '3879520000112345', '160-42910618998733-02', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'Vuk Vukovic', '147', 4000.00),
-- ('2f3e4d5c-6b7a-8f9e-0123-456789abcd12', '3879520000123456', '160-41898313465951-95', false, '387952', NULL, NULL, NOW() + INTERVAL '1 year', 'Tamara Tamic', '369', 3200.75),
-- ('3f4e5d6c-7b8a-9f0e-1234-56789abcde23', '3879520000134567', '160-99216092923556-56', false, '387952', NULL, NULL, NOW() + INTERVAL '5 years', 'Aleksa Aleksic', '852', 23000.00),
-- ('4f5e6d7c-8b9a-0f1e-2345-6789abcdef34', '3879520000145678', '160-19861636440986-12', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'Nina Nikolic', '951', 103200.00),
-- ('5f6e7d8c-9b0a-1f2e-3456-789abcdef045', '3879520000156789', '160-78054280842909-24', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'Dusan Dusanic', '142', 1000.00),
--
-- -- More merchant accounts
-- ('6f7e8d9c-0b1a-2f3e-4567-89abcdef0567', '3879520000167890', '160-17747416133912-09', true, '387952', 'merchant6!', 'MERCH-1006', NOW() + INTERVAL '2 years', 'Kristina Krstic', '777', 265000.00),
-- ('7f8e9d0c-1b2a-3f4e-5678-9abcdef06789', '3879520000178901', '160-65094541400868-04', true, '387952', 'merchant7!', 'MERCH-1007', NOW() + INTERVAL '3 years', 'Filip Filipovic', '543', 170900.00),
-- ('8f9e0d1c-2b3a-4f5e-6789-abcdef07890a', '3879520000189012', '160-18488188011365-57', true, '387952', 'merchant8!', 'MERCH-1008', NOW() + INTERVAL '1 year', 'Natalija Novak', '456', 201000.00),
-- ('9f0e1d2c-3b4a-5f6e-7890-bcdef08901ab', '3879520000190123', '160-85473844459373-60', true, '387952', 'merchant9!', 'MERCH-1009', NOW() + INTERVAL '2 years', 'Sava Savic', '234', 300000.00),
-- ('0f1e2d3c-4b5a-6f7e-8901-cdef0123456b', '3879520000201234', '160-32957845108126-78', true, '387952', 'merchant10!','MERCH-1010', NOW() + INTERVAL '4 years', 'Marija Maric', '999', 188000.00);

INSERT INTO accounts (
    id, primary_account_number, number, merchant_account, bank_identifier_code, merchant_password,
    merchant_id, expiration_date, card_holder_name, balance
) VALUES
-- Merchant accounts
('1e2f3a4b-5c6d-7e8f-9012-3456789abc01', 'a1ZAtuU9iC+9ltKo/eEI+MVPY0eMUWvC93owm9Yx/v0FJuyphtok86/UCRw=', '160-13646185005023-13', true, '387952', 'm3rchant!', 'MERCH-1001', NOW() + INTERVAL '2 years', 'v7qt/J89sLXYZgTBfbVYIS0C6WD/tYnklSaEGzSEFXeh45PC8OeEVFy9', 234500.75),
('2b3c4d5e-6f7a-8b9c-0123-456789abcdef', '7U+SeawZlN+wV0wTp0O5s7VTJdhtp/5P7AGwlTyPZdbggdxBLbbx8D96+8U=', '160-44811483878583-30', true, '387952', 's3cur3Pa$$', 'MERCH-1002', NOW() + INTERVAL '3 years', 'MLeMl+jYFscc8iUoIcwk3gYt8gPnKl0xb/nCvHPtX+qpMAiNXT7nPeT4UTcy', 152000.00),
('3c4d5e6f-7a8b-9c0d-1234-56789abcdef0', 'nLKw5GaBUDXuui+VXdoRv/Op5+B890QbfDaXGDN4egqR07NYLU+C2TFShw8=', '160-33752021114200-38', true, '387952', 'P@ssw0rd', 'MERCH-1003', NOW() + INTERVAL '1 year', '9OFDNHjGomIjIdqBomDB2g91zVy+prHBaGGdDfii/9domScLcHYBmUNe', 199990.99),
('4d5e6f7a-8b9c-0d1e-2345-6789abcdef01', 'jHa9N3duobQ69K1VigJ+HRQONw5p2iM5FDjfzHrrS8iQSw/ena9/ffyNQD8=', '160-12188161996163-87', true, '387952', 'safe1234', 'MERCH-1004', NOW() + INTERVAL '4 years', 'nah59lccSduqg3MpCOMBC12VnFrq6xsw8jgDRnU8q22hPyd4MHR9IDiKxQ==', 182500.10),
('5e6f7a8b-9c0d-1e2f-3456-789abcdef012', 'z+dEQBK+Xq4HOD/mqbtau+hOjYV4u8Cl/3a0WfDk4wvbkPssDqQkWVIOSrE=', '160-32450648275404-70', true, '387952', 'merchant5!', 'MERCH-1005', NOW() + INTERVAL '2 years', 'QIvX+JFsNeA1Fhm/KN+KCddhEatts7oXgJ+nhXtTB04PHTSWq4J3nxKk', 225000.00),

-- Regular personal accounts
('6f7a8b9c-0d1e-2f3a-4567-89abcdef0123', 'Y8USA/trDQBBqSxkl5QRasZGPrfk/fTlilUlm/ClzhzsLtJ+nRtbiWc2El8=', '160-65935980837512-98', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'bAVE+Q8onTKsIYuqJPOAA/izTYkhNvYLtr7IDW8++SiWyrX0FFjTequL', 62000.50),
('7a8b9c0d-1e2f-3a4b-5678-9abcdef01234', 'Y/OeSk8aIHcil1SVxj6qVnqaJ+IU29CuQuW/pP5VIsJbRZ4Mk+bQZrVvyk0=', '160-78897075011199-57', false, '387952', NULL, NULL, NOW() + INTERVAL '1 year', 'lafYbuNapq1o6kkOCDcJ3HxKakyNYRr+Xn/6URCAxoI7FdPP', 430000.00),
('8b9c0d1e-2f3a-4b5c-6789-abcdef012345', 'jMtgWIOhBJL8YFk3dXHqZKGj3Pf1t95ZZEurTlXm11y+Mhy9kwKpVxN7Xbc=', '160-50846027517296-38', false, '387952', NULL, NULL, NOW() + INTERVAL '3 years', 'iZE7cJi6JwefkH9E39YXr23E6oz0AI8mCCcBH3kbEmtaxF7LNHhY', 5500.00),
('9c0d1e2f-3a4b-5c6d-789a-bcdef0123456', 'OsGIfdOpXFdYRtEOKckqqSb8HCYigzIB9vF8D00/Nh8Z/Dyc7lzGcgbHEsM=', '160-40040111647409-81', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'CecNcRVVWlfj0AxoAcv5GWiXoAQSTquN1i3xWIP5gX2L5k/t3gz6k1k8', 88000.90),
('0d1e2f3a-4b5c-6d7e-89ab-cdef01234567', 'AAV8YlbkEfYzbMuxzgSQDitBx7jx5wYmIbCMQ2VUEqxeT+VAVuhQo7xhdko=', '160-19574843957189-95', false, '387952', NULL, NULL, NOW() + INTERVAL '4 years', 'VGm1CtYTxE3Z4dLkIcaDkEdZQ+aAR+p2Uj44SrqEgq8w7JF3Sb3mxHk=',127500.25),
('1f2e3d4c-5b6a-7f8e-9012-3456789abc11', 'GGqbnj+MaFCSACMkp8DNgOTCfGtWdDartx3TYfXnTwtw/WxsQRwxYyeWicc=', '160-42910618998733-02', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'reeCqoQ8Scu8xHpwjezyRe6QBU/3YTYbK11esKZn3GX4op5J3Kn0xQ1mk3aI', 4000.00),
('2f3e4d5c-6b7a-8f9e-0123-456789abcd12', 'psbd8/F8Q10izPbnlnLAMG8OoXj66nRxC8+lFiKhJCbQTdGMxMYipS/3wXU=', '160-41898313465951-95', false, '387952', NULL, NULL, NOW() + INTERVAL '1 year', 'jGq2ppI4+DT0xaBhi9nBv10BQWbfkYGc/tE2XPVYI2CMJVaB3auL',3200.75),
('3f4e5d6c-7b8a-9f0e-1234-56789abcde23', '7H3aul84qds5g1hwMRlwvvic8CkoXZjcfDfrOex/P9uoiEcBgzZzL1+2a7Q=', '160-99216092923556-56', false, '387952', NULL, NULL, NOW() + INTERVAL '5 years', '3AyZcafeCID3yvNU5+UWKc4tmGaEfXF95wAu7uQnwbpwU77TO443CQ==', 23000.00),
('4f5e6d7c-8b9a-0f1e-2345-6789abcdef34', 'z/rJPrrTEIYphdFW7ppLSKtGzrEXYlx0FOrSjALXzYESCWbu4JNRI1pEXPc=', '160-19861636440986-12', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'Bu9HFj0NcsT0072gC3mtGloX6uCmh3FdETfFKmoQ32TsQcAJ0BbGrsk2', 103200.00),
('5f6e7d8c-9b0a-1f2e-3456-789abcdef045', 'pFq/2SieMlqaQAJbIujZ44gVzMcNRutWj4GpTfF8ifePklDhHcWzL7865HM=', '160-78054280842909-24', false, '387952', NULL, NULL, NOW() + INTERVAL '2 years', 'F/v+a56Qo1FGUFh2YOgk1tYeEw/Pjj0S0hpSBzXMn7uU7HJsj/AT3g==', 1000.00),

-- More merchant accounts
('6f7e8d9c-0b1a-2f3e-4567-89abcdef0567', '5OjS5g1J8ODu9itLlRsxEwHRQrYJz67EJnqIuTEFm/fOr+oWWVWk7Z9/llo=', '160-17747416133912-09', true, '387952', 'merchant6!', 'MERCH-1006', NOW() + INTERVAL '2 years', 'hBI2MRp4nja8UwLpt9dxK0oSfQYc4//w/N7YPyg3sRrjSL7KC8gBAxw=', 265000.00),
('7f8e9d0c-1b2a-3f4e-5678-9abcdef06789', 'ONYdlAUpitXqK12qwEXIhKs9eO1Dbe2RXxGZVwyQpKAGAIr0R7oKvRw7BYU=', '160-65094541400868-04', true, '387952', 'merchant7!', 'MERCH-1007', NOW() + INTERVAL '3 years', 'GiRIz2mPiN7RQt2hQSLdfFJxnelRad+EXiNCeQEkzBpgP3720u39+yJ5cQ==', 170900.00),
('8f9e0d1c-2b3a-4f5e-6789-abcdef07890a', 'zs93OVsFULbtaH/xnnJp7DFY6PvgZezoUfC6Jtrpx58sN+BFCf5JZbzg/10=', '160-18488188011365-57', true, '387952', 'merchant8!', 'MERCH-1008', NOW() + INTERVAL '1 year', 'cnYyVOedLWCxyFLB+XLiM7S66cLC58mORZjwRjkllYt+3vQ9VwYwbdav9Q==', 201000.00),
('9f0e1d2c-3b4a-5f6e-7890-bcdef08901ab', 'ZtkE+ZkyFxR3xDtorL71wbZ4MDFiGpZXWX0tJBeHrpLn6OHD1MQKCS1+VT0=', '160-85473844459373-60', true, '387952', 'merchant9!', 'MERCH-1009', NOW() + INTERVAL '2 years', 'U7PT3i43ITX10kmngmtLLQU/Y4r3VXHF0Ok0cVxkNbqdDabKFjU=', 300000.00),
('0f1e2d3c-4b5a-6f7e-8901-cdef0123456b', 'lb0ka1EzGv+ZnK27LVEaXid577MBQo9dgqxGtSwsNGOokPpbuBkp4QBrgbk=', '160-32957845108126-78', true, '387952', 'merchant10!','MERCH-1010', NOW() + INTERVAL '4 years', '0RHvCvjUJDVu7MIAPiIXEA5Y76dk81FUEi4hrcFqQtdVuHtgMLc/cQ==', 188000.00);
