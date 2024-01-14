CREATE TABLE
    IF NOT EXISTS ipfs_metadata (
        cid character varying(255) PRIMARY KEY,
        token JSONB
    );