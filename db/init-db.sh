#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting database initialization...${NC}"


# Database 1: Orchestrator Service
echo -e "${GREEN}Creating orchestrator database and user...${NC}"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    -- Create user for orchestrator service
    CREATE USER ${ORCHESTRATOR_DB_USER} WITH PASSWORD '${ORCHESTRATOR_DB_PASS}';
    CREATE DATABASE orchestrator_db OWNER ${ORCHESTRATOR_DB_USER};
    GRANT ALL PRIVILEGES ON DATABASE orchestrator_db TO ${ORCHESTRATOR_DB_USER};
EOSQL

# Connect to orchestrator_db and set up schema permissions
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname=orchestrator_db <<-EOSQL
    -- Grant schema permissions
    GRANT ALL ON SCHEMA public TO ${ORCHESTRATOR_DB_USER};

    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ${ORCHESTRATOR_DB_USER};
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO ${ORCHESTRATOR_DB_USER};
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO ${ORCHESTRATOR_DB_USER};
EOSQL

echo -e "${GREEN}Orchestrator database setup complete!${NC}"

# Database 2: Video Processing Service
echo -e "${GREEN}Creating video processing database and user...${NC}"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE USER ${VPS_DB_USER} WITH PASSWORD '${VPS_DB_PASS}';
    CREATE DATABASE video_processing_db OWNER ${VPS_DB_USER};
    GRANT ALL PRIVILEGES ON DATABASE video_processing_db TO ${VPS_DB_USER};
EOSQL

# Connect to video_processing_db and set up schema permissions
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname=video_processing_db <<-EOSQL
    GRANT ALL ON SCHEMA public TO ${VPS_DB_USER};
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ${VPS_DB_USER};
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO ${VPS_DB_USER};
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO ${VPS_DB_USER};
EOSQL

echo -e "${GREEN}Video processing database setup complete!${NC}"
echo -e "${BLUE}Database initialization finished successfully!${NC}"
