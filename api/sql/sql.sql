CREATE DATABASE IF NOT EXISTS visualmodasbykel;
USE visualmodasbykel;

DROP TABLE IF EXISTS favoritos;
DROP TABLE IF EXISTS produtos;
DROP TABLE IF EXISTS usuarios;
DROP TABLE IF EXISTS carrinho;
DROP TABLE IF EXISTS pedidos;
DROP TABLE IF EXISTS itens_pedido;

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
    secao VARCHAR(50) NOT NULL DEFAULT 'Geral',
    genero VARCHAR(20) NOT NULL DEFAULT 'Unissex',
    foto_url LONGTEXT,
    usuario_id INT NOT NULL,
    ativo BOOLEAN DEFAULT TRUE,
    criadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    atualizadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    INDEX idx_categoria (categoria),
    INDEX idx_secao (secao),
    INDEX idx_genero (genero),
    INDEX idx_ativo (ativo),
    INDEX idx_usuario (usuario_id),
    INDEX idx_tamanho (tamanho)
) ENGINE=InnoDB;

CREATE TABLE favoritos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    usuario_id INT NOT NULL,
    produto_id INT NOT NULL,
    criadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    FOREIGN KEY (produto_id) REFERENCES produtos(id) ON DELETE CASCADE,
    UNIQUE KEY unique_favorito (usuario_id, produto_id),
    INDEX idx_usuario (usuario_id),
    INDEX idx_produto (produto_id)
) ENGINE=InnoDB;

CREATE TABLE carrinho (
    id INT AUTO_INCREMENT PRIMARY KEY,
    usuario_id INT NOT NULL,
    produto_id INT NOT NULL,
    quantidade INT NOT NULL DEFAULT 1,
    criadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    atualizadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    FOREIGN KEY (produto_id) REFERENCES produtos(id) ON DELETE CASCADE,
    UNIQUE KEY unique_item_carrinho (usuario_id, produto_id),
    INDEX idx_usuario (usuario_id),
    INDEX idx_produto (produto_id)
) ENGINE=InnoDB;

CREATE TABLE pedidos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    usuario_id INT NOT NULL,
    nome_completo VARCHAR(200) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telefone VARCHAR(20) NOT NULL,
    endereco VARCHAR(500) NOT NULL,
    numero VARCHAR(20) NOT NULL,
    complemento VARCHAR(200),
    bairro VARCHAR(100) NOT NULL,
    cidade VARCHAR(100) NOT NULL,
    estado VARCHAR(2) NOT NULL,
    cep VARCHAR(10) NOT NULL,
    forma_pagamento VARCHAR(50) NOT NULL DEFAULT 'pendente',
    status VARCHAR(50) NOT NULL DEFAULT 'pendente',
    codigo_rastreio VARCHAR(100),
    total DECIMAL(10,2) NOT NULL,
    criadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    atualizadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    INDEX idx_usuario (usuario_id),
    INDEX idx_status (status),
    INDEX idx_data (criadoEm)
) ENGINE=InnoDB;

CREATE TABLE itens_pedido (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pedido_id INT NOT NULL,
    produto_id INT NOT NULL,
    nome_produto VARCHAR(255) NOT NULL,
    preco_unitario DECIMAL(10,2) NOT NULL,
    quantidade INT NOT NULL,
    tamanho VARCHAR(10) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    foto_url LONGTEXT,
    criadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pedido_id) REFERENCES pedidos(id) ON DELETE CASCADE,
    FOREIGN KEY (produto_id) REFERENCES produtos(id) ON DELETE CASCADE,
    INDEX idx_pedido (pedido_id),
    INDEX idx_produto (produto_id)
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
