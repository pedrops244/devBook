CREATE DATABASE IF NOT EXISTS AppRep;
USE AppRep;

DROP TABLE IF EXISTS usuarios;

CREATE TABLE usuarios (
    id INT IDENTITY PRIMARY KEY,
    username NVARCHAR(100) NOT NULL,
    senha NVARCHAR(255) NOT NULL,
    role NVARCHAR(50) NOT NULL DEFAULT 'comprador',
    created_at DATETIME DEFAULT GETDATE(),
    CONSTRAINT chk_role CHECK (role IN ('comprador', 'repositor', 'admin', 'gerente'))
);

INSERT INTO usuarios (Username, Senha, Role)
VALUES ('amigaoadmin', '$2a$10$3KQfTjX69Hu8Xvksp94Z3O7BSZp4ZmLb/NIqqRp8yqMHPomCa31bS', 'admin');



CREATE TABLE Pedidos (
    ID INT IDENTITY PRIMARY KEY,
    Status NVARCHAR(50) NOT NULL, -- Ex.: 'created', 'received', 'checked'
    CriadoEm DATETIME DEFAULT GETDATE(), -- Data e hora da criação do pedido
    RecebidoEm DATETIME NULL, -- Data e hora em que o pedido foi recebido
    ConferidoEm DATETIME NULL, -- Data e hora em que o pedido foi conferido
    UsuarioID INT NOT NULL, -- ID do usuário que criou o pedido
    CONSTRAINT FK_Pedidos_Usuario FOREIGN KEY (UsuarioID) REFERENCES Usuarios(ID)
);



CREATE TABLE ItensPedidos (
    ID INT IDENTITY PRIMARY KEY,
    PedidoID INT NOT NULL, -- Relacionamento com a tabela Pedidos
    Codigo NVARCHAR(50) NOT NULL, -- Código de barras do produto
    QuantidadeSolicitada INT NOT NULL,
    QuantidadeRecebida INT DEFAULT 0, 
    QuantidadeConferida INT DEFAULT 0, 
    CONSTRAINT FK_ItensPedidos_Pedidos FOREIGN KEY (PedidoID) REFERENCES Pedidos(ID)
);