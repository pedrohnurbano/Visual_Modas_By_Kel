$(document).ready(function() {
    buscarProdutos();
    
    $("#busca-produto").on("input", function() {
        buscarProdutos($(this).val());
    });
    
    $("#filtro-categoria").on("change", function() {
        // Implementar filtro por categoria
        buscarProdutos();
    });
});

function buscarProdutos(filtro = "") {
    $.ajax({
        url: `/api/produtos?filtro=${filtro}`,
        method: "GET",
        success: function(produtos) {
            exibirProdutos(produtos);
        },
        error: function(xhr) {
            console.error("Erro ao buscar produtos:", xhr);
            $("#lista-produtos").html("<p>Erro ao carregar produtos.</p>");
        }
    });
}

function exibirProdutos(produtos) {
    if (produtos.length === 0) {
        $("#lista-produtos").html("<p>Nenhum produto encontrado.</p>");
        return;
    }

    let html = '<div class="produtos-grid">';
    
    produtos.forEach(function(produto) {
        html += `
            <div class="produto-card">
                <img src="http://localhost:5000${produto.foto_url}" alt="${produto.nome}">
                <h3>${produto.nome}</h3>
                <p>${produto.descricao}</p>
                <p class="preco">R$ ${produto.preco.toFixed(2)}</p>
                <p class="tamanho">Tamanho: ${produto.tamanho}</p>
                <p class="categoria">${produto.categoria}</p>
                <button onclick="verDetalhes(${produto.id})">Ver Detalhes</button>
            </div>
        `;
    });
    
    html += '</div>';
    $("#lista-produtos").html(html);
}

function verDetalhes(produtoId) {
    window.location.href = `/produto/${produtoId}`;
}

// MUDAR AS VARIAVEIS
