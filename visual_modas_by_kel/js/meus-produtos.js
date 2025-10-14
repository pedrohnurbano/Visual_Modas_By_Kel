$(document).ready(function() {
    buscarMeusProdutos();
});

function buscarMeusProdutos() {
    $.ajax({
        url: "/api/meus-produtos",
        method: "GET",
        success: function(produtos) {
            exibirMeusProdutos(produtos);
        },
        error: function(xhr) {
            console.error("Erro ao buscar produtos:", xhr);
            $("#lista-meus-produtos").html("<p>Erro ao carregar produtos.</p>");
        }
    });
}

function exibirMeusProdutos(produtos) {
    if (produtos.length === 0) {
        $("#lista-meus-produtos").html("<p>Você ainda não cadastrou nenhum produto.</p>");
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
                <div class="acoes">
                    <button onclick="editarProduto(${produto.id})">Editar</button>
                    <button onclick="deletarProduto(${produto.id})" class="btn-deletar">Deletar</button>
                </div>
            </div>
        `;
    });
    
    html += '</div>';
    $("#lista-meus-produtos").html(html);
}

function deletarProduto(produtoId) {
    Swal.fire({
        title: 'Tem certeza?',
        text: "Deseja deletar este produto? Esta ação não pode ser desfeita!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#370400',
        cancelButtonColor: '#6c757d',
        confirmButtonText: 'Sim, deletar!',
        cancelButtonText: 'Cancelar'
    }).then((result) => {
        if (result.isConfirmed) {
            $.ajax({
                url: `/api/produtos/${produtoId}`,
                method: "DELETE",
                success: function() {
                    Swal.fire({
                        icon: 'success',
                        title: 'Deletado!',
                        text: 'Produto deletado com sucesso!',
                        confirmButtonColor: '#370400'
                    });
                    buscarMeusProdutos();
                },
                error: function(xhr) {
                    console.error("Erro ao deletar:", xhr);
                    Swal.fire({
                        icon: 'error',
                        title: 'Erro',
                        text: 'Erro ao deletar produto.',
                        confirmButtonColor: '#370400'
                    });
                }
            });
        }
    });
}

function editarProduto(produtoId) {
    // Seu amigo implementa a página de edição
    window.location.href = `/editar-produto/${produtoId}`;
}

// MUDAR AS VARIAVEIS
