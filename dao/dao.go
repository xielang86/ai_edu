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

type EduKnowledge struct {
	level1        string
	level2        string
	level3        string
	level4        string
	level5_prompt string // maybe it's json
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

func QueryEduKnowledge(dao *UserDAO, query string, result *EduKnowledge) error {
	row := dao.db.QueryRow(query)

	err := row.Scan(&result.level4, &result.level5_prompt)
	if err != nil {
		return err
	}
	return nil
}

func CreateKnowledgeTable(dao *UserDAO) {
	var sql string = `
  CREATE TABLE IF NOT EXISTS  knowledge_edu.en_knowledge_point (
    id INT AUTO_INCREMENT PRIMARY KEY,
    level1 VARCHAR(255),
    level2 VARCHAR(255),
    level3 VARCHAR(255),
    level4 VARCHAR(255),
    level5 VARCHAR(255),
    level5_prompt VARCHAR(1024)
  );
  `
	_, err := dao.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("succ create")
}

// func main() {
// 	dao := &UserDAO{nil, kEduKnowledgeDB}
// 	ConnectDB(dao)
// 	var query string = "select level4, level5 from knowledge_edu.en_knowledge_point where level4=?"
// 	result := &EduKnowledge{}
// 	CreateKnowledgeTable(dao)
// 	QueryEduKnowledge(dao, query, result)
// 	CloseDB(dao)
// }
