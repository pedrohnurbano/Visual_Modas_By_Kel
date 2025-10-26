
// Criar um novo pedido
function criarPedidoAPI(dadosPedido, callback) {
    $.ajax({
        url: "/api/pedidos",
        method: "POST",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        data: JSON.stringify(dadosPedido),
        contentType: "application/json",
        success: function(response) {
            console.log("Pedido criado com sucesso:", response);
            if (callback) callback(true, response);
        },
        error: function(xhr) {
            console.error("Erro ao criar pedido:", xhr);
            let mensagem = "Erro ao criar pedido. Tente novamente.";
            
            if (xhr.status === 401) {
                mensagem = "Sessão expirada. Faça login novamente.";
                window.location.href = '/login';
            } else if (xhr.responseJSON && xhr.responseJSON.erro) {
                mensagem = xhr.responseJSON.erro;
            }
            
            if (callback) callback(false, mensagem);
        }
    });
}

// Buscar pedidos do usuário
function buscarPedidosUsuarioAPI(callback) {
    $.ajax({
        url: "/api/pedidos",
        method: "GET",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        success: function(pedidos) {
            console.log("Pedidos do usuário:", pedidos);
            if (callback) callback(pedidos || []);
        },
        error: function(xhr) {
            console.error("Erro ao buscar pedidos:", xhr);
            if (xhr.status === 401) {
                window.location.href = '/login';
            }
            if (callback) callback([]);
        }
    });
}

// Buscar um pedido específico
function buscarPedidoAPI(pedidoId, callback) {
    $.ajax({
        url: `/api/pedidos/${pedidoId}`,
        method: "GET",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        success: function(pedido) {
            console.log("Pedido encontrado:", pedido);
            if (callback) callback(pedido);
        },
        error: function(xhr) {
            console.error("Erro ao buscar pedido:", xhr);
            if (xhr.status === 401) {
                window.location.href = '/login';
            }
            if (callback) callback(null);
        }
    });
}

// Listar todos os pedidos (admin)
function listarTodosPedidosAPI(callback) {
    $.ajax({
        url: "/api/admin/pedidos",
        method: "GET",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        success: function(pedidos) {
            console.log("Todos os pedidos:", pedidos);
            if (callback) callback(pedidos || []);
        },
        error: function(xhr) {
            console.error("Erro ao listar pedidos:", xhr);
            if (xhr.status === 401 || xhr.status === 403) {
                alert("Acesso negado. Apenas administradores podem acessar.");
                window.location.href = '/home';
            }
            if (callback) callback([]);
        }
    });
}

// Atualizar status do pedido (admin)
function atualizarStatusPedidoAPI(pedidoId, status, callback) {
    $.ajax({
        url: `/api/admin/pedidos/${pedidoId}/status`,
        method: "PUT",
        headers: {
            "Authorization": `Bearer ${localStorage.getItem("token")}`
        },
        data: JSON.stringify({
            status: status
        }),
        contentType: "application/json",
        success: function(response) {
            console.log("Status atualizado:", response);
            if (callback) callback(true);
        },
        error: function(xhr) {
            console.error("Erro ao atualizar status:", xhr);
            if (xhr.status === 401 || xhr.status === 403) {
                alert("Acesso negado.");
            }
            if (callback) callback(false);
        }
    });
}

