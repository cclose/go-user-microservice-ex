set -e
#export PGPASSWORD=$POSTGRES_PASSWORD;
#psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
#  CREATE USER $PG_USER_USER WITH PASSWORD '$PG_USER_PASSWORD';
#  CREATE DATABASE $PG_USER_DB OWNER $PG_USER_USER;
#  \connect $PG_USER_DB $PG_USER_USER
#  BEGIN;
#    CREATE TABLE IF NOT EXISTS user (
#      id INT NOT NULL PRIMARY KEY,
#      name VARCHAR(355)
#    );
#  COMMIT;
EOSQL