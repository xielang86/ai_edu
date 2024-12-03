package model

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// 这里添加key、secret和url
var (
	appKeyStr = "1312055085303332864"
	secretStr = "dcf112d1049f407babd1b027f5ae8227"
	your_url  = "https://openai.100tal.com/aiimage/"
	your_api  = "handrecognition-pro"
)

const request_body = "request_body"
const Application_json = "application/json"
const Application_x_www_form_urlencoded = "application/x-www-form-urlencoded"
const Multipart_formdata = "multipart/form-data"
const Multipart_formdata_body = "multipartformDataBody"
const Binary = "binary"
const BinaryBody = "BinaryBody"

func bodyFormat(bodyParams map[string]interface{}) (result string) {
	params := url.Values{}
	for k, v := range bodyParams {
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			params.Add(k, v.(string))
			break
		default:
			vJson, _ := json.Marshal(v)
			params.Add(k, string(vJson))
			break
		}
	}
	return params.Encode()
}

// 生成uuid
func GenUUIDv4() (idStr string) {
	return uuid.NewV4().String()
}

func DoPost(url string, contentType string, bodyParams map[string]interface{}) (*http.Response, error) {
	var body io.Reader
	if contentType == Application_x_www_form_urlencoded {
		body = bytes.NewBufferString(bodyFormat(bodyParams))
	} else if contentType == Multipart_formdata {
		//body := bodyParams[Multipart_formdata_body]
	} else if contentType == Binary {
		//body := bytes.NewBuffer(bodyParams[BinaryBody].([]byte))
	} else {
		bytesData, err := json.Marshal(bodyParams)
		if err != nil {
			return nil, errors.New("json.Marshal body_params error")
		}
		body = bytes.NewReader(bytesData)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	return client.Do(req)
}

func StrIsEmpty(str string) (result bool) {
	if len(str) == 0 {
		return true
	}
	return false
}

// 使用HmacSha1计算签名
func HmacSha1(secret string, query string) string {
	secret = secret + "&"
	key := []byte(secret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(query))
	query = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return query
}

// 格式化入参，并计算签名
func GetSignature(
	urlParams map[string]string,
	bodyParams map[string]interface{},
	requestMethod string,
	contentType string,
	accessKeySecret string) (signature string, signatureNonce string) {

	signatureNonce = GenUUIDv4()

	signParams := make(map[string]interface{})
	signParams["signature_nonce"] = signatureNonce

	//只有Application_x_www_form_urlencoded和Application_x_www_form_urlencoded，且是POST/PATCH/PUT时，body才参与鉴权
	if bodyParams != nil && len(bodyParams) != 0 && (requestMethod == "POST" || requestMethod == "PATCH" || requestMethod == "PUT") && (contentType == Application_x_www_form_urlencoded || contentType == Application_json) {
		if contentType == Application_x_www_form_urlencoded {
			bodyParamsEncode := url.Values{}
			for k, v := range bodyParams {
				//str, _ := json.Marshal(v)
				//bodyParamsEncode.Add(k, string(str))

				switch reflect.TypeOf(v).Kind() {
				case reflect.String:
					bodyParamsEncode.Add(k, v.(string))
					break
				default:
					vJson, _ := json.Marshal(v)
					bodyParamsEncode.Add(k, string(vJson))
					break
				}
			}
			//对body进行format，并不是URLEncode
			body := bodyParamsEncode.Encode()
			signParams[request_body] = body
		} else {
			bodyJson, _ := json.Marshal(bodyParams)
			signParams[request_body] = string(bodyJson)
		}
	}

	for k, v := range urlParams {
		signParams[k] = v
	}

	sortKeys := SortMapKey(signParams)

	stringToSign := SingFormat(sortKeys, signParams)
	signature = HmacSha1(accessKeySecret, stringToSign)
	return signature, signatureNonce
}

func GetInterfaceToBytes(key interface{}) (result []byte, err error) {
	var rawRoomIdBuffer bytes.Buffer
	enc := gob.NewEncoder(&rawRoomIdBuffer)
	if err = enc.Encode(key); err != nil {
		return nil, err
	}
	return rawRoomIdBuffer.Bytes(), nil

}

// 计算签名参数格式化
func SingFormat(sortKeys []string, parameters map[string]interface{}) (result string) {
	var paramList []string
	for _, k := range sortKeys {
		v, _ := parameters[k]

		var buffer bytes.Buffer
		buffer.WriteString(k)
		buffer.WriteString("=")
		//vByte,_  := GetInterfaceToBytes(v)
		//println(string(vByte))
		//buffer.Write(vByte)

		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			buffer.WriteString(v.(string))
			break
		default:
			vJson, _ := json.Marshal(v)
			buffer.WriteString(string(vJson))
			break

		}
		paramList = append(paramList, buffer.String())
	}
	return strings.Join(paramList, "&")
}

// 入参格式化为URL参数形式
func UrlFormat(parameters map[string]string) (result string) {
	params := url.Values{}
	for k, v := range parameters {
		params.Add(k, v)
	}
	return params.Encode()
}

// 排序sourceMap，升序
func SortMapKey(sourceMap map[string]interface{}) (sortKeys []string) {
	for key := range sourceMap {
		sortKeys = append(sortKeys, key)
	}
	sort.Strings(sortKeys)
	return sortKeys
}

// 计算签名并发送http请求
func SendRequest(
	accessKeyId string,
	accessKeySecret string,
	timestamp string,
	requestUrl string,
	urlParams map[string]string,
	bodyParams map[string]interface{},
	requestMethod string,
	contentType string) (*http.Response, error) {

	if StrIsEmpty(accessKeyId) {
		return nil, errors.New("参数access_key_id不能为空")
	}
	if StrIsEmpty(accessKeySecret) {
		return nil, errors.New("参数access_key_secret不能为空")
	}
	if StrIsEmpty(timestamp) {
		return nil, errors.New("参数timestamp不能为空")
	}
	if StrIsEmpty(requestUrl) {
		return nil, errors.New("参数requestUrl不能为空")
	}
	if urlParams == nil {
		return nil, errors.New("参数urlParams不能为null,会带回签名，至少做初始化")
	}
	if bodyParams == nil {
		bodyParams = make(map[string]interface{})
	}
	if StrIsEmpty(requestMethod) {
		return nil, errors.New("参数requestMethod不能为空")
	}
	if StrIsEmpty(contentType) {
		return nil, errors.New("参数contentType不能为空")
	}

	urlParams["access_key_id"] = accessKeyId
	urlParams["timestamp"] = timestamp

	signature, signatureNonce := GetSignature(urlParams,
		bodyParams,
		requestMethod,
		contentType,
		accessKeySecret)

	urlParams["signature"] = signature
	urlParams["signature_nonce"] = signatureNonce
	urlParams["timestamp"] = timestamp

	requestUrl = requestUrl + "?" + UrlFormat(urlParams)
	fmt.Println(requestUrl)

	return DoPost(requestUrl, contentType, bodyParams)
}

// 获取东8区时间
func GetCurrentDate() (date string) {
	t := "2006-01-02T15:04:05"
	now := time.Now()
	location, error := time.LoadLocation("Asia/Chongqing")
	if error != nil {
		fmt.Println(error)
	}
	return now.In(location).Format(t)
}

func Post(url string, body interface{}, timeout int64) (*http.Response, error) {
	var err error
	url = your_url + url

	var pbody map[string]interface{}

	p, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(p, &pbody); err != nil {
		return nil, err
	}
	urls := make(map[string]string)
	return SendRequest(
		appKeyStr, secretStr, GetCurrentDate(), url, urls,
		pbody, "POST", "application/json",
	)

}

func request(url string, body interface{}, timeout time.Duration) (error, []byte) {
	var r *http.Response
	var err error

	r, err = Post(url, body, int64(timeout))

	if err != nil {
		return err, []byte("")
	}

	defer r.Body.Close()
	if r.StatusCode != 200 {
		return fmt.Errorf("http code = %d", r.StatusCode), []byte("")
	}

	resp, err := io.ReadAll(r.Body)
	// reader := bufio.NewReader(r.Body)
	// io.Copy(os.Stdout, reader)

	return err, resp
}

type TalOcrReqBody struct {
	ImageUrl    string `json:"image_url"`
	ImageBase64 []byte `json:"image_base64"`
	Function    int    `json:"function"`
}

func demo() {
	image_path := "D:\\work\\edu\\api_server\\test_ocr.png"
	imageBytes, err := os.ReadFile(image_path)
	if err != nil {
		fmt.Println("读取文件错误:", err)
		return
	}
	msg := &TalOcrReqBody{
		ImageBase64: imageBytes,
		Function:    1,
	}

	err, resp := request(your_api, &msg, time.Second*3)

	if err != nil {
		fmt.Printf("requestId : request metre err : %s\n", err.Error())
	}

	fmt.Printf("resp: %v\n", resp)
}

type TalHandText struct {
	Poses [][]int `json:"poses"`
	Probs float64 `json:"probs"`
	Texts string  `json:"texts"`
}

type TalRespData struct {
	HandFormula interface{}   `json:"hand_formula"`
	HandText    []TalHandText `json:"hand_text"`
	LineResults []TalHandText `json:"line_results"`
	SingleBoxes []TalHandText `json:"single_boxes"`
}

type TalRespBody struct {
	Code      int         `json:"code"`
	Data      TalRespData `json:"data"`
	Msg       string      `json:"msg"`
	RequestId string      `json:"requestId"`
}

func GenContent(data *TalRespData, buffer *strings.Builder) {
	for _, text_slice := range data.LineResults {
		buffer.WriteString(text_slice.Texts)
	}
}

// haoweilai api
func GetTalCompResp(lesson string, fileBytes []byte, buffer *strings.Builder) error {
	msg := TalOcrReqBody{
		ImageBase64: fileBytes,
	}

	err, resp := request(your_api, &msg, time.Second*3)
	if err != nil {
		fmt.Println("request tal err:", err)
		return err
	}
	var resp_body TalRespBody
	err = json.Unmarshal(resp, &resp_body)
	if err != nil {
		fmt.Println("parse resp data JSON failed:", err)
		return err
	}
	// 访问和使用解析后的 JSON 数据
	GenContent(&resp_body.Data, buffer)
	// fmt.Println(buffer.String())
	return nil
}
