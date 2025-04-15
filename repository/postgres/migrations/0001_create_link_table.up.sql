CREATE TABLE "link" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    original TEXT NOT NULL,
    short VARCHAR(16) UNIQUE NOT NULL,
    max_uses INTEGER DEFAULT NULL,
    click_count INTEGER DEFAULT 0,
    expires_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT now()
);
