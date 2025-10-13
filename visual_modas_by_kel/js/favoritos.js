// Gerenciamento de Favoritos com API

let favoritosCache = [];
let produtosCache = [];

// Carregar favoritos do usuário
function carregarFavoritosDaAPI() {
    console.log('Carregando favoritos da API...');
    
    $.ajax({
        url: "/api/favoritos",
        method: "GET",
        success: function(favoritos) {
            console.log('Favoritos recebidos:', favoritos);
            // Garantir que favoritos seja um array
            favoritosCache = favoritos || [];
            console.log('Total de favoritos:', favoritosCache.length);
            exibirFavoritos(favoritosCache);
            atualizarContadorFavoritos();
        },
        error: function(xhr) {
            console.error("Erro ao buscar favoritos:", xhr);
            console.error("Status:", xhr.status);
            console.error("Response:", xhr.responseText);
            
            // Se não estiver autenticado, mostrar tela vazia
            if (xhr.status === 401 || xhr.status === 403) {
                console.log('Não autenticado - mostrando mensagem');
                exibirMensagemNaoAutenticado();
            } else {
                // Qualquer outro erro, mostrar vazio
                console.log('Outro erro - mostrando tela vazia');
                exibirFavoritos([]);
            }
        }
    });
}

// Buscar apenas IDs dos favoritos (mais leve)
function buscarIDsFavoritos(callback) {
    $.ajax({
        url: "/api/favoritos/ids",
        method: "GET",
        success: function(ids) {
            // Garantir que ids seja um array
            if (callback) callback(ids || []);
        },
        error: function(xhr) {
            console.error("Erro ao buscar IDs dos favoritos:", xhr);
            if (callback) callback([]);
        }
    });
}

// Toggle favorito (adiciona ou remove)
function toggleFavoritoAPI(produtoId, callback) {
    console.log('Tentando alternar favorito do produto:', produtoId);
    
    $.ajax({
        url: `/api/favoritos/toggle/${produtoId}`,
        method: "POST",
        success: function(response) {
            console.log('Resposta do servidor:', response);
            mostrarNotificacao(response.mensagem);
            if (callback) callback(response.isFavorito);
            // Atualizar contador
            buscarIDsFavoritos(function(ids) {
                console.log('IDs dos favoritos atualizados:', ids);
                atualizarContadorFavoritos(ids.length);
            });
        },
        error: function(xhr) {
            console.error('Erro ao alternar favorito:', xhr);
            console.error('Status:', xhr.status);
            console.error('Response:', xhr.responseText);
            
            if (xhr.status === 401 || xhr.status === 403) {
                console.log('Usuário não autenticado - mostrando modal de login');
                mostrarModalLoginFavoritos();
            } else {
                let mensagem = "Erro ao processar favorito";
                try {
                    const resposta = JSON.parse(xhr.responseText);
                    mensagem = resposta.erro || resposta.message || mensagem;
                } catch (e) {
                    console.error('Erro ao fazer parse da resposta:', e);
                }
                mostrarNotificacao(mensagem);
            }
        }
    });
}

// Exibir favoritos na página
function exibirFavoritos(favoritos) {
    const grid = document.getElementById('favoritosGrid');
    const empty = document.getElementById('favoritosEmpty');
    const total = document.getElementById('favoritosTotal');

    if (!grid || !empty || !total) return;

    // Garantir que favoritos seja um array
    if (!favoritos || !Array.isArray(favoritos) || favoritos.length === 0) {
        empty.style.display = 'flex';
        grid.style.display = 'none';
        total.textContent = '0 itens salvos';
        return;
    }

    empty.style.display = 'none';
    grid.style.display = 'grid';
    total.textContent = `${favoritos.length} ${favoritos.length === 1 ? 'item salvo' : 'itens salvos'}`;

    grid.innerHTML = favoritos.map(fav => criarCardFavorito(fav)).join('');
}

// Criar card de produto favorito
function criarCardFavorito(favorito) {
    const produto = favorito.produto;
    const preco = parseFloat(produto.preco);
    const installments = Math.floor(preco / 100); // Exemplo de parcelas
    const installmentValue = (preco / (installments || 1)).toFixed(2);

    return `
        <div class="produto-card" data-product-id="${produto.id}">
            <button class="produto-favorito favorito-active" 
                    title="Remover dos favoritos" 
                    onclick="removerFavoritoClick(${produto.id})">
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
                    <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41 0.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
                        stroke="#370400" stroke-width="1.5" fill="#370400"/>
                </svg>
            </button>
            
            <div class="produto-img-container">
                <img src="http://localhost:5000${produto.foto_url}" 
                     alt="${produto.nome}" 
                     class="produto-img img-primaria"
                     onerror="this.src='design/cabide.png'">
            </div>
            <div class="produto-info">
                <span class="produto-nome">${produto.nome}</span>
                <span class="produto-preco">R$ ${preco.toFixed(2).replace('.', ',')}</span>
                <span class="produto-parcela">Tamanho: ${produto.tamanho}</span>
                <span class="produto-categoria" style="font-size: 12px; color: #666; margin-top: 4px; display: block;">
                    ${produto.categoria}
                </span>
            </div>
        </div>
    `;
}

// Remover favorito ao clicar
function removerFavoritoClick(produtoId) {
    toggleFavoritoAPI(produtoId, function(isFavorito) {
        if (!isFavorito) {
            // Recarregar lista de favoritos
            carregarFavoritosDaAPI();
        }
    });
}

// Atualizar contador de favoritos no header
function atualizarContadorFavoritos(count) {
    const contador = document.getElementById('favoritosCount');
    if (!contador) return;

    if (count === undefined) {
        // Buscar contagem da API
        buscarIDsFavoritos(function(ids) {
            if (ids.length > 0) {
                contador.textContent = ids.length;
                contador.style.display = 'inline';
            } else {
                contador.style.display = 'none';
            }
        });
    } else {
        if (count > 0) {
            contador.textContent = count;
            contador.style.display = 'inline';
        } else {
            contador.style.display = 'none';
        }
    }
}

// Mostrar mensagem quando não autenticado
function exibirMensagemNaoAutenticado() {
    const grid = document.getElementById('favoritosGrid');
    const empty = document.getElementById('favoritosEmpty');
    const total = document.getElementById('favoritosTotal');

    if (grid && empty && total) {
        empty.style.display = 'flex';
        grid.style.display = 'none';
        total.textContent = '0 itens salvos';
    }
}

// Mostrar modal de login para favoritos
function mostrarModalLoginFavoritos() {
    const modal = document.getElementById('favoriteLoginModal');
    if (modal) {
        modal.classList.add('active');
        document.body.style.overflow = 'hidden';
    }
}

// Fechar modal de login
function fecharModalLoginFavoritos() {
    const modal = document.getElementById('favoriteLoginModal');
    if (modal) {
        modal.classList.remove('active');
        document.body.style.overflow = '';
    }
}

// Mostrar notificação
function mostrarNotificacao(mensagem) {
    const notification = document.createElement('div');
    notification.style.cssText = `
        position: fixed;
        top: 80px;
        right: 20px;
        background: #370400;
        color: white;
        padding: 15px 25px;
        border-radius: 4px;
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        z-index: 10000;
        animation: slideIn 0.3s ease-out;
    `;
    notification.textContent = mensagem;
    document.body.appendChild(notification);

    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => notification.remove(), 300);
    }, 2000);
}

// Inicialização
document.addEventListener('DOMContentLoaded', function() {
    // Se estiver na página de favoritos, carregar favoritos
    if (document.getElementById('favoritosGrid')) {
        carregarFavoritosDaAPI();
    } else {
        // Em outras páginas, apenas atualizar contador
        atualizarContadorFavoritos();
    }
});

