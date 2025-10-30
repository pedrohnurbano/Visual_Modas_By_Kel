// Dados do usuário (serão carregados da API)
let userData = {};
let pedidos = [];

// Inicialização
document.addEventListener('DOMContentLoaded', function () {
    carregarDadosUsuario();
    configurarMenu();
    configurarFormularios();
});

// Carregar dados do usuário da API
async function carregarDadosUsuario() {
    try {
        const response = await fetch('/api/usuarios/dados');
        
        if (!response.ok) {
            throw new Error('Erro ao carregar dados do usuário');
        }
        
        const data = await response.json();
        userData = data;

        // Atualizar cabeçalho
        document.getElementById('userName').textContent = userData.nome || 'Usuário';

        // Atualizar visualização dos dados
        document.getElementById('viewNome').textContent = userData.nome && userData.sobrenome 
            ? `${userData.nome} ${userData.sobrenome}` 
            : userData.nome || '-';
        document.getElementById('viewCpf').textContent = formatarCPF(userData.cpf) || '-';
        document.getElementById('viewEmail').textContent = userData.email || '-';
        document.getElementById('viewTelefone').textContent = formatarTelefone(userData.telefone) || '-';
        
        // Carregar dados relacionados
        carregarPedidos();
        carregarFavoritos();
    } catch (error) {
        console.error('Erro ao carregar dados do usuário:', error);
        showNotification('Erro ao carregar seus dados. Tente novamente.');
    }
}

// Configurar navegação do menu
function configurarMenu() {
    const menuItems = document.querySelectorAll('.menu-item');
    const sections = document.querySelectorAll('.usuario-section');

    menuItems.forEach(item => {
        item.addEventListener('click', function () {
            const targetSection = this.getAttribute('data-section');

            // Remover active de todos
            menuItems.forEach(mi => mi.classList.remove('active'));
            sections.forEach(s => s.classList.remove('active'));

            // Adicionar active no clicado
            this.classList.add('active');
            document.getElementById(targetSection).classList.add('active');

            // Scroll para o topo em mobile
            if (window.innerWidth <= 1024) {
                window.scrollTo({ top: 0, behavior: 'smooth' });
            }
        });
    });
}

// Configurar formulários
function configurarFormularios() {
    // Botão editar dados
    document.getElementById('btnEditarDados').addEventListener('click', function () {
        document.getElementById('dadosView').classList.add('hidden');
        document.getElementById('dadosEdit').classList.remove('hidden');

        // Preencher formulário
        document.getElementById('editNome').value = userData.nome || '';
        document.getElementById('editSobrenome').value = userData.sobrenome || '';
        document.getElementById('editCpf').value = formatarCPF(userData.cpf) || '';
        document.getElementById('editTelefone').value = formatarTelefone(userData.telefone) || '';
    });

    // Form de edição de dados
    document.getElementById('dadosEdit').addEventListener('submit', function (e) {
        e.preventDefault();
        salvarDadosPessoais();
    });

    // Form de alteração de senha
    document.getElementById('formAlterarSenha').addEventListener('submit', function (e) {
        e.preventDefault();
        alterarSenha();
    });

    // Máscaras
    aplicarMascaras();
}

// Aplicar máscaras nos inputs
function aplicarMascaras() {
    const cpfInput = document.getElementById('editCpf');
    if (cpfInput) {
        cpfInput.addEventListener('input', function () {
            this.value = maskCPF(this.value);
        });
    }

    const telefoneInput = document.getElementById('editTelefone');
    if (telefoneInput) {
        telefoneInput.addEventListener('input', function () {
            this.value = maskPhone(this.value);
        });
    }
}

function maskCPF(value) {
    return value
        .replace(/\D/g, '')
        .replace(/(\d{3})(\d)/, '$1.$2')
        .replace(/(\d{3})(\d)/, '$1.$2')
        .replace(/(\d{3})(\d{1,2})$/, '$1-$2')
        .slice(0, 14);
}

function maskPhone(value) {
    return value
        .replace(/\D/g, '')
        .replace(/^(\d{2})(\d)/g, '($1) $2')
        .replace(/(\d{5})(\d{1,4})$/, '$1-$2')
        .slice(0, 15);
}

// Salvar dados pessoais
async function salvarDadosPessoais() {
    const dadosAtualizados = {
        nome: document.getElementById('editNome').value.trim(),
        sobrenome: document.getElementById('editSobrenome').value.trim(),
        cpf: document.getElementById('editCpf').value.trim(),
        telefone: document.getElementById('editTelefone').value.trim(),
        email: userData.email // Mantém o email original (não editável)
    };

    try {
        const response = await fetch(`/api/usuarios/${userData.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(dadosAtualizados)
        });

        if (response.ok || response.status === 204) {
            await carregarDadosUsuario();
            document.getElementById('dadosEdit').classList.add('hidden');
            document.getElementById('dadosView').classList.remove('hidden');
            showNotification('Dados atualizados com sucesso!');
        } else {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.erro || 'Erro ao atualizar dados');
        }
    } catch (error) {
        console.error('Erro ao salvar dados:', error);
        showNotification(error.message || 'Erro ao salvar dados. Tente novamente.');
    }
}

function cancelarEdicaoDados() {
    document.getElementById('dadosEdit').classList.add('hidden');
    document.getElementById('dadosView').classList.remove('hidden');
}

// Pedidos
async function carregarPedidos() {
    try {
        const response = await fetch('/api/pedidos');
        if (!response.ok) {
            throw new Error('Erro ao carregar pedidos');
        }
        
        pedidos = await response.json() || [];
        
        const lista = document.getElementById('pedidosLista');
        const empty = document.getElementById('pedidosEmpty');

        if (!pedidos || pedidos.length === 0) {
            lista.style.display = 'none';
            empty.style.display = 'flex';
            return;
        }

        lista.style.display = 'flex';
        empty.style.display = 'none';

        lista.innerHTML = pedidos.map(pedido => {
            const dataFormatada = new Date(pedido.criadoEm).toLocaleDateString('pt-BR');
            const statusClass = getStatusClassUsuario(pedido.status);
            
            return `
                <div class="pedido-card" data-status="${pedido.status}">
                    <div class="pedido-header">
                        <div>
                            <div class="pedido-numero">Pedido #${pedido.id}</div>
                            <div class="pedido-data">${dataFormatada}</div>
                        </div>
                        <span class="pedido-status badge badge-${statusClass}">${getStatusText(pedido.status)}</span>
                    </div>
                    <div class="pedido-info">
                        <p><strong>Total:</strong> R$ ${pedido.total.toFixed(2).replace('.', ',')}</p>
                        <p><strong>Forma de Pagamento:</strong> ${pedido.formaPagamento}</p>
                        ${pedido.codigoRastreio ? `<p><strong>Código de Rastreio:</strong> <span style="color: #370400; font-weight: 600;">${pedido.codigoRastreio}</span></p>` : ''}
                    </div>
                    <div class="pedido-footer">
                        <div class="pedido-acoes">
                            <button class="btn-small primary" onclick="verDetalhesPedidoUsuario(${pedido.id})">
                                Ver Detalhes
                            </button>
                            ${pedido.codigoRastreio ? `
                                <a href="https://www.correios.com.br/rastreamento" target="_blank" class="btn-small secondary">
                                    Rastrear Entrega
                                </a>
                            ` : ''}
                            ${pedido.status === 'enviado' && pedido.codigoRastreio ? `
                                <button class="btn-small success" onclick="confirmarEntregaPedido(${pedido.id})">
                                    Confirmar Recebimento
                                </button>
                            ` : ''}
                        </div>
                    </div>
                </div>
            `;
        }).join('');
    } catch (error) {
        console.error('Erro ao carregar pedidos:', error);
        const lista = document.getElementById('pedidosLista');
        const empty = document.getElementById('pedidosEmpty');
        if (lista && empty) {
            lista.style.display = 'none';
            empty.style.display = 'flex';
        }
    }
}

function getStatusText(status) {
    const statusMap = {
        'pendente': 'Pendente',
        'processamento': 'Em Processamento',
        'enviado': 'Enviado',
        'recebido': 'Recebido',
        'cancelado': 'Cancelado'
    };
    return statusMap[status] || status;
}

function getStatusClassUsuario(status) {
    const classMap = {
        'pendente': 'warning',
        'processamento': 'info',
        'enviado': 'primary',
        'recebido': 'success',
        'cancelado': 'danger'
    };
    return classMap[status] || 'secondary';
}

function verDetalhesPedidoUsuario(pedidoId) {
    const pedido = pedidos.find(p => p.id === pedidoId);
    if (!pedido) {
        showNotification('Pedido não encontrado');
        return;
    }

    // Buscar detalhes completos do pedido com itens
    fetch(`/api/pedidos/${pedidoId}`)
        .then(response => response.json())
        .then(pedidoCompleto => {
            const itens = pedidoCompleto.itens || [];
            let itensHTML = '';
            
            itens.forEach(item => {
                const fotoUrl = item.fotoUrl ? 
                    (item.fotoUrl.startsWith('http') ? item.fotoUrl : `http://localhost:5000${item.fotoUrl}`) : 
                    'design/cabide.png';
                    
                itensHTML += `
                    <div style="display: flex; gap: 15px; padding: 10px; border-bottom: 1px solid #eee;">
                        <img src="${fotoUrl}" style="width: 60px; height: 60px; object-fit: cover; border-radius: 4px;" onerror="this.src='design/cabide.png'">
                        <div style="flex: 1;">
                            <div style="font-weight: 600;">${item.nomeProduto}</div>
                            <div style="color: #666; font-size: 13px;">
                                Tamanho: ${item.tamanho} | Quantidade: ${item.quantidade}
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
                            <h3 style="color: #370400; margin-bottom: 10px;">Status</h3>
                            <span class="badge badge-${getStatusClassUsuario(pedido.status)}">${getStatusText(pedido.status)}</span>
                            <p style="margin-top: 10px;"><strong>Data do Pedido:</strong> ${new Date(pedido.criadoEm).toLocaleDateString('pt-BR')}</p>
                            ${pedido.codigoRastreio ? `<p><strong>Código de Rastreio:</strong> ${pedido.codigoRastreio}</p>` : ''}
                        </div>

                        <div style="margin-bottom: 20px;">
                            <h3 style="color: #370400; margin-bottom: 10px;">Endereço de Entrega</h3>
                            <p>${pedido.endereco}, ${pedido.numero}${pedido.complemento ? ' - ' + pedido.complemento : ''}</p>
                            <p>${pedido.bairro} - ${pedido.cidade}/${pedido.estado}</p>
                            <p>CEP: ${pedido.cep}</p>
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
                width: '700px',
                confirmButtonColor: '#370400',
                confirmButtonText: 'Fechar'
            });
        })
        .catch(error => {
            console.error('Erro ao carregar detalhes do pedido:', error);
            showNotification('Erro ao carregar detalhes do pedido');
        });
}

function confirmarEntregaPedido(pedidoId) {
    Swal.fire({
        title: 'Confirmar Recebimento',
        text: 'Você confirma que recebeu este pedido?',
        icon: 'question',
        showCancelButton: true,
        confirmButtonColor: '#370400',
        cancelButtonColor: '#6c757d',
        confirmButtonText: 'Sim, recebi',
        cancelButtonText: 'Cancelar'
    }).then(async (result) => {
        if (result.isConfirmed) {
            try {
                const response = await fetch(`/api/pedidos/${pedidoId}/confirmar-entrega`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });

                if (response.ok) {
                    Swal.fire({
                        icon: 'success',
                        title: 'Entrega Confirmada!',
                        text: 'Obrigado por confirmar o recebimento do seu pedido.',
                        confirmButtonColor: '#370400'
                    }).then(() => {
                        carregarPedidos(); // Recarregar lista de pedidos
                    });
                } else {
                    throw new Error('Erro ao confirmar entrega');
                }
            } catch (error) {
                console.error('Erro ao confirmar entrega:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Erro',
                    text: 'Não foi possível confirmar o recebimento. Tente novamente.',
                    confirmButtonColor: '#370400'
                });
            }
        }
    });
}

// Favoritos
async function carregarFavoritos() {
    try {
        const grid = document.getElementById('favoritosGrid');
        const empty = document.getElementById('favoritosEmpty');
        
        // Buscar favoritos da API
        const response = await fetch('/api/favoritos');
        if (!response.ok) {
            throw new Error('Erro ao carregar favoritos');
        }
        
        const favoritos = await response.json();

        if (!favoritos || favoritos.length === 0) {
            grid.style.display = 'none';
            empty.style.display = 'block';
            return;
        }

        grid.style.display = 'grid';
        empty.style.display = 'none';

        grid.innerHTML = favoritos.map(produto => {
            const fotoUrl = produto.foto_url && !produto.foto_url.startsWith('http') 
                ? `http://localhost:5000${produto.foto_url.startsWith('/') ? '' : '/'}${produto.foto_url}`
                : produto.foto_url || 'design/ex-roupa1.png';
                
            return `
                <div class="produto-card">
                    <button class="produto-favorito favorito-active" onclick="removerFavorito(${produto.id})">
                        <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
                            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41 0.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
                                stroke="#370400" stroke-width="1.5" fill="#370400"/>
                        </svg>
                    </button>
                    <div class="produto-img-container">
                        <img src="${fotoUrl}" alt="${produto.nome}" class="produto-img">
                    </div>
                    <div class="produto-info">
                        <span class="produto-nome">${produto.nome}</span>
                        <span class="produto-preco">R$ ${produto.preco.toFixed(2).replace('.', ',')}</span>
                    </div>
                </div>
            `;
        }).join('');
    } catch (error) {
        console.error('Erro ao carregar favoritos:', error);
        const grid = document.getElementById('favoritosGrid');
        const empty = document.getElementById('favoritosEmpty');
        grid.style.display = 'none';
        empty.style.display = 'block';
    }
}

async function removerFavorito(produtoId) {
    try {
        const response = await fetch(`/api/favoritos/${produtoId}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            await carregarFavoritos();
            showNotification('Produto removido dos favoritos');
        } else {
            throw new Error('Erro ao remover favorito');
        }
    } catch (error) {
        console.error('Erro ao remover favorito:', error);
        showNotification('Erro ao remover favorito. Tente novamente.');
    }
}

// Segurança
async function alterarSenha() {
    const senhaAtual = document.getElementById('senhaAtual').value;
    const novaSenha = document.getElementById('novaSenha').value;
    const confirmarSenha = document.getElementById('confirmarNovaSenha').value;

    if (!senhaAtual || !novaSenha || !confirmarSenha) {
        showNotification('Preencha todos os campos!');
        return;
    }

    if (novaSenha !== confirmarSenha) {
        showNotification('As senhas não coincidem!');
        return;
    }

    try {
        const response = await fetch(`/api/usuarios/${userData.id}/atualizar-senha`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                atual: senhaAtual,
                nova: novaSenha
            })
        });

        if (response.ok || response.status === 204) {
            showNotification('Senha alterada com sucesso!');
            document.getElementById('formAlterarSenha').reset();
        } else {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.erro || 'Erro ao alterar senha');
        }
    } catch (error) {
        console.error('Erro ao alterar senha:', error);
        showNotification(error.message || 'Erro ao alterar senha. Verifique a senha atual e tente novamente.');
    }
}

function togglePasswordField(fieldId, element) {
    const field = document.getElementById(fieldId);
    const icon = element.querySelector('i');

    if (field.type === 'password') {
        field.type = 'text';
        icon.classList.remove('fa-eye');
        icon.classList.add('fa-eye-slash');
    } else {
        field.type = 'password';
        icon.classList.remove('fa-eye-slash');
        icon.classList.add('fa-eye');
    }
}

// Exclusão de conta
async function confirmarExclusaoConta() {
    const { value: senha } = await Swal.fire({
        title: 'Excluir Conta',
        html: `
            <div style="text-align: left; margin: 20px 0;">
                <div style="background: #fff3cd; padding: 15px; border-radius: 5px; margin-bottom: 20px; border-left: 4px solid #ffc107;">
                    <strong>⚠️ Atenção!</strong> Esta ação não pode ser desfeita.
                </div>
                <p style="margin-bottom: 15px;">Ao excluir sua conta, você perderá:</p>
                <ul style="margin-left: 20px; margin-bottom: 20px;">
                    <li>Histórico de pedidos</li>
                    <li>Lista de favoritos</li>
                    <li>Todos os seus dados pessoais</li>
                </ul>
                <p style="margin-bottom: 10px;"><strong>Digite sua senha para confirmar:</strong></p>
            </div>
        `,
        input: 'password',
        inputPlaceholder: 'Digite sua senha',
        inputAttributes: {
            autocapitalize: 'off',
            autocorrect: 'off'
        },
        showCancelButton: true,
        confirmButtonText: 'Excluir Conta',
        cancelButtonText: 'Cancelar',
        confirmButtonColor: '#dc3545',
        cancelButtonColor: '#6c757d',
        inputValidator: (value) => {
            if (!value) {
                return 'Você precisa digitar sua senha!'
            }
        }
    });

    if (senha) {
        // Confirmação final
        const confirmacao = await Swal.fire({
            title: 'Tem certeza absoluta?',
            text: 'Esta ação é IRREVERSÍVEL!',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#dc3545',
            cancelButtonColor: '#6c757d',
            confirmButtonText: 'Sim, excluir permanentemente!',
            cancelButtonText: 'Cancelar'
        });

        if (confirmacao.isConfirmed) {
            try {
                const response = await fetch(`/api/usuarios/${userData.id}`, {
                    method: 'DELETE',
                    headers: { 
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ senha: senha })
                });

                if (response.ok || response.status === 204) {
                    // Limpar dados locais e cookies
                    localStorage.clear();
                    sessionStorage.clear();
                    
                    // Limpar todos os cookies
                    document.cookie.split(";").forEach(function(c) { 
                        document.cookie = c.replace(/^ +/, "").replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/"); 
                    });
                    
                    await Swal.fire({
                        icon: 'success',
                        title: 'Conta excluída',
                        text: 'Sua conta foi excluída com sucesso. Até breve!',
                        confirmButtonColor: '#370400',
                        timer: 2500,
                        showConfirmButton: false
                    });
                    
                    window.location.href = '/login';
                } else {
                    const errorData = await response.json().catch(() => ({}));
                    throw new Error(errorData.erro || 'Senha incorreta ou erro ao excluir conta');
                }
            } catch (error) {
                console.error('Erro ao excluir conta:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: error.message || 'Não foi possível excluir a conta. Verifique sua senha e tente novamente.',
                    confirmButtonColor: '#370400'
                });
            }
        }
    }
}

// Logout
function logout() {
    Swal.fire({
        title: 'Sair da conta?',
        text: 'Deseja realmente sair da sua conta?',
        icon: 'question',
        showCancelButton: true,
        confirmButtonColor: '#370400',
        cancelButtonColor: '#6c757d',
        confirmButtonText: 'Sim, sair',
        cancelButtonText: 'Cancelar'
    }).then(async (result) => {
        if (result.isConfirmed) {
            try {
                const response = await fetch('/logout', { method: 'POST' });
                
                Swal.fire({
                    icon: 'success',
                    title: 'Até logo!',
                    text: 'Saindo...',
                    confirmButtonColor: '#370400',
                    timer: 1000,
                    showConfirmButton: false
                }).then(() => {
                    window.location.href = '/login';
                });
            } catch (error) {
                console.error('Erro ao fazer logout:', error);
                // Mesmo com erro, redirecionar para o login
                window.location.href = '/login';
            }
        }
    });
}

// Funções auxiliares
function formatarCPF(cpf) {
    if (!cpf) return '';
    const apenasNumeros = cpf.replace(/\D/g, '');
    if (apenasNumeros.length === 11) {
        return apenasNumeros.replace(/(\d{3})(\d{3})(\d{3})(\d{2})/, '$1.$2.$3-$4');
    }
    return cpf;
}

function formatarTelefone(telefone) {
    if (!telefone) return '';
    const apenasNumeros = telefone.replace(/\D/g, '');
    if (apenasNumeros.length === 11) {
        return apenasNumeros.replace(/(\d{2})(\d{5})(\d{4})/, '($1) $2-$3');
    } else if (apenasNumeros.length === 10) {
        return apenasNumeros.replace(/(\d{2})(\d{4})(\d{4})/, '($1) $2-$3');
    }
    return telefone;
}

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
        z-index: 10001;
        animation: slideIn 0.3s ease-out;
        max-width: 300px;
    `;
    notification.textContent = message;
    document.body.appendChild(notification);

    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

// Fechar modais ao clicar fora
document.addEventListener('click', function (e) {
    const modals = document.querySelectorAll('.modal');
    modals.forEach(modal => {
        if (e.target === modal) {
            modal.classList.remove('active');
        }
    });
});

// Fechar modais com ESC
document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') {
        document.querySelectorAll('.modal').forEach(modal => {
            modal.classList.remove('active');
        });
    }
});

// Animações de entrada
const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.style.animation = 'fadeIn 0.5s ease-out';
        }
    });
}, { threshold: 0.1 });

document.querySelectorAll('.endereco-card, .pedido-card, .produto-card').forEach(el => {
    observer.observe(el);
});