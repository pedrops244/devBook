package controllers

import (
	"api/src/banco"
	"api/src/repositorios"
	"api/src/respostas"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func VerificarProduto(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}
	defer r.Body.Close()

	var produto struct {
		Codigo     string `json:"codigo"`
		Quantidade int    `json:"quantidade"`
	}
	if err := json.Unmarshal(body, &produto); err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("dados inválidos"))
		return
	}

	if produto.Codigo == "" || produto.Quantidade <= 0 {
		respostas.Erro(w, http.StatusBadRequest, errors.New("código e quantidade devem ser informados"))
		return
	}

	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao conectar com o banco de dados"))
		return
	}

	repositorioEstoque := repositorios.NovoRepositorioDeEstoque(db)

	estoqueDisponivel, err := repositorioEstoque.ObterEstoqueDisponivel(produto.Codigo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respostas.Erro(w, http.StatusNotFound, errors.New("produto não encontrado no estoque"))
		} else {
			respostas.Erro(w, http.StatusInternalServerError, err)
		}
		return
	}

	if produto.Quantidade > estoqueDisponivel {
		respostas.Erro(w, http.StatusConflict, errors.New("estoque insuficiente para o produto"))
		return
	}

	respostas.JSON(w, http.StatusOK, map[string]interface{}{
		"codigo":            produto.Codigo,
		"quantidadePedido":  produto.Quantidade,
		"quantidadeEstoque": estoqueDisponivel,
	})
}
