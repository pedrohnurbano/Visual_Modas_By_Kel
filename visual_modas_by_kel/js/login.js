$("#formulario-login").on("submit", fazerLogin);

function fazerLogin(evento) {
    evento.preventDefault();

    // Adicionar loading state
    const submitBtn = $("#formulario-login button[type='submit']");
    const originalText = submitBtn.text();
    submitBtn.text("Entrando...").prop("disabled", true);

    $.ajax({
        url: "/login",
        method: "POST",
        data: {
            email: $("#loginEmail").val(),
            senha: $("#loginPassword").val(),
        },
        success: function(response) {
            // Armazenar a role no sessionStorage
            if (response && response.role) {
                sessionStorage.setItem("userRole", response.role);
                sessionStorage.setItem("userId", response.id);
                
                // Redirecionar baseado na role
                if (response.role === "admin") {
                    alert("Login de administrador efetuado com sucesso!");
                    window.location.href = "/painel-admin";
                } else {
                    alert("Login efetuado com sucesso!");
                    window.location.href = "/home";
                }
            } else {
                // Se não retornar role, assume user comum
                alert("Login efetuado com sucesso!");
                window.location.href = "/home";
            }
        },
        error: function(xhr, status, error) {
            console.error("Erro ao fazer o login:", error);
            console.log("Response:", xhr.responseText);
            
            let mensagem = "Erro ao fazer o login. Tente novamente.";
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