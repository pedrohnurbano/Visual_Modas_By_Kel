<?php
// cadastro_fornecedor.php

// Conexão com o banco de dados (ajuste conforme necessário)
$host = 'localhost';
$user = 'root';
$pass = '';
$db = 'visual_modas';

$conn = new mysqli($host, $user, $pass, $db);

if ($conn->connect_error) {
    die('Erro de conexão: ' . $conn->connect_error);
}

// Processa o formulário
$mensagem = '';
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $nome = $conn->real_escape_string($_POST['nome']);
    $email = $conn->real_escape_string($_POST['email']);
    $telefone = $conn->real_escape_string($_POST['telefone']);
    $cnpj = $conn->real_escape_string($_POST['cnpj']);
    $endereco = $conn->real_escape_string($_POST['endereco']);

    $sql = "INSERT INTO fornecedores (nome, email, telefone, cnpj, endereco)
            VALUES ('$nome', '$email', '$telefone', '$cnpj', '$endereco')";

    if ($conn->query($sql) === TRUE) {
        $mensagem = "Fornecedor cadastrado com sucesso!";
    } else {
        $mensagem = "Erro ao cadastrar: " . $conn->error;
    }
}
?>

<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <title>Cadastro de Fornecedor</title>
    <link rel="stylesheet" href="estilos.css">
</head>
<body>
    <h2>Cadastro de Fornecedor</h2>
    <?php if ($mensagem): ?>
        <div class="mensagem"><?= $mensagem ?></div>
    <?php endif; ?>
    <form method="post" action="">
        <label for="nome">Nome do Fornecedor:</label>
        <input type="text" name="nome" id="nome" required>

        <label for="email">E-mail:</label>
        <input type="email" name="email" id="email" required>

        <label for="telefone">Telefone:</label>
        <input type="text" name="telefone" id="telefone" required>

        <label for="cnpj">CNPJ:</label>
        <input type="text" name="cnpj" id="cnpj" required>

        <label for="endereco">Endereço:</label>
        <input type="text" name="endereco" id="endereco" required>

        <input type="submit" value="Cadastrar">
    </form>
</body>
</html>