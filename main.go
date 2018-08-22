package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"smpNVR/cmd"
	"smpNVR/ffmpeg"
	"smpNVR/lib"
	"smpNVR/server"
	"strconv"
	//	"smpNVR/utils"
	"strings"
	"time"
)

var queueName string
var mode string
var uuid string
var ip string
var port int
var breath int
var vdpath string
var jpgpath string
var senable string

var Processor []ffmpeg.C
var Client cmd.Cli
var staStream []string

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	flag.Parse()
	//read config file
	myConfig := new(Config)
	myConfig.InitConfig("./conf/config.conf")
	mode = myConfig.Read("mode", "mode")
	uuid = myConfig.Read("mode", "uuid")
	ip = myConfig.Read("mode", "ip")
	port, _ = strconv.Atoi(myConfig.Read("mode", "port"))
	breath, _ = strconv.Atoi(myConfig.Read("mode", "breath"))
	if mode != "" {
		queueName = myConfig.Read(mode, "queueName")
		senable = myConfig.Read(mode, "server")
	} else {
		queueName = myConfig.Read("default", "queueName")
		senable = myConfig.Read("default", "server")
	}
	//是否开启流媒体
	if strings.EqualFold(senable, "on") {
		go server.Start()
	}
	//等待10秒，开启默认转流
	time.Sleep(10 * time.Second)

	count := myConfig.Read("stream", "count") //转流路数
	num, err := strconv.Atoi(count)
	if err != nil { //
		fmt.Println("will not stream")
	} else {
		for ii := 1; ii <= num; ii++ {
			sname := myConfig.Read(fmt.Sprintf("route%d", ii), "cam")
			srtmp := myConfig.Read(fmt.Sprintf("route%d", ii), "rtmp")
			srtsp := myConfig.Read(fmt.Sprintf("route%d", ii), "rtsp")
			m, err := lib.StreamLiveVideo(srtsp, srtmp, "", sname)
			if err == nil {
				Processor = append(Processor, m)
				staStream = append(staStream, sname)
			} else {
				fmt.Println(err)
			}
		}
	}

	vdpath = GetAppPath() + "\\video\\"
	jpgpath = GetAppPath() + "\\pic\\"
}

func main() {
	//TCP
	Client = cmd.Cli{Uid: uuid, Ip: ip, Port: port, Breath: 5, Mhandler: handle}
	for {
		if Client.Dial() {
			//开启循环读取
			go Client.Read(Client.Mhandler)

			//发送auth
			Client.Cinit()
			Client.Auth()
			for {
				if ok := Client.Ready(); ok != nil {
					fmt.Println(ok)
					break
				}
				time.Sleep(5 * time.Second)
			}
		}
		time.Sleep(5 * time.Second)
		fmt.Println("ygg 重连")
		continue
	}
}

//配置
type Config struct {
	Mymap  map[string]string
	strcet string
}

func (c *Config) InitConfig(path string) {
	c.Mymap = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(c.strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := c.strcet + "." + frist
		c.Mymap[key] = strings.TrimSpace(second)
	}
}

func (c Config) Read(node, key string) string {
	key = node + "." + key
	v, found := c.Mymap[key]
	if !found {
		return ""
	}
	return v
}

func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

func handle(b []byte) {
	for {
		m := cmd.BytesToInt32(b[:4])
		//解析四位，如果长度大于原长度，认为解析错误，丢弃
		if m > uint32(len(b)) {
			break
		}
		n := b[4 : m+4]
		//调用解析
		fmt.Println(string(n))
		go func(bb []byte) {
			pp := lib.Pas{}
			Parser(pp, bb)
		}(n)
		b = b[m+3 : len(b)]
		//嵌套，用于分离粘连
		if len(b) < 4 {
			break
		}
		handle(b)
	}
}

func Parser(pmp lib.Pas, b []byte) {
	if strings.Contains(string(b), "code") {
		var rr lib.Result = lib.Result{}
		json.Unmarshal(b, &rr)
	} else {
		json.Unmarshal(b, &pmp.Request)
		if strings.EqualFold(pmp.Action, "play") {
			fl := true
			for _, tst := range Processor {
				if strings.EqualFold(tst.Name, pmp.CameraCode) {
					fl = false //已经开启则不再操作
				}
			}
			if fl {
				streamName := fmt.Sprintf("%s/%s%s.flv", pmp.CameraCode, pmp.ImosCode, pmp.Resolution)
				m, err := lib.StreamLiveVideo(pmp.RtspUrl, queueName, streamName, pmp.CameraCode)
				if err != nil {
					//返回失败
					Client.RespFailed(pmp.Rspuuid, err)
				}
				//wait the stream for cache
				time.Sleep(3 * time.Second)
				Client.RespSuccess(pmp.Rspuuid)
				Processor = append(Processor, m)
				return
			}
			Client.RespSuccess(pmp.Rspuuid)
			return
		}
		if strings.EqualFold(pmp.Action, "stopTransform") {
			fl := true
			for _, tst := range staStream {
				if strings.EqualFold(tst, pmp.CameraCode) {
					fl = false //不再操作
				}
			}
			if fl {
				for index, m := range Processor {
					//摄像机编码作为唯一标识
					if strings.EqualFold(m.Name, pmp.CameraCode) {
						lib.ClosefmgVideo(m)
						Client.RespSuccess(pmp.Rspuuid)
						//delete one element from the slice
						Processor = append(Processor[:index], Processor[index+1:]...)
						return
					}
				}
			}
			Client.RespSuccess(pmp.Rspuuid)
			return
		}
		if strings.EqualFold(pmp.Action, "capture") {
			lib.GetfmgFrame(pmp.RtspUrl, jpgpath, fmt.Sprintf("%s%d", pmp.CameraCode, pmp.SendTime))
			Client.RespSuccess(pmp.Rspuuid)
		}
	}

}
