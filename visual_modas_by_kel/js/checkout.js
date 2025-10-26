// Dados do checkout
let checkoutData = {
    cliente: {},
    endereco: {},
    pagamento: {},
    items: []
};

let cartItems = [];

// INICIALIZAÇÃO
$(document).ready(function() {
    // Carregar carrinho da API
    buscarResumoCarrinhoAPI(function(resumo) {
        if (!resumo || !resumo.itens || resumo.itens.length === 0) {
            Swal.fire({
                icon: 'warning',
                title: 'Sacola vazia',
                text: 'Adicione produtos à sacola antes de finalizar a compra.',
                confirmButtonColor: '#370400'
            }).then(() => {
                window.location.href = '/sacola';
            });
            return;
        }

        cartItems = resumo.itens;
        checkoutData.items = cartItems;
        carregarResumo();
        inicializarFormularios();
        aplicarMascaras();
        preencherDadosUsuario();
    });
});

// Preencher dados do usuário logado
function preencherDadosUsuario() {
    const usuarioStr = localStorage.getItem('usuario');
    if (usuarioStr) {
        try {
            const usuario = JSON.parse(usuarioStr);
            document.getElementById('emailCliente').value = usuario.email || '';
            document.getElementById('nomeCliente').value = `${usuario.nome || ''} ${usuario.sobrenome || ''}`.trim();
            document.getElementById('telefoneCliente').value = usuario.telefone || '';
        } catch (e) {
            console.error('Erro ao parsear dados do usuário:', e);
        }
    }
}

// RESUMO DO PEDIDO
function carregarResumo() {
    const subtotal = cartItems.reduce((sum, item) => {
        return sum + (parseFloat(item.produto.preco) * item.quantidade);
    }, 0);
    
    const frete = subtotal >= 500 ? 0 : 30;
    const total = subtotal + frete;

    const resumoHTML = `
        <h3 class="summary-title">Resumo do Pedido</h3>
        <div class="summary-items">
            ${cartItems.map(item => {
                const preco = parseFloat(item.produto.preco);
                let fotoUrl = item.produto.foto_url;
                if (fotoUrl && !fotoUrl.startsWith('http') && !fotoUrl.startsWith('data:')) {
                    fotoUrl = fotoUrl.startsWith('/') ? `http://localhost:5000${fotoUrl}` : `http://localhost:5000/${fotoUrl}`;
                }
                
                return `
                <div class="summary-item">
                    <img src="${fotoUrl || 'design/cabide.png'}" alt="${item.produto.nome}" onerror="this.src='design/cabide.png'">
                    <div class="summary-item-info">
                        <span class="summary-item-name">${item.produto.nome}</span>
                        <span class="summary-item-details">Tam: ${item.produto.tamanho} | Qtd: ${item.quantidade}</span>
                    </div>
                    <span class="summary-item-price">R$ ${(preco * item.quantidade).toFixed(2).replace('.', ',')}</span>
                </div>
                `;
            }).join('')}
        </div>
        <div class="summary-totals">
            <div class="summary-line">
                <span>Subtotal</span>
                <span>R$ ${subtotal.toFixed(2).replace('.', ',')}</span>
            </div>
            <div class="summary-line">
                <span>Frete</span>
                <span>${frete === 0 ? 'GRÁTIS' : 'R$ ' + frete.toFixed(2).replace('.', ',')}</span>
            </div>
            <div class="summary-line total">
                <span>Total</span>
                <span>R$ ${total.toFixed(2).replace('.', ',')}</span>
            </div>
        </div>
    `;

    document.querySelectorAll('[id^="orderSummary"]').forEach(el => {
        el.innerHTML = resumoHTML;
    });
}

// FORMULÁRIOS
function inicializarFormularios() {
    // Formulário de Identificação
    document.getElementById('formIdentificacao').addEventListener('submit', (e) => {
        e.preventDefault();
        
        checkoutData.cliente = {
            email: document.getElementById('emailCliente').value,
            nome: document.getElementById('nomeCliente').value,
            cpf: document.getElementById('cpfCliente').value.replace(/\D/g, ''),
            telefone: document.getElementById('telefoneCliente').value.replace(/\D/g, ''),
            nascimento: document.getElementById('nascimentoCliente').value
        };

        avancarEtapa(3);
    });

    // Formulário de Entrega
    document.getElementById('formEntrega').addEventListener('submit', (e) => {
        e.preventDefault();
        
        checkoutData.endereco = {
            cep: document.getElementById('cep').value.replace(/\D/g, ''),
            endereco: document.getElementById('endereco').value,
            numero: document.getElementById('numero').value,
            complemento: document.getElementById('complemento').value,
            bairro: document.getElementById('bairro').value,
            cidade: document.getElementById('cidade').value,
            estado: document.getElementById('estado').value
        };

        avancarEtapa(4);
    });
}

// BUSCAR CEP
async function buscarCep() {
    const cep = document.getElementById('cep').value.replace(/\D/g, '');
    
    if (cep.length !== 8) {
        Swal.fire({
            icon: 'warning',
            title: 'CEP inválido',
            text: 'Por favor, informe um CEP válido com 8 dígitos.',
            confirmButtonColor: '#370400'
        });
        return;
    }

    try {
        const response = await fetch(`https://viacep.com.br/ws/${cep}/json/`);
        const data = await response.json();

        if (data.erro) {
            Swal.fire({
                icon: 'error',
                title: 'CEP não encontrado',
                text: 'Não foi possível encontrar o CEP informado. Verifique e tente novamente.',
                confirmButtonColor: '#370400'
            });
            return;
        }

        document.getElementById('endereco').value = data.logradouro;
        document.getElementById('bairro').value = data.bairro;
        document.getElementById('cidade').value = data.localidade;
        document.getElementById('estado').value = data.uf;
        document.getElementById('numero').focus();

    } catch (error) {
        Swal.fire({
            icon: 'error',
            title: 'Erro',
            text: 'Erro ao buscar CEP. Verifique sua conexão e tente novamente.',
            confirmButtonColor: '#370400'
        });
    }
}

// MÁSCARAS
function aplicarMascaras() {
    // CPF
    document.getElementById('cpfCliente').addEventListener('input', (e) => {
        let value = e.target.value.replace(/\D/g, '');
        value = value.replace(/(\d{3})(\d)/, '$1.$2');
        value = value.replace(/(\d{3})(\d)/, '$1.$2');
        value = value.replace(/(\d{3})(\d{1,2})$/, '$1-$2');
        e.target.value = value;
    });

    // Telefone
    document.getElementById('telefoneCliente').addEventListener('input', (e) => {
        let value = e.target.value.replace(/\D/g, '');
        value = value.replace(/^(\d{2})(\d)/g, '($1) $2');
        value = value.replace(/(\d)(\d{4})$/, '$1-$2');
        e.target.value = value;
    });

    // CEP
    document.getElementById('cep').addEventListener('input', (e) => {
        let value = e.target.value.replace(/\D/g, '');
        value = value.replace(/^(\d{5})(\d)/, '$1-$2');
        e.target.value = value;
    });
}

// NAVEGAÇÃO ENTRE ETAPAS
function avancarEtapa(etapa) {
    // Esconde todas as seções
    document.querySelectorAll('.checkout-section').forEach(section => {
        section.classList.add('hidden');
    });

    // Remove estados ativos
    document.querySelectorAll('.checkout-step').forEach(step => {
        step.classList.remove('active');
    });

    // Mostra seção atual
    if (etapa === 2) {
        document.getElementById('sectionIdentificacao').classList.remove('hidden');
        document.getElementById('step2').classList.add('active');
    } else if (etapa === 3) {
        document.getElementById('sectionEntrega').classList.remove('hidden');
        document.getElementById('step3').classList.add('active');
        document.getElementById('step2').classList.add('completed');
    } else if (etapa === 4) {
        document.getElementById('sectionPagamento').classList.remove('hidden');
        document.getElementById('step4').classList.add('active');
        document.getElementById('step3').classList.add('completed');
    } else if (etapa === 5) {
        document.getElementById('sectionConfirmacao').classList.remove('hidden');
        document.getElementById('step5').classList.add('active');
        document.getElementById('step4').classList.add('completed');
        mostrarConfirmacao();
    }

    window.scrollTo(0, 0);
}

function voltarEtapa(etapa) {
    avancarEtapa(etapa);
}

// SELEÇÃO DE PAGAMENTO
function selecionarPagamento(tipo) {
    document.querySelectorAll('.payment-option').forEach(opt => {
        opt.classList.remove('selected');
    });
    event.currentTarget.classList.add('selected');
    
    checkoutData.pagamento.tipo = tipo;
}

// FINALIZAR PAGAMENTO
async function finalizarPagamento() {
    if (!checkoutData.pagamento.tipo) {
        Swal.fire({
            icon: 'warning',
            title: 'Forma de pagamento',
            text: 'Por favor, selecione uma forma de pagamento.',
            confirmButtonColor: '#370400'
        });
        return;
    }

    const btnFinalizar = document.getElementById('btnFinalizarPagamento');
    btnFinalizar.disabled = true;
    btnFinalizar.textContent = 'Processando...';

    // Preparar dados do pedido
    const dadosPedido = {
        nomeCompleto: checkoutData.cliente.nome,
        email: checkoutData.cliente.email,
        telefone: checkoutData.cliente.telefone,
        endereco: checkoutData.endereco.endereco,
        numero: checkoutData.endereco.numero,
        complemento: checkoutData.endereco.complemento || '',
        bairro: checkoutData.endereco.bairro,
        cidade: checkoutData.endereco.cidade,
        estado: checkoutData.endereco.estado,
        cep: checkoutData.endereco.cep,
        formaPagamento: checkoutData.pagamento.tipo
    };

    // Criar pedido via API
    criarPedidoAPI(dadosPedido, function(sucesso, response) {
        if (sucesso) {
            checkoutData.numeroPedido = response.pedidoId;
            avancarEtapa(5);
        } else {
            Swal.fire({
                icon: 'error',
                title: 'Erro ao criar pedido',
                text: response || 'Não foi possível criar o pedido. Tente novamente.',
                confirmButtonColor: '#370400'
            });
            btnFinalizar.disabled = false;
            btnFinalizar.textContent = 'Finalizar Compra';
        }
    });
}

// MOSTRAR CONFIRMAÇÃO
function mostrarConfirmacao() {
    const numeroPedido = checkoutData.numeroPedido || Math.floor(Math.random() * 1000000);
    document.getElementById('numeroPedido').textContent = `#${numeroPedido}`;

    const subtotal = cartItems.reduce((sum, item) => {
        return sum + (parseFloat(item.produto.preco) * item.quantidade);
    }, 0);
    const frete = subtotal >= 500 ? 0 : 30;
    const total = subtotal + frete;

    const detalhesHTML = `
        <div class="confirmation-section">
            <h3>Dados do Cliente</h3>
            <p><strong>Nome:</strong> ${checkoutData.cliente.nome}</p>
            <p><strong>E-mail:</strong> ${checkoutData.cliente.email}</p>
            <p><strong>CPF:</strong> ${formatarCPF(checkoutData.cliente.cpf)}</p>
        </div>

        <div class="confirmation-section">
            <h3>Endereço de Entrega</h3>
            <p>${checkoutData.endereco.endereco}, ${checkoutData.endereco.numero}</p>
            ${checkoutData.endereco.complemento ? `<p>${checkoutData.endereco.complemento}</p>` : ''}
            <p>${checkoutData.endereco.bairro} - ${checkoutData.endereco.cidade}/${checkoutData.endereco.estado}</p>
            <p>CEP: ${formatarCEP(checkoutData.endereco.cep)}</p>
        </div>

        <div class="confirmation-section">
            <h3>Itens do Pedido</h3>
            ${cartItems.map(item => {
                const preco = parseFloat(item.produto.preco);
                return `<p>${item.produto.nome} - Tamanho ${item.produto.tamanho} - Qtd: ${item.quantidade} - R$ ${(preco * item.quantidade).toFixed(2).replace('.', ',')}</p>`;
            }).join('')}
        </div>

        <div class="confirmation-section">
            <h3>Total do Pedido</h3>
            <p class="confirmation-total">R$ ${total.toFixed(2).replace('.', ',')}</p>
        </div>

        <div class="confirmation-section">
            <p style="text-align: center; color: #666; margin-top: 20px;">
                ${checkoutData.pagamento.tipo === 'pix' ? 'Aguardando confirmação do pagamento via PIX...' : 'Aguardando confirmação do pagamento...'}
            </p>
            <p style="text-align: center; color: #666;">
                Você receberá um e-mail com os detalhes do seu pedido.
            </p>
        </div>
    `;

    document.getElementById('confirmationDetails').innerHTML = detalhesHTML;
}

function formatarCPF(cpf) {
    return cpf.replace(/(\d{3})(\d{3})(\d{3})(\d{2})/, '$1.$2.$3-$4');
}

function formatarCEP(cep) {
    return cep.replace(/(\d{5})(\d{3})/, '$1-$2');
}
