package api

import (
	"api_server/dao"
	"api_server/model"
	"bytes"
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	// "github.com/tencentcloud/cos-go-sdk-v5"
)

var (
	local_file_root          = "./local_files" // just for testing
	tencent_cloud_cos_root   = ""
	tencent_cloud_cos_key    = ""
	tencent_cloud_cos_secret = ""
)

func FileCenterHandler(w http.ResponseWriter, r *http.Request) {
	PageHandler(w, "./pages/file_center.html")
}

func FilePageHandler(w http.ResponseWriter, r *http.Request) {
	PageHandler(w, "./pages/article.html")
}

func FileCenterPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}
	responseData := ResponseData{
		Status:  "success",
		Message: "操作成功",
	}

	PostResponse(w, responseData)
}

func SaveFileToLocal(fileData []byte, filePath string) error {
	err := os.WriteFile(filePath, fileData, 0644)
	if err != nil {
		log.Printf("保存文件到本地失败: %v", err)
		return err
	}
	return nil
}

// UploadFileToTencentCOS 将文件上传到腾讯云COS
// func UploadFileToTencentCOS(fileData []byte, bucket, objectKey string) error {
// 	// 填入你的腾讯云 SecretId 和 SecretKey
// 	secretId := "your_secret_id"
// 	secretKey := "your_secret_key"
// 	// 填入你的腾讯云COS所在的地域
// 	region := "your_region"
//
// 	u, _ := url.Parse(fmt.Sprintf("https://cos-%s.myqcloud.com", region))
// 	b := &cos.BaseURL{BucketURL: u}
// 	client := cos.NewClient(b, &http.Client{
// 		Transport: &cos.AuthorizationTransport{
// 			SecretID:  secretId,
// 			SecretKey: secretKey,
// 		},
// 	})
//
// 	fileData, err := ioutil.ReadFile(localFilePath)
// 	if err != nil {
// 		log.Printf("读取本地文件失败: %v", err)
// 		return err
// 	}
//
// 	_, err = client.Object.Put(context.Background(), objectKey, fileData, nil)
// 	if err != nil {
// 		log.Printf("上传文件到腾讯云COS失败: %v", err)
// 		return err
// 	}
//
// 	return nil
// }

type OneFile struct {
	FileBytes []byte
	Filename  string
}

type UploadFileInfo struct {
	Username    string
	Title       string
	Lesson      string // language, trick fu
	Grade       string
	ContentType string
	Files       []OneFile
}

func UploadFile(upload_info *UploadFileInfo) error {
	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	// query for user_id
	// TODO(*): just for test username, risk for attacking
	var info dao.UserInfo
	info.Name = upload_info.Username
	info.Id = 0
	if upload_info.Username != "default" {
		err := dao.QueryUser(mydao, upload_info.Username, "", &info)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("user not exist or login eror for user %s", upload_info.Username)
			}
			return err
		}
	} else {
		fmt.Println("default user for test")
	}

	// cal md5
	for _, one_file := range upload_info.Files {
		hash := md5.Sum(one_file.FileBytes)
		md5_str := fmt.Sprintf("%x", hash)
		var file_info dao.FileInfo
		err := dao.GetFileByMd5(mydao, md5_str, &file_info)
		if err == nil {
			fmt.Printf("do not save dedup file %s\n", one_file.Filename)
			return nil
		}

		// save to local	path
		uniq_name := fmt.Sprintf("%d_%s", time.Now().Unix(), one_file.Filename)
		local_file_path := fmt.Sprintf("%s/%s", local_file_root, uniq_name)
		err = SaveFileToLocal(one_file.FileBytes, local_file_path)
		if err != nil {
			fmt.Printf("store file %s failed\n", one_file.Filename)
			return err
		}
		// insert record to db
		// strick use lesson to store the info for same article
		lesson_str := upload_info.Lesson
		if len(lesson_str) < 1 || lesson_str == "中文" || lesson_str == "英文" {
			lesson_str = upload_info.Files[0].Filename
		}
		err = dao.AddFile(mydao, info.Id, upload_info.Username, local_file_path, md5_str, lesson_str)
		if err != nil {
			fmt.Printf("add file info %s to db failed for user %s\n", one_file.Filename, upload_info.Username)
		}
	}
	return nil
}

func ParseFileData(r *http.Request) (UploadFileInfo, error) {
	var info UploadFileInfo
	info.ContentType = r.FormValue("content_type")
	info.Username = r.FormValue("username")
	info.Lesson = r.FormValue("language")
	info.Grade = r.FormValue("grade")
	if len(info.Grade) < 3 {
		info.Grade = fmt.Sprintf("%s年级", info.Grade)
		fmt.Printf("change grade to %s", info.Grade)
	}

	err := r.ParseMultipartForm(32 << 20) // 32MB 大小限制
	if err != nil {
		return info, err
	}
	// file, header, err := r.FormFile("file")
	// defer file.Close()
	files := r.MultipartForm.File["images"]

	// 遍历图片文件并保存
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		var one_file OneFile
		one_file.Filename = fileHeader.Filename
		one_file.FileBytes, err = io.ReadAll(file)
		if err != nil {
			fmt.Printf("read file %s failed for user %s", one_file.Filename, info.Username)
			continue
		}
		info.Files = append(info.Files, one_file)
	}
	return info, err
}

type CompData struct {
	Title                string      `json:"title"`
	Name                 string      `json:"name"`
	Grade                string      `json:"grade"`
	Type                 string      `json:"type"`
	Requirements         string      `json:"requirements"`
	ContentType          string      `json:"content_type"`
	Date                 string      `json:"date"`
	Content              string      `json:"content"`
	Example              string      `json:"example"`
	OverallScore         string      `json:"overall_score"`
	OverallGrade         string      `json:"overall_grade"`         // "合格", 等级和分数只给一个结果
	OverallAdvantages    string      `json:"overall_Advantages"`    //优点
	OverallDisadvantages string      `json:"overall_Disadvantages"` //缺点
	OverallRemarks       string      `json:"overall_remarks"`       // "总体评价的得分和等级，表示本次作文在所有维度上的综合表现。"
	LanguageCriteria     interface{} `json:"language_criteria"`
	WritingCriteria      interface{} `json:"writing_criteria"`
	TextContent          interface{} `json:"text_content"`
}

func UploadAndOcrHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	upload_info, err := ParseFileData(r)
	if err != nil {
		http.Error(w, "parse file info failed ", http.StatusInternalServerError)
		return
	}

	err = UploadFile(&upload_info)
	if err != nil {
		fmt.Printf("upload file failed for user=%s would continue do ocr", upload_info.Username)
	}

	content := "燕子去了，还有再来的时候，桃花谢了，还有再开的时候，杨柳枯了，还有再青的时候，花瓶碎了，还可以修补。可，燕子再来，还是从前的那一群吗？桃花再开，还是从前的那一朵吗？杨柳再青，也不会是从前的那一缕了，花瓶修补，那不也有裂痕吗？友情就如同花瓶，破碎了还可以修补，可却再也不是曾经那份纯粹，真挚的友谊了。\n        公园里，杨柳依依，如同纱帘般垂在水面上，偶尔有小鸟轻点下水面，漾开淡淡波纹，半天透明半天云，一切都美好的不像话，风缱绻而温柔，拂过我的发梢，温温软软的，可我却面色铁青，生气的坐在椅子上，身旁坐着不断道歉的依依。我们原本约定10:30在公园见面，可这都11:46了，依依才姗姗来迟，我多少有些愤怒，但我这脾气来的快去得也快，不一会儿就冷静下来了。我拉着依依去到湖边的亭子里坐下，迫不及待的打开书包，没错，这才是我此行的真正目的，和依依交流自己最近读的书，忽然，我察觉身侧有些不对劲，我抬了抬头，便见依依嘴巴抿了抿，张开，又闭上，像要说什么，又没说，脸红红的，像熟透了的西红柿，手指搓着衣角，掌心有些亮亮的。看她这样，我心中已猜出了八九分，我神色一凛，试探的开口:“你……是不是……没带书。”她迟疑了一下，然后有些愧疚的点了点头，刹那间，我的心头升起一团火焰，先前的怒火叠加此时的愤怒，我质问到:“你怎么能这样呢？分享书是你提出来的，时间是你定的，迟到就算了，还把书忘了，你真的太过分了！”依依赶忙上来抓我的胳膊，我侧身躲开，闪到一旁，她一怔，然后哭着跑开了，我也没管，直接回家了。一周后，她来跟我道歉，我已经消气，可话到嘴边又变了，我拒绝了他，以为她会来找我，可她再也没来过。\n       我们的友情就如同夕阳的最后一抹余晖，我想抓却抓不住。后来，我也曾去过那个公园，那里承载着我们的不少回忆，从天真烂漫的奔跑到娴静的漫步，我们终究结束了这段友谊。杨柳依然温柔，小鸟依旧俏皮，一切依旧美得不像话，可湖边再也没有两个小姑娘了。我和依依，也只能渐行渐远。"
	buf := bytes.NewBufferString("")
	for _, onefile := range upload_info.Files {
		// do ocr
		result, ocr_err := model.GetTalCompResp(upload_info.Lesson, onefile.FileBytes)
		if ocr_err != nil {
			fmt.Printf("do ocr failed1 for %s\n", onefile.Filename)
			continue
		}
		buf.WriteString(result)
	}
	if buf.Len() > 8 {
		content = buf.String()
	}
	// do ana, test case
	fmt.Printf("info: %s, %s, %s, filenum=%d, content_len=%d",
		upload_info.Grade, upload_info.Lesson, upload_info.ContentType, len(upload_info.Files), len(content))
	evaluation, example := model.GetCompAna("test title", upload_info.Grade, upload_info.Lesson, upload_info.ContentType, content)

	var comp_data CompData
	comp_data.OverallRemarks = evaluation
	comp_data.Content = content
	comp_data.Type = upload_info.Lesson
	comp_data.ContentType = upload_info.ContentType
	comp_data.Grade = upload_info.Grade
	comp_data.Example = example
	comp_data.Name = upload_info.Username
	comp_data.Title = upload_info.Title

	responseData := ResponseData{
		Status:  "success",
		Message: fmt.Sprintf("comp ana succ"),
	}

	responseData.Data = comp_data
	// reformat the result to json
	PostResponse(w, responseData)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}
	upload_info, err := ParseFileData(r)
	if err != nil {
		http.Error(w, "parse file info failed ", http.StatusInternalServerError)
		return
	}

	err = UploadFile(&upload_info)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := ResponseData{
		Status:  "success",
		Message: "upload failed!",
	}

	if len(upload_info.Files) > 0 {
		buf := bytes.NewBufferString("")
		for _, onefile := range upload_info.Files {
			buf.WriteString(onefile.Filename)
			buf.WriteString(",")
		}
		responseData.Message = fmt.Sprintf("upload succ for %s", buf.String())
	}

	PostResponse(w, responseData)
}

type FileData struct {
	Folder  interface{} `json:"folder"`
	AllFile interface{} `json:"file"`
}

func GetAllFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	username := r.FormValue("username")
	var result []dao.FileInfo
	dao.QueryAllFileByUsername(mydao, username, &result)

	responseData := ResponseData{
		Status:  "success",
		Message: fmt.Sprintf("get files succ for user %s", username),
	}

	var data FileData
	data.AllFile = result
	responseData.Data = data
	PostResponse(w, responseData)
}
