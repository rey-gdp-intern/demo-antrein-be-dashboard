CREATE TABLE IF NOT EXISTS tenants (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(75) UNIQUE NOT NULL,
    password VARCHAR(155) NOT NULL,
    name VARCHAR(155) NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS projects (
    id VARCHAR(75) PRIMARY KEY,
    name VARCHAR(155) NOT NULL,
    tenant_id uuid NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp
);

DO $$ BEGIN
    CREATE TYPE style AS ENUM ('base', 'custom');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS configurations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id VARCHAR(75) NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    threshold INTEGER DEFAULT 0,
    session_time INTEGER DEFAULT 5,
    host VARCHAR(155),
    base_url VARCHAR(155),
    max_users_in_queue INTEGER DEFAULT 0,
    queue_start TIMESTAMP,
    queue_end TIMESTAMP,
    queue_page_style style DEFAULT 'base',
    queue_html_page VARCHAR(155),
    queue_page_base_color VARCHAR(10),
    queue_page_title VARCHAR(155),
    queue_page_logo VARCHAR(155),
    is_configure boolean DEFAULT FALSE,
    updated_at timestamp
);