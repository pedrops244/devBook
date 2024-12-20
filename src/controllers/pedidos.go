package controllers

import (
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"fmt"
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
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao conectar com o banco de dados"))
		return
	}

	repositorio := repositorios.NovoRepositorioDePedidos(db)
	repositorioEstoque := repositorios.NovoRepositorioDeEstoque(db)

	// Verifica se há estoque disponível para cada item do pedido
	for _, item := range pedido.Itens {
		estoqueDisponivel, err := repositorioEstoque.ObterEstoqueDisponivel(item.Codigo)
		if err != nil {
			respostas.Erro(w, http.StatusInternalServerError, err)
			return
		}

		if item.QuantidadeSolicitada > estoqueDisponivel {
			respostas.Erro(w, http.StatusConflict, errors.New("estoque insuficiente para o produto: "+item.Codigo))
			return
		}

		// Reserva o item no estoque
		if err := repositorioEstoque.ReservarItens(item.Codigo, item.QuantidadeSolicitada); err != nil {
			respostas.Erro(w, http.StatusInternalServerError, errors.New("estoque insuficiente para a reserva"))
			return
		}
	}

	// Cria o pedido no banco de dados
	pedidoID, erro := repositorio.Criar(pedido)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, map[string]string{"mensagem": fmt.Sprintf("Pedido criado com sucesso! ID: %d", pedidoID)})
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
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao conectar com o banco de dados"))
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

// BuscarPedidos retorna todos os pedidos com seus itens
func BuscarPedidos(w http.ResponseWriter, r *http.Request) {
	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao conectar com o banco de dados"))
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

// DeletarPedido deleta o pedido pelo ID
func DeletarPedido(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	pedidoID, erro := strconv.ParseUint(parametros["pedidoID"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao conectar com o banco de dados"))
		return
	}

	repositorio := repositorios.NovoRepositorioDePedidos(db)
	if erro = repositorio.DeletarPedido(uint(pedidoID)); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao deletar o pedido"))
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// ConfirmarRecebimento atualiza a QuantidadeRecebida no banco
func ConfirmarRecebimento(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	pedidoID, err := strconv.ParseUint(parametros["pedidoID"], 10, 64)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}
	defer r.Body.Close()

	var pedido modelos.Pedido
	if err := json.Unmarshal(body, &pedido); err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("dados inválidos"))
		return
	}

	if err := pedido.Validar(); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao conectar com o banco de dados"))
		return
	}

	repositorio := repositorios.NovoRepositorioDePedidos(db)
	repositorioEstoque := repositorios.NovoRepositorioDeEstoque(db)

	// Verifica o status atual do pedido
	statusAtual, err := repositorio.VerificarStatus(uint(pedidoID))
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if statusAtual == "recebido" {
		respostas.JSON(w, http.StatusConflict, map[string]string{
			"erro": "Pedido já foi recebido ou está em status não permitido para alteração",
		})
		return
	}

	// Processar cada item do pedido
	for _, item := range pedido.Itens {

		sobraNegativo := item.QuantidadeRecebida - item.QuantidadeSolicitada
		// Verifica se o checkbox da flag IsInsufficientStock foi marcado
		if item.IsInsufficientStock {
			// Se flag foi marcada, zera o estoque físico
			err := repositorioEstoque.ZerarEstoque(item.Codigo, sobraNegativo)
			if err != nil {
				respostas.Erro(w, http.StatusInternalServerError, err)
				return
			}
			continue // Não tenta mais confirmar a saída para este item
		}
		// Se houver sobra, devolve ao estoque
		sobra := item.QuantidadeSolicitada - item.QuantidadeRecebida
		fmt.Println(sobra)
		if sobra > 0 {
			err := repositorioEstoque.DevolverSobraEstoque(item.Codigo, sobra)
			if err != nil {
				respostas.Erro(w, http.StatusInternalServerError, err)
				return
			}
		}

		// Confirma a saída do estoque
		err = repositorioEstoque.ConfirmarSaida(item.Codigo, item.QuantidadeRecebida)
		if err != nil {
			respostas.Erro(w, http.StatusInternalServerError, err)
			return
		}
	}

	// Atualiza o recebimento no pedido
	err = repositorio.AtualizarRecebimento(uint(pedidoID), pedido)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, map[string]string{"mensagem": "Recebimento confirmado com sucesso"})
}

// ConfirmarConferencia atualiza a QuantidadeConferida no banco
func ConfirmarConferencia(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	pedidoID, err := strconv.ParseUint(parametros["pedidoID"], 10, 64)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("id do pedido inválido"))
		return
	}

	// Lê o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, err)
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
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao conectar com o banco de dados"))
		return
	}

	repositorio := repositorios.NovoRepositorioDePedidos(db)

	statusAtual, err := repositorio.VerificarStatus(uint(pedidoID))
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if statusAtual == "conferido" {
		respostas.JSON(w, http.StatusConflict, map[string]string{
			"erro": "Pedido já foi conferido ou está em status não permitido para alteração",
		})
		return
	}

	err = repositorio.AtualizarConferencia(uint(pedidoID), pedido)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, map[string]string{"mensagem": "Conferência confirmada com sucesso"})
}
