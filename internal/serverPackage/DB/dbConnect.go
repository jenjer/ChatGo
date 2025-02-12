package DBConnect

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

type UserDB struct {
	db *sql.DB
}

func NewUserDB(dbPath string) (*UserDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("database open error: %v", err)
	}

	//사용자 수가 많지 않을것으로 예상해 넘버를 제외, 아이디로만 구분
	createTableSQL :=`
	CREATE TABLE IF NOT EXISTS users(
		ID TEXT NOT NULL UNIQUE,
		PW TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL) 
	if err != nil {
		return nil, fmt.Errorf("table creation error %v", err)
	}
	return &UserDB{db: db}, nil
}

func (udb *UserDB) Close() error {
	return udb.db.Close()
}

func (udb *UserDB) AddUser(ID, PW string) error{
	stmt, err := udb.db.Prepare("INSERT INTO users(ID, PW) VALUES(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID,PW)
	return err
}

func (udb *UserDB) ValidateUser(ID, password string) (bool, error) {
	var storedPassword string
	err := udb.db.QueryRow("SELECT PW FROM users WHERE ID =?", ID).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return storedPassword == password, nil
}

func (udb *UserDB) DeleteUser(username string) error {
	stmt,err := udb.db.Prepare("DELETE FROM users WHERE ID = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username)
	return err
}

func (udb *UserDB) ListUsers() ([]string, error) {
	rows, err := udb.db.Query("SELECT ID FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var ID string
		if err := rows.Scan(&ID); err != nil {
			return nil, err
		}
		users = append(users, ID)
	}
	return users, nil
}
/*
func CompareLoginDB(ID string, PW string) (bool) {

}
*/
