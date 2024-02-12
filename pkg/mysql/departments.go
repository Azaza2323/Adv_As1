package mysql

import (
	models "asik1/pkg"
	"database/sql"
)

type DepoModel struct {
	DB *sql.DB
}

func (d *DepoModel) InsertDepo(id, depoID, staffQuantity int) (int, error) {
	stmt := `INSERT INTO departments (id, dep_id, staff_quantity) VALUES (?, ?, ?)`
	result, err := d.DB.Exec(stmt, id, depoID, staffQuantity)
	if err != nil {
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertID), nil
}
func (d *DepoModel) Deps() ([]*models.Depo, error) {
	stmt := `SELECT id, dep_id,staff_quantity FROM departments LIMIT 10`
	rows, err := d.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []*models.Depo

	for rows.Next() {
		n := &models.Depo{}

		err := rows.Scan(&n.ID, &n.DepId, &n.Staff)
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
