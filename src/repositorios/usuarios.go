package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

type usuarios struct {
	db *sql.DB
}

func NovoRepositorioDeUsuario(db *sql.DB) *usuarios {
	return &usuarios{db}
}

func (repositorio usuarios) Criar(usuario modelos.Usuario) (uint, error) {
	statement, erro := repositorio.db.Prepare("insert into usuarios(nome, nick, email, senha) values(?, ?, ?, ?)")
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if erro != nil {
		return 0, erro
	}

	ultimoIdInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}
	return uint(ultimoIdInserido), nil
}

func (repositorio usuarios) Buscar(nomeOuNick string) ([]modelos.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick)

	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where nome LIKE ? or nick LIKE ?",
		nomeOuNick, nomeOuNick,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []modelos.Usuario

	for linhas.Next() {
		var usuario modelos.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}

func (repositorio usuarios) BuscarPorID(id uint64) (modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoEm form usuarios where id = ?", id,
	)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linhas.Close()

	var usuario modelos.Usuario

	if linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nick,
			&usuario.Nome,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return modelos.Usuario{}, erro
		}
	}
	return usuario, nil
}

func (repositorio usuarios) Atualizar(ID uint64, usuario modelos.Usuario) error {
	statement, erro := repositorio.db.Prepare(
		"update usuarios set nome = ?, nick = ?, email = ? where id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, ID); erro != nil {
		return erro
	}
	return nil
}
