package Models

//*Senha representa um formato de senha
type Senha struct {
	Nova string `json:"nova"`
	Atual string `json:"atual"`
}