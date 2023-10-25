$("#form-sign-in").on("submit", fazerLogin);
function fazerLogin(evento) {
  evento.preventDefault();

  $.ajax({
    url: "/login",
    method: "POST",
    data: {
      email: $("#emailS").val(),
      senha: $("#senhaS").val(),
    },
  })
    .done(function () {
      Swal.fire("Sucesso", "Login efetuado com sucesso!", "success")
      .then(function () {
        window.location = "/home";
      })
    })
    .fail(function () {
      Swal.fire("Ops...", "Erro ao realizar o login!", "error")
    });
}
