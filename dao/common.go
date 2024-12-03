package dao

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// u can change dsn for different db

type UserDAO struct {
	db  *sql.DB
	dsn string
}

func NewUserDAO(db *sql.DB, dsn string) *UserDAO {
	return &UserDAO{
		db:  db,
		dsn: dsn,
	}
}

func CloseDB(dao *UserDAO) {
	defer dao.db.Close()
	fmt.Printf("Close MySQL db: %s\n", dao.dsn)
}

func ConnectDB(dao *UserDAO) error {
	// open db connection
	var err error
	dao.db, err = sql.Open("mysql", dao.dsn)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// check the connection
	err = dao.db.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Connected to the database successfully!")

	// check the version, just for a query test
	var version string
	err = dao.db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("MySQL version: %s\n", version)
	return nil
}

func CreateTable(dao *UserDAO, sql_str string) {
	_, err := dao.db.Exec(sql_str)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("succ create")
}
