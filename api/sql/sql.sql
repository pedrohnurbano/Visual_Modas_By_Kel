CREATE DATABASE IF NOT EXISTS visualmodasbykel;
USE visualmodasbykel;

DROP TABLE IF EXISTS usuarios;

CREATE TABLE usuarios (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    sobrenome VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha VARCHAR(255) NOT NULL,
    telefone VARCHAR(20) NOT NULL,
    cpf VARCHAR(14) NOT NULL UNIQUE,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    criadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

DROP TABLE IF EXISTS produtos;

CREATE TABLE produtos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT NOT NULL,
    preco DECIMAL(10,2) NOT NULL,
    tamanho VARCHAR(10) NOT NULL,
    categoria VARCHAR(50) NOT NULL,
    foto_url VARCHAR(500),
    usuario_id INT NOT NULL,
    ativo BOOLEAN DEFAULT TRUE,
    criadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    atualizadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    INDEX idx_categoria (categoria),
    INDEX idx_ativo (ativo)
) ENGINE=InnoDB;

ALTER TABLE produtos MODIFY foto_url LONGTEXT;

INSERT INTO usuarios (nome, sobrenome, email, senha, telefone, cpf, role) 
VALUES (
    'Admin', 
    'Sistema', 
    'admin@visualmodasbykel.com', 
    '123456',
    '00000000000',
    '00000000000',
    'admin'
);

-- Trocar a senha acima por uma senha segura em produção