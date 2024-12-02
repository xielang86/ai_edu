package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestModelApi(t *testing.T) {
	req_body := SiliconFlowReqBody{
		Model:    "alibaba/Qwen1.5-110B-Chat",
		Messages: [3]GPTMessage{},
	}
	req_body.Messages[0] = GPTMessage{"user", "如何走出情劫"}
	var resp_body SiliconFlowRespBody
	err := GetRawResp(req_body, &resp_body)
	if err != nil {
		t.Errorf("failed to test raw resp for %v", req_body)
	}
	fmt.Println(resp_body)
}

func TestSingleAnswer(t *testing.T) {
	var query = "如何走出情劫"
	var ans = GetSingleAnswer(query, "Qwen/Qwen2-7B-Instruct")
	if len(ans) == 0 {
		t.Errorf("empty result for %s", query)
	}
	fmt.Println(ans)
}

func TestSiliconFlow(t *testing.T) {
	// 目标URL
	url := "https://api.siliconflow.cn/v1/chat/completions"

	req_body := SiliconFlowReqBody{
		Model:    "alibaba/Qwen1.5-110B-Chat",
		Messages: [3]GPTMessage{},
	}
	req_body.Messages[0] = GPTMessage{"user", "如何走出情劫"}

	json_data, err := json.Marshal(req_body)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
		return
	}

	// data := []byte(`{"model":"alibaba/Qwen1.5-110B-Chat", "messages":[{"role": "user", "content": "抛砖引玉是什么意思呀"}]}`)
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))

	// 创建一个POST请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 设置请求头，例如设置Content-Type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", "Bearer sk-acgzjedfzicxmxzuprvduikxfaoenzdmrkxyyyimxqvesppj")
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}
	fmt.Printf("响应状态码: %d\n响应体: %s\n", resp.StatusCode, body)

	var resp_body SiliconFlowRespBody
	json.Unmarshal(body, &resp_body)
	if len(resp_body.Choices) > 0 {
		fmt.Printf("answer: %s\n", resp_body.Choices[0].Message.Content)
	}
}

func TestCompAna(t *testing.T) {
	title := "我和她的故事"
	grade := "6年级"
	language := "中文"
	content_type := "记叙文"
	content := "燕子去了，还有再来的时候，桃花谢了，还有再开的时候，杨柳枯了，还有再青的时候，花瓶碎了，还可以修补。可，燕子再来，还是从前的那一群吗？桃花再开，还是从前的那一朵吗？杨柳再青，也不会是从前的那一缕了，花瓶修补，那不也有裂痕吗？友情就如同花瓶，破碎了还可以修补，可却再也不是曾经那份纯粹，真挚的友谊了。\n        公园里，杨柳依依，如同纱帘般垂在水面上，偶尔有小鸟轻点下水面，漾开淡淡波纹，半天透明半天云，一切都美好的不像话，风缱绻而温柔，拂过我的发梢，温温软软的，可我却面色铁青，生气的坐在椅子上，身旁坐着不断道歉的依依。我们原本约定10:30在公园见面，可这都11:46了，依依才姗姗来迟，我多少有些愤怒，但我这脾气来的快去得也快，不一会儿就冷静下来了。我拉着依依去到湖边的亭子里坐下，迫不及待的打开书包，没错，这才是我此行的真正目的，和依依交流自己最近读的书，忽然，我察觉身侧有些不对劲，我抬了抬头，便见依依嘴巴抿了抿，张开，又闭上，像要说什么，又没说，脸红红的，像熟透了的西红柿，手指搓着衣角，掌心有些亮亮的。看她这样，我心中已猜出了八九分，我神色一凛，试探的开口:“你……是不是……没带书。”她迟疑了一下，然后有些愧疚的点了点头，刹那间，我的心头升起一团火焰，先前的怒火叠加此时的愤怒，我质问到:“你怎么能这样呢？分享书是你提出来的，时间是你定的，迟到就算了，还把书忘了，你真的太过分了！”依依赶忙上来抓我的胳膊，我侧身躲开，闪到一旁，她一怔，然后哭着跑开了，我也没管，直接回家了。一周后，她来跟我道歉，我已经消气，可话到嘴边又变了，我拒绝了他，以为她会来找我，可她再也没来过。\n       我们的友情就如同夕阳的最后一抹余晖，我想抓却抓不住。后来，我也曾去过那个公园，那里承载着我们的不少回忆，从天真烂漫的奔跑到娴静的漫步，我们终究结束了这段友谊。杨柳依然温柔，小鸟依旧俏皮，一切依旧美得不像话，可湖边再也没有两个小姑娘了。我和依依，也只能渐行渐远。"
	result := GetCompAna(title, grade, language, content_type, content)
	fmt.Println(result)
}
