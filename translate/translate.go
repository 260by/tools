package translate
// https://github.com/zmj1368/GoogleTranslateApi
import (
	"net/http"
	"log"
	"errors"
	"strings"
	"io/ioutil"
)

type TranslateInfo struct {
	ch string
	en string
}

// TranslateChToEnApi 中文到英文
func TranslateChToEnApi(str, tk string) (TranslateInfo, error) {
	var info TranslateInfo

	// 抓取翻译信息
	baseURL := "https://translate.google.cn/translate_a/single?client=t&sl=zh-CN&tl=en&hl=zh-CN&dt=at&dt=bd&dt=ex&dt=ld&dt=md&dt=qca&dt=rw&dt=rm&dt=ss&dt=t&ie=UTF-8&oe=UTF-8&source=btn&ssel=5&tsel=5&kc=0&tk=" + tk + "&q="
	url := baseURL + FormatStr(str)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return info, err
	}
	// 关闭链接
	defer resp.Body.Close()
	// 如果状态码不等于200
	if resp.StatusCode != http.StatusOK {
		return info, errors.New(resp.Status) // 例如: 404 Not Found
	}

	body, err := ioutil.ReadAll(resp.Body) //此处可增加输入过滤
	if err != nil {
		//fmt.Println("错误:读取body", err)
		return info, err
	}
	// 分割字符串
	info.en, info.ch = DealContentFunc(string(body))

	return info, nil
}

// DealContentFunc 获取字符内容处理函数 按照顺序输出处理结果
func DealContentFunc(content string) (string, string) {
	var (
		res1 string
		res2 string
	)
	tmpArr := strings.Split(content[3:], "],[")

	for _, v := range tmpArr {
		// 如果前4个字母是null 就可以break
		if v[:4] == "null" {
			break
		}
		// 分割字符串拼接结果
		tmpArr := strings.Split(v, ",")
		res1 += string(tmpArr[0][1:len(tmpArr[0])-1])
		res2 += string(tmpArr[1][1:len(tmpArr[1])-1])
	}
	return res1, res2
}

func TranslateEnToChApi(str, tk string) (TranslateInfo, error) {
	var info TranslateInfo

	// 抓取翻译信息
	baseURL := "https://translate.google.cn/translate_a/single?client=t&sl=en&tl=zh-CN&hl=zh-CN&dt=at&dt=bd&dt=ex&dt=ld&dt=md&dt=qca&dt=rw&dt=rm&dt=ss&dt=t&ie=UTF-8&oe=UTF-8&source=btn&ssel=0&tsel=0&kc=0&tk=" + tk + "&q="
	url := baseURL + FormatStr(str)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return info, err
	}
	// 关闭链接
	defer resp.Body.Close()
	// 如果状态码不等于200
	if resp.StatusCode != http.StatusOK {
		log.Fatal(errors.New(resp.Status))
		return info, errors.New(resp.Status) // 例如: 404 Not Found
	}

	body, err := ioutil.ReadAll(resp.Body) //此处可增加输入过滤
	if err != nil {
		//fmt.Println("错误:读取body", err)
		return info, err
	}
	// 分割字符串
	info.ch, info.en = DealContentFunc(string(body))

	return info, nil
}

func FormatStr(str string) string {
	// 将url空格转换为 %20
	str = strings.TrimSpace(str)
	str = strings.Replace(str, " ", "%20", -1)
	return str
}

// 中文翻译成英文
func GetTranslateChToEnContent(str string) (string, error) {
	var content = ""

	tk := GetGoogleToken(str)
	info, err := TranslateChToEnApi(str, tk)
	if err != nil {
		return "", nil
	}
	content = info.en
	return content, nil
}

// 英文翻译成中文
func GetTranslateEnToChContent(str string) (string, error) {
	var content = ""

	tk := GetGoogleToken(str)
	info, err := TranslateEnToChApi(str, tk)
	if err != nil {
		return "", nil
	}
	content = info.ch
	return content, nil
}

func GoGetTranslate(lineArr map[int]string, start, end int, ch chan map[int]string) {
	var tk = make(map[int]string)
	var res = make(map[int]string)
	for {
		if start > end {
			break
		}
		tk[start] = GetGoogleToken(lineArr[start])
		info, err := TranslateEnToChApi(lineArr[start], tk[start])
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println("key is ", start, "ch is ", info.ch)
		res[start] = info.ch
		start++
	}
	ch <- res
}
