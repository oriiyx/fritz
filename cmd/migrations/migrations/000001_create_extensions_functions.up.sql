CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;


-- Function to hash password
CREATE OR REPLACE FUNCTION hash_password(password TEXT)
    RETURNS TEXT AS
$$
BEGIN
    RETURN crypt(password, gen_salt('bf', 10)); -- bcrypt with cost factor 10
END;
$$ LANGUAGE plpgsql;

-- Function to verify password
CREATE OR REPLACE FUNCTION verify_password(password TEXT, hash TEXT)
    RETURNS BOOLEAN AS
$$
BEGIN
    RETURN hash = crypt(password, hash);
END;
$$ LANGUAGE plpgsql;