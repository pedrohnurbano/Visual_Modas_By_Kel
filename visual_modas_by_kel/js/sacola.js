// Sistema de Gerenciamento da Sacola
class SacolaManager {
    constructor() {
        this.sacola = this.carregarSacola();
        this.cupomAplicado = null;
        this.descontoPercentual = 0;
    }

    // Carregar sacola do armazenamento
    carregarSacola() {
        const dados = sessionStorage.getItem('sacola_visual_modas');
        return dados ? JSON.parse(dados) : [];
    }

    // Salvar sacola no armazenamento
    salvarSacola() {
        sessionStorage.setItem('sacola_visual_modas', JSON.stringify(this.sacola));
    }

    // Adicionar item à sacola
    adicionarItem(produto) {
        const itemExistente = this.sacola.find(
            item => item.id === produto.id && item.tamanho === produto.tamanho
        );

        if (itemExistente) {
            itemExistente.quantidade += 1;
        } else {
            this.sacola.push({
                ...produto,
                quantidade: 1
            });
        }

        this.salvarSacola();
        this.atualizarContador();
        return true;
    }

    // Remover item da sacola
    removerItem(id, tamanho) {
        this.sacola = this.sacola.filter(
            item => !(item.id === id && item.tamanho === tamanho)
        );
        this.salvarSacola();
        this.atualizarInterface();
    }

    // Atualizar quantidade
    atualizarQuantidade(id, tamanho, novaQuantidade) {
        const item = this.sacola.find(
            item => item.id === id && item.tamanho === tamanho
        );

        if (item && novaQuantidade > 0) {
            item.quantidade = novaQuantidade;
            this.salvarSacola();
            this.atualizarInterface();
        } else if (novaQuantidade === 0) {
            this.removerItem(id, tamanho);
        }
    }

    // Calcular subtotal
    calcularSubtotal() {
        return this.sacola.reduce((total, item) => {
            return total + (item.preco * item.quantidade);
        }, 0);
    }

    // Calcular desconto
    calcularDesconto(subtotal) {
        return subtotal * (this.descontoPercentual / 100);
    }

    // Calcular total
    calcularTotal() {
        const subtotal = this.calcularSubtotal();
        const desconto = this.calcularDesconto(subtotal);
        return subtotal - desconto;
    }

    // Formatar moeda
    formatarMoeda(valor) {
        return valor.toLocaleString('pt-BR', {
            style: 'currency',
            currency: 'BRL'
        });
    }

    // Atualizar contador no header
    atualizarContador() {
        const contador = document.getElementById('sacola-contador');
        const totalItens = this.sacola.reduce((total, item) => total + item.quantidade, 0);
        
        if (contador) {
            if (totalItens > 0) {
                contador.textContent = totalItens;
                contador.style.display = 'flex';
            } else {
                contador.style.display = 'none';
            }
        }
    }

    // Renderizar interface da sacola
    renderizarSacola() {
        const container = document.getElementById('sacola-items');
        const resumo = document.getElementById('sacola-resumo');
        const subtitulo = document.getElementById('sacola-subtitulo');
        const vazia = document.getElementById('sacola-vazia');

        if (this.sacola.length === 0) {
            vazia.style.display = 'block';
            resumo.style.display = 'none';
            subtitulo.textContent = '0 itens';
            return;
        }

        vazia.style.display = 'none';
        resumo.style.display = 'block';

        const totalItens = this.sacola.reduce((total, item) => total + item.quantidade, 0);
        subtitulo.textContent = `${totalItens} ${totalItens === 1 ? 'item' : 'itens'}`;

        container.innerHTML = this.sacola.map(item => `
            <div class="sacola-item" data-id="${item.id}" data-tamanho="${item.tamanho}">
                <div class="item-imagem">
                    <img src="${item.imagem}" alt="${item.nome}">
                </div>
                <div class="item-detalhes">
                    <h3 class="item-nome">${item.nome}</h3>
                    <p class="item-info">Tamanho: ${item.tamanho}</p>
                    <p class="item-info">Cor: ${item.cor || 'Padrão'}</p>
                    <p class="item-preco">${this.formatarMoeda(item.preco)}</p>
                    <div class="item-acoes">
                        <div class="item-quantidade">
                            <button class="btn-quantidade" onclick="sacolaManager.diminuirQuantidade('${item.id}', '${item.tamanho}')">−</button>
                            <span class="quantidade-valor">${item.quantidade}</span>
                            <button class="btn-quantidade" onclick="sacolaManager.aumentarQuantidade('${item.id}', '${item.tamanho}')">+</button>
                        </div>
                        <button class="btn-remover" onclick="sacolaManager.removerItem('${item.id}', '${item.tamanho}')">Remover</button>
                    </div>
                </div>
            </div>
        `).join('');

        this.atualizarResumo();
    }

    // Atualizar resumo do pedido
    atualizarResumo() {
        const subtotal = this.calcularSubtotal();
        const desconto = this.calcularDesconto(subtotal);
        const total = this.calcularTotal();

        document.getElementById('resumo-subtotal').textContent = this.formatarMoeda(subtotal);
        document.getElementById('resumo-total').textContent = this.formatarMoeda(total);

        const descontoLinha = document.getElementById('desconto-linha');
        if (desconto > 0) {
            descontoLinha.style.display = 'flex';
            document.getElementById('resumo-desconto').textContent = '- ' + this.formatarMoeda(desconto);
        } else {
            descontoLinha.style.display = 'none';
        }
    }

    // Aumentar quantidade
    aumentarQuantidade(id, tamanho) {
        const item = this.sacola.find(item => item.id === id && item.tamanho === tamanho);
        if (item) {
            this.atualizarQuantidade(id, tamanho, item.quantidade + 1);
        }
    }

    // Diminuir quantidade
    diminuirQuantidade(id, tamanho) {
        const item = this.sacola.find(item => item.id === id && item.tamanho === tamanho);
        if (item) {
            this.atualizarQuantidade(id, tamanho, item.quantidade - 1);
        }
    }

    // Atualizar toda interface
    atualizarInterface() {
        this.renderizarSacola();
        this.atualizarContador();
    }

    // Aplicar cupom
    aplicarCupom(codigo) {
        const cupons = {
            'BEMVINDOVM': 10,
            'PRIMEIRACOMPRA': 15,
            'DESCONTO20': 20,
            'FRETE5': 5
        };

        const codigoUpper = codigo.toUpperCase().trim();
        
        if (cupons[codigoUpper]) {
            this.cupomAplicado = codigoUpper;
            this.descontoPercentual = cupons[codigoUpper];
            this.atualizarResumo();
            alert(`Cupom "${codigoUpper}" aplicado com sucesso! Desconto de ${this.descontoPercentual}%`);
            return true;
        } else {
            alert('Cupom inválido ou expirado.');
            return false;
        }
    }
}

// Inicializar gerenciador
const sacolaManager = new SacolaManager();

// Funções globais
function aplicarCupom() {
    const input = document.getElementById('cupom-input');
    if (input.value) {
        sacolaManager.aplicarCupom(input.value);
        input.value = '';
    }
}

function finalizarCompra() {
    if (sacolaManager.sacola.length === 0) {
        alert('Adicione produtos à sacola antes de finalizar a compra!');
        return;
    }
    
    alert('Redirecionando para o checkout...');
    // Aqui você implementaria a navegação para página de checkout
    // window.location.href = 'checkout.html';
}

// Inicializar ao carregar página
document.addEventListener('DOMContentLoaded', function() {
    sacolaManager.atualizarInterface();
});