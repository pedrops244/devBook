package controllers

import (
	"api/src/banco"
	"api/src/config"
	"api/src/respostas" // Pacote de respostas
	"encoding/json"
	"errors"
	"net/http"
)

func ConfigurarBanco(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respostas.Erro(w, http.StatusMethodNotAllowed, errors.New("método não permitido"))
		return
	}

	var dbConfig config.DBConfig
	if err := json.NewDecoder(r.Body).Decode(&dbConfig); err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("dados inválidos"))
		return
	}

	stringConexao := config.ConstruirStringConexao(dbConfig)
	if err := banco.Conectar(stringConexao); err != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao conectar ao banco: "+err.Error()))
		return
	}

	respostas.JSON(w, http.StatusOK, map[string]string{"mensagem": "Pool de conexões configurado com sucesso"})
}
