package mysql

import (
	models "asik1/pkg"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, email, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, role) VALUES (?, ?, ?, ?)`
	_, err = u.DB.Exec(stmt, name, email, hashedPassword, role)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	row := u.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}
func (u *UserModel) GetByEmail(email string) (*models.User, error) {
	query := "SELECT * FROM users WHERE email = ?"
	row := u.DB.QueryRow(query, email)

	user := &models.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (u *UserModel) GetRoleByEmail(email string) (string, error) {
	query := "SELECT role FROM users WHERE email = ?"
	row := u.DB.QueryRow(query, email)

	var role string
	err := row.Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil // No rows found, return nil without error
		}
		return "", err // Return error for other cases
	}
	return role, nil
}
func (u *UserModel) GetNameByEmail(email string) (string, error) {
	query := "SELECT name FROM users WHERE email = ?"
	row := u.DB.QueryRow(query, email)

	var name string
	err := row.Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil // No rows found, return nil without error
		}
		return "", err // Return error for other cases
	}
	return name, nil
}
func (u *UserModel) GetAllUsers() ([]*models.User, error) {
	rows, err := u.DB.Query("SELECT id,name,email,role FROM users WHERE role!='admin'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		n := &models.User{}
		err := rows.Scan(&n.Id, &n.Name, &n.Email, &n.Role)
		if err != nil {
			return nil, err
		}

		users = append(users, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
func (u *UserModel) ChangeUserRole(userID, newRole string) error {
	_, err := u.DB.Exec("UPDATE users SET role = ? WHERE id = ?", newRole, userID)
	if err != nil {
		return err
	}
	return nil
}
