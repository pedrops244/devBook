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


CREATE TABLE Pedidos (
    ID INT IDENTITY PRIMARY KEY,
    Status NVARCHAR(50) NOT NULL, -- Ex.: 'Em andamento', 'Enviado', 'Conferido'
    DataCriacao DATETIME DEFAULT GETDATE(), -- Data e hora da criação do pedido
    UsuarioID INT NOT NULL, -- ID do usuário que criou o pedido
    CONSTRAINT FK_Pedidos_Usuario FOREIGN KEY (UsuarioID) REFERENCES Usuarios(ID)
);


CREATE TABLE ItensPedidos (
    ID INT IDENTITY PRIMARY KEY,
    PedidoID INT NOT NULL, -- Relacionamento com a tabela Pedidos
    Codigo NVARCHAR(50) NOT NULL, -- Código de barras do produto
    QuantidadeSolicitada INT NOT NULL,
    QuantidadeConferida INT DEFAULT 0, 
    QuantidadeAprovada INT DEFAULT 0, 
    CONSTRAINT FK_ItensPedidos_Pedidos FOREIGN KEY (PedidoID) REFERENCES Pedidos(ID)
);