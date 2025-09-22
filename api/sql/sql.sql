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

// trocar a senha acima por uma senha segura em produção
