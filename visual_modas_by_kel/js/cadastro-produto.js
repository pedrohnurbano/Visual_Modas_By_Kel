// Variável global para armazenar a imagem em base64
let imagemBase64 = "";

// Aguardar DOM carregar
$(document).ready(function () {
    // Handler para quando o usuário seleciona uma imagem
    $("#foto-produto").on("change", function (e) {
        const arquivo = e.target.files[0];

        if (!arquivo) {
            imagemBase64 = "";
            $("#preview-imagem").html("");
            return;
        }

        // Validar tipo de arquivo
        const tiposPermitidos = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp'];
        if (!tiposPermitidos.includes(arquivo.type)) {
            alert("Por favor, selecione apenas arquivos JPG, PNG ou WEBP!");
            $(this).val('');
            imagemBase64 = "";
            $("#preview-imagem").html("");
            return;
        }

        // Validar tamanho (máximo 5MB)
        const tamanhoMaximo = 5 * 1024 * 1024; // 5MB
        if (arquivo.size > tamanhoMaximo) {
            alert("A imagem deve ter no máximo 5MB!");
            $(this).val('');
            imagemBase64 = "";
            $("#preview-imagem").html("");
            return;
        }

        // Converter para base64
        const reader = new FileReader();
        reader.onload = function (event) {
            imagemBase64 = event.target.result;

            // Mostrar preview da imagem
            $("#preview-imagem").html(`
                <img src="${imagemBase64}" alt="Preview" style="max-width: 200px; max-height: 200px; margin-top: 10px; border-radius: 8px;">
            `);

            console.log("Imagem carregada com sucesso!");
        };

        reader.onerror = function () {
            alert("Erro ao carregar a imagem. Tente novamente.");
            $("#foto-produto").val('');
            imagemBase64 = "";
            $("#preview-imagem").html("");
        };

        reader.readAsDataURL(arquivo);
    });

    // Handler do formulário
    $("#produto-form").on("submit", function (evento) {
        evento.preventDefault();
        criarProduto();
    });
});

function criarProduto() {
    // Validações
    const nome = $("#nome-produto").val().trim();
    const descricao = $("#descricao-produto").val().trim();
    const precoString = $("#preco-produto").val();
    const tamanho = $("#tamanho-produto").val();
    const categoria = $("#categoria-produto").val();
    const secao = $("#secao-produto").val();

    if (!nome) {
        alert("O nome do produto é obrigatório!");
        return;
    }

    if (!descricao) {
        alert("A descrição do produto é obrigatória!");
        return;
    }

    const preco = parseFloat(precoString);
    if (!preco || preco <= 0 || isNaN(preco)) {
        alert("Por favor, informe um preço válido!");
        return;
    }

    if (!tamanho) {
        alert("Por favor, selecione um tamanho!");
        return;
    }

    if (!categoria) {
        alert("Por favor, selecione uma categoria!");
        return;
    }

    if (!imagemBase64) {
        alert("Por favor, selecione uma foto do produto!");
        return;
    }

    if (!secao) {
        alert("Por favor, selecione uma seção!");
        return;
    }

    // Adicionar loading state
    const submitBtn = $('.painel-btn-primary[form="produto-form"]');
    const originalText = submitBtn.text();
    submitBtn.text("Cadastrando...").prop("disabled", true);

    // Preparar dados para envio
    const dadosProduto = {
        nome: nome,
        descricao: descricao,
        preco: preco, // Agora é número
        tamanho: tamanho,
        categoria: categoria,
        secao: secao,  // NOVO CAMPO
        foto_url: imagemBase64
    };

    console.log("Enviando produto:", dadosProduto);

    $.ajax({
        url: "/api/produtos",
        method: "POST",
        contentType: "application/json",
        data: JSON.stringify(dadosProduto),
        success: function (response) {
            alert("Produto cadastrado com sucesso!");

            // Limpar formulário
            $("#produto-form")[0].reset();
            imagemBase64 = "";
            $("#preview-imagem").html("");

            // Fechar modal
            closeModal('produto-modal');

            // Atualizar tabela de produtos se estiver na página de produtos
            if (typeof carregarProdutos === 'function') {
                carregarProdutos();
            }
        },
        error: function (xhr, status, error) {
            console.error("Erro ao cadastrar produto:", error);
            console.log("Response:", xhr.responseText);

            let mensagem = "Erro ao cadastrar produto. Tente novamente.";
            if (xhr.responseText) {
                try {
                    const resposta = JSON.parse(xhr.responseText);
                    mensagem = resposta.erro || resposta.message || mensagem;
                } catch (e) {
                    mensagem = xhr.responseText || mensagem;
                }
            }
            alert(mensagem);
        },
        complete: function () {
            submitBtn.text(originalText).prop("disabled", false);
        }
    });
}

// Função para buscar produtos (para outras páginas)
function buscarProdutos(filtro = "") {
    $.ajax({
        url: `/api/produtos?filtro=${filtro}`,
        method: "GET",
        success: function (produtos) {
            exibirProdutos(produtos);
        },
        error: function (xhr, status, error) {
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
        success: function (produtos) {
            exibirMeusProdutos(produtos);
        },
        error: function (xhr, status, error) {
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
        success: function () {
            alert("Produto deletado com sucesso!");
            buscarMeusProdutos(); // Recarregar lista
        },
        error: function (xhr, status, error) {
            console.error("Erro ao deletar produto:", error);
            let mensagem = "Erro ao deletar produto.";
            if (xhr.responseText) {
                try {
                    const resposta = JSON.parse(xhr.responseText);
                    mensagem = resposta.erro || mensagem;
                } catch (e) { }
            }
            alert(mensagem);
        }
    });
}

// Função auxiliar para exibir produtos
function exibirProdutos(produtos) {
    const tbody = document.getElementById('produtos-list');
    if (!tbody) return;

    tbody.innerHTML = '';

    if (!produtos || produtos.length === 0) {
        tbody.innerHTML = '<tr><td colspan="11" style="text-align: center;">Nenhum produto encontrado</td></tr>';
        return;
    }

    produtos.forEach(produto => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>${produto.id || ''}</td>
            <td>${produto.secao || '-'}</td>  
            <td>${produto.nome || ''}</td>
            <td>R$ ${(produto.preco || 0).toFixed(2)}</td>
            <td>${produto.tamanho || ''}</td>
            <td>-</td>
            <td>${produto.categoria || ''}</td>
            <td>-</td>
            <td>${produto.descricao ? produto.descricao.substring(0, 50) + '...' : ''}</td>
            <td>${produto.foto_url ? '<img src="' + (produto.foto_url.startsWith('http') ? produto.foto_url : (produto.foto_url.startsWith('/') ? 'http://localhost:5000' + produto.foto_url : 'http://localhost:5000/' + produto.foto_url)) + '" style="width: 50px; height: 50px; object-fit: cover; border-radius: 4px;">' : '-'}</td>
            <td>
                <button class="painel-btn-secondary" style="padding: 5px 10px; font-size: 12px;" onclick="editarProduto(${produto.id})">Editar</button>
                <button class="painel-btn-danger" style="padding: 5px 10px; font-size: 12px; background: #dc3545; color: white; border: none; border-radius: 4px; cursor: pointer;" onclick="deletarProduto(${produto.id})">Excluir</button>
            </td>
        `;
        tbody.appendChild(tr);
    });
}

// Função auxiliar para exibir meus produtos
function exibirMeusProdutos(produtos) {
    // Esta função será customizada conforme necessário
    console.log("Meus produtos:", produtos);
    exibirProdutos(produtos); // Por enquanto, usa a mesma função
}

// Função para editar produto (placeholder)
function editarProduto(produtoId) {
    alert('Funcionalidade de edição em desenvolvimento. ID: ' + produtoId);
}

// Carregar produtos ao entrar na seção de produtos
$(document).ready(function () {
    // Se estiver na página de produtos, carregar automaticamente
    if (window.location.pathname.includes('painel-admin')) {
        // Adicionar listener para quando clicar em Produtos no menu
        $('[data-section="produtos"]').on('click', function () {
            setTimeout(carregarProdutos, 100);
        });
    }
});
