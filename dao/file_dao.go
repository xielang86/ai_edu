package dao

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type FileInfo struct {
	Id                string
	Name              string `json:"name"`
	CloudPath         string `json:"cloud_path"`
	Md5               string
	RelatedLessonList string `json:"related_lessson_list"`
	UploadUserId      int64
	AuthUserList      string
	CreateTime        int64
}

func CreateFileTable(dao *UserDAO) {
	var sql_str string = `
	CREATE TABLE IF NOT EXISTS  user_file (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
  	name VARCHAR(255),
  	cloud_path	VARCHAR(255),
		md5                 VARCHAR(255),
		related_lesson_list VARCHAR(255),
		upload_user_id BIGINT,
		auth_user_list VARCHAR(255),
		create_time BIGINT);`
	CreateTable(dao, sql_str)
}

func QueryAllFileByUsername(dao *UserDAO, username string, result *[]FileInfo) error {
	var info StudentBaseInfo
	err := QueryStudent(dao, username, "", &info)
	var upload_user_id int64
	upload_user_id = 0

	if err == nil {
		upload_user_id = info.Id
	}
	return QueryAllFileByUserId(dao, upload_user_id, result)
}

func QueryAllFileByUserId(dao *UserDAO, upload_user_id int64, result *[]FileInfo) error {
	query := fmt.Sprintf("select name, cloud_path, md5, related_lesson_list, upload_user_id, auth_user_list from user_file where upload_user_id=%d", upload_user_id)
	rows, err := dao.db.Query(query)
	if err != nil {
		fmt.Println("查询错误:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var file_info FileInfo
		err := rows.Scan(&file_info.Name, &file_info.CloudPath, &file_info.Md5, &file_info.RelatedLessonList,
			file_info.UploadUserId, &file_info.AuthUserList)
		if err != nil {
			fmt.Println("映射错误:", err)
			continue
		}

		// 将结构体添加到数组
		*result = append(*result, file_info)
	}

	return nil
}

func GetFileByMd5(dao *UserDAO, md5 string, info *FileInfo) error {
	query := fmt.Sprintf("select name, cloud_path, md5, related_lesson_list, upload_user_id, auth_user_list from user_file where md5=\"%s\"", md5)
	row := dao.db.QueryRow(query)
	err := row.Scan(&info.Name, &info.CloudPath, &info.Md5, &info.RelatedLessonList,
		&info.UploadUserId, &info.AuthUserList)
	if err == sql.ErrNoRows {
		fmt.Println("没有找到匹配的行")
	} else if err != nil {
		fmt.Println("发生其他错误:", err)
	}
	return err
}

func AddFile(dao *UserDAO, upload_user_id int64, name string, file_path string, md5 string, related_lesson string) error {
	// check by md5
	var info FileInfo
	err := GetFileByMd5(dao, md5, &info)
	if err != sql.ErrNoRows {
		fmt.Println("add failed for", name, err)
		return err
	}

	info.Name = name
	info.UploadUserId = upload_user_id
	info.CloudPath = file_path
	info.Md5 = md5
	info.RelatedLessonList = related_lesson
	info.AuthUserList = fmt.Sprintf("%d", upload_user_id)
	info.CreateTime = time.Now().Unix()
	// TODO(*): add auth to lesson's teacher

	insert_sql := `INSERT INTO user_file 
	(name,cloud_path, md5, related_lesson_list,upload_user_id,auth_user_list,create_time)
	VALUES (?,?,?,?,?,?,?)`
	_, insert_err := dao.db.Exec(insert_sql, info.Name, info.CloudPath, info.Md5, info.RelatedLessonList,
		upload_user_id, info.AuthUserList, info.CreateTime)
	return insert_err
}
