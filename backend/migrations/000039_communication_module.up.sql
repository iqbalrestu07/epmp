-- Migration: 000039_communication_module.up.sql
-- Tables for WhatsApp gateway devices and blast message history

-- WA Gateway Devices
CREATE TABLE IF NOT EXISTS wa_devices (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id      UUID NOT NULL,
    label       VARCHAR(100) NOT NULL,
    phone       VARCHAR(30),
    status      VARCHAR(20) NOT NULL DEFAULT 'disconnected', -- connected | disconnected | qr_pending
    session_data TEXT,
    last_seen   TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Blast Messages (campaigns)
CREATE TABLE IF NOT EXISTS blast_messages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id          UUID NOT NULL,
    device_id       UUID REFERENCES wa_devices(id) ON DELETE SET NULL,
    title           VARCHAR(200),
    template        TEXT NOT NULL,
    target_type     VARCHAR(50) NOT NULL DEFAULT 'all_tenants', -- all_tenants | overdue | building | manual
    target_filter   JSONB,
    status          VARCHAR(20) NOT NULL DEFAULT 'draft', -- draft | sending | done | failed
    total_recipients INT NOT NULL DEFAULT 0,
    sent_count      INT NOT NULL DEFAULT 0,
    failed_count    INT NOT NULL DEFAULT 0,
    scheduled_at    TIMESTAMPTZ,
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    created_by      UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Individual message logs per recipient
CREATE TABLE IF NOT EXISTS blast_message_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blast_id        UUID NOT NULL REFERENCES blast_messages(id) ON DELETE CASCADE,
    tenant_id       UUID,
    phone           VARCHAR(30) NOT NULL,
    recipient_name  VARCHAR(200),
    message         TEXT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending | sent | failed | read
    error_message   TEXT,
    sent_at         TIMESTAMPTZ,
    read_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Message templates
CREATE TABLE IF NOT EXISTS message_templates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id      UUID NOT NULL,
    name        VARCHAR(100) NOT NULL,
    content     TEXT NOT NULL,
    variables   TEXT[], -- e.g. {tenant_name, room_name, invoice_amount}
    category    VARCHAR(50) DEFAULT 'general', -- general | invoice | reminder | announcement
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed default templates
INSERT INTO message_templates (org_id, name, content, variables, category) VALUES
  ('00000000-0000-0000-0000-000000000000', 'Tagihan Bulanan', 
   'Halo {{tenant_name}}, kami ingin menginformasikan bahwa tagihan sewa bulan ini sebesar *Rp {{invoice_amount}}* telah jatuh tempo pada {{due_date}}. Mohon segera lakukan pembayaran. Terima kasih.', 
   ARRAY['tenant_name','invoice_amount','due_date'], 'invoice'),
  ('00000000-0000-0000-0000-000000000000', 'Pengumuman Umum', 
   'Halo {{tenant_name}}, kami ingin menyampaikan pengumuman penting terkait hunian Anda di {{room_name}}. Harap perhatikan informasi berikut: {{message}}. Terima kasih atas perhatiannya.', 
   ARRAY['tenant_name','room_name','message'], 'announcement'),
  ('00000000-0000-0000-0000-000000000000', 'Pengingat Pembayaran', 
   'Halo {{tenant_name}}, ini adalah pengingat bahwa tagihan Anda sebesar *Rp {{invoice_amount}}* akan jatuh tempo dalam {{days_left}} hari. Silakan lakukan pembayaran sebelum tanggal {{due_date}}. Info: {{room_name}}.', 
   ARRAY['tenant_name','invoice_amount','days_left','due_date','room_name'], 'reminder'),
  ('00000000-0000-0000-0000-000000000000', 'Selamat Datang', 
   'Selamat datang {{tenant_name}}! Kami senang Anda bergabung. Kamar Anda adalah {{room_name}} dengan harga sewa *Rp {{rent_amount}}*/bulan. Kontrak berlaku mulai {{start_date}}. Jangan ragu untuk menghubungi kami jika ada pertanyaan.', 
   ARRAY['tenant_name','room_name','rent_amount','start_date'], 'general');
