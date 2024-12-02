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

	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

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

	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

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
	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	repositorio := repositorios.NovoRepositorioDePedidos(db)
	pedidos, erro := repositorio.Listar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, pedidos)
}

// ConfirmarRecebimento atualiza a QuandidadeConferida no banco
func ConfirmarRecebimento(w http.ResponseWriter, r *http.Request) {
	// Captura o ID do pedido
	parametros := mux.Vars(r)
	pedidoID, err := strconv.ParseUint(parametros["pedidoID"], 10, 64)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("id do pedido inválido"))
		return
	}

	// Lê o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, errors.New("erro ao ler o corpo da requisição"))
		return
	}
	defer r.Body.Close()

	// Mapeia os dados recebidos
	var pedido modelos.Pedido
	if err := json.Unmarshal(body, &pedido); err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("dados inválidos"))
		return
	}

	if err := pedido.Validar(); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	// Conecta ao banco
	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Cria o repositório e atualiza os dados
	repositorio := repositorios.NovoRepositorioDePedidos(db)
	err = repositorio.AtualizarRecebimento(uint(pedidoID), pedido)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Retorna sucesso
	respostas.JSON(w, http.StatusOK, map[string]string{"mensagem": "Recebimento confirmado com sucesso"})
}

func ConfirmarConferencia(w http.ResponseWriter, r *http.Request) {
	// Captura o ID do pedido
	parametros := mux.Vars(r)
	pedidoID, err := strconv.ParseUint(parametros["pedidoID"], 10, 64)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("id do pedido inválido"))
		return
	}

	// Lê o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, errors.New("erro ao ler o corpo da requisição"))
		return
	}
	defer r.Body.Close()

	// Mapeia os dados recebidos
	var pedido modelos.Pedido
	if err := json.Unmarshal(body, &pedido); err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("dados inválidos"))
		return
	}

	if err := pedido.Validar(); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	// Conecta ao banco
	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Cria o repositório e atualiza os dados
	repositorio := repositorios.NovoRepositorioDePedidos(db)
	err = repositorio.AtualizarConferencia(uint(pedidoID), pedido)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Retorna sucesso
	respostas.JSON(w, http.StatusOK, map[string]string{"mensagem": "Recebimento confirmado com sucesso"})
}
