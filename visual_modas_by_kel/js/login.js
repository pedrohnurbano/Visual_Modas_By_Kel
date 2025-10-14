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
                    Swal.fire({
                        icon: 'success',
                        title: 'Bem-vindo!',
                        text: 'Login de administrador efetuado com sucesso!',
                        confirmButtonColor: '#370400',
                        timer: 1500,
                        showConfirmButton: false
                    }).then(() => {
                        window.location.href = "/painel-admin";
                    });
                } else {
                    Swal.fire({
                        icon: 'success',
                        title: 'Bem-vindo!',
                        text: 'Login efetuado com sucesso!',
                        confirmButtonColor: '#370400',
                        timer: 1500,
                        showConfirmButton: false
                    }).then(() => {
                        window.location.href = "/home";
                    });
                }
            } else {
                // Se não retornar role, assume user comum
                Swal.fire({
                    icon: 'success',
                    title: 'Bem-vindo!',
                    text: 'Login efetuado com sucesso!',
                    confirmButtonColor: '#370400',
                    timer: 1500,
                    showConfirmButton: false
                }).then(() => {
                    window.location.href = "/home";
                });
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
            Swal.fire({
                icon: 'error',
                title: 'Erro no login',
                text: mensagem,
                confirmButtonColor: '#370400'
            });
        },
        complete: function() {
            submitBtn.text(originalText).prop("disabled", false);
        }
    });
}