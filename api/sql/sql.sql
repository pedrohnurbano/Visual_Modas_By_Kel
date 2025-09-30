CREATE DATABASE IF NOT EXISTS visualmodasbykel;
USE visualmodasbykel;

DROP TABLE IF EXISTS produtos;
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
    criadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_cpf (cpf),
    INDEX idx_role (role)
) ENGINE=InnoDB;

CREATE TABLE produtos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT NOT NULL,
    preco DECIMAL(10,2) NOT NULL,
    tamanho VARCHAR(10) NOT NULL,
    categoria VARCHAR(50) NOT NULL,
    foto_url LONGTEXT,
    usuario_id INT NOT NULL,
    ativo BOOLEAN DEFAULT TRUE,
    criadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    atualizadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    INDEX idx_categoria (categoria),
    INDEX idx_ativo (ativo),
    INDEX idx_usuario (usuario_id),
    INDEX idx_tamanho (tamanho)
) ENGINE=InnoDB;

-- Inserir usuário admin padrão
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

 -- Trocar por senha hasheada em produção
 