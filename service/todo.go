package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	res, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var todo model.TODO
	todo.ID = id
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, err
	case err != nil:
		return nil, err
	default:
		return &todo, nil
	}
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		defaultSQL = `SELECT id, subject, description, created_at, updated_at FROM todos`
		readALl    = defaultSQL + ` ORDER BY id DESC`
		read       = defaultSQL + ` ORDER BY id DESC LIMIT ?`
		readWithID = defaultSQL + ` WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	var err error
	var rows *sql.Rows

	if size == 0 {
		rows, err = s.db.QueryContext(ctx, readALl)
	} else {
		if prevID == 0 {
			rows, err = s.db.QueryContext(ctx, read, size)
		} else {
			rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
		}
	}

	if err != nil {
		return []*model.TODO{}, err
	}

	todos := []*model.TODO{}

	for rows.Next() {
		todo := model.TODO{}
		err = rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return []*model.TODO{}, err
		}
		todos = append(todos, &todo)
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	res, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if row == 0 {
		return nil, &model.ErrNotFound{}
	}

	var todo model.TODO
	todo.ID = id
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, err
	case err != nil:
		return nil, err
	default:
		return &todo, nil
	}
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const delete = `DELETE FROM todos WHERE id = ?`

	if len(ids) == 0 {
		return nil
	}

	for _, id := range ids {
		res, err := s.db.ExecContext(ctx, delete, fmt.Sprint(id))
		if err != nil {
			return err
		}
		row, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if row == 0 {
			return &model.ErrNotFound{}
		}
	}

	return nil
}
