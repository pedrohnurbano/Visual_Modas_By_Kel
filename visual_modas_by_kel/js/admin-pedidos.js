// Gerenciamento de Pedidos no Painel Admin

// Carregar pedidos quando a aba for ativada
function carregarPedidosAdmin() {
    listarTodosPedidosAPI(function(pedidos) {
        if (!pedidos || pedidos.length === 0) {
            const tbody = document.querySelector('#pedidos-section .painel-data-table tbody');
            if (tbody) {
                tbody.innerHTML = '<tr><td colspan="6" style="text-align: center;">Nenhum pedido encontrado</td></tr>';
            }
            return;
        }

        renderizarPedidosAdmin(pedidos);
    });
}

function renderizarPedidosAdmin(pedidos) {
    const tbody = document.querySelector('#pedidos-section .painel-data-table tbody');
    if (!tbody) return;

    tbody.innerHTML = '';

    pedidos.forEach(pedido => {
        const tr = document.createElement('tr');
        const dataFormatada = new Date(pedido.criadoEm).toLocaleDateString('pt-BR');
        const statusClass = getStatusClass(pedido.status);
        const statusText = getStatusText(pedido.status);

        tr.innerHTML = `
            <td>#${pedido.id}</td>
            <td>${pedido.nomeCompleto}</td>
            <td>${dataFormatada}</td>
            <td>R$ ${pedido.total.toFixed(2).replace('.', ',')}</td>
            <td><span class="badge badge-${statusClass}">${statusText}</span></td>
            <td>
                <button class="painel-btn-secondary" style="padding: 5px 10px; font-size: 12px;" onclick="visualizarPedido(${pedido.id})">
                    Ver Detalhes
                </button>
                ${pedido.status === 'processamento' ? `
                    <button class="painel-btn-primary" style="padding: 5px 10px; font-size: 12px;" onclick="abrirModalRastreio(${pedido.id})">
                        Adicionar Rastreio
                    </button>
                ` : ''}
            </td>
        `;
        tbody.appendChild(tr);
    });
}

function getStatusClass(status) {
    const classes = {
        'pendente': 'warning',
        'processamento': 'info',
        'enviado': 'primary',
        'recebido': 'success',
        'cancelado': 'danger'
    };
    return classes[status] || 'secondary';
}

function getStatusText(status) {
    const texts = {
        'pendente': 'Pendente',
        'processamento': 'Em Processamento',
        'enviado': 'Enviado',
        'recebido': 'Recebido',
        'cancelado': 'Cancelado'
    };
    return texts[status] || status;
}

// Visualizar detalhes do pedido (admin)
function visualizarPedido(pedidoId) {
    // Usar ajax direto para admin buscar qualquer pedido
    $.ajax({
        url: `/api/pedidos/${pedidoId}`,
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        success: function(pedidoCompleto) {
            if (!pedidoCompleto) {
                Swal.fire({
                    icon: 'error',
                    title: 'Erro',
                    text: 'Não foi possível carregar os detalhes do pedido.',
                    confirmButtonColor: '#370400'
                });
                return;
            }

            const pedido = pedidoCompleto.pedido;
            const itens = pedidoCompleto.itens || [];
            
            mostrarDetalhesPedido(pedido, itens);
        },
        error: function(xhr) {
            console.error('Erro ao buscar pedido:', xhr);
            Swal.fire({
                icon: 'error',
                title: 'Erro',
                text: 'Não foi possível carregar os detalhes do pedido.',
                confirmButtonColor: '#370400'
            });
        }
    });
}

// Função para mostrar detalhes do pedido
function mostrarDetalhesPedido(pedido, itens) {

    let itensHTML = '';
    itens.forEach(item => {
        itensHTML += `
            <div style="display: flex; gap: 15px; padding: 10px; border-bottom: 1px solid #eee;">
                ${item.fotoUrl ? `<img src="${item.fotoUrl.startsWith('http') ? item.fotoUrl : 'http://localhost:5000' + item.fotoUrl}" style="width: 60px; height: 60px; object-fit: cover; border-radius: 4px;">` : ''}
                <div style="flex: 1;">
                    <div style="font-weight: 600;">${item.nomeProduto}</div>
                    <div style="color: #666; font-size: 13px;">
                        Tamanho: ${item.tamanho} | Quantidade: ${item.quantidade} | 
                        Preço unitário: R$ ${item.precoUnitario.toFixed(2).replace('.', ',')}
                    </div>
                </div>
                <div style="font-weight: 600;">R$ ${item.subtotal.toFixed(2).replace('.', ',')}</div>
            </div>
        `;
    });

    Swal.fire({
        title: `Pedido #${pedido.id}`,
        html: `
            <div style="text-align: left; max-height: 500px; overflow-y: auto;">
                <div style="margin-bottom: 20px;">
                    <h3 style="color: #370400; margin-bottom: 10px;">Dados do Cliente</h3>
                    <p><strong>Nome:</strong> ${pedido.nomeCompleto}</p>
                    <p><strong>E-mail:</strong> ${pedido.email}</p>
                    <p><strong>Telefone:</strong> ${pedido.telefone}</p>
                </div>

                <div style="margin-bottom: 20px;">
                    <h3 style="color: #370400; margin-bottom: 10px;">Endereço de Entrega</h3>
                    <p>${pedido.endereco}, ${pedido.numero}${pedido.complemento ? ' - ' + pedido.complemento : ''}</p>
                    <p>${pedido.bairro} - ${pedido.cidade}/${pedido.estado}</p>
                    <p>CEP: ${pedido.cep}</p>
                </div>

                <div style="margin-bottom: 20px;">
                    <h3 style="color: #370400; margin-bottom: 10px;">Status e Pagamento</h3>
                    <p><strong>Status:</strong> <span class="badge badge-${getStatusClass(pedido.status)}">${getStatusText(pedido.status)}</span></p>
                    <p><strong>Forma de Pagamento:</strong> ${pedido.formaPagamento}</p>
                    ${pedido.codigoRastreio ? `<p><strong>Código de Rastreio:</strong> ${pedido.codigoRastreio}</p>` : ''}
                </div>

                <div style="margin-bottom: 20px;">
                    <h3 style="color: #370400; margin-bottom: 10px;">Itens do Pedido</h3>
                    ${itensHTML}
                </div>

                <div style="margin-top: 20px; padding-top: 15px; border-top: 2px solid #370400;">
                    <div style="display: flex; justify-content: space-between; font-size: 18px; font-weight: 700; color: #370400;">
                        <span>Total:</span>
                        <span>R$ ${pedido.total.toFixed(2).replace('.', ',')}</span>
                    </div>
                </div>
            </div>
        `,
        width: '800px',
        confirmButtonColor: '#370400',
        confirmButtonText: 'Fechar'
    });
}

// Abrir modal para adicionar código de rastreio
function abrirModalRastreio(pedidoId) {
    Swal.fire({
        title: 'Adicionar Código de Rastreio',
        html: `
            <div style="text-align: left;">
                <p style="margin-bottom: 10px;">Digite o código de rastreio dos Correios:</p>
                <input id="codigoRastreio" class="swal2-input" placeholder="Ex: AA123456789BR" style="width: 90%;">
            </div>
        `,
        showCancelButton: true,
        confirmButtonColor: '#370400',
        cancelButtonColor: '#6c757d',
        confirmButtonText: 'Salvar',
        cancelButtonText: 'Cancelar',
        preConfirm: () => {
            const codigo = document.getElementById('codigoRastreio').value.trim();
            if (!codigo) {
                Swal.showValidationMessage('Por favor, digite o código de rastreio');
                return false;
            }
            return codigo;
        }
    }).then((result) => {
        if (result.isConfirmed) {
            salvarCodigoRastreio(pedidoId, result.value);
        }
    });
}

// Salvar código de rastreio
function salvarCodigoRastreio(pedidoId, codigoRastreio) {
    $.ajax({
        url: `/api/admin/pedidos/${pedidoId}/rastreio`,
        method: 'PUT',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
            'Content-Type': 'application/json'
        },
        data: JSON.stringify({ codigoRastreio: codigoRastreio }),
        success: function(response) {
            Swal.fire({
                icon: 'success',
                title: 'Sucesso!',
                text: 'Código de rastreio adicionado com sucesso! O pedido foi marcado como enviado.',
                confirmButtonColor: '#370400'
            }).then(() => {
                carregarPedidosAdmin();
            });
        },
        error: function(xhr) {
            console.error('Erro ao salvar código de rastreio:', xhr);
            let mensagem = 'Erro ao salvar código de rastreio. Tente novamente.';
            if (xhr.responseJSON && xhr.responseJSON.erro) {
                mensagem = xhr.responseJSON.erro;
            }
            Swal.fire({
                icon: 'error',
                title: 'Erro',
                text: mensagem,
                confirmButtonColor: '#370400'
            });
        }
    });
}

// Inicializar quando o documento estiver pronto
$(document).ready(function() {
    // Adicionar listener para quando a aba de pedidos for clicada
    const pedidosMenuItem = document.querySelector('.menu-item[data-section="pedidos"]');
    if (pedidosMenuItem) {
        pedidosMenuItem.addEventListener('click', function() {
            setTimeout(carregarPedidosAdmin, 100);
        });
    }
});

