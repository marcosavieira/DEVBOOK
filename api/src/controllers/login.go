package controllers

import (
	"api/src/Models"
	models "api/src/Models"
	"api/src/autenticacao"
	db "api/src/banco"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/security"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func Login(w http.ResponseWriter, r *http.Request){
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		return
	}

	db, erro := db.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
			return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarioSalvoNobanco, erro := repositorio.BuscarPorEmail(usuario.Email)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	if erro = security.VerificarSenha(usuarioSalvoNobanco.Senha, usuario.Senha); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	token, erro := autenticacao.CriarToken(usuarioSalvoNobanco.ID)
	 if erro != nil{
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	 }

	 usuarioID := strconv.FormatUint(usuarioSalvoNobanco.ID, 10)

	 
	 respostas.JSON(w, http.StatusOK, Models.DadosAutenticacao{ID: usuarioID, Token: token})
}