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
    // Se estiver na página de confirmação, recuperar dados do localStorage
    if (window.location && window.location.pathname === '/checkout/confirmacao') {
        const savedCheckoutData = localStorage.getItem('checkoutData');
        if (savedCheckoutData) {
            checkoutData = JSON.parse(savedCheckoutData);
            cartItems = checkoutData.items || [];
            
            // Confirmar pagamento e criar pedido
            confirmarPagamentoECriarPedido();
            return;
        }
    }

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
        // Selecionar PIX automaticamente
        checkoutData.pagamento.tipo = 'pix';
    } else if (etapa === 5) {
        document.getElementById('sectionConfirmacao').classList.remove('hidden');
        document.getElementById('step5').classList.add('active');
        document.getElementById('step4').classList.add('completed');
        // Atualiza a rota do navegador para a rota dedicada de confirmação
        if (window.location && window.location.pathname !== '/checkout/confirmacao') {
            window.history.pushState({}, '', '/checkout/confirmacao');
        }
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

    try {
        // Preparar dados para o AbacatePay
        const dadosAbacatePay = {
            cliente: checkoutData.cliente,
            endereco: checkoutData.endereco
        };

        // Chamar API para criar cobrança no AbacatePay
        const response = await fetch('/api/abacatepay/criar-cobranca', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: JSON.stringify(dadosAbacatePay)
        });

        const data = await response.json();
        
        console.log('=== Resposta da API ===');
        console.log('Status:', response.status);
        console.log('Data:', data);
        console.log('paymentUrl:', data.paymentUrl);
        console.log('=====================');

        if (response.ok && data.success && data.paymentUrl) {
            // Verificar se há itens no carrinho
            if (!cartItems || cartItems.length === 0) {
                Swal.fire({
                    icon: 'warning',
                    title: 'Carrinho vazio',
                    text: 'Seu carrinho está vazio. Adicione produtos antes de finalizar a compra.',
                    confirmButtonColor: '#370400'
                }).then(() => {
                    window.location.href = '/sacola';
                });
                btnFinalizar.disabled = false;
                btnFinalizar.textContent = 'Finalizar Compra';
                return;
            }

            // Criar o pedido antes de ir para o AbacatePay
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
                formaPagamento: 'PIX'
            };

            // Criar pedido
            const pedidoResponse = await fetch('/api/pedidos', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify(dadosPedido)
            });

            const pedidoData = await pedidoResponse.json();

            if (pedidoResponse.ok) {
                // Salvar dados para a confirmação
                checkoutData.numeroPedido = pedidoData.pedidoId;
                localStorage.setItem('checkoutData', JSON.stringify(checkoutData));
                localStorage.setItem('abacateExternalId', data.externalId);
                
                console.log('Pedido criado! ID:', pedidoData.pedidoId);
                console.log('Redirecionando para:', data.paymentUrl);
                
                // Redirecionar para o gateway de pagamento do AbacatePay
                window.location.href = data.paymentUrl;
            } else {
                // Se o erro for de carrinho vazio, dar mensagem específica
                if (pedidoData.erro && pedidoData.erro.includes('carrinho vazio')) {
                    Swal.fire({
                        icon: 'warning',
                        title: 'Carrinho vazio',
                        html: 'Seu carrinho está vazio ou já foi processado.<br><br>Por favor, adicione produtos novamente antes de finalizar a compra.',
                        confirmButtonColor: '#370400'
                    }).then(() => {
                        window.location.href = '/sacola';
                    });
                } else {
                    throw new Error(pedidoData.erro || 'Erro ao criar pedido');
                }
            }
        } else {
            console.error('Erro: resposta inválida', data);
            throw new Error(data.erro || 'Erro ao processar pagamento: URL de pagamento não foi retornada');
        }
    } catch (error) {
        console.error('Erro ao criar cobrança:', error);
        
        let mensagemErro = error.message || 'Não foi possível processar o pagamento. Tente novamente.';
        
        // Mensagens específicas para erros comuns
        if (mensagemErro.includes('Invalid taxId') || mensagemErro.includes('CPF inválido')) {
            mensagemErro = 'CPF inválido. Por favor, verifique o CPF informado e tente novamente.';
        }
        
        Swal.fire({
            icon: 'error',
            title: 'Erro ao processar pagamento',
            text: mensagemErro,
            confirmButtonColor: '#370400'
        });
        btnFinalizar.disabled = false;
        btnFinalizar.textContent = 'Finalizar Compra';
    }
}

// MOSTRAR CONFIRMAÇÃO
function mostrarConfirmacao() {
    const numeroPedido = checkoutData.numeroPedido;
    const numeroPedidoEl = document.getElementById('numeroPedido');
    if (numeroPedido && numeroPedidoEl) {
        numeroPedidoEl.textContent = `#${numeroPedido}`;
    } else if (numeroPedidoEl) {
        // Se não houver número real do pedido, esconder a linha
        const parentParagraph = numeroPedidoEl.parentElement;
        if (parentParagraph) {
            parentParagraph.style.display = 'none';
        }
    }

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
            <p style="text-align: center; color: #370400; margin-top: 20px; font-weight: 600;">
                Pedido criado com sucesso!
            </p>
            <p style="text-align: center; color: #666;">
                Em breve você será redirecionado para o gateway de pagamento Abacate Pay para realizar o pagamento via PIX.
            </p>
            <p style="text-align: center; color: #666; font-size: 13px; margin-top: 10px;">
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

// Mostrar confirmação após retorno do AbacatePay
async function confirmarPagamentoECriarPedido() {
    try {
        // Limpar dados salvos
        const savedData = localStorage.getItem('checkoutData');
        if (savedData) {
            const data = JSON.parse(savedData);
            checkoutData.numeroPedido = data.numeroPedido;
        }
        
        localStorage.removeItem('checkoutData');
        localStorage.removeItem('abacateExternalId');
        
        // Mostrar confirmação
        carregarResumo();
        avancarEtapa(5);
        
    } catch (error) {
        console.error('Erro ao mostrar confirmação:', error);
        Swal.fire({
            icon: 'error',
            title: 'Erro',
            text: 'Houve um erro ao processar seu pedido.',
            confirmButtonColor: '#370400'
        }).then(() => {
            window.location.href = '/home';
        });
    }
}
