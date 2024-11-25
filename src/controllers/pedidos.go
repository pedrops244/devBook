package controllers

import (
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CriarPedido cria um novo pedido com seus itens
func CriarPedido(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}
	defer r.Body.Close()

	var pedido modelos.Pedido
	if err := json.Unmarshal(body, &pedido); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	if err := pedido.Validar(); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePedidos(db)
	pedidoID, erro := repositorio.Criar(pedido)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, map[string]uint{"pedidoID": pedidoID})
}

// BuscarPedido busca um pedido pelo ID
func BuscarPedido(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	pedidoID, err := strconv.ParseUint(parametros["pedidoID"], 10, 64)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePedidos(db)
	pedido, erro := repositorio.BuscarPorID(uint(pedidoID))
	if erro != nil {
		respostas.Erro(w, http.StatusNotFound, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, pedido)
}

// ListarPedidos retorna todos os pedidos com seus itens
func ListarPedidos(w http.ResponseWriter, r *http.Request) {
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePedidos(db)
	pedidos, erro := repositorio.Listar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, pedidos)
}

// AtualizarStatusPedido atualiza o status de um pedido
func AtualizarStatusPedido(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	pedidoID, err := strconv.ParseUint(parametros["pedidoID"], 10, 64)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	var dados struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &dados); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	if dados.Status == "" {
		respostas.Erro(w, http.StatusBadRequest, errors.New("o status não pode ser vazio"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePedidos(db)
	if erro := repositorio.AtualizarStatus(uint(pedidoID), dados.Status); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
