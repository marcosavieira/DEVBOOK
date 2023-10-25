package controllers

import (
	"net/http"
	"webapp/src/cookies"
)

// *Fazer logout remove os dados de autenticação do usuario salvo no browser
func FazerLogout(w http.ResponseWriter, r *http.Request) {
	cookies.Deletar(w)
	http.Redirect(w, r, "/login", 302)
}
