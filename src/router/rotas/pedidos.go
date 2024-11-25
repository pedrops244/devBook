package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasPedidos = []Rota{
	{
		Uri:                "/pedidos",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarPedido,
		RequerAutenticacao: true,
	},
	{
		Uri:                "/pedidos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.ListarPedidos,
		RequerAutenticacao: true,
	},
	{
		Uri:                "/pedidos/{pedidoID}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarPedido,
		RequerAutenticacao: true,
	},
	{
		Uri:                "/pedidos/{pedidoID}/status",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarStatusPedido,
		RequerAutenticacao: true,
	},
}
