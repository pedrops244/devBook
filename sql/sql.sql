-- Verifique se o banco de dados já existe antes de criar
IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'app_rep')
BEGIN
    CREATE DATABASE app_rep;
END;

-- Selecionar o banco de dados
USE app_rep;

-- Remover tabelas existentes, caso já existam
DROP TABLE IF EXISTS usuarios;
DROP TABLE IF EXISTS pedidos;
DROP TABLE IF EXISTS itens_pedidos;

-- Criação da tabela 'usuarios'
CREATE TABLE usuarios (
    id INT IDENTITY PRIMARY KEY,
    username NVARCHAR(100) NOT NULL,
    senha NVARCHAR(255) NOT NULL,
    is_deleted BIT DEFAULT 0,  -- ALTEREI BOOLEAN para BIT, pois no SQL Server não existe o tipo BOOLEAN
    role NVARCHAR(50) NOT NULL DEFAULT 'comprador',
    created_at DATETIME DEFAULT GETDATE(),
    CONSTRAINT chk_role CHECK (role IN ('comprador', 'repositor', 'admin', 'gerente'))
);

-- Inserção de dados na tabela 'usuarios'
INSERT INTO usuarios (username, senha, role)
VALUES ('amigaoadmin', '$2a$10$Sp3T.233Ouy5EGv9lzbOA.PU0aknljai/VFkUqEVaz5L.zSudeKNe', 'admin');

-- Criação da tabela 'pedidos'
CREATE TABLE pedidos (
    id INT IDENTITY PRIMARY KEY,
    status NVARCHAR(50) NOT NULL, -- Ex.: 'criado', 'recebido', 'conferido'
    usuario_id INT NOT NULL, -- ID do usuário que criou o pedido
    criado_em DATETIME DEFAULT GETDATE(), -- Data e hora da criação do pedido
    recebido_em DATETIME NULL, -- Data e hora em que o pedido foi recebido
    conferido_em DATETIME NULL, -- Data e hora em que o pedido foi conferido
    CONSTRAINT FK_pedidos_usuario FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
);

-- Criação da tabela 'itens_pedidos'
CREATE TABLE itens_pedidos (
    id INT IDENTITY PRIMARY KEY,
    pedido_id INT NOT NULL, -- Relacionamento com a tabela Pedidos
    codigo NVARCHAR(50) NOT NULL, -- Código de barras do produto
    quantidade_solicitada INT NOT NULL,
    quantidade_recebida INT DEFAULT 0, 
    quantidade_conferida INT DEFAULT 0, 
    CONSTRAINT FK_itens_pedidos_pedidos FOREIGN KEY (pedido_id) REFERENCES pedidos(id) ON DELETE CASCADE
);

CREATE TABLE estoque (
    id INT IDENTITY PRIMARY KEY,
    codigo NVARCHAR(50) NOT NULL UNIQUE, -- Código único do produto
    quantidade INT NOT NULL,            -- Quantidade total no estoque
    reservado INT DEFAULT 0,            -- Quantidade reservada para pedidos
    faltas INT DEFAULT 0,            -- Quantidade em falta para pedidos
    criado_em DATETIME DEFAULT GETDATE(), -- Data de criação
);
