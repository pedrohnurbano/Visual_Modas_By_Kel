// Gerenciamento da página da sacola usando API

let cartItems = [];

// Carregar itens do carrinho da API
function carregarSacola() {
    buscarCarrinhoAPI(function(itens) {
        cartItems = itens || [];
        renderizarSacola();
    });
}

// Renderizar sacola
function renderizarSacola() {
    const itemsWrapper = document.getElementById('sacolaItems');
    const emptyState = document.getElementById('sacolaEmpty');
    const resumo = document.getElementById('sacolaResumo');
    const totalElement = document.getElementById('sacolaTotal');
    const sacolaContent = document.querySelector('.sacola-content');

    if (cartItems.length === 0) {
        emptyState.style.display = 'flex';
        itemsWrapper.style.display = 'none';
        resumo.style.display = 'none';
        totalElement.textContent = '0 itens';
        if (sacolaContent) {
            sacolaContent.classList.add('sacola-vazia');
        }
    } else {
        emptyState.style.display = 'none';
        itemsWrapper.style.display = 'block';
        resumo.style.display = 'block';

        const totalItens = cartItems.reduce((sum, item) => sum + item.quantidade, 0);
        totalElement.textContent = `${totalItens} ${totalItens === 1 ? 'item' : 'itens'}`;

        itemsWrapper.innerHTML = cartItems.map(item => criarItemSacola(item)).join('');
        atualizarResumo();
        
        if (sacolaContent) {
            sacolaContent.classList.remove('sacola-vazia');
        }
    }
}

// Criar HTML de um item da sacola
function criarItemSacola(item) {
    const produto = item.produto;
    const preco = parseFloat(produto.preco);
    const quantidade = item.quantidade;
    
    // Construir URL da foto
    let fotoUrl = produto.foto_url;
    if (fotoUrl && !fotoUrl.startsWith('http') && !fotoUrl.startsWith('data:')) {
        fotoUrl = fotoUrl.startsWith('/') ? `http://localhost:5000${fotoUrl}` : `http://localhost:5000/${fotoUrl}`;
    }

    return `
        <div class="sacola-item" data-item-id="${produto.id}">
            <div class="item-image">
                <img src="${fotoUrl || 'design/cabide.png'}" alt="${produto.nome}" onerror="this.src='design/cabide.png'">
            </div>
            <div class="item-details">
                <h3 class="item-name">${produto.nome}</h3>
                
                <!-- Informação de Tamanho -->
                <div class="item-size-info">
                    <span>Tamanho: ${produto.tamanho}</span>
                </div>
                
                <p class="item-price">R$ ${preco.toFixed(2).replace('.', ',')}</p>
                
                <div class="item-actions">
                    <div class="quantity-control">
                        <button class="btn-qty" onclick="alterarQuantidade(${produto.id}, ${quantidade}, -1)">−</button>
                        <span class="qty-value">${quantidade}</span>
                        <button class="btn-qty" onclick="alterarQuantidade(${produto.id}, ${quantidade}, 1)">+</button>
                    </div>
                    <button class="btn-remove" onclick="removerItemSacola(${produto.id})">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <polyline points="3 6 5 6 21 6"></polyline>
                            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                        </svg>
                        Remover
                    </button>
                </div>
            </div>
            <div class="item-total">
                <span class="item-total-label">Total</span>
                <span class="item-total-value">R$ ${(preco * quantidade).toFixed(2).replace('.', ',')}</span>
            </div>
        </div>
    `;
}

// Alterar quantidade de um item
function alterarQuantidade(produtoId, quantidadeAtual, mudanca) {
    const novaQuantidade = quantidadeAtual + mudanca;
    
    if (novaQuantidade <= 0) {
        removerItemSacola(produtoId);
        return;
    }

    atualizarQuantidadeCarrinhoAPI(produtoId, novaQuantidade, function(sucesso) {
        if (sucesso) {
            carregarSacola();
        } else {
            showNotification('Erro ao atualizar quantidade');
        }
    });
}

// Remover item da sacola
function removerItemSacola(produtoId) {
    removerDoCarrinhoAPI(produtoId, function(sucesso) {
        if (sucesso) {
            showNotification('Produto removido da sacola');
            carregarSacola();
        } else {
            showNotification('Erro ao remover produto');
        }
    });
}

// Atualizar resumo do pedido
function atualizarResumo() {
    const subtotal = cartItems.reduce((sum, item) => {
        return sum + (parseFloat(item.produto.preco) * item.quantidade);
    }, 0);

    const total = subtotal;

    document.getElementById('resumoSubtotal').textContent = `R$ ${subtotal.toFixed(2).replace('.', ',')}`;
    document.getElementById('resumoTotal').textContent = `R$ ${total.toFixed(2).replace('.', ',')}`;

    // Frete grátis acima de R$ 500
    const freteElement = document.getElementById('resumoFrete');
    if (subtotal >= 500) {
        freteElement.textContent = 'GRÁTIS';
        freteElement.style.color = '#28a745';
        freteElement.style.fontWeight = '700';
    } else {
        freteElement.textContent = 'A calcular';
        freteElement.style.color = '';
        freteElement.style.fontWeight = '';
    }
}

// Finalizar compra
function finalizarCompra() {
    if (cartItems.length === 0) {
        showNotification('Adicione produtos à sacola primeiro!');
        return;
    }

    // Redireciona para o checkout
    window.location.href = '/checkout';
}

// Mostrar notificação
function showNotification(message) {
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
    notification.textContent = message;
    document.body.appendChild(notification);

    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => notification.remove(), 300);
    }, 2000);
}

// Inicialização
$(document).ready(function() {
    carregarSacola();
});
