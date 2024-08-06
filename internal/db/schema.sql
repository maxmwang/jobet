CREATE TYPE site AS ENUM (
    'ashby',
    'greenhouse',
    'lever'
);

CREATE TABLE companies (
    name VARCHAR NOT NULL,
    alias VARCHAR NOT NULL,
    site site NOT NULL,
    priority INT4 NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,

    CONSTRAINT u_listing PRIMARY KEY (
         name, site
   )
);

CREATE TABLE discord_channels (
    id VARCHAR PRIMARY KEY,
    created_at timestamptz DEFAULT NOW() NOT NULL
)
