$("#form-sign-up").on("submit", criarUsuario);
function criarUsuario(e) {
  e.preventDefault();

  if ($("#senha").val() !== $("#confirmar-senha").val()) {
    alert("As senhas não conferem");
    return;
  }

  $.ajax({
    url: "/usuarios",
    method: "POST",
    data: {
      nome: $("#nome").val(),
      email: $("#email").val(),
      nick: $("#nick").val(),
      senha: $("#senha").val(),
    },
  })
    .done(function () {
      alert("Usuário criado com sucesso!");
    })
    .fail(function (erro) {
      console.log(erro);
      alert("Erro ao criar usuário!");
    });
}
