package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotaDatabase = Rota{
	Uri:                "/conectar-database",
	Metodo:             http.MethodPost,
	Funcao:             controllers.ConfigurarBanco,
	RequerAutenticacao: false,
}
