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
		Funcao:             controllers.BuscarPedidos,
		RequerAutenticacao: true,
	},
	{
		Uri:                "/pedidos/{pedidoID}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarPedido,
		RequerAutenticacao: true,
	},
	{
		Uri:                "/pedidos/{pedidoID}/receber-produtos",
		Metodo:             http.MethodPost,
		Funcao:             controllers.ConfirmarRecebimento,
		RequerAutenticacao: true,
	},
	{
		Uri:                "/pedidos/{pedidoID}/conferir-produtos",
		Metodo:             http.MethodPost,
		Funcao:             controllers.ConfirmarConferencia,
		RequerAutenticacao: true,
	},
}
