package controllers

import (
	"api/src/auth"
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

// CriarUsuario cria um novo usuário
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("cadastro"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao conectar com o banco de dados"))
		return
	}

	repositorio := repositorios.NovoRepositorioDeUsuario(db)

	usuarioExistente, erro := repositorio.BuscarPorUsername(usuario.Username)
	if erro != nil {
		if erro.Error() != "usuário não encontrado" {
			respostas.Erro(w, http.StatusInternalServerError, erro)
			return
		}
	} else if usuarioExistente.ID != 0 {
		respostas.Erro(w, http.StatusConflict, errors.New("nome de usuário já está em uso"))
		return
	}

	usuario.ID, erro = repositorio.Criar(usuario)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, map[string]string{"mensagem": "Usuário criado com sucesso"})
}

// BuscarUsuarios retorna todos os usuários
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {

	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao conectar com o banco de dados"))
		return
	}

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	usuarios, erro := repositorio.Buscar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}

// AtualizarUsuario atualiza o usuário pelo ID
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametro := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametro["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIdNoToken, erro := auth.ExtrairUsuarioId(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIdNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possivel atualizar um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {

		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	if erro = repositorio.Atualizar(usuarioID, usuario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarUsuario deleta o usuário pelo ID
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Verifica se o usuário é admin
	roleLogado, erro := auth.ExtrairRole(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if roleLogado != "admin" && roleLogado != "gerente" {
		respostas.Erro(w, http.StatusConflict, errors.New("apenas administradores ou gerentes podem deletar usuários"))
		return
	}

	// Evitar que o admin/gerente delete a si mesmo
	usuarioLogadoID, erro := auth.ExtrairUsuarioId(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID == usuarioLogadoID {
		respostas.Erro(w, http.StatusConflict, errors.New("você não pode deletar seu próprio usuário"))
		return
	}

	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao conectar com o banco de dados"))
		return
	}

	repositorio := repositorios.NovoRepositorioDeUsuario(db)

	roleUsuario, erro := repositorio.ObterRoleUsuario(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao verificar role do usuário"))
		return
	}

	if roleUsuario == "admin" {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é permitido deletar um usuário administrador"))
		return
	}

	if erro = repositorio.Deletar(usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao remover usuário"))
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
