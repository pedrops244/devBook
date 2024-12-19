package repositorios

import (
	"database/sql"
	"errors"
)

type estoque struct {
	db *sql.DB
}

// NovoRepositorioDeEstoque cria um novo repositório de estoque
func NovoRepositorioDeEstoque(db *sql.DB) *estoque {
	return &estoque{db}
}

func (repositorio estoque) ObterEstoqueDisponivel(codigo string) (int, error) {
	var disponivel int
	query := `
        SELECT (quantidade - reservado) AS Disponivel
        FROM estoque
        WHERE codigo = @codigo
    `
	err := repositorio.db.QueryRow(query, sql.Named("codigo", codigo)).Scan(&disponivel)
	if err != nil {
		return 0, err
	}
	return disponivel, nil
}

// ReservarItens atualiza a quantidade reservada no estoque
func (repositorio estoque) ReservarItens(codigo string, quantidade int) error {
	query := `
        UPDATE estoque
        SET reservado = reservado + @quantidade
        WHERE codigo = @codigo AND (quantidade - reservado) >= @quantidade
    `
	result, err := repositorio.db.Exec(query, sql.Named("quantidade", quantidade), sql.Named("codigo", codigo))
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("estoque insuficiente para a reserva")
	}
	return err
}

// ConfirmarSaida atualiza o estoque real após a confirmação do pedido
func (repositorio estoque) ConfirmarSaida(codigo string, quantidade int) error {
	query := `
        UPDATE estoque
        SET quantidade = quantidade - @quantidade, reservado = reservado - @quantidade
        WHERE codigo = @codigo AND reservado >= @quantidade
    `
	result, err := repositorio.db.Exec(query, sql.Named("quantidade", quantidade), sql.Named("codigo", codigo))
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("falha ao confirmar a saída do estoque ou quantidade reservada insuficiente")
	}
	return nil
}

func (repositorio estoque) DevolverSobraEstoque(codigo string, sobra int) error {
	query := `
        UPDATE estoque
        SET reservado = reservado - @sobra
        WHERE codigo = @codigo
    `
	_, err := repositorio.db.Exec(query, sql.Named("sobra", sobra), sql.Named("codigo", codigo))
	return err
}

func (repositorio estoque) ZerarEstoque(codigo string, sobra int) error {
	query := `
        UPDATE estoque
        SET quantidade = 0, reservado = 0, faltas = @sobra
        WHERE codigo = @codigo
    `
	_, err := repositorio.db.Exec(query, sql.Named("codigo", codigo), sql.Named("sobra", sobra))
	return err
}
