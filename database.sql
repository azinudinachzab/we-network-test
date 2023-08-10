/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE users (
                       id bigint PRIMARY KEY,
                       full_name VARCHAR (60) NOT NULL,
                       phone_number VARCHAR(13) UNIQUE NOT NULL,
                       password TEXT NOT NULL,
                       created_at timestamp DEFAULT now() NOT NULL,
                       updated_at timestamp DEFAULT now() NOT NULL
);

CREATE TABLE user_counts (
                        user_id bigint PRIMARY KEY,
                        count bigint NOT NULL DEFAULT 0,
                        created_at timestamp DEFAULT now() NOT NULL,
                        updated_at timestamp DEFAULT now() NOT NULL
);
