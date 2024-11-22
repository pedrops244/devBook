package modelos

type DadosAutenticacao struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Role  string `json:"role"`
}
