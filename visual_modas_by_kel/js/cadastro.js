$("#formulario-cadastro").on("submit", criarUsuario);

function criarUsuario(evento) {
    evento.preventDefault();

    if ($("#senha").val() != $("#confirmar-senha").val()) {
        alert("As senhas não coincidem!");
        return;
    }

    // Funções para limpar os campos
    function limparCPF(cpf) {
        return cpf.replace(/[^\d]/g, ''); // Remove tudo exceto dígitos
    }

    function limparTelefone(telefone) {
        return telefone.replace(/[^\d]/g, ''); // Remove tudo exceto dígitos
    }

    // Validações básicas
    const cpfLimpo = limparCPF($("#cpf").val());
    const telefoneLimpo = limparTelefone($("#telefone").val());

    if (cpfLimpo.length !== 11) {
        alert("CPF deve ter 11 dígitos!");
        return;
    }

    if (telefoneLimpo.length < 10 || telefoneLimpo.length > 11) {
        alert("Telefone deve ter 10 ou 11 dígitos!");
        return;
    }

    // Adicionar loading state
    const submitBtn = $("#formulario-cadastro button[type='submit']");
    const originalText = submitBtn.text();
    submitBtn.text("Criando conta...").prop("disabled", true);

    $.ajax({
        url: "/usuarios",
        method: "POST",
        data: {
            nome: $("#nome").val().trim(),
            sobrenome: $("#sobrenome").val().trim(),
            email: $("#email").val().trim(),
            senha: $("#senha").val(),
            telefone: telefoneLimpo,
            cpf: cpfLimpo,
        },
        success: function(response) {
            alert("Usuário criado com sucesso!");
            window.location.href = "/login";
        },
        error: function(xhr, status, error) {
            console.error("Erro ao criar usuário:", error);
            console.log("Response:", xhr.responseText);
            
            let mensagem = "Erro ao criar usuário. Tente novamente.";
            if (xhr.responseText) {
                try {
                    const resposta = JSON.parse(xhr.responseText);
                    mensagem = resposta.message || mensagem;
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
};