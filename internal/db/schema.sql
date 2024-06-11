CREATE TABLE companies (
    name TEXT NOT NULL,
    alias TEXT NOT NULL,
    site TEXT NOT NULL,
    priority INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT u_listing UNIQUE (
        name, site
    ),
    CONSTRAINT u_alias UNIQUE (
        alias
    )
);
