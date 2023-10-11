package controllers

import (
	models "api/src/Models"
	"api/src/autenticacao"
	"api/src/banco"
	db "api/src/banco"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/security"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request){
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	var usuario models.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil{
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("cadastro"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := db.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario.ID, erro = repositorio.Criar(usuario)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusCreated, usuario)
}
func BuscarUsuarios(w http.ResponseWriter, r *http.Request){
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))
	db, erro := db.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.Buscar(nomeOuNick)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}
func BuscarUsuario(w http.ResponseWriter, r *http.Request){
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := db.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	usuario, erro := repositorio.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuario)
}
// AtualizarUsuario altera as informações de um usuário no banco
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Atualizar(usuarioID, usuario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
func DeletarUsuario(w http.ResponseWriter, r *http.Request){
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro!= nil {
        respostas.Erro(w, http.StatusUnauthorized, erro)
        return
    }
	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível deletar um usuário que não seja o seu"))
        return
	}

	db, erro := db.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Deletar(usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

//*SeguirUsuario permite seguir um usuario
func SeguirUsuario(w http.ResponseWriter, r *http.Request){
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro!= nil {
        respostas.Erro(w, http.StatusUnauthorized, erro)
        return
    }
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }
	if seguidorID == usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível seguir o mesmo usuário"))
        return
	}

	db, erro := db.Conectar()
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	defer db.Close()
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Seguir(seguidorID, usuarioID); erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }

	respostas.JSON(w, http.StatusNoContent, nil)
}

//*Parar de seguir usuario deixa de seguir um usuario
func PararDeSeguirUsuario(w http.ResponseWriter, r *http.Request){
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro!= nil {
        respostas.Erro(w, http.StatusUnauthorized, erro)
        return
    }
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }
	if seguidorID == usuarioID {
        respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível parar de seguir o mesmo usuário"))
        return
    }
	db, erro := db.Conectar()
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	defer db.Close()
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.PararDeSeguir(seguidorID, usuarioID); erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	respostas.JSON(w, http.StatusNoContent, nil)
}

//*BuscarSeguidores lista os seguidores de um usuario 
func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuario_id"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
        return
	}
	db, erro := db.Conectar()
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	seguidores, erro := repositorio.BuscarSeguidores(usuarioID)
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	respostas.JSON(w, http.StatusOK, seguidores)
}

//*BuscarSeguindo lista todos os seguidores de um usuario 
func BuscarSeguindo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuario_id"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
        return
	}
	db, erro := db.Conectar()
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.BuscarSeguindo(usuarioID)
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	respostas.JSON(w, http.StatusOK, usuarios)
}

// AtualizarSenha permite alterar a senha de um usuário
func AtualizarSenha(w http.ResponseWriter, r *http.Request) {
	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if usuarioIDNoToken != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar a senha de um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, _ := io.ReadAll(r.Body)
	var senha models.Senha
	if erro = json.Unmarshal(corpoRequisicao, &senha); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	senhaSalvaNoBanco, erro := repositorio.BuscarSenha(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = security.VerificarSenha(senhaSalvaNoBanco, senha.Atual); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("A senha atual não condiz com a que está salva no banco"))
		return
	}

	senhaComHash, erro := security.Hash(senha.Nova)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repositorio.AtualizarSenha(usuarioID, string(senhaComHash)); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}