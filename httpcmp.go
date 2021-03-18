package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/jl2005/image-compare/compare"
)

var SIZE = 2048
var buf = make(chan *http.Request, SIZE)

type Info struct {
	Base   string
	Update string
	Next   string
}

func diffHandle(c *gin.Context) {
	info := &Info{
		Base:   c.Query("base"),
		Update: c.Query("update"),
	}
	i, _ := strconv.Atoi(info.Base)
	if len(info.Update) == 0 {
		info.Update = fmt.Sprintf("%05d", i+1)
	}
	info.Base = "images/" + info.Base
	info.Update = "images/" + info.Update
	info.Next = fmt.Sprintf("%05d", i+2)

	c.HTML(200, "diff.tmpl", info)
}

func handle(c *gin.Context) {
	buf <- c.Request
	c.JSON(200, gin.H{
		"message": "success",
	})
}

func handle1(c *gin.Context) {
	print(c.Request)
	c.JSON(200, gin.H{
		"message": "this is handle1",
	})
}

func handle2(c *gin.Context) {
	print(c.Request)
	c.JSON(200, gin.H{
		"message": "this is handle2",
	})
}

func print(req *http.Request) {
	fmt.Printf("Host: %s\n", req.Host)
	for k, v := range req.Header {
		fmt.Printf("%s: ", k)
		for i := range v {
			if i == 0 {
				fmt.Printf("%s", v[i])
			} else {
				fmt.Printf(",%s", v[i])
			}
		}
		fmt.Println()
	}
}

var fileID int64
var num int
var addr string
var base string
var update string
var path string
var fileListen string

func main() {
	flag.IntVar(&num, "num", 10, "client num")
	flag.StringVar(&addr, "l", ":9090", "listen addr address")
	flag.StringVar(&base, "base", "http://127.0.0.1:8081", "base machine address")
	flag.StringVar(&update, "update", "http://127.0.0.1:8082", "update machine address")
	flag.StringVar(&path, "path", "./diff", "save diff file path")
	flag.StringVar(&fileListen, "fl", ":8080", "show diff image compare")
	flag.Parse()

	go fileServer(fileListen, path)

	for i := 0; i < num; i++ {
		go worker(base, update)
	}

	/*
		go runServer(":8081", handle1)
		go runServer(":8082", handle2)
	*/
	runServer(addr, handle)
}

func runServer(addr string, handle func(*gin.Context)) {
	// Creates a router without any middleware by default
	r := gin.New()
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.Any("/*proxyPath", handle)
	r.Run(addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func fileServer(addr string, path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("can not get abs path for %s\n", path)
		return
	}
	log.Printf("get path: %s, absPaht: %s\n", path, absPath)

	router := gin.Default()
	router.Static("/images", "./diff/")
	router.LoadHTMLGlob("templates/*")
	router.Any("/", diffHandle)

	router.Run(addr)
	log.Printf("diff file listen %s\n", addr)
}

func worker(base string, update string) {
	for src := range buf {
		baseInfo, baseErr := sendRequest(base, src)
		updateInfo, updateErr := sendRequest(update, src)
		if baseErr != nil && updateErr != nil {
			log.Printf("same get error baseErr=%v, updateErr=%v, url=%s", baseErr, updateErr, src.RequestURI)
			continue
		}
		result := make(map[string]comare.Diff)
		baseInfo.Compare(updateInfo, result)
		if len(result) == 0 {
			log.Printf("same %v", src.RequestURI)
		} else {
			var out string
			for k, v := range result {
				out += fmt.Sprintf("%s: %s ", k, v.String())
			}
			id := getFileId()
			basePath = save(fmt.Sprintf("%s/%05d", path, id), baseInfo.Data)
			updatePath = save(fmt.Sprintf("%s/%05d", path, id+1), updateInfo.Data)
			log.Printf("not same %s file: %s,%s %s", out, basePath, updatePath, src.RequestURI)
		}
	}
}

func sendRequest(addr string, src *http.Request) (*compare.Info, error) {
	client := &http.Client{}
	req, err := http.NewRequest(src.Method, addr+src.RequestURI, nil)
	if err != nil {
		return nil, err
	}
	req.Host = src.Host
	req.Header = src.Header
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("status %v body len %d err=%v, url=%v\n", resp.Status, len(data), err, src.RequestURI)
	}
	info := compare.ParseInfo(resp.StatusCode, resp.Header, data)
	return info, nil
}

func getFileID() int {
	return atomic.AddInt64(&fileID, 2)
}

func save(path string, data []byte) string {
	if data == nil || len(data) == 0 {
		return "null"
	}
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		log.Printf("wirte error %s, %v", path, err)
	}
	return path
}
