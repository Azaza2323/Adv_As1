package mysql

import (
	models "asik1/pkg"
	"database/sql"
	_ "fmt"
	_ "strconv"
)

type NewsModel struct {
	DB *sql.DB
}

func (m *NewsModel) Insert(audience, author, title, description, content string) (int, error) {
	stmt := `INSERT INTO news (audience, author, title, description, content) VALUES (?, ?, ?, ?, ?)`
	result, err := m.DB.Exec(stmt, audience, author, title, description, content)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *NewsModel) Latest() ([]*models.News, error) {
	stmt := `SELECT id, audience, author, title, description, content, created FROM news ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	newsList := []*models.News{}

	for rows.Next() {
		n := &models.News{}

		err := rows.Scan(&n.ID, &n.Audience, &n.Author, &n.Title, &n.Description, &n.Content, &n.Created)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return newsList, nil
}

func (m *NewsModel) GetByAudience(audience string) ([]*models.News, error) {
	stmt := `SELECT id, audience, author, title, description, content, created FROM news WHERE audience = ? ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(stmt, audience)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	newsList := []*models.News{}

	for rows.Next() {
		n := &models.News{}

		err := rows.Scan(&n.ID, &n.Audience, &n.Author, &n.Title, &n.Description, &n.Content, &n.Created)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return newsList, nil
}

func (m *NewsModel) Delete(id int) (bool, error) {
	stmt := `DELETE FROM news WHERE id = ?`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *NewsModel) Update(id int, audience, author, title, description, content string) (bool, error) {
	stmt := `UPDATE news SET audience = ?, author = ?, title = ?, description = ?, content = ? WHERE id = ?`
	_, err := m.DB.Exec(stmt, audience, author, title, description, content, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
