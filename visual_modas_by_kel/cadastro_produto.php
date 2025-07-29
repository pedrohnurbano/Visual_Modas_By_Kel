<?php
$host = 'localhost';
$user = 'root';
$pass = '';
$db = 'loja_roupas';

$conn = new mysqli($host, $user, $pass, $db);

if ($conn->connect_error) {
    die('Erro de conexão: ' . $conn->connect_error);
}

$mensagem = '';

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $nome       = $conn->real_escape_string($_POST['nome']     );
    $descricao  = $conn->real_escape_string($_POST['descricao']);
    $preco      = floatval($_POST['preco']                      );
    $tamanho    = $conn->real_escape_string($_POST['tamanho']  );
    $quantidade = intval($_POST['quantidade']                   );

    $sql = "INSERT INTO produtos (nome, descricao, preco, tamanho, quantidade)
            VALUES ('$nome', '$descricao', $preco, '$tamanho', $quantidade)";

    if ($conn->query($sql) === TRUE) {
        $mensagem = "Produto cadastrado com sucesso!";
    } else {
        $mensagem = "Erro ao cadastrar produto: " . $conn->error;
    }
}
?>

<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <title> Cadastro de Produto </title>
    <link rel="stylesheet" href="estilos.css">
</head>
<body>
    <h2>Cadastro de Produto</h2>
    <?php if ($mensagem): ?>
        <div class="mensagem"><?= $mensagem ?></div>
    <?php endif; ?>
    <form method="post">
        <label for="nome"> Nome do Produto: </label>
        <input type="text" name="nome" id="nome" required>

        <label for="descricao"> Descrição: </label>
        <textarea name="descricao" id="descricao" required></textarea>

        <label for="preco"> Preço (R$): </label>
        <input type="number" step="0.01" name="preco" id="preco" required>

        <label for="tamanho"> Tamanho: </label>
        <select name="tamanho" id="tamanho" required>
            <option value = ""   > Selecione </option>
            <option value = "PP" > PP        </option>
            <option value = "P"  > P         </option>
            <option value = "M"  > M         </option>
            <option value = "G"  > G         </option>
            <option value = "GG" > GG        </option>
        </select>

        <label for="quantidade">Quantidade em Estoque:</label>
        <input type="number" name="quantidade" id="quantidade" min="1" required>

        <button type="submit" style="margin-top: 15px;">Cadastrar Produto</button>
    </form>
</body>
</html>