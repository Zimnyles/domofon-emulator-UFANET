CREATE TABLE intercom (
    id SERIAL PRIMARY KEY,
    mac_address VARCHAR(17) UNIQUE NOT NULL,
    domofon_status BOOLEAN DEFAULT TRUE,
    door_status BOOLEAN DEFAULT FALSE,
    address TEXT NOT NULL,
    number_of_apartments INTEGER NOT NULL,
    is_calling BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);