CREATE DATABASE IF NOT EXISTS AppRep;
USE AppRep;

DROP TABLE IF EXISTS usuarios;

CREATE TABLE usuarios (
    id INT IDENTITY PRIMARY KEY,
    username NVARCHAR(100) NOT NULL,
    senha NVARCHAR(255) NOT NULL,
    role NVARCHAR(50) NOT NULL DEFAULT 'vendedor',
    created_at DATETIME DEFAULT GETDATE(),
    CONSTRAINT chk_role CHECK (role IN ('vendedor', 'repositor', 'admin'))
);
