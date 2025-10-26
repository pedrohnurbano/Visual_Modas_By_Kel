// Funções para gerenciar o carrinho através da API

// Adicionar produto ao carrinho via API
function adicionarAoCarrinhoAPI(produtoId, quantidade = 1, callback) {
    $.ajax({
        url: "/api/carrinho",
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        data: JSON.stringify({
            produtoId: produtoId,
            quantidade: quantidade
        }),
        contentType: "application/json",
        success: function(response) {
            console.log("Produto adicionado ao carrinho:", response);
            atualizarContadorCarrinho();
            if (callback) callback(true);
        },
        error: function(xhr) {
            console.error("Erro ao adicionar produto ao carrinho:", xhr);
            if (xhr.status === 401) {
                // Não autenticado, redirecionar para login
                window.location.href = '/login';
            } else {
                alert("Erro ao adicionar produto ao carrinho. Tente novamente.");
            }
            if (callback) callback(false);
        }
    });
}

// Atualizar quantidade de um item no carrinho
function atualizarQuantidadeCarrinhoAPI(produtoId, quantidade, callback) {
    $.ajax({
        url: `/api/carrinho/${produtoId}`,
        method: "PUT",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        data: JSON.stringify({
            quantidade: quantidade
        }),
        contentType: "application/json",
        success: function(response) {
            console.log("Quantidade atualizada:", response);
            if (callback) callback(true);
        },
        error: function(xhr) {
            console.error("Erro ao atualizar quantidade:", xhr);
            if (callback) callback(false);
        }
    });
}

// Remover produto do carrinho
function removerDoCarrinhoAPI(produtoId, callback) {
    $.ajax({
        url: `/api/carrinho/${produtoId}`,
        method: "DELETE",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        success: function(response) {
            console.log("Produto removido do carrinho:", response);
            atualizarContadorCarrinho();
            if (callback) callback(true);
        },
        error: function(xhr) {
            console.error("Erro ao remover produto do carrinho:", xhr);
            if (callback) callback(false);
        }
    });
}

// Buscar itens do carrinho
function buscarCarrinhoAPI(callback) {
    $.ajax({
        url: "/api/carrinho",
        method: "GET",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        success: function(itens) {
            console.log("Itens do carrinho carregados:", itens);
            if (callback) callback(itens || []);
        },
        error: function(xhr) {
            console.error("Erro ao buscar carrinho:", xhr);
            if (xhr.status === 401) {
                // Não autenticado
                if (callback) callback([]);
            } else {
                if (callback) callback([]);
            }
        }
    });
}

// Buscar resumo do carrinho com totais
function buscarResumoCarrinhoAPI(callback) {
    $.ajax({
        url: "/api/carrinho/resumo",
        method: "GET",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        success: function(resumo) {
            console.log("Resumo do carrinho:", resumo);
            if (callback) callback(resumo);
        },
        error: function(xhr) {
            console.error("Erro ao buscar resumo do carrinho:", xhr);
            if (callback) callback({itens: [], quantidadeTotal: 0, valorTotal: 0});
        }
    });
}

// Limpar carrinho
function limparCarrinhoAPI(callback) {
    $.ajax({
        url: "/api/carrinho/limpar",
        method: "DELETE",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        success: function(response) {
            console.log("Carrinho limpo:", response);
            atualizarContadorCarrinho();
            if (callback) callback(true);
        },
        error: function(xhr) {
            console.error("Erro ao limpar carrinho:", xhr);
            if (callback) callback(false);
        }
    });
}

// Atualizar contador de itens no carrinho (header)
function atualizarContadorCarrinho() {
    $.ajax({
        url: "/api/carrinho/contar",
        method: "GET",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        success: function(response) {
            const count = document.getElementById('sacolaCount');
            if (count) {
                if (response.total > 0) {
                    count.textContent = response.total;
                    count.style.display = 'inline';
                } else {
                    count.style.display = 'none';
                }
            }
        },
        error: function(xhr) {
            console.error("Erro ao contar itens do carrinho:", xhr);
            const count = document.getElementById('sacolaCount');
            if (count) {
                count.style.display = 'none';
            }
        }
    });
}

// Inicializar contador ao carregar a página
$(document).ready(function() {
    // Verificar se o usuário está logado antes de carregar o contador
    const token = localStorage.getItem("token");
    if (token) {
        atualizarContadorCarrinho();
    }
});

