-- Create AdUnits table
CREATE TABLE IF NOT EXISTS AdUnits (
    ID TEXT PRIMARY KEY,
    Format TEXT,
    Width INT,
    Height INT
);

-- Create Creatives table
CREATE TABLE IF NOT EXISTS Creatives (
    ID TEXT PRIMARY KEY,
    Format TEXT,
    Width INT,
    Height INT,
    Content TEXT,
    Price REAL
);

-- Insert sample AdUnits
INSERT INTO AdUnits (ID, Format, Width, Height)
VALUES
    ('adunit1', 'banner', 300, 250),
    ('adunit2', 'interstitial', 1024, 768);

-- Insert sample Creatives
INSERT INTO Creatives (ID, Format, Width, Height, Content, Price)
VALUES
    ('creative1', 'banner', 300, 250, 'Sample Banner Ad', 15),
    ('creative2', 'interstitial', 1024, 768, 'Sample Interstitial Ad', 3.0);