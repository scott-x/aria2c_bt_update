package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

func main() {
	tracker_url := flag.String("url", "https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_best.txt", "text url for tracker")
	file := flag.String("c", path.Join(os.Getenv("HOME"), ".aria2/aria2.conf"), "aria2 config file")
	flag.Parse()
	res, err := http.Get(*tracker_url)
	if err != nil {
		log.Println("网络不可达，请开启代理再试")
		return
	}
	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	content := string(bs)
	if len(content) == 0 {
		log.Println("[更新失败]: 网络出错，tracker url 不可达")
		return
	}
	bs2, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Printf("%s doesn't exist\n", *file)
		return
	}
	//replace
	btRe := regexp.MustCompile(`bt-tracker=(.*)`)
	result := btRe.FindAllStringSubmatch(string(bs2), -1)
	if len(result) == 0 {
		//未配置bt-tracker
		ioutil.WriteFile(*file, []byte(string(bs2)+"\nbt-tracker="+formatRes(content)), 0755)
		log.Println("success!")
		return
	}
	new_conetent := strings.ReplaceAll(string(bs2), result[0][0], "bt-tracker="+formatRes(content))
	ioutil.WriteFile(*file, []byte(new_conetent), 0755)
	log.Println("success!")
}

func formatRes(str string) string {
	re := regexp.MustCompile(`.+`)
	var result string
	res := re.FindAllString(str, -1)
	if len(res) > 0 {
		var new_arr []string
		for _, y := range res {
			new_arr = append(new_arr, strings.TrimSpace(y))
		}
		result = strings.Join(new_arr, ",")
	}
	return result
}
