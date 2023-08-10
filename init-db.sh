set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    use database;
    CREATE TABLE user (
    	id bigint(20) PRIMARY KEY,
    	full_name VARCHAR (60) NOT NULL,
    	phone_number VARCHAR(13) NOT NULL,
    	password TEXT NOT NULL,
    	created_at timestamp DEFAULT now() NOT NULL,
    	updated_at timestamp DEFAULT now() NOT NULL
    );
EOSQL
