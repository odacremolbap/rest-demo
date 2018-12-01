package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/odacremolbap/rest-demo/pkg/db/clauses"
	"github.com/odacremolbap/rest-demo/pkg/log"
	"github.com/odacremolbap/rest-demo/pkg/types"
	"github.com/pkg/errors"
)

// SelectTasks executes a tasks query at the database
func (p PersistenceManager) SelectTasks(q *clauses.Query) ([]types.Task, error) {

	query := `
		select
			id, 
			name, 
			description, 
			category, 
			status, 
			duedate,
			created
		from tasks`

	if len(q.Where) != 0 {
		query = fmt.Sprintf("%s where %s", query, q.Where)
	}
	if len(q.Pagination) != 0 {
		query = fmt.Sprintf("%s %s", query, q.Pagination)
	}
	if len(q.OrderByClause) != 0 {
		query = fmt.Sprintf("%s order by %s", query, q.OrderByClause)
	}

	log.V(10).Info("Executing query",
		"query", query,
		"parameters", q.WhereParams)

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "error preparing SelectTasks statement")
	}

	rows, err := stmt.Query(q.WhereParams...)

	if err != nil {
		return nil, errors.Wrap(err, "error retrieving Tasks")
	}
	defer rows.Close()

	items := []types.Task{}
	for rows.Next() {
		item := types.Task{}
		if rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Category,
			&item.Status,
			&item.DueDate,
			&item.Created) != nil {
			return nil, errors.Wrap(err, "error scanning Tasks")
		}
		items = append(items, item)
	}
	return items, nil
}

// GetTask from the database
// If object by ID doesn't exists, nil is returned
func (p *PersistenceManager) GetTask(ID int) (*types.Task, error) {
	query := `
		select
			name, 
			description, 
			category, 
			status, 
			duedate,
			created
		from tasks
		where id = $1`
	item := &types.Task{ID: ID}

	log.V(10).Info("Executing query",
		"query", query,
		"ID", ID)
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "error preparing GetTask statement")
	}

	err = stmt.QueryRow().Scan(&item.ID,
		&item.Name,
		&item.Description,
		&item.Category,
		&item.Status,
		&item.DueDate,
		&item.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "error scanning Task")
	}
	return item, nil
}

// CreateTask at the database
func (p *PersistenceManager) CreateTask(item *types.Task) (*types.Task, error) {
	query := `
		insert into tasks
		(
			name, 
			description, 
			category, 
			status, 
			duedate,
		)
		values
			($1, $2, $3, $4, $5)
		returning
			id, created`
	log.V(10).Info("Executing query",
		"query", query,
		"parameters", item)

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "error preparing CreateTask statement")
	}

	err = stmt.QueryRow(
		item.Name,
		item.Description,
		item.Category,
		strings.ToLower(item.Status),
		*item.DueDate,
		item.Status).
		Scan(
			&item.ID,
			&item.Created)

	if err != nil {
		return nil, errors.Wrap(err, "error creating Task")
	}
	return item, nil
}

// UpdateOneTask object at the database
func (p *PersistenceManager) UpdateOneTask(item *types.Task) (*types.Task, error) {
	query := `
		update tasks set
			name = $1,
			description = $2,
			category = $3,
			status = $4,
			duedate = $5
		where
			id = $6`
	log.V(10).Info("Executing query",
		"query", query,
		"parameters", item)

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "error preparing UpdateOneTask statement")
	}

	_, err = stmt.Exec(
		item.Name,
		item.Description,
		item.Category,
		strings.ToLower(item.Status),
		*item.DueDate,
		item.ID)

	if err != nil {
		return nil, errors.Wrap(err, "error updating Task")
	}
	return item, nil
}

// DeleteOneTask object at the database
func (p *PersistenceManager) DeleteOneTask(ID int) error {
	query := `
		delete from tasks
		where
		id = $1`
	log.V(10).Info("Executing query",
		"query", query,
		"ID", ID)

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return errors.Wrap(err, "error preparing DeleteOneTask statement")
	}

	_, err = stmt.Exec(query, ID)
	if err != nil {
		return errors.Wrap(err, "error deleting Tasks")
	}
	return nil
}
