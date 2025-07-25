#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
    CREATE DATABASE identity_db;
    CREATE DATABASE order_db;
EOSQL


##!/bin/bash
#set -e
#
## Connect to the default 'postgres' database to run our setup commands
#psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
#    -- Check if 'identity_db' exists, and if not, create it
#    SELECT 'CREATE DATABASE identity_db'
#    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'identity_db')\gexec
#
#    -- Check if 'order_db' exists, and if not, create it
#    SELECT 'CREATE DATABASE order_db'
#    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'order_db')\gexec
#EOSQL