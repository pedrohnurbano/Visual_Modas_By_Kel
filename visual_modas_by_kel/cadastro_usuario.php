<?php
$host = 'localhost';
$user = 'root';
$pass = '';
$db = 'visual_modas';

$conn = new mysqli($host, $user, $pass, $db);

if ($conn->connect_error) {
    die('Erro de conexão: ' . $conn->connect_error);
}

$mensagem = '';

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $nome = $conn->real_escape_string($_POST['nome']);
    $email = $conn->real_escape_string($_POST['email']);
    $senha = password_hash($_POST['senha'], PASSWORD_DEFAULT);

    $sql = "INSERT INTO usuarios (nome, email, senha) VALUES ('$nome', '$email', '$senha')";

    if ($conn->query($sql) === TRUE) {
        $mensagem = "Usuário cadastrado com sucesso!";
    } else {
        $mensagem = "Erro ao cadastrar: " . $conn->error;
    }
}
?>

<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <title>Cadastro de Usuário - Visual Modas</title>
    <link rel="stylesheet" href="estilos.css">
</head>
<body>
    <h2>Cadastro de Usuário</h2>
    <?php if ($mensagem) echo "<div class='mensagem'>$mensagem</div>"; ?>
    <form method="post" action="">
        <label>Nome:</label>
        <input type="text" name="nome" required>
        <label>Email:</label>
        <input type="email" name="email" required>
        <label>Senha:</label>
        <input type="password" name="senha" required>
        <button type="submit">Cadastrar</button>
    </form>
</body>
</html>