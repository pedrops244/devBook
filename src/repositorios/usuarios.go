package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

type usuarios struct {
	db *sql.DB
}

func NovoRepositorioDeUsuario(db *sql.DB) *usuarios {
	return &usuarios{db}
}

func (repositorio usuarios) Criar(usuario modelos.Usuario) (uint, error) {
	statement, erro := repositorio.db.Prepare(
		"INSERT INTO usuarios (username, senha, role) OUTPUT Inserted.id VALUES (@username, @senha, @role)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	var ultimoIdInserido uint
	if erro := statement.QueryRow(
		sql.Named("username", usuario.Username),
		sql.Named("senha", usuario.Senha),
		sql.Named("role", usuario.Role),
	).Scan(&ultimoIdInserido); erro != nil {
		return 0, erro
	}
	return ultimoIdInserido, nil
}

func (repositorio usuarios) Buscar() ([]modelos.Usuario, error) {

	linhas, erro := repositorio.db.Query(
		"SELECT id, username, role, created_at, is_deleted FROM usuarios",
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
			&usuario.Username,
			&usuario.Role,
			&usuario.CriadoEm,
			&usuario.IsDeleted,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}

func (repositorio usuarios) BuscarPorUsername(username string) (modelos.Usuario, error) {
	linha, erro := repositorio.db.Query(
		"SELECT id, senha, role FROM usuarios WHERE username = @username",
		sql.Named("username", username),
	)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linha.Close()

	var usuario modelos.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Senha, &usuario.Role); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

func (repositorio usuarios) BuscarPorID(id uint64) (modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"SELECT id, username, role FROM usuarios where id = @id", sql.Named("id", id),
	)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linhas.Close()

	var usuario modelos.Usuario

	if linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Username,
			&usuario.Role,
		); erro != nil {
			return modelos.Usuario{}, erro
		}
	}
	return usuario, nil
}

func (repositorio usuarios) Atualizar(id uint64, usuario modelos.Usuario) error {
	statement, erro := repositorio.db.Prepare(
		"UPDATE usuarios SET username = @username WHERE id = @id",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(
		sql.Named("username", usuario.Username),
		sql.Named("id", id),
	); erro != nil {
		return erro
	}
	return nil
}

func (repositorio usuarios) Deletar(id uint64) error {
	statement, erro := repositorio.db.Prepare("DELETE FROM usuarios WHERE id = @id")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(sql.Named("id", id)); erro != nil {
		return erro
	}

	return nil
}

func (repositorio usuarios) ObterRoleUsuario(usuarioID uint64) (string, error) {
	var role string

	linha := repositorio.db.QueryRow(
		"SELECT role FROM usuarios WHERE id = @id",
		sql.Named("id", usuarioID),
	)

	if erro := linha.Scan(&role); erro != nil {
		if erro == sql.ErrNoRows {
			return "", erro
		}
	}

	return role, nil
}

func (repositorio usuarios) ExisteUsuarioAdmin() (bool, error) {
	var exists bool

	linha := repositorio.db.QueryRow(
		"SELECT CASE WHEN EXISTS (SELECT 1 FROM usuarios WHERE role = @role) THEN 1 ELSE 0 END",
		sql.Named("role", "admin"),
	)

	if erro := linha.Scan(&exists); erro != nil {
		return false, erro
	}

	return exists, nil
}

func (repositorio usuarios) InativarUsuario(id uint64) error {
	statement, erro := repositorio.db.Prepare("UPDATE usuarios SET is_deleted = 1 WHERE id = @id")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(sql.Named("id", id)); erro != nil {
		return erro
	}
	return nil
}
