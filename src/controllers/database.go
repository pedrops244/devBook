package controllers

import (
	"api/src/banco"
	"api/src/config"
	"encoding/json"
	"net/http"
)

func ConfigurarBanco(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var dbConfig config.DBConfig
	err := json.NewDecoder(r.Body).Decode(&dbConfig)
	if err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	stringConexao := config.ConstruirStringConexao(dbConfig)

	_, err = banco.Conectar(stringConexao)
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Conexão configurada com sucesso"))
}
