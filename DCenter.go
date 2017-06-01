package main

import (
	"colorballs"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"util"
)

var (
	program = filepath.Base(os.Args[0])
)

const (
	OPCODE_NEW_BALLS = iota
	OPCODE_ALL_BALLS
	OPCODE_HOT_NEWS
)

func init() {
	cwd, _ := os.Getwd()
	file := fmt.Sprintf("%s/%s.golang.log", cwd, program)
	logFile, logErr := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Printf("Fail to open log file<%s>", file)
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	log.Println("starting...")
	go dataCollector()
	time.Sleep(10)

	http.HandleFunc("/", handleRequest)
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	log.Printf("listening on %s...\n", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Panicln("err:", err)
	}
}

type collector struct {
	phaseDate string //colorballs
}

func initDC() *collector {
	c := new(collector)
	c.phaseDate = colorballs.InitPhaseDate()

	return c

}
func dataCollector() {
	c := initDC()
	for {
		colorballs.Update(c.phaseDate)

		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
	}
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	m := util.Msg{}
	var data interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("[ERROR]: failed to read request<", err, ">")
		return
	}
	util.Decoder(body, m)
	switch m.OpCode {
	case OPCODE_NEW_BALLS:
		data = colorballs.GetNewBalls()
	case OPCODE_ALL_BALLS:
		data = colorballs.GetAllBalls()
	case OPCODE_HOT_NEWS:
	default:
		data = "welcome to tomato world"
	}
	msg := util.Encoder(data)
	fmt.Fprintf(res, "%s", msg)

}
