<?php
session_start();

// Inicializa o carrinho se não existir
if (!isset($_SESSION['carrinho'])) {
    $_SESSION['carrinho'] = [];
}

// Adiciona produto ao carrinho
if (isset($_POST['adicionar'])) {
    $id = $_POST['id'];
    $nome = $_POST['nome'];
    $preco = floatval($_POST['preco']);
    $quantidade = intval($_POST['quantidade']);

    // Se já existe, soma a quantidade
    if (isset($_SESSION['carrinho'][$id])) {
        $_SESSION['carrinho'][$id]['quantidade'] += $quantidade;
    } else {
        $_SESSION['carrinho'][$id] = [
            'nome' => $nome,
            'preco' => $preco,
            'quantidade' => $quantidade
        ];
    }
}

// Remove produto do carrinho
if (isset($_GET['remover'])) {
    $id = $_GET['remover'];
    unset($_SESSION['carrinho'][$id]);
}

// Lista de produtos de exemplo
$produtos = [
    1 => ['nome' => 'Camiseta', 'preco' => 49.90],
    2 => ['nome' => 'Calça Jeans', 'preco' => 129.90],
    3 => ['nome' => 'Vestido', 'preco' => 89.90]
];
?>
<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <title>Carrinho de Compras</title>
    <link rel="stylesheet" href="estilos.css">
</head>
<body>
    <h1>Produtos</h1>
    <form method="post">
        <select name="id">
            <?php foreach ($produtos as $id => $produto): ?>
                <option value="<?= $id ?>"><?= $produto['nome'] ?> - R$ <?= number_format($produto['preco'], 2, ',', '.') ?></option>
            <?php endforeach; ?>
        </select>
        <input type="hidden" name="nome" value="<?= $produtos[1]['nome'] ?>">
        <input type="hidden" name="preco" value="<?= $produtos[1]['preco'] ?>">
        <input type="number" name="quantidade" value="1" min="1" style="width: 60px;">
        <button type="submit" name="adicionar">Adicionar ao Carrinho</button>
    </form>

    <h2>Carrinho de Compras</h2>
    <?php if (empty($_SESSION['carrinho'])): ?>
        <p>Seu carrinho está vazio.</p>
    <?php else: ?>
        <table>
            <tr>
                <th>Produto</th>
                <th>Preço</th>
                <th>Quantidade</th>
                <th>Total</th>
                <th>Ação</th>
            </tr>
            <?php
            $total = 0;
            foreach ($_SESSION['carrinho'] as $id => $item):
                $subtotal = $item['preco'] * $item['quantidade'];
                $total += $subtotal;
            ?>
            <tr>
                <td><?= htmlspecialchars($item['nome']) ?></td>
                <td>R$ <?= number_format($item['preco'], 2, ',', '.') ?></td>
                <td><?= $item['quantidade'] ?></td>
                <td>R$ <?= number_format($subtotal, 2, ',', '.') ?></td>
                <td><a href="?remover=<?= $id ?>">Remover</a></td>
            </tr>
            <?php endforeach; ?>
            <tr>
                <th colspan="3">Total</th>
                <th colspan="2">R$ <?= number_format($total, 2, ',', '.') ?></th>
            </tr>
        </table>
    <?php endif; ?>
</body>
</html>