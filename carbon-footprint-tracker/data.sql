-- Drop tables if they exist to allow for clean re-creation (optional, for development)
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS electric_consumption CASCADE;
DROP TABLE IF EXISTS population CASCADE;
DROP TABLE IF EXISTS transport CASCADE;
DROP TABLE IF EXISTS water_consumption CASCADE;
DROP TABLE IF EXISTS waste CASCADE;
DROP TABLE IF EXISTS accommodation CASCADE;
DROP TABLE IF EXISTS goods_purchased CASCADE;

-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL -- e.g., 'admin', 'staff', 'viewer'
);

-- Electric Consumption Table
CREATE TABLE electric_consumption (
    id SERIAL PRIMARY KEY,
    source VARCHAR(255) NOT NULL, -- 'main board', 'generator'
    kwh DECIMAL(10, 2), -- Kilowatt-hours (for main board)
    fuel_liters DECIMAL(10, 2), -- Liters of fuel (for generators)
    hours DECIMAL(10, 2), -- Hours of operation (for generators)
    date DATE NOT NULL,
    location VARCHAR(255) NOT NULL DEFAULT 'Overall'
);

-- Population Table
CREATE TABLE population (
    id SERIAL PRIMARY KEY,
    registered_count INT NOT NULL,
    floating_count INT NOT NULL,
    date DATE NOT NULL,
    location VARCHAR(255) NOT NULL DEFAULT 'Overall'
);

-- Transport Table
CREATE TABLE transport (
    id SERIAL PRIMARY KEY,
    vehicle_type VARCHAR(255) NOT NULL,
    fuel_type VARCHAR(255) NOT NULL,
    distance_km DECIMAL(10, 2) NOT NULL,
    fuel_liters DECIMAL(10, 2) NOT NULL,
    date DATE NOT NULL,
    location VARCHAR(255) NOT NULL DEFAULT 'Overall'
);

-- Water Consumption Table
CREATE TABLE water_consumption (
    id SERIAL PRIMARY KEY,
    meter_reading DECIMAL(10, 2) NOT NULL,
    date DATE NOT NULL,
    location VARCHAR(255) NOT NULL DEFAULT 'Overall'
);

-- Waste Generation Table
CREATE TABLE waste (
    id SERIAL PRIMARY KEY,
    spot_name VARCHAR(255) NOT NULL,
    waste_type VARCHAR(255) NOT NULL, -- e.g., 'biodegradable', 'recyclable', 'landfill'
    weight_kg DECIMAL(10, 2) NOT NULL,
    date DATE NOT NULL,
    location VARCHAR(255) NOT NULL DEFAULT 'Overall'
);

-- Accommodation Table
CREATE TABLE accommodation (
    id SERIAL PRIMARY KEY,
    people_count INT NOT NULL,
    nights INT NOT NULL,
    date DATE NOT NULL,
    location VARCHAR(255) NOT NULL DEFAULT 'Overall'
);

-- Goods Purchased Table (Procurement)
CREATE TABLE goods_purchased (
    id SERIAL PRIMARY KEY,
    item_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    cost DECIMAL(10, 2) NOT NULL,
    date DATE NOT NULL,
    location VARCHAR(255) NOT NULL DEFAULT 'Overall'
);

select * from goods_purchased