-- +goose Up
-- Migration 0001: Initial Infrastructure, Extensions, Roles & Base Types
-- PayBridge AI — PostgreSQL 16 Foundation

-- 1. Enable Required Extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

-- 2. Provision Database Roles (Idempotent)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'paybridge_admin') THEN
        CREATE ROLE paybridge_admin WITH LOGIN NOBYPASSRLS;
    END IF;
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'paybridge_app') THEN
        CREATE ROLE paybridge_app WITH LOGIN NOBYPASSRLS;
    END IF;
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'paybridge_worker') THEN
        CREATE ROLE paybridge_worker WITH LOGIN NOBYPASSRLS;
    END IF;
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'paybridge_connector') THEN
        CREATE ROLE paybridge_connector WITH LOGIN NOBYPASSRLS;
    END IF;
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'paybridge_auth') THEN
        CREATE ROLE paybridge_auth WITH LOGIN NOBYPASSRLS;
    END IF;
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'paybridge_webhook') THEN
        CREATE ROLE paybridge_webhook WITH LOGIN NOBYPASSRLS;
    END IF;
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'paybridge_retention') THEN
        CREATE ROLE paybridge_retention WITH LOGIN NOBYPASSRLS;
    END IF;
END $$;

-- 3. Schema Security Lockdown
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO paybridge_admin, paybridge_app, paybridge_worker, paybridge_connector, paybridge_auth, paybridge_webhook, paybridge_retention;
GRANT ALL ON SCHEMA public TO paybridge_admin;

-- 4. Schema Migrations Table (Goose Version Tracking / Infrastructure Table 1/40)
CREATE TABLE IF NOT EXISTS schema_migrations (
    version_id BIGINT PRIMARY KEY,
    is_applied BOOLEAN NOT NULL DEFAULT TRUE,
    tstamp TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE schema_migrations OWNER TO paybridge_admin;
REVOKE ALL ON schema_migrations FROM PUBLIC, paybridge_app, paybridge_worker, paybridge_connector, paybridge_auth, paybridge_webhook, paybridge_retention;

-- +goose Down
DROP TABLE IF EXISTS schema_migrations;
