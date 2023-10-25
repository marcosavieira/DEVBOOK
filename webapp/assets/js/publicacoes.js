$("#nova-publicacao").on("submit", criarPublicacao);
$(document).on("click", ".curtir-publicacao", curtirPublicacao);
$(document).on("click", ".descurtir-publicacao", descurtirPublicacao);
$('#atualizar-publicacao').on("click", atualizarPublicacao);
$('.deletar-publicacao').on("click", deletarPublicacao);
function criarPublicacao(evento) {
  evento.preventDefault();

  $.ajax({
    url: "/publicacoes",
    method: "POST",
    data: {
      titulo: $("#titulo").val(),
      conteudo: $("#conteudo").val(),
    },
  })
    .done(function () {
      window.location = "/home";
    })
    .fail(function () {
      Swal.fire("Erro", "Erro ao criar publicação!Tente novamente", "error");
    });
}

function curtirPublicacao(evento) {
  evento.preventDefault();

  const elementoClicado = $(evento.target)
  const publicacaoID = elementoClicado.closest('div').data('publicacao-id');

  elementoClicado.prop('disabled', true)
  $.ajax({
    url: `/publicacoes/${publicacaoID}/curtir`,
    method: 'POST',
  }).done(function() {
    const conatadorDeCurtidas = elementoClicado.next('span');
    const quantidadeDeCurtidas = parseInt(conatadorDeCurtidas.text());

    conatadorDeCurtidas.text(quantidadeDeCurtidas + 1);

    elementoClicado.addClass('descurtir-publicacao');
    elementoClicado.addClass('text-danger');
    elementoClicado.removeClass('curtir-publicacao');
  }).fail(function() {
    Swal.fire("Erro", "Erro ao curtir a publicação!Tente novamente...", "error");
  }).always(function() {
    elementoClicado.prop('disabled', false)
  })
}
function descurtirPublicacao(evento) {
  evento.preventDefault();

  const elementoClicado = $(evento.target)
  const publicacaoID = elementoClicado.closest('div').data('publicacao-id');

  elementoClicado.prop('disabled', true);
  $.ajax({
    url: `/publicacoes/${publicacaoID}/descurtir`,
    method: 'POST',
  }).done(function() {
    const conatadorDeCurtidas = elementoClicado.next('span');
    const quantidadeDeCurtidas = parseInt(conatadorDeCurtidas.text());

    conatadorDeCurtidas.text(quantidadeDeCurtidas - 1);

    elementoClicado.removeClass('descurtir-publicacao');
    elementoClicado.removeClass('text-danger');
    elementoClicado.addClass('curtir-publicacao');
  }).fail(function() {
    Swal.fire("Erro", "Erro ao descurtir!Tente novamente...", "error");
  }).always(function() {
    elementoClicado.prop('disabled', false)
  })
}

function atualizarPublicacao(evento){
  $(this).prop('disabled', true);

  const publicacaoID = $(this).data('publicacao-id');
  
  $.ajax({
    url: `/publicacoes/${publicacaoID}`,
    method: 'PUT',
    data: {
      titulo: $('#titulo').val(),
      conteudo: $('#conteudo').val()
    }
  }).done(function() {
    Swal.fire('Sucesso!','Publicação atualizada com sucesso','success')
    .then(function() {
     window.location = "/home"; 
    })
  }).fail(function() {
    Swal.fire("Erro", "Erro ao editar publicação!Tente novamente...", "error");
  }).always(function() {
   $('#atualizar-publicacao').prop('disabled', false); 
  })
}

function deletarPublicacao(evento) {
  evento.preventDefault();
  Swal.fire({
    title: "Atenção!",
    text: "Tem certeza que deseja excluir a publicação? Essa ação é irreversível!",
    showCancelButton: true,
    cancelButtonText: "Cancelar",
    icon: "warning"
  }).then(function(confirmacao) {
    confirmacao.value;
    if(!confirmacao.value) return;
    const elementoClicado = $(evento.target)
  const publicacao = elementoClicado.closest('div')
  const publicacaoID = publicacao.data('publicacao-id');

  elementoClicado.prop('disabled', true);

  $.ajax({
    url: `/publicacoes/${publicacaoID}`,
    method: 'DELETE',
  }).done(function() {
    Swal.fire("Sucesso!", "Publicação Excluída!", "success");
    publicacao.fadeOut("slow", function() {
      $(this).remove();
    });
  }).fail( function() {
    Swal.fire("Erro!", "Erro ao excluir a publicação!", "error");
  });
  })
  
}