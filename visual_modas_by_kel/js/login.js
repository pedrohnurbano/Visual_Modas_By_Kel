$("#formulario-login").on("submit", fazerLogin);

function fazerLogin(evento) {
    evento.preventDefault();

    $.ajax({
        url: "/login",
        method: "POST",
        data: {
            email: $("#loginEmail").val(),
            senha: $("#loginPassword").val(),
        },
        success: function(response) {
            alert("Login efetuado com sucesso!");
            window.location.href = "/home";
        },
        error: function(xhr, status, error) {
            console.error("Erro ao fazer o login:", error);
            console.log("Response:", xhr.responseText);
            
            let mensagem = "Erro fazer o login. Tente novamente.";
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
