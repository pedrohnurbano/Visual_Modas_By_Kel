// Dados do usuário (simulação - em produção viriam do backend)
let userData = {
    nome: 'Maria',
    sobrenome: 'Silva',
    cpf: '123.456.789-00',
    telefone: '(48) 99999-9999',
    email: 'maria.silva@email.com',
    dataNasc: '1990-05-15',
    genero: 'feminino'
};

let enderecos = [
    {
        id: 1,
        cep: '88000-000',
        rua: 'Rua das Flores',
        numero: '123',
        complemento: 'Apto 501',
        bairro: 'Centro',
        cidade: 'Criciúma',
        estado: 'SC',
        principal: true
    }
];

let pedidos = [
    {
        id: '2025001',
        data: '15/01/2025',
        status: 'entregue',
        items: [
            {
                nome: 'Vestido Midi Floral',
                tamanho: 'M',
                quantidade: 1,
                preco: 899.90,
                imagem: 'design/ex-roupa1.png'
            }
        ],
        total: 899.90
    },
    {
        id: '2025002',
        data: '20/01/2025',
        status: 'enviado',
        items: [
            {
                nome: 'Blusa Cropped Básica',
                tamanho: 'P',
                quantidade: 2,
                preco: 449.90,
                imagem: 'design/ex-roupa2.png'
            }
        ],
        total: 899.80
    }
];

// Inicialização
document.addEventListener('DOMContentLoaded', function () {
    verificarLogin();
    carregarDadosUsuario();
    configurarMenu();
    configurarFormularios();
    carregarEnderecos();
    carregarPedidos();
    carregarFavoritos();
});

// Verificar se usuário está logado
function verificarLogin() {
    const isLoggedIn = localStorage.getItem('userLoggedIn') === 'true';
    if (!isLoggedIn) {
        window.location.href = 'login.html';
    }
}

// Carregar dados do usuário
function carregarDadosUsuario() {
    const savedData = localStorage.getItem('userData');
    if (savedData) {
        userData = JSON.parse(savedData);
    }

    // Atualizar cabeçalho
    document.getElementById('userName').textContent = userData.nome;

    // Atualizar visualização dos dados
    document.getElementById('viewNome').textContent = `${userData.nome} ${userData.sobrenome}`;
    document.getElementById('viewCpf').textContent = userData.cpf;
    document.getElementById('viewEmail').textContent = userData.email;
    document.getElementById('viewTelefone').textContent = userData.telefone;
    document.getElementById('viewDataNasc').textContent = userData.dataNasc ? formatarData(userData.dataNasc) : '-';
    document.getElementById('viewGenero').textContent = formatarGenero(userData.genero);
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
        document.getElementById('editNome').value = userData.nome;
        document.getElementById('editSobrenome').value = userData.sobrenome;
        document.getElementById('editCpf').value = userData.cpf;
        document.getElementById('editTelefone').value = userData.telefone;
        document.getElementById('editEmail').value = userData.email;
        document.getElementById('editDataNasc').value = userData.dataNasc || '';
        document.getElementById('editGenero').value = userData.genero || '';
    });

    // Form de edição de dados
    document.getElementById('dadosEdit').addEventListener('submit', function (e) {
        e.preventDefault();
        salvarDadosPessoais();
    });

    // Form de endereço
    document.getElementById('formEndereco').addEventListener('submit', function (e) {
        e.preventDefault();
        salvarEndereco();
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
    const cpfInputs = document.querySelectorAll('#editCpf, #enderecoCep');
    cpfInputs.forEach(input => {
        if (input.id.includes('Cpf')) {
            input.addEventListener('input', function () {
                this.value = maskCPF(this.value);
            });
        } else if (input.id.includes('Cep')) {
            input.addEventListener('input', function () {
                this.value = maskCEP(this.value);
            });
        }
    });

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

function maskCEP(value) {
    return value
        .replace(/\D/g, '')
        .replace(/(\d{5})(\d)/, '$1-$2')
        .slice(0, 9);
}

function maskPhone(value) {
    return value
        .replace(/\D/g, '')
        .replace(/^(\d{2})(\d)/g, '($1) $2')
        .replace(/(\d{5})(\d{1,4})$/, '$1-$2')
        .slice(0, 15);
}

// Salvar dados pessoais
function salvarDadosPessoais() {
    userData = {
        nome: document.getElementById('editNome').value,
        sobrenome: document.getElementById('editSobrenome').value,
        cpf: document.getElementById('editCpf').value,
        telefone: document.getElementById('editTelefone').value,
        email: document.getElementById('editEmail').value,
        dataNasc: document.getElementById('editDataNasc').value,
        genero: document.getElementById('editGenero').value
    };

    localStorage.setItem('userData', JSON.stringify(userData));
    carregarDadosUsuario();

    document.getElementById('dadosEdit').classList.add('hidden');
    document.getElementById('dadosView').classList.remove('hidden');

    showNotification('Dados atualizados com sucesso!');
}

function cancelarEdicaoDados() {
    document.getElementById('dadosEdit').classList.add('hidden');
    document.getElementById('dadosView').classList.remove('hidden');
}

// Endereços
function carregarEnderecos() {
    const savedEnderecos = localStorage.getItem('enderecos');
    if (savedEnderecos) {
        enderecos = JSON.parse(savedEnderecos);
    }

    const grid = document.getElementById('enderecosGrid');
    const empty = document.getElementById('enderecosEmpty');

    if (enderecos.length === 0) {
        grid.style.display = 'none';
        empty.style.display = 'block';
        return;
    }

    grid.style.display = 'grid';
    empty.style.display = 'none';

    grid.innerHTML = enderecos.map(end => `
        <div class="endereco-card ${end.principal ? 'principal' : ''}">
            ${end.principal ? '<span class="endereco-badge">Principal</span>' : ''}
            <div class="endereco-info">
                <strong>Endereço</strong>
                <p>${end.rua}, ${end.numero}</p>
                ${end.complemento ? `<p>${end.complemento}</p>` : ''}
                <p>${end.bairro} - ${end.cidade}/${end.estado}</p>
                <p>CEP: ${end.cep}</p>
            </div>
            <div class="endereco-actions">
                <button class="btn-icon" onclick="editarEndereco(${end.id})">
                    <i class="fas fa-edit"></i> Editar
                </button>
                <button class="btn-icon delete" onclick="excluirEndereco(${end.id})">
                    <i class="fas fa-trash"></i> Excluir
                </button>
            </div>
        </div>
    `).join('');
}

function abrirModalEndereco(id = null) {
    const modal = document.getElementById('modalEndereco');
    const titulo = document.getElementById('modalEnderecoTitulo');

    if (id) {
        const endereco = enderecos.find(e => e.id === id);
        titulo.textContent = 'Editar Endereço';
        preencherFormEndereco(endereco);
    } else {
        titulo.textContent = 'Adicionar Endereço';
        document.getElementById('formEndereco').reset();
    }

    modal.classList.add('active');
}

function fecharModalEndereco() {
    document.getElementById('modalEndereco').classList.remove('active');
    document.getElementById('formEndereco').reset();
}

function preencherFormEndereco(endereco) {
    document.getElementById('enderecoId').value = endereco.id;
    document.getElementById('enderecoCep').value = endereco.cep;
    document.getElementById('enderecoRua').value = endereco.rua;
    document.getElementById('enderecoNumero').value = endereco.numero;
    document.getElementById('enderecoComplemento').value = endereco.complemento || '';
    document.getElementById('enderecoBairro').value = endereco.bairro;
    document.getElementById('enderecoCidade').value = endereco.cidade;
    document.getElementById('enderecoEstado').value = endereco.estado;
    document.getElementById('enderecoPrincipal').checked = endereco.principal;
}

function buscarCep() {
    const cep = document.getElementById('enderecoCep').value.replace(/\D/g, '');

    if (cep.length !== 8) {
        showNotification('CEP inválido!');
        return;
    }

    // Aqui você faria uma chamada para API de CEP (ViaCEP, por exemplo)
    // Simulação:
    fetch(`https://viacep.com.br/ws/${cep}/json/`)
        .then(response => response.json())
        .then(data => {
            if (data.erro) {
                showNotification('CEP não encontrado!');
                return;
            }

            document.getElementById('enderecoRua').value = data.logradouro;
            document.getElementById('enderecoBairro').value = data.bairro;
            document.getElementById('enderecoCidade').value = data.localidade;
            document.getElementById('enderecoEstado').value = data.uf;
        })
        .catch(() => {
            showNotification('Erro ao buscar CEP!');
        });
}

function salvarEndereco() {
    const id = document.getElementById('enderecoId').value;
    const endereco = {
        id: id ? parseInt(id) : Date.now(),
        cep: document.getElementById('enderecoCep').value,
        rua: document.getElementById('enderecoRua').value,
        numero: document.getElementById('enderecoNumero').value,
        complemento: document.getElementById('enderecoComplemento').value,
        bairro: document.getElementById('enderecoBairro').value,
        cidade: document.getElementById('enderecoCidade').value,
        estado: document.getElementById('enderecoEstado').value,
        principal: document.getElementById('enderecoPrincipal').checked
    };

    if (endereco.principal) {
        enderecos.forEach(e => e.principal = false);
    }

    if (id) {
        const index = enderecos.findIndex(e => e.id === parseInt(id));
        enderecos[index] = endereco;
    } else {
        enderecos.push(endereco);
    }

    localStorage.setItem('enderecos', JSON.stringify(enderecos));
    carregarEnderecos();
    fecharModalEndereco();
    showNotification('Endereço salvo com sucesso!');
}

function editarEndereco(id) {
    abrirModalEndereco(id);
}

function excluirEndereco(id) {
    Swal.fire({
        title: 'Excluir endereço?',
        text: 'Deseja realmente excluir este endereço?',
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#370400',
        cancelButtonColor: '#6c757d',
        confirmButtonText: 'Sim, excluir!',
        cancelButtonText: 'Cancelar'
    }).then((result) => {
        if (result.isConfirmed) {
            enderecos = enderecos.filter(e => e.id !== id);
            localStorage.setItem('enderecos', JSON.stringify(enderecos));
            carregarEnderecos();
            showNotification('Endereço excluído!');
        }
    });
}

// Pedidos
function carregarPedidos() {
    const savedPedidos = localStorage.getItem('pedidos');
    if (savedPedidos) {
        pedidos = JSON.parse(savedPedidos);
    }

    const lista = document.getElementById('pedidosLista');
    const empty = document.getElementById('pedidosEmpty');

    if (pedidos.length === 0) {
        lista.style.display = 'none';
        empty.style.display = 'block';
        return;
    }

    lista.style.display = 'flex';
    empty.style.display = 'none';

    lista.innerHTML = pedidos.map(pedido => `
        <div class="pedido-card">
            <div class="pedido-header">
                <div>
                    <div class="pedido-numero">Pedido #${pedido.id}</div>
                    <div class="pedido-data">${pedido.data}</div>
                </div>
                <span class="pedido-status status-${pedido.status}">${getStatusText(pedido.status)}</span>
                </div>
            <div class="pedido-items">
                ${pedido.items.map(item => `
                    <div class="pedido-item">
                        <img src="${item.imagem}" alt="${item.nome}">
                        <div class="pedido-item-info">
                            <div class="pedido-item-nome">${item.nome}</div>
                            <div class="pedido-item-detalhes">
                                Tamanho: ${item.tamanho} | Qtd: ${item.quantidade} | 
                                R$ ${item.preco.toFixed(2).replace('.', ',')}
                            </div>
                        </div>
                    </div>
                `).join('')}
            </div>
            <div class="pedido-footer">
                <div class="pedido-total">Total: R$ ${pedido.total.toFixed(2).replace('.', ',')}</div>
                <div class="pedido-acoes">
                    <button class="btn-small primary" onclick="verDetalhesPedido('${pedido.id}')">
                        Ver Detalhes
                    </button>
                    ${pedido.status === 'enviado' ? `
                        <button class="btn-small secondary" onclick="rastrearPedido('${pedido.id}')">
                            Rastrear
                        </button>
                    ` : ''}
                </div>
            </div>
        </div>
    `).join('');

    // Filtro de pedidos
    document.getElementById('filterPedidos').addEventListener('change', function () {
        filtrarPedidos(this.value);
    });
}

function getStatusText(status) {
    const statusMap = {
        'em-processamento': 'Em Processamento',
        'enviado': 'Enviado',
        'entregue': 'Entregue',
        'cancelado': 'Cancelado'
    };
    return statusMap[status] || status;
}

function filtrarPedidos(filtro) {
    const cards = document.querySelectorAll('.pedido-card');
    cards.forEach(card => {
        const status = card.querySelector('.pedido-status').classList[1].replace('status-', '');
        if (filtro === 'todos' || status === filtro) {
            card.style.display = 'block';
        } else {
            card.style.display = 'none';
        }
    });
}

function verDetalhesPedido(id) {
    showNotification(`Ver detalhes do pedido #${id}`);
    // Aqui você implementaria a navegação para página de detalhes
}

function rastrearPedido(id) {
    showNotification(`Rastreamento do pedido #${id}`);
    // Aqui você implementaria o rastreamento
}

// Favoritos
function carregarFavoritos() {
    const favorites = JSON.parse(localStorage.getItem('favorites')) || [];
    const grid = document.getElementById('favoritosGrid');
    const empty = document.getElementById('favoritosEmpty');

    if (favorites.length === 0) {
        grid.style.display = 'none';
        empty.style.display = 'block';
        return;
    }

    grid.style.display = 'grid';
    empty.style.display = 'none';

    // Aqui você buscaria os produtos favoritos
    // Por enquanto vou usar dados mockados
    grid.innerHTML = favorites.map(id => `
        <div class="produto-card">
            <button class="produto-favorito favorito-active" onclick="removerFavorito(${id})">
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
                    <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41 0.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
                        stroke="#370400" stroke-width="1.5" fill="#370400"/>
                </svg>
            </button>
            <div class="produto-img-container">
                <img src="design/ex-roupa${id}.png" alt="Produto" class="produto-img">
            </div>
            <div class="produto-info">
                <span class="produto-nome">Produto ${id}</span>
                <span class="produto-preco">R$ 899,90</span>
            </div>
        </div>
    `).join('');
}

function removerFavorito(id) {
    let favorites = JSON.parse(localStorage.getItem('favorites')) || [];
    favorites = favorites.filter(fav => fav !== id);
    localStorage.setItem('favorites', JSON.stringify(favorites));
    carregarFavoritos();
    showNotification('Produto removido dos favoritos');
}

// Segurança
function alterarSenha() {
    const senhaAtual = document.getElementById('senhaAtual').value;
    const novaSenha = document.getElementById('novaSenha').value;
    const confirmarSenha = document.getElementById('confirmarNovaSenha').value;

    if (novaSenha.length < 8) {
        showNotification('A senha deve ter no mínimo 8 caracteres!');
        return;
    }

    if (novaSenha !== confirmarSenha) {
        showNotification('As senhas não coincidem!');
        return;
    }

    // Aqui você validaria a senha atual com o backend
    // Simulação:
    showNotification('Senha alterada com sucesso!');
    document.getElementById('formAlterarSenha').reset();
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

// Configurações
document.addEventListener('DOMContentLoaded', function () {
    const emailPromocional = document.getElementById('emailPromocional');
    const emailPedidos = document.getElementById('emailPedidos');

    if (emailPromocional) {
        emailPromocional.addEventListener('change', function () {
            localStorage.setItem('emailPromocional', this.checked);
            showNotification('Preferência salva!');
        });
    }

    if (emailPedidos) {
        emailPedidos.addEventListener('change', function () {
            localStorage.setItem('emailPedidos', this.checked);
            showNotification('Preferência salva!');
        });
    }

    // Carregar preferências salvas
    const savedEmailProm = localStorage.getItem('emailPromocional');
    const savedEmailPed = localStorage.getItem('emailPedidos');

    if (savedEmailProm !== null) {
        emailPromocional.checked = savedEmailProm === 'true';
    }
    if (savedEmailPed !== null) {
        emailPedidos.checked = savedEmailPed === 'true';
    }
});

// Exclusão de conta
function confirmarExclusaoConta() {
    document.getElementById('modalExcluirConta').classList.add('active');
}

function fecharModalExcluirConta() {
    document.getElementById('modalExcluirConta').classList.remove('active');
    document.getElementById('senhaExclusao').value = '';
}

function excluirConta() {
    const senha = document.getElementById('senhaExclusao').value;

    if (!senha) {
        showNotification('Digite sua senha para confirmar!');
        return;
    }

    // Aqui você validaria a senha com o backend
    // Simulação:
    Swal.fire({
        title: 'ATENÇÃO!',
        text: 'Tem certeza absoluta? Esta ação é IRREVERSÍVEL!',
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#dc3545',
        cancelButtonColor: '#6c757d',
        confirmButtonText: 'Sim, excluir minha conta!',
        cancelButtonText: 'Cancelar'
    }).then((result) => {
        if (result.isConfirmed) {
            // Limpar todos os dados
            localStorage.clear();
            Swal.fire({
                icon: 'success',
                title: 'Conta excluída',
                text: 'Conta excluída com sucesso. Até breve!',
                confirmButtonColor: '#370400',
                timer: 2000,
                showConfirmButton: false
            }).then(() => {
                window.location.href = 'home.html';
            });
        }
    });
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
    }).then((result) => {
        if (result.isConfirmed) {
            localStorage.setItem('userLoggedIn', 'false');
            Swal.fire({
                icon: 'success',
                title: 'Até logo!',
                text: 'Saindo...',
                confirmButtonColor: '#370400',
                timer: 1000,
                showConfirmButton: false
            }).then(() => {
                window.location.href = 'login.html';
            });
        }
    });
}

// Funções auxiliares
function formatarData(data) {
    if (!data) return '-';
    const [ano, mes, dia] = data.split('-');
    return `${dia}/${mes}/${ano}`;
}

function formatarGenero(genero) {
    const generos = {
        'feminino': 'Feminino',
        'masculino': 'Masculino',
        'outro': 'Outro',
        'prefiro-nao-informar': 'Prefiro não informar'
    };
    return generos[genero] || '-';
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