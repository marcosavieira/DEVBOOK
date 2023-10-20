package controllers

import (
	models "api/src/Models"
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//*CriarPublicao permite o usuario criar uma publicacao
func CriarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }

	var publicacao models.Publicacao
	if erro = json.Unmarshal(corpoRequisicao, &publicacao); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
        return
	}

	publicacao.AutorID = usuarioID

	if erro = publicacao.Preparar(); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacao.ID, erro = repositorio.Criar(publicacao)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
        return
	}

	respostas.JSON(w, http.StatusCreated, publicacao)
}

//*BuscarPublicacoes permite o usuario buscar todas publicacoes no feed do usuario
func BuscarPublicacoes(w http.ResponseWriter, r *http.Request)   {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro!= nil {
        respostas.Erro(w, http.StatusUnauthorized, erro)
        return
    }

	db, erro := banco.Conectar()
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacoes, erro := repositorio.Buscar(usuarioID)
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }

	respostas.JSON(w, http.StatusOK, publicacoes)
}

//*BuscarPublicacao permite o usuario buscar uma publicacao
func BuscarPublicacao(w http.ResponseWriter, r *http.Request)    {
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10 , 64)
	if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }

	db, erro := banco.Conectar()
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacao, erro := repositorio.BuscarPorID(publicacaoID)
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	respostas.JSON(w, http.StatusOK, publicacao)
}

//*AtualizarPublicacao permite o usuario atualizar uma determinada publicacao
func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro!= nil {
        respostas.Erro(w, http.StatusUnauthorized, erro)
        return
    }
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }
	db, erro := banco.Conectar()
    if erro!= nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacaoSalvaNoBanco, erro := repositorio.BuscarPorID(publicacaoID)
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	if publicacaoSalvaNoBanco.AutorID!= usuarioID {
        respostas.Erro(w, http.StatusForbidden, errors.New(`Não é possivel atualizar uma publicação que não é sua!`))

        return
    }
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }
	var publicacao models.Publicacao
	if erro = json.Unmarshal(corpoRequisicao, &publicacao); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	if erro = publicacao.Preparar(); erro!= nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
        return
	}
	if erro = repositorio.Atualizar(publicacaoID, publicacao); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)
}
//*DeletarPublicacao permite o usuario deletar uma publicacao
func DeletarPublicacao(w http.ResponseWriter, r *http.Request)   {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro!= nil {
        respostas.Erro(w, http.StatusUnauthorized, erro)
        return
    }
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }
	db, erro := banco.Conectar()
    if erro!= nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacaoSalvaNoBanco, erro := repositorio.BuscarPorID(publicacaoID)
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	if publicacaoSalvaNoBanco.AutorID!= usuarioID {
        respostas.Erro(w, http.StatusForbidden, errors.New(`Não é possivel deletar uma publicação que não é sua!`))

        return
    }
	if erro = repositorio.Deletar(publicacaoID); erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	respostas.JSON(w, http.StatusNoContent, nil)
}

//*BuscarPublicacoesPorUsuario permite o usuario buscar todas publicacoes do usuario
func BuscarPublicacoesPorUsuario(w http.ResponseWriter, r *http.Request){
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro!= nil {
        respostas.Erro(w, http.StatusBadRequest, erro)
        return
    }
	db, erro := banco.Conectar()
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	defer db.Close()
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacoes, erro := repositorio.BuscarPorUsuario(usuarioID)
	if erro!= nil {
        respostas.Erro(w, http.StatusInternalServerError, erro)
        return
    }
	respostas.JSON(w, http.StatusOK, publicacoes)

}

//*CurtirPublicacao permite ao usuario curtir uma publicacao no feed
func CurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
        return
	}
	db, erro := banco.Conectar()
    if erro!= nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	if erro = repositorio.Curtir(publicacaoID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
        return
	}
	respostas.JSON(w, http.StatusNoContent, nil)
}

//*DescurtirPublicacao permite ao usuario descurtir uma publicacao no feed
func DescurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
        return
	}
	db, erro := banco.Conectar()
    if erro!= nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	if erro = repositorio.Descurtir(publicacaoID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
        return
	}
	respostas.JSON(w, http.StatusNoContent, nil)

}