
-- Updated SQL script to create the uploads table with 'username' field
CREATE TABLE IF NOT EXISTS uploads (
    id INT AUTO_INCREMENT PRIMARY KEY,
    filename VARCHAR(255),
    filedate DATETIME,
    sha256hash CHAR(64),
    vesselid INT CHECK (vesselid >= 0 AND vesselid <= 65535),
    username VARCHAR(100)
);
