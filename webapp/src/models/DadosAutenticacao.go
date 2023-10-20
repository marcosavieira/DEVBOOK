package models

//*DadosAutenticacao contem o id e o token do usuario logado
type DadosAutenticacao struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}