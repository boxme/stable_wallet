CREATE TABLE IF NOT EXISTS users {
    id bigserial PRIMARY KEY,
    phone_number integer NOT NULL,
    country_code integer NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    email citext,
    CONSTRAINT fk_country_code FOREIGN KEY(country_code) REFERENCES countries(country_code)
};