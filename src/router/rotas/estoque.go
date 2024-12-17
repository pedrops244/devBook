package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotaEstoque = Rota{
	Uri:                "/estoque/validar-produto",
	Metodo:             http.MethodPost,
	Funcao:             controllers.VerificarProduto,
	RequerAutenticacao: true,
}
