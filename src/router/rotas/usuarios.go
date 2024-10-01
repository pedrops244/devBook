package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasUsuarios = []Rota{
	{
		Uri:              "/usuarios",
		Metodo:           http.MethodPost,
		Funcao:           controllers.CriarUsuario,
		RequerAutenticao: false,
	},
	{
		Uri:              "/usuarios",
		Metodo:           http.MethodGet,
		Funcao:           controllers.BuscarUsuarios,
		RequerAutenticao: false,
	},
	{
		Uri:              "/usuarios/{usuarioId}",
		Metodo:           http.MethodGet,
		Funcao:           controllers.BuscarUsuario,
		RequerAutenticao: false,
	},
	{
		Uri:              "/usuarios/{usuarioId}",
		Metodo:           http.MethodPut,
		Funcao:           controllers.AtualizarUsuario,
		RequerAutenticao: false,
	},
	{
		Uri:              "/usuarios/{usuarioId}",
		Metodo:           http.MethodDelete,
		Funcao:           controllers.DeletarUsuario,
		RequerAutenticao: false,
	},
}
