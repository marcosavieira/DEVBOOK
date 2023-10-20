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
      window.location = "/home";
    })
    .fail(function () {
      alert("Falha ao realizar login");
    });
}
