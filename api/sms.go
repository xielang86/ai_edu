package api

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// generateVerificationCode 生成6位数字验证码
func generateVerificationCode() string {
	min := 100000
	max := 999999
	return strconv.Itoa(min + rand.Intn(max-min))
}

// SendSMS 使用腾讯云短信服务发送短信
// phoneNumbers := []string{"+86138xxxx5678"} // 替换为真实的手机号码，支持多个手机号码同时发送
// templateID := "your_template_id"         // 替换为你在腾讯云短信服务配置的短信模板ID
// templateParam := []string{`{"code":"123456"}`}      // 替换为真实的模板参数，格式要与短信模板定义的参数格式一致
func SendSMS(phoneNumbers []string, templateID string, templateParam []string) error {
	// 替换为你自己的腾讯云短信服务配置信息
	cred := common.NewCredential("your_secret_id", "your_secret_key")
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, err := sms.NewClient(cred, "ap-guangzhou", cpf)
	if err != nil {
		return err
	}

	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr("your_app_id")
	request.SignName = common.StringPtr("your_sign_name")
	request.TemplateId = common.StringPtr(templateID)
	request.TemplateParamSet = common.StringPtrs(templateParam)
	request.PhoneNumberSet = common.StringPtrs(phoneNumbers)

	response, err := client.SendSms(request)
	if err != nil {
		return err
	}

	if response.Response.SendStatusSet == nil || len(response.Response.SendStatusSet) == 0 {
		return errors.NewTencentCloudSDKError("Empty response", "", "")
	}

	for _, status := range response.Response.SendStatusSet {
		if status.Code == nil || *status.Code != "Ok" {
			log.Printf("发送短信到 %s 失败，错误码：%s，错误信息：%s\n", *status.PhoneNumber, *status.Code, *status.Message)
		}
	}

	return nil
}

// VerificationCodeInfo 存储验证码相关信息
type VerificationCodeInfo struct {
	Code           string
	ExpirationTime int64
}

var (
	verificationCodes = make(map[string]*VerificationCodeInfo)
	mutex             sync.Mutex
)

func SendVerificationCodeHandler(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "请求参数错误"})
		return
	}
	phone := req.Phone
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "手机号不能为空"})
		return
	}
	code := generateVerificationCode()
	expirationTime := time.Now().Unix() + 120 // 验证码有效期120秒（2分钟），可按需调整

	// 调用阿里云短信服务发送验证码（此处需替换为真实有效的配置）
	// err := SendSMS({phone}, "", {code})
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "验证码发送失败"})
	// 	return
	// }

	mutex.Lock()
	verificationCodes[phone] = &VerificationCodeInfo{
		Code:           code,
		ExpirationTime: expirationTime,
	}
	mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "success", "code": code})
}

func VerificationCodeHandler(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "请求参数错误"})
		return
	}
	phone := req.Phone
	inputCode := req.Code
	if phone == "" || inputCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "手机号和验证码不能为空"})
		return
	}

	mutex.Lock()
	codeInfo, ok := verificationCodes[phone]
	mutex.Unlock()
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "未发送验证码或验证码已过期"})
		return
	}

	if codeInfo.Code == inputCode && time.Now().Unix() < codeInfo.ExpirationTime {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "验证码验证成功"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "验证码错误或已过期"})
	}
}
