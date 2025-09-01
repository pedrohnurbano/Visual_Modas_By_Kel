$("#formulario-cadastro").on("submit", criarUsuario); //quando rolar um submit nesse formulario, a funcao criarUsuario é chamada

function criarUsuario(evento) {
    evento.preventDefault(); //previne o comportamento padrao do formulario ao ser enviado

    if ($("#senha").val() != $("#confirmar-senha").val()) {
        console.log("As senhas não coincidem!");
        return;
    }

    $.ajax({
        url: "/usuarios",
        method: "POST",
        data: {
            nome: $("#nome").val(),
            sobrenome: $("#sobrenome").val(),
            email: $("#email").val(),
            senha: $("#senha").val(),
            telefone: $("#telefone").val(),
            cpf: $("#cpf").val(),
        }
    });
};

