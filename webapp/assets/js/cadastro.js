$("#form-sign-up").on("submit", criarUsuario);
function criarUsuario(e) {
  e.preventDefault();

  if ($("#senha").val() !== $("#confirmar-senha").val()) {
    Swal.fire("Ops...", "As senhas não conferem!", "error");
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
      Swal.fire("Sucesso", "Usuário criado com sucesso!", "success")
      .then(function () {
        $.ajax({
          url: "/login",
          method: "POST",
          data: {
            email: $("#email").val(),
            senha: $("#senha").val()
          }
        }).done(function() {
          window.location = "/home";
        }).fail(function() {
          Swal.fire("Erro", "Erro ao autenticar o usuario!, tente fazer o login novamente", "error");
          window.location = "/login"
        })
      })
    })
    .fail(function () {
      Swal.fire("Erro", "Erro ao criar usuário!", "error");
    });
}
