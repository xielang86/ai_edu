package api

import (
	"api_server/dao"
	"fmt"
	"net/http"
)

var kEduKnowledgeDB string = "xielang:lang.xie86@(127.0.0.1:3306)/knowledge_edu"

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	// get knowledge
	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	var query string = "select level4, level5 from knowledge_edu.en_knowledge_point where level4=?"
	result := &dao.EduKnowledge{}
	dao.CreateKnowledgeTable(mydao)
	dao.QueryEduKnowledge(mydao, query, result)
	dao.CloseDB(mydao)

	fmt.Println("Hello, World!")
}
