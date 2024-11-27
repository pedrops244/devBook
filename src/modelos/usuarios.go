package modelos

import (
	"api/src/seguranca"
	"errors"
	"strings"
	"time"
)

type Usuario struct {
	ID       uint      `json:"id,omitempty"`
	Username string    `json:"username,omitempty"`
	Senha    string    `json:"senha,omitempty"`
	Role     string    `json:"role,omitempty"`
	CriadoEm time.Time `json:"criadoEm,omitempty"`
}

func (usuario *Usuario) Preparar(etapa string) error {
	if erro := usuario.validar(etapa); erro != nil {
		return erro
	}
	if erro := usuario.formatar(etapa); erro != nil {
		return erro
	}
	return nil
}

func (usuario *Usuario) validar(etapa string) error {
	if usuario.Username == "" {
		return errors.New("o username é obrigatório e não pode estar em branco")
	}
	if usuario.Role != "comprador" && usuario.Role != "repositor" && usuario.Role != "gerente" {
		return errors.New("a role é obrigatória e não pode ser diferente de comprador, repositor ou gerente")
	}

	if etapa == "cadastro" && usuario.Senha == "" {
		return errors.New("a senha é obrigatório e não pode estar em branco")
	}

	return nil
}

func (usuario *Usuario) formatar(etapa string) error {
	usuario.Username = strings.TrimSpace(usuario.Username)

	if etapa == "cadastro" {
		senhaComHash, erro := seguranca.Hash(usuario.Senha)
		if erro != nil {
			return erro
		}
		usuario.Senha = string(senhaComHash)
	}
	return nil

}
