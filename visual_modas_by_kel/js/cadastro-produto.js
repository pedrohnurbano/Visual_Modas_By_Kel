$("#produto-form").on("submit", criarProduto);

// Variável global para armazenar a imagem em base64
let imagemBase64 = "";

// Handler para quando o usuário seleciona uma imagem
$("#foto-produto").on("change", function(e) {
    const arquivo = e.target.files[0];
    
    if (!arquivo) {
        return;
    }

    // Validar tipo de arquivo
    const tiposPermitidos = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp'];
    if (!tiposPermitidos.includes(arquivo.type)) {
        alert("Por favor, selecione apenas arquivos JPG, PNG ou WEBP!");
        $(this).val('');
        return;
    }

    // Validar tamanho (máximo 5MB)
    const tamanhoMaximo = 5 * 1024 * 1024; // 5MB
    if (arquivo.size > tamanhoMaximo) {
        alert("A imagem deve ter no máximo 5MB!");
        $(this).val('');
        return;
    }

    // Converter para base64
    const reader = new FileReader();
    reader.onload = function(event) {
        imagemBase64 = event.target.result;
        
        // Mostrar preview da imagem
        $("#preview-imagem").html(`
            <img src="${imagemBase64}" alt="Preview" style="max-width: 200px; max-height: 200px; margin-top: 10px; border-radius: 8px;">
        `);
        
        console.log("Imagem carregada com sucesso!");
    };
    
    reader.onerror = function() {
        alert("Erro ao carregar a imagem. Tente novamente.");
        $("#foto-produto").val('');
    };
    
    reader.readAsDataURL(arquivo);
});

function criarProduto(evento) {
    evento.preventDefault();

    // Validações
    if (!$("#nome-produto").val().trim()) {
        alert("O nome do produto é obrigatório!");
        return;
    }

    if (!$("#descricao-produto").val().trim()) {
        alert("A descrição do produto é obrigatória!");
        return;
    }

    const preco = parseFloat($("#preco-produto").val());
    if (!preco || preco <= 0) {
        alert("Por favor, informe um preço válido!");
        return;
    }

    if (!$("#tamanho-produto").val()) {
        alert("Por favor, selecione um tamanho!");
        return;
    }

    if (!$("#categoria-produto").val()) {
        alert("Por favor, selecione uma categoria!");
        return;
    }

    if (!imagemBase64) {
        alert("Por favor, selecione uma foto do produto!");
        return;
    }

    // Adicionar loading state
    const submitBtn = $("#produto-form").closest('.painel-modal-content').find('.painel-btn-primary[type="submit"]');
    const originalText = submitBtn.text();
    submitBtn.text("Cadastrando...").prop("disabled", true);

    $.ajax({
        url: "/api/produtos",
        method: "POST",
        data: {
            nome: $("#nome-produto").val().trim(),
            descricao: $("#descricao-produto").val().trim(),
            preco: preco,
            tamanho: $("#tamanho-produto").val(),
            categoria: $("#categoria-produto").val(),
            foto_url: imagemBase64
        },
        success: function(response) {
            alert("Produto cadastrado com sucesso!");
            
            // Limpar formulário
            $("#produto-form")[0].reset();
            imagemBase64 = "";
            $("#preview-imagem").html("");
            
            // Fechar modal
            closeModal('produto-modal');
            
            // Atualizar tabela de produtos
            if (typeof carregarProdutos === 'function') {
                carregarProdutos();
            }
        },
        error: function(xhr, status, error) {
            console.error("Erro ao cadastrar produto:", error);
            console.log("Response:", xhr.responseText);
            
            let mensagem = "Erro ao cadastrar produto. Tente novamente.";
            if (xhr.responseText) {
                try {
                    const resposta = JSON.parse(xhr.responseText);
                    mensagem = resposta.erro || resposta.message || mensagem;
                } catch (e) {
                    // Se não conseguir fazer parse do JSON, usa mensagem padrão
                }
            }
            alert(mensagem);
        },
        complete: function() {
            submitBtn.text(originalText).prop("disabled", false);
        }
    });
}

// Função para buscar produtos (para outras páginas)
function buscarProdutos(filtro = "") {
    $.ajax({
        url: `/api/produtos?filtro=${filtro}`,
        method: "GET",
        success: function(produtos) {
            exibirProdutos(produtos);
        },
        error: function(xhr, status, error) {
            console.error("Erro ao buscar produtos:", error);
            alert("Erro ao carregar produtos.");
        }
    });
}

// Função para buscar meus produtos
function buscarMeusProdutos() {
    $.ajax({
        url: "/api/meus-produtos",
        method: "GET",
        success: function(produtos) {
            exibirMeusProdutos(produtos);
        },
        error: function(xhr, status, error) {
            console.error("Erro ao buscar produtos:", error);
            alert("Erro ao carregar seus produtos.");
        }
    });
}

// Função para deletar produto
function deletarProduto(produtoId) {
    if (!confirm("Tem certeza que deseja deletar este produto?")) {
        return;
    }

    $.ajax({
        url: `/api/produtos/${produtoId}`,
        method: "DELETE",
        success: function() {
            alert("Produto deletado com sucesso!");
            buscarMeusProdutos(); // Recarregar lista
        },
        error: function(xhr, status, error) {
            console.error("Erro ao deletar produto:", error);
            let mensagem = "Erro ao deletar produto.";
            if (xhr.responseText) {
                try {
                    const resposta = JSON.parse(xhr.responseText);
                    mensagem = resposta.erro || mensagem;
                } catch (e) {}
            }
            alert(mensagem);
        }
    });
}

// Função auxiliar para exibir produtos (seu amigo vai implementar o HTML)
function exibirProdutos(produtos) {
    // Esta função será customizada pelo seu amigo para exibir os produtos na tela
    console.log("Produtos carregados:", produtos);
}

// Função auxiliar para exibir meus produtos
function exibirMeusProdutos(produtos) {
    // Esta função será customizada pelo seu amigo para exibir os produtos na tela
    console.log("Meus produtos:", produtos);
}
