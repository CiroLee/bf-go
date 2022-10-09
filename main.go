package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	// custom config
	"github.com/CiroLee/bf-go/config"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

type TransProps struct {
	from  string
	to    string
	query string
}
type TransResult struct {
	Src string
	Dst string
}
type TransResponse struct {
	From        string
	To          string
	TransResult []TransResult `json:"trans_result"`
	ErrorCode   string        `json:"error_code"`
	ErrorMsg    string        `json:"error_msg"`
}

var versonStr = "0.0.1"

var greenPrint = color.New(color.FgGreen, color.Bold)
var tipPrint = color.New(color.FgCyan)
var wrongPrint = color.New(color.FgRed)

var loading = spinner.New(spinner.CharSets[9], 100*time.Millisecond)

func md5Str(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

// 生成签名和随机值
func getSign(query string) (salt, sign string) {
	appid := config.Secrets.Appid
	key := config.Secrets.Key
	salt = fmt.Sprintf("%x", time.Now().Unix())
	sign = appid + query + salt + key

	return salt, md5Str(sign)
}

func getTranalation(props TransProps) {
	salt, sign := getSign(props.query)
	payload := url.Values{
		"q":     {props.query},
		"from":  {props.from},
		"to":    {props.to},
		"salt":  {salt},
		"sign":  {sign},
		"appid": {config.Secrets.Appid},
	}
	loading.Start()
	res, err := http.PostForm("https://fanyi-api.baidu.com/api/trans/vip/translate", payload)

	if err != nil {
		log.Fatal("Error:", err)
	} else {
		result := &TransResponse{}
		renderResut(res, result)
	}

	loading.Stop()

	defer res.Body.Close()

}

// 渲染查询结果
func renderResut(response *http.Response, result *TransResponse) {
	body, _ := io.ReadAll(response.Body)
	jsonError := json.Unmarshal(body, result)
	if jsonError != nil {
		log.Fatal("json.Unmarshal", jsonError)
	} else {
		if result.ErrorCode != "" {
			wrongPrint.Println(result.ErrorMsg)
		} else {
			fmt.Print("\n")
			tipPrint.Printf("- from: %v to: %v\n", result.From, result.To)
			greenPrint.Println("-", result.TransResult[0].Dst)
		}
	}
}

// 自定义帮助信息
func usage() {
	fmt.Println("Usage bf-go [options] <query>")
	fmt.Println("-v                  output version number")
	fmt.Println("-h                  output help inf")
	fmt.Println("-from <lang>        source language")
	fmt.Println("-to <lang>          target language")
	fmt.Println("-lang               output the list of supported languages")
}

func showLanguages() {
	var langs = `zh(Chinese|简体中文), cht(Traditional Chinese|繁体中文),en(English|英文), yue(粤语), wyw(文言文), jp(Japanese|日语),
 kor(Korean|韩语), fra(France|法语),spa(Spanish|西班牙语), th(Thai|泰语), ara(Arabic|阿拉伯语), ru(Russian|俄语), 
 pt(Portuguese|葡萄牙语), de(German|德语), it(Italian|意大利语),el(Greek|希腊语), nl(Dutch|荷兰语), pl(Polish|波兰语), 
 bul(Bulgarian|保加利亚语), est(Estonian|爱沙尼亚语), dan(Danish|丹麦语), fin(Finnish|芬兰语), cs(Czech|捷克语), 
 rom(Romanian|罗马尼亚语), slo(Slovenian|斯洛文尼亚语), swe(Swedish|瑞典语), hu(Hungarian|匈牙利语), vie(Vietnamese|越南语)`
	fmt.Println("Languages:\n", langs)
	// fmt.Println(langs)
}
func getArgv() (from, to, query string) {
	var _from string
	var _to string
	var _query string
	var _version bool
	var _lang bool
	var isZh = regexp.MustCompile("[\u4E00-\u9FA5]+")
	flag.Usage = usage
	flag.StringVar(&_from, "from", "", "")
	flag.StringVar(&_to, "to", "", "")
	flag.BoolVar(&_version, "v", false, "")
	flag.BoolVar(&_lang, "lang", false, "")

	flag.Parse()
	// flag.Args() 要放在parse后面
	_query = strings.Join(flag.Args(), " ")

	// 帮助信息有限
	if len(os.Args) == 1 {
		usage()
		os.Exit(0)
	}
	if _version {
		fmt.Println(versonStr)
		os.Exit(0)
	}
	if _lang {
		showLanguages()
		os.Exit(0)
	}

	if _from == "" {
		if isZh.MatchString(query) {
			_from = "zh"
		} else {
			_from = "auto"
		}
	}

	if _to == "" {
		if _from == "zh" {
			_to = "en"
		} else {
			_to = "auto"
		}
	}
	return _from, _to, _query
}

func main() {
	from, to, query := getArgv()
	getTranalation(TransProps{from, to, query})
}
