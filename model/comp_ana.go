package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// {
//     "inputs": {
//         "title": "我和她的故事",
//         "grade": "6年级",
//         "type": "中文",
//         "content_type": "记叙文",
//         "requirements": "",
//         "content": "燕子去了，还有再来的时候，桃花谢了，还有再开的时候，杨柳枯了，还有再青的时候，花瓶碎了，还可以修补。可，燕子再来，还是从前的那一群吗？桃花再开，还是从前的那一朵吗？杨柳再青，也不会是从前的那一缕了，花瓶修补，那不也有裂痕吗？友情就如同花瓶，破碎了还可以修补，可却再也不是曾经那份纯粹，真挚的友谊了。\n        公园里，杨柳依依，如同纱帘般垂在水面上，偶尔有小鸟轻点下水面，漾开淡淡波纹，半天透明半天云，一切都美好的不像话，风缱绻而温柔，拂过我的发梢，温温软软的，可我却面色铁青，生气的坐在椅子上，身旁坐着不断道歉的依依。我们原本约定10:30在公园见面，可这都11:46了，依依才姗姗来迟，我多少有些愤怒，但我这脾气来的快去得也快，不一会儿就冷静下来了。我拉着依依去到湖边的亭子里坐下，迫不及待的打开书包，没错，这才是我此行的真正目的，和依依交流自己最近读的书，忽然，我察觉身侧有些不对劲，我抬了抬头，便见依依嘴巴抿了抿，张开，又闭上，像要说什么，又没说，脸红红的，像熟透了的西红柿，手指搓着衣角，掌心有些亮亮的。看她这样，我心中已猜出了八九分，我神色一凛，试探的开口:“你……是不是……没带书。”她迟疑了一下，然后有些愧疚的点了点头，刹那间，我的心头升起一团火焰，先前的怒火叠加此时的愤怒，我质问到:“你怎么能这样呢？分享书是你提出来的，时间是你定的，迟到就算了，还把书忘了，你真的太过分了！”依依赶忙上来抓我的胳膊，我侧身躲开，闪到一旁，她一怔，然后哭着跑开了，我也没管，直接回家了。一周后，她来跟我道歉，我已经消气，可话到嘴边又变了，我拒绝了他，以为她会来找我，可她再也没来过。\n       我们的友情就如同夕阳的最后一抹余晖，我想抓却抓不住。后来，我也曾去过那个公园，那里承载着我们的不少回忆，从天真烂漫的奔跑到娴静的漫步，我们终究结束了这段友谊。杨柳依然温柔，小鸟依旧俏皮，一切依旧美得不像话，可湖边再也没有两个小姑娘了。我和依依，也只能渐行渐远。"
//     },
//     "response_mode": "blocking",
//     "user": "zhizhi-user"
// }

type CompInput struct {
	Title        string `json:"title"`
	Grade        string `json:"grade"`
	Type         string `json:"type"`
	ContentType  string `json:"content_type"`
	Requirements string `json:"requirements"`
	Content      string `json:"content"`
}

type CompReqBody struct {
	Inputs       CompInput `json:"inputs"`
	ResponseMode string    `json:"response_mode"`
	User         string    `json:"user"`
}

// {
// "task_id": "9f87907e-7214-4d17-9bd1-4fe7fcc8465c",
// "workflow_run_id": "83a6f024-937d-4d61-ba06-483b0f5b2903",
// "data": {
// "id": "83a6f024-937d-4d61-ba06-483b0f5b2903",
// "workflow_id": "f799862e-8f45-4add-a7d0-881d0036a979",
// "status": "succeeded",
// "outputs": {
// "result": "好的，让我们开始对这篇作文进行点评吧。\n\n## 作文要求\n（由于原文未提供作文要求的具体内容，这里假设作文要求为记叙文，题目为《我的一天》，字数不少于300字，文体要求内容连贯、情节完整、主题鲜明。）\n\n## 学生作文\n（由于原文未提供学生作文内容，这里假设学生作文如下：）\n\n今天是一个阳光明媚的日子。早上，我睡了一个懒觉，起床后吃了妈妈做的早餐，然后坐着公交车去学校。在公交车上，我看到了一个老爷爷拿着沉重的行李上车，于是我就主动让座。到了学校，我参加了早上的晨跑活动，然后开始了一天的学习。下午放学后，我和同学们去了公园玩，还拍了许多照片。晚上，我回家后写完作业，看了一会儿书，就早早地睡觉了。通过这一天，我感到十分充实和快乐。\n\n## 亮点解读\n\n这篇作文主要讲述了作者的一天，从早上起床到晚上睡觉，记录了一天的生活琐事和感受。作文内容连贯自然，能够清晰地反映出一天的经历。亮点如下：\n1. 整体内容连贯，叙述了一天的经历，逻辑性较强。\n2. 通过描述让座的细节，展现了作者的善良品质。\n3. 作文开头的“今天是一个阳光明媚的日子”很好地营造了氛围，引出了下文的描写。\n4. 描述了上下学路上、学校生活、放学后的活动，内容充实。\n5. 文章结尾“通过这一天，我感到十分充实和快乐”表达了作者的内心感受，升华了主题。\n\n## 改进建议\n\n1. 写作要求为记叙文，题目为《我的一天》，内容要求连贯、情节完整、主题鲜明。这篇作文基本上符合要求，但内容可以更加丰富，包括更多的细节描写。\n2. 作文中的过渡衔接比较通顺，但可以在关键情节间添加适当的过渡句，如“早晨我吃完早餐，走进了公交车”。\n3. 作文中的思想积极健康，但可以进一步深化，深入描述每个细节背后的感受，如“看到老爷爷拿着沉重的行李上车，我心里很不是滋味，于是主动让座给他，那一刻我感到非常开心”。\n4. 可以在作文中增加一些具体的描写手法，如用比喻、拟人等修辞手法，进一步丰富内容，如“阳光像一位慈祥的母亲，温柔地抚摸着大地”。\n\n## 量规评价\n\n| 评价维度     | 评价指标         | 评价级别 | 评价依据                                         |\n|-------------|----------------|--------|------------------------------------------------|\n| 错别字      |                | 良好   | 作文中有少量的错别字，主要为常见的错别字，少量复杂词汇的笔误，不影响整体阅读和理解。 |\n| 标点        |                | 良好   | 作文中能正确、规范使用常用标点符号，包括逗号、句号、感叹号等。                 |\n| 词汇运用    |                | 良好   | 作文中能运用较为丰富的词语搭配，有少量重复使用现象；大部分选词恰当，存在少量误用现象。   |\n| 句式结构    |                | 良好   | 作文中能运用较丰富的句式结构，但有部分句式较为简单，稍显重复；能够运用一些基础的修辞手法。 |\n| 语法        |                | 良好   | 作文中掌握常用语法，语篇中存在一些高级语法的使用不当，但对语篇理解和可读性几乎无影响。 |\n| 语言风格    |                | 良好   | 作文的语言风格较好地符合文体要求，整体表达力较强，能够有效传达作者的观点和意图。   |\n| 结构组织    |                | 良好   | 作文结构清晰、完整且符合文体要求，段落衔接连贯，整体逻辑性较强。                 |\n| 语义连接    |                | 良好   | 作文能使用连接词来连接简单句、复合句和段落，表达较复杂的语义关系，使语篇结构更加清晰。      |\n| 语篇结构    |                | 良好   | 作文结构清晰完整，有较好的连贯性，有较为完整的开头、情节发展和结尾，内容组织合理。       |\n| 语篇主题    |                | 良好   | 作文主题较明确，具有较高的思考价值。                                             |\n| 记叙主体    |                | 良好   | 记叙主体较为明确，具有较好的辨识度；记叙主体的特征和行为与语篇主题相符。                     |\n| 故事情节    |                | 良好   | 故事情节完整，有较充分的细节支撑；情节发展逻辑清晰，连贯性好。                                 |\n| 综合得分     |                | 80%    | 作文整体表现良好，有较多的优点，但仍有提升空间。                                             |\n\n## 总结\n\n这篇作文整体上内容丰富而连贯，结构清晰，表现出了一个普通但充实的一天。记叙主体明确，主题也较为鲜明。不过，作文中的一些细节描写和情感表达可以更加深入，丰富内容，使作文更加生动。另外，适当增加一些写作技巧，如描写手法和修辞手法，能够进一步提升文章的感染力。继续保持优点，努力改进不足，相信你的作文会越来越精彩！"
// },
// "error": null,
// "elapsed_time": 17.987086632987484,
// "total_tokens": 3727,
// "total_steps": 7,
// "created_at": 1732888517,
// "finished_at": 1732888535
// }
// }

type CompOutput struct {
	// Result string `json:"result"`
	Evaluation string `json:"evaluation"`
	Example    string `json:"example"`
}

type CompRespData struct {
	Id         string     `json:"id"`
	WorkflowId string     `json:"workflow_id"`
	Status     string     `json:"status"`
	Outputs    CompOutput `json:"outputs"`
}

type CompRespBody struct {
	TaskId        string       `json:"task_id"`
	WorkflowRunId string       `json:"workflow_run_id"`
	Data          CompRespData `json:""`
}

func GetCompRawResp(req_body *CompReqBody, resp_body *CompRespBody) error {
	comp_url := "https://demo.eye8.top/v1/workflows/run"
	// secrete_key := "app-EC4t4gDjaaakzWFpbQCImYOe"
	json_data, err := json.Marshal(*req_body)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
		return err
	}

	// 创建一个POST请求
	req, err := http.NewRequest("POST", comp_url, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return err
	}

	// 设置请求头，例如设置Content-Type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	// TODO(xl): mv to config
	req.Header.Set("authorization", "Bearer app-EC4t4gDjaaakzWFpbQCImYOe")
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read resp failed:%v", err)
		return err
	}
	// fmt.Printf("响应状态码: %d\n响应体: %s\n", resp.StatusCode, body)

	err = json.Unmarshal(body, resp_body)
	if err != nil {
		fmt.Println("unmarshal failed!:", err)
		return err
	}
	return nil
}

func GetCompAna(title string, grade string, language string, content_type string, content string) (string, string) {
	comp_input := CompInput{Title: title, Grade: grade, Type: language, ContentType: content_type, Requirements: ""}
	req_body := CompReqBody{Inputs: comp_input, ResponseMode: "blocking", User: "zhizhi-user"}
	req_body.Inputs.Content = content

	var resp_body CompRespBody
	err := GetCompRawResp(&req_body, &resp_body)
	if err != nil {
		fmt.Printf("failed to get raw resp for %v", req_body)
		return "", ""
	}

	return resp_body.Data.Outputs.Evaluation, resp_body.Data.Outputs.Example
}
