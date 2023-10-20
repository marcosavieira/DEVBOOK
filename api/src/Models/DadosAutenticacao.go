package Models

//*DadosAutenticacao contem o ID e o Token do usuario autenticado
type DadosAutenticacao struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}