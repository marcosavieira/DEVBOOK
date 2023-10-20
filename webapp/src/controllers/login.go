package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"webapp/src/config"
	"webapp/src/cookies"
	"webapp/src/models"
	"webapp/src/respostas"
)

// FazerLogin utiliza o e-mail e senha do usuário para autenticar na aplicação
func FazerLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	usuario, erro := json.Marshal(map[string]string{
		"email": r.FormValue("email"),
		"senha": r.FormValue("senha"),
	})
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/login", config.APIURL)
	usuarioJson := bytes.NewBuffer(usuario)
	response, erro := http.Post(url, "application/json", usuarioJson)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeError(w, response)
		return
	}
	/* body, err := io.ReadAll(response.Body)
	if err != nil {
    // Handle the error reading the response body
	} else {
    fmt.Println(string(body))
	} */

	decoder := json.NewDecoder(response.Body)
var dadosAutenticacao models.DadosAutenticacao
if err := decoder.Decode(&dadosAutenticacao); err != nil {
    respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
/* } else {
    fmt.Println("ID:", dadosAutenticacao.ID)
    fmt.Println("Token:", dadosAutenticacao.Token)
} */
}

	//var dadosAutenticacao models.DadosAutenticacao
	/* if erro = json.NewDecoder(response.Body).Decode(&dadosAutenticacao); erro != nil {
		fmt.Println("Esta entrando neste erro!")
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	} */

	if erro = cookies.Salvar(w, dadosAutenticacao.ID, dadosAutenticacao.Token); erro != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	response.Body.Close()
	respostas.JSON(w, http.StatusOK, nil)

}