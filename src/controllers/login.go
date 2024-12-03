package controllers

import (
	"api/src/auth"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func Login(w http.ResponseWriter, r *http.Request) {
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

	db, erro := banco.ObterConexao()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("erro ao conectar com o banco de dados"))
		return
	}

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	usuarioSalvoNoBanco, erro := repositorio.BuscarPorUsername(usuario.Username)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(usuarioSalvoNoBanco.Senha, usuario.Senha); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	token, erro := auth.CriarToken(usuarioSalvoNoBanco.ID, usuarioSalvoNoBanco.Role)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	usuarioId := strconv.FormatUint(uint64(usuarioSalvoNoBanco.ID), 10)

	respostas.JSON(w, http.StatusOK, modelos.DadosAutenticacao{ID: usuarioId, Token: token, Role: usuarioSalvoNoBanco.Role})
}
