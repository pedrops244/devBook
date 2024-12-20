package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"errors"
)

// pedidos representa o repositório de pedidos
type pedidos struct {
	db *sql.DB
}

// NovoRepositorioDePedidos cria um repositório de pedidos
func NovoRepositorioDePedidos(db *sql.DB) *pedidos {
	return &pedidos{db}
}

// Criar adiciona um novo pedido com seus itens
func (repositorio pedidos) Criar(pedido modelos.Pedido) (uint, error) {
	tx, err := repositorio.db.Begin()
	if err != nil {
		return 0, err
	}

	// Inserir o pedido
	queryPedido := `
		INSERT INTO pedidos (status, usuario_id, criado_em)
		OUTPUT INSERTED.ID
		VALUES (@status, @usuario_id, @criado_em)
	`
	stmtPedido, err := tx.Prepare(queryPedido)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmtPedido.Close()

	var pedidoID uint
	err = stmtPedido.QueryRow(
		sql.Named("status", pedido.Status),
		sql.Named("usuario_id", pedido.UsuarioId),
		sql.Named("criado_em", pedido.CriadoEm),
	).Scan(&pedidoID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Inserir os itens do pedido
	queryItem := `
		INSERT INTO itens_pedidos 
		(pedido_id, quantidade_solicitada, quantidade_recebida, quantidade_conferida, codigo)
		VALUES (@pedido_id, @quantidade_solicitada, @quantidade_recebida, @quantidade_conferida, @codigo)
	`
	stmtItem, err := tx.Prepare(queryItem)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmtItem.Close()

	for _, item := range pedido.Itens {
		_, err := stmtItem.Exec(
			sql.Named("pedido_id", pedidoID),
			sql.Named("quantidade_solicitada", item.QuantidadeSolicitada),
			sql.Named("quantidade_recebida", 0),
			sql.Named("quantidade_conferida", 0),
			sql.Named("codigo", item.Codigo),
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Commit da transação
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	// Retornar o ID do pedido criado
	return pedidoID, nil
}

// BuscarPorID busca o pedido pelo ID
func (repositorio pedidos) BuscarPorID(id uint) (modelos.Pedido, error) {
	queryPedido := `
		SELECT id, status, usuario_id, criado_em
		FROM pedidos
		WHERE id = @id
	`
	var pedido modelos.Pedido
	err := repositorio.db.QueryRow(queryPedido, sql.Named("id", id)).
		Scan(&pedido.ID, &pedido.Status, &pedido.UsuarioId, &pedido.CriadoEm)
	if err == sql.ErrNoRows {
		return pedido, errors.New("pedido não encontrado")
	} else if err != nil {
		return pedido, err
	}

	// Buscar os itens associados ao pedido
	queryItens := `
		SELECT id, pedido_id, quantidade_solicitada, quantidade_recebida, quantidade_conferida, codigo
		FROM itens_pedidos
		WHERE pedido_id = @pedido_id
	`
	rows, err := repositorio.db.Query(queryItens, sql.Named("pedido_id", pedido.ID))
	if err != nil {
		return pedido, err
	}
	defer rows.Close()

	var itens []modelos.ItensPedido
	for rows.Next() {
		var item modelos.ItensPedido
		if err := rows.Scan(
			&item.ID,
			&item.PedidoID,
			&item.QuantidadeSolicitada,
			&item.QuantidadeRecebida,
			&item.QuantidadeConferida,
			&item.Codigo,
		); err != nil {
			return pedido, err
		}
		itens = append(itens, item)
	}
	pedido.Itens = itens
	return pedido, nil
}

// Listar lista todos os pedidos
func (repositorio pedidos) Listar() ([]modelos.Pedido, error) {
	queryPedidos := `
		SELECT id, status, usuario_id, criado_em, recebido_em, conferido_em
		FROM pedidos
	`
	rows, err := repositorio.db.Query(queryPedidos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pedidos []modelos.Pedido
	for rows.Next() {
		var pedido modelos.Pedido
		if err := rows.Scan(&pedido.ID, &pedido.Status, &pedido.UsuarioId, &pedido.CriadoEm, &pedido.RecebidoEm, &pedido.ConferidoEm); err != nil {
			return nil, err
		}

		// Buscar os itens do pedido
		itens, err := repositorio.BuscarItensDoPedido(pedido.ID)
		if err != nil {
			return nil, err
		}
		pedido.Itens = itens
		pedidos = append(pedidos, pedido)
	}
	return pedidos, nil
}

// BuscarItensDoPedido busca todos os itens dentro de um pedido especifico pelo ID
func (repositorio pedidos) BuscarItensDoPedido(pedidoID uint) ([]modelos.ItensPedido, error) {
	query := `
		SELECT id, pedido_id, quantidade_solicitada, quantidade_recebida, quantidade_conferida, codigo
		FROM itens_pedidos
		WHERE pedido_id = @pedido_id
	`
	rows, err := repositorio.db.Query(query, sql.Named("pedido_id", pedidoID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itens []modelos.ItensPedido
	for rows.Next() {
		var item modelos.ItensPedido
		if err := rows.Scan(
			&item.ID,
			&item.PedidoID,
			&item.QuantidadeSolicitada,
			&item.QuantidadeRecebida,
			&item.QuantidadeConferida,
			&item.Codigo,
		); err != nil {
			return nil, err
		}
		itens = append(itens, item)
	}
	return itens, nil
}

// DeletarPedido remove um pedido pelo ID
func (repositorio pedidos) DeletarPedido(pedidoID uint) error {
	tx, err := repositorio.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		DELETE FROM pedidos
		WHERE id = @pedido_id
	`

	_, err = tx.Exec(query, sql.Named("pedido_id", pedidoID))
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// AtualizarRecebimento atualiza a QuantidadeRecebida no banco
func (repositorio pedidos) AtualizarRecebimento(pedidoID uint, pedido modelos.Pedido) error {
	tx, err := repositorio.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		UPDATE pedidos
		SET status = @status, recebido_em = @recebido_em, repositor_id = @repositor_id
		WHERE id = @pedido_id
	`

	_, err = tx.Exec(
		query,
		sql.Named("status", pedido.Status),
		sql.Named("recebido_em", pedido.RecebidoEm),
		sql.Named("repositor_id", pedido.RepositorId),
		sql.Named("pedido_id", pedidoID),
	)
	if err != nil {
		return err
	}

	queryItem := `
	UPDATE itens_pedidos
	SET quantidade_recebida = @quantidade_recebida
	WHERE pedido_id = @pedido_id AND codigo = @codigo
`
	for _, itemRecebido := range pedido.Itens {
		_, err = tx.Exec(
			queryItem,
			sql.Named("quantidade_recebida", itemRecebido.QuantidadeRecebida),
			sql.Named("pedido_id", pedidoID),
			sql.Named("codigo", itemRecebido.Codigo),
		)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// AtualizarConferencia atualiza a QuantidadeConferida no banco
func (repositorio pedidos) AtualizarConferencia(pedidoID uint, pedido modelos.Pedido) error {
	tx, err := repositorio.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		UPDATE pedidos
		SET status = @status, conferido_em = @conferido_em, conferente_id = @conferente_id
		WHERE id = @pedido_id
	`
	_, err = tx.Exec(
		query,
		sql.Named("status", pedido.Status),
		sql.Named("conferido_em", pedido.ConferidoEm),
		sql.Named("conferente_id", pedido.ConferenteId),
		sql.Named("pedido_id", pedidoID),
	)
	if err != nil {
		return err
	}

	queryItem := `
	UPDATE itens_pedidos
	SET quantidade_conferida = @quantidade_conferida
	WHERE pedido_id = @pedido_id AND codigo = @codigo
`
	for _, itemRecebido := range pedido.Itens {
		_, err = tx.Exec(
			queryItem,
			sql.Named("quantidade_conferida", itemRecebido.QuantidadeConferida),
			sql.Named("pedido_id", pedidoID),
			sql.Named("codigo", itemRecebido.Codigo),
		)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// VerificarStatus retorna o status atual do pedido no banco
func (repositorio pedidos) VerificarStatus(pedidoID uint) (string, error) {
	query := `
		SELECT status
		FROM pedidos
		WHERE id = @id
	`

	var status string
	err := repositorio.db.QueryRow(query, sql.Named("id", pedidoID)).Scan(&status)
	if err == sql.ErrNoRows {
		return "", errors.New("pedido não encontrado")
	} else if err != nil {
		return "", err
	}

	return status, nil
}
