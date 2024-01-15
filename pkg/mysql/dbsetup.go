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
	stmt := `SELECT id, audience, author, title, description, content FROM news LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []*models.News

	for rows.Next() {
		n := &models.News{}

		err := rows.Scan(&n.ID, &n.Audience, &n.Author, &n.Title, &n.Description, &n.Content)
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
	stmt := `SELECT id, audience, author, title, description, content FROM news WHERE audience = ? DESC LIMIT 10`
	rows, _ := m.DB.Query(stmt, audience)
	var s []*models.News

	for rows.Next() {
		n := &models.News{}
		err := rows.Scan(&n.ID, &n.Audience, &n.Author, &n.Title, &n.Description, &n.Content)
		if err != nil {
			return nil, err
		}
		s = append(s, n)
	}
	return s, nil
}
