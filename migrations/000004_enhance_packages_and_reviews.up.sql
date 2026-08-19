-- Migration 000004: Enhance packages with badge/tag and ensure reviews & landing sections support full CRUD
ALTER TABLE packages ADD COLUMN IF NOT EXISTS badge VARCHAR(50) DEFAULT '';

-- Pastikan seeder default packages ada di database jika tabel kosong
INSERT INTO packages (id, name, description, price, benefits, is_active, display_order, badge, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    'Silver',
    'Paket ekonomis dengan fitur esensial untuk acara pernikahan yang intim.',
    99000,
    '["1 Pilihan Desain Elegan", "Masa Aktif 3 Bulan", "Hitung Mundur Waktu Acara", "Galeri Foto (Hingga 5 Foto)", "Google Maps Navigasi Lokasi", "RSVP & Kolom Ucapan Tamu"]'::jsonb,
    TRUE,
    1,
    'Hemat',
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM packages WHERE name = 'Silver');

INSERT INTO packages (id, name, description, price, benefits, is_active, display_order, badge, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    'Gold',
    'Paket paling populer dengan kustomisasi lengkap dan fitur multimedia interaktif.',
    199000,
    '["Semua Fitur Paket Silver", "Masa Aktif 6 Bulan", "Galeri Foto (Hingga 15 Foto) & Background Music", "Amplop Digital & QRIS Pembayaran", "Kirim Undangan Nama Tamu Otomatis (WhatsApp Generator)", "Cerita Cinta (Love Story Timeline)", "Live Streaming Link Integration"]'::jsonb,
    TRUE,
    2,
    'Paling Populer',
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM packages WHERE name = 'Gold');

INSERT INTO packages (id, name, description, price, benefits, is_active, display_order, badge, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    'Platinum Custom Domain',
    'Paket eksklusif dengan Custom Domain pribadi (.com / .id) untuk pengalaman tak terlupakan.',
    349000,
    '["Semua Fitur Paket Gold", "Custom Domain Pribadi (contoh: rani-dan-budi.com)", "Gratis SSL / HTTPS Otomatis", "Masa Aktif 1 Tahun Penuh", "Galeri Foto Tanpa Batas & Video Embed", "Prioritas Dukungan VIP & Bebas Revisi", "Statistik Tamu & Export Data RSVP Excel"]'::jsonb,
    TRUE,
    3,
    'Eksklusif Domain',
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM packages WHERE name = 'Platinum Custom Domain');
