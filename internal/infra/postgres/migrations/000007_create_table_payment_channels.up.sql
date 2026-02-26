CREATE TABLE payment_channels (
  id SERIAL PRIMARY KEY,

  code VARCHAR(50) UNIQUE NOT NULL,
  name VARCHAR(100) NOT NULL,

  -- gateway fee (dari midtrans)
  gateway_fee_fixed BIGINT NOT NULL DEFAULT 0,
  gateway_fee_percentage_bps INT NOT NULL DEFAULT 0,

  -- default platform fee
  platform_fee_fixed BIGINT NOT NULL DEFAULT 0,
  platform_fee_percentage_bps INT NOT NULL DEFAULT 200,

  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


INSERT INTO payment_channels 
(code, name, gateway_fee_fixed, gateway_fee_percentage_bps, platform_fee_fixed, platform_fee_percentage_bps, is_active, created_at) 
VALUES

-- QRIS (MDR 0.7%)
('qris', 'QRIS (All providers)', 0, 70, 0, 200, TRUE, now()),

-- QRIS detail
('qris_gopay', 'QRIS GoPay', 0, 70, 0, 200, TRUE, now()),
('qris_shopeepay', 'QRIS ShopeePay', 0, 70, 0, 200, TRUE, now()),
('qris_dana', 'QRIS DANA', 0, 70, 0, 200, TRUE, now()),

-- E-Wallet native (Midtrans pricing umum)
('ewallet_gopay', 'GoPay (DeepLink)', 0, 200, 0, 200, TRUE, now()),
('ewallet_shopeepay', 'ShopeePay', 0, 200, 0, 200, TRUE, now()),
('ewallet_dana', 'DANA', 0, 150, 0, 200, TRUE, now());

