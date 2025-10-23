// CONFIGURAÇÃO DO MERCADO PAGO
const MP_PUBLIC_KEY = 'TEST-18a74325-f90e-4f35-8056-4af663c6482f'; // Substitua pela sua chave pública
let mp;
let cardPaymentBrickController;

// Dados do checkout
let checkoutData = {
    cliente: {},
    endereco: {},
    pagamento: {},
    items: []
};

let cart = JSON.parse(localStorage.getItem('cart')) || [];
const allProducts = [
    // Cole aqui o mesmo array de produtos das outras páginas
    {
        id: 1,
        name: "Vestido Midi Floral",
        price: 899.90,
        installments: 7,
        image1: "design/ex-roupa1.png",
        sizes: ["PP", "P", "M", "G", "GG", "XG", "XGG"]
    },
    // ... resto dos produtos
];

// INICIALIZAÇÃO
document.addEventListener('DOMContentLoaded', () => {
    if (cart.length === 0) {
        window.location.href = '/sacola';
        return;
    }

    checkoutData.items = cart;
    carregarResumo();
    inicializarFormularios();
    aplicarMascaras();
});

// RESUMO DO PEDIDO
function carregarResumo() {
    const subtotal = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
    const desconto = subtotal >= 500 ? subtotal * 0.05 : 0;
    const frete = subtotal >= 500 ? 0 : 30;
    const total = subtotal - desconto + frete;

    const resumoHTML = `
        <h3 class="summary-title">Resumo do Pedido</h3>
        <div class="summary-items">
            ${cart.map(item => `
                <div class="summary-item">
                    <img src="${item.image}" alt="${item.name}">
                    <div class="summary-item-info">
                        <span class="summary-item-name">${item.name}</span>
                        <span class="summary-item-details">Tam: ${item.size} | Qtd: ${item.quantity}</span>
                    </div>
                    <span class="summary-item-price">R$ ${(item.price * item.quantity).toFixed(2).replace('.', ',')}</span>
                </div>
            `).join('')}
        </div>
        <div class="summary-totals">
            <div class="summary-line">
                <span>Subtotal</span>
                <span>R$ ${subtotal.toFixed(2).replace('.', ',')}</span>
            </div>
            ${desconto > 0 ? `
            <div class="summary-line discount">
                <span>Desconto</span>
                <span>- R$ ${desconto.toFixed(2).replace('.', ',')}</span>
            </div>
            ` : ''}
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
        inicializarMercadoPago();
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

    if (tipo === 'mercadopago') {
        document.getElementById('mercadoPagoContainer').classList.remove('hidden');
    } else {
        document.getElementById('mercadoPagoContainer').classList.add('hidden');
    }
}

// INICIALIZAR MERCADO PAGO
function inicializarMercadoPago() {
    if (!mp) {
        mp = new MercadoPago(MP_PUBLIC_KEY, {
            locale: 'pt-BR'
        });
    }
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

    try {
        if (checkoutData.pagamento.tipo === 'mercadopago') {
            await processarPagamentoMercadoPago();
        } else if (checkoutData.pagamento.tipo === 'pix') {
            await processarPagamentoPix();
        }
    } catch (error) {
        console.error('Erro ao processar pagamento:', error);
        Swal.fire({
            icon: 'error',
            title: 'Erro no pagamento',
            text: 'Erro ao processar pagamento. Por favor, tente novamente.',
            confirmButtonColor: '#370400'
        });
        btnFinalizar.disabled = false;
        btnFinalizar.textContent = 'Finalizar Compra';
    }
}

// PROCESSAR PAGAMENTO MERCADO PAGO
async function processarPagamentoMercadoPago() {
    // Aqui você fará a chamada para o seu backend
    const pedido = {
        cliente: checkoutData.cliente,
        endereco: checkoutData.endereco,
        items: checkoutData.items,
        total: calcularTotal()
    };

    // Exemplo de chamada ao backend
    const response = await fetch('/api/criar-pagamento-mp', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(pedido)
    });

    const data = await response.json();
    
    if (data.success) {
        checkoutData.pagamento.transactionId = data.transactionId;
        limparCarrinho();
        avancarEtapa(5);
    } else {
        throw new Error(data.message);
    }
}

// PROCESSAR PAGAMENTO PIX
async function processarPagamentoPix() {
    const pedido = {
        cliente: checkoutData.cliente,
        endereco: checkoutData.endereco,
        items: checkoutData.items,
        total: calcularTotal() * 0.95 // 5% de desconto
    };

    const response = await fetch('/api/criar-pagamento-pix', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(pedido)
    });

    const data = await response.json();
    
    if (data.success) {
        checkoutData.pagamento.qrCode = data.qrCode;
        checkoutData.pagamento.qrCodeBase64 = data.qrCodeBase64;
        limparCarrinho();
        avancarEtapa(5);
    } else {
        throw new Error(data.message);
    }
}

// FUNÇÕES AUXILIARES
function calcularTotal() {
    const subtotal = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
    const frete = subtotal >= 500 ? 0 : 30;
    return subtotal + frete;
}

function limparCarrinho() {
    localStorage.removeItem('cart');
    cart = [];
}

function mostrarConfirmacao() {
    const numeroPedido = Math.floor(Math.random() * 1000000).toString().padStart(6, '0');
    document.getElementById('numeroPedido').textContent = `#${numeroPedido}`;

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
            <p>${checkoutData.endereco.bairro} - ${checkoutData.endereco.cidade}/${checkoutData.endereco.estado}</p>
            <p>CEP: ${formatarCEP(checkoutData.endereco.cep)}</p>
        </div>

        <div class="confirmation-section">
            <h3>Itens do Pedido</h3>
            ${cart.map(item => `
                <p>${item.name} - Tamanho ${item.size} - Qtd: ${item.quantity} - R$ ${(item.price * item.quantity).toFixed(2).replace('.', ',')}</p>
            `).join('')}
        </div>

        <div class="confirmation-section">
            <h3>Total do Pedido</h3>
            <p class="confirmation-total">R$ ${calcularTotal().toFixed(2).replace('.', ',')}</p>
        </div>

        ${checkoutData.pagamento.tipo === 'pix' && checkoutData.pagamento.qrCodeBase64 ? `
        <div class="confirmation-section pix-section">
            <h3>Pagamento via PIX</h3>
            <img src="${checkoutData.pagamento.qrCodeBase64}" alt="QR Code PIX" style="max-width: 300px; margin: 20px auto; display: block;">
            <p style="text-align: center; word-break: break-all; padding: 10px; background: #f5f5f5; border-radius: 4px;">
                ${checkoutData.pagamento.qrCode}
            </p>
            <p style="text-align: center; color: #666; margin-top: 10px;">
                Escaneie o QR Code ou copie o código acima para realizar o pagamento
            </p>
        </div>
        ` : ''}
    `;

    document.getElementById('confirmationDetails').innerHTML = detalhesHTML;
}

function formatarCPF(cpf) {
    return cpf.replace(/(\d{3})(\d{3})(\d{3})(\d{2})/, '$1.$2.$3-$4');
}

function formatarCEP(cep) {
    return cep.replace(/(\d{5})(\d{3})/, '$1-$2');
}