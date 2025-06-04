<?php
session_start();

$error = '';
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    // Exemplo de usuário e senha (substitua por consulta ao banco de dados)
    $usuario_correto = 'cliente@exemplo.com';
    $senha_correta = '123456';

    $email = $_POST['email'] ?? '';
    $senha = $_POST['senha'] ?? '';

    if ($email === $usuario_correto && $senha === $senha_correta) {
        $_SESSION['usuario'] = $email;
        header('Location: pagina_principal.php');
        exit;
    } else {
        $error = 'E-mail ou senha inválidos.';
    }
}
?>

<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <title>Login - Visual Modas By Kel</title>
    <link rel="stylesheet" href="estilos.css">
</head>
<body>
    <div class="login-container">
        <h2>Login</h2>
        <?php if ($error): ?>
            <div class="error"><?= htmlspecialchars($error) ?></div>
        <?php endif; ?>
        <form method="post" action="">
            <label for="email">E-mail:</label>
            <input type="email" id="email" name="email" required>

            <label for="senha">Senha:</label>
            <input type="password" id="senha" name="senha" required>

            <button type="submit">Entrar</button>
        </form>
        <a class="register-link" href="pagina_cadastro.php">Criar nova conta</a>
    </div>
</body>
</html>