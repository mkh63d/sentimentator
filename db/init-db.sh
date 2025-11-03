#!/bin/bash
set -euo pipefail

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting database initialization...${NC}"


setup_db() {
  local db_name="$1"
  local app_user="$2"
  local app_pass="$3"

  echo -e "${GREEN}Creating ${db_name} database and user...${NC}"

  # Part 1: Create user and database if not exist; grant DB privileges

  # Create user if not exists
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<EOSQL
DO \$do\$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = '$app_user') THEN
    CREATE USER "$app_user" WITH PASSWORD '$app_pass';
  END IF;
END
\$do\$;
EOSQL

  # Create database if not exists (must be outside transaction)
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<EOSQL
SELECT 'CREATE DATABASE "$db_name" OWNER "$app_user"'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$db_name')\gexec
EOSQL

  # Grant privileges on the database
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<EOSQL
GRANT ALL PRIVILEGES ON DATABASE "$db_name" TO "$app_user";
EOSQL

  # Part 2: Grants for schema and objects in the target DB
  psql -v ON_ERROR_STOP=1 \
       --username "$POSTGRES_USER" --dbname="$db_name" <<EOSQL
-- Schema privileges
GRANT USAGE, CREATE ON SCHEMA public TO "$app_user";

-- Existing objects
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "$app_user";
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "$app_user";
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO "$app_user";

-- Future objects (effective for objects created by the role running this block)
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO "$app_user";
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO "$app_user";
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO "$app_user";
-- Optionally (PG11+):
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON ROUTINES TO "$app_user";
EOSQL

  echo -e "${GREEN}${db_name} database setup complete!${NC}"
}

# Orchestrator
setup_db "orchestrator_db" "${ORCHESTRATOR_DB_USER}" "${ORCHESTRATOR_DB_PASS}"

# Video Processing
setup_db "video_processing_db" "${VPS_DB_USER}" "${VPS_DB_PASS}"

echo -e "${BLUE}Database initialization finished successfully!${NC}"
