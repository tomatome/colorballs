package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	nilTest()
}

var (
	dbHostIP   = "127.9.26.130:3306" //IP地址
	dbUserName = "admin2g2HTwv"      //用户名
	dbPassword = "fP_DVT_QPiUq"      //密码
	dbName     = "golang"            //表名
	program    = filepath.Base(os.Args[0])
	logger     *log.Logger
)

type Balls struct {
	phase    string
	redBalls []string
	blueBall string
	date     string
}

func init() {
	dir := os.Getenv("$OPENSHIFT_GO_LOG_DIR")
	fileName := fmt.Sprintf("%s/%s.golang.log", dir, program)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("fail to create test.log file!")
	}
	logger = log.New(file, "", log.LstdFlags|log.Llongfile)
}

func openDataBase() *sql.DB {
	info := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUserName, dbPassword, dbHostIP, dbName)
	db, err := sql.Open("mysql", info)
	if err != nil {
		logger.Println("[ERROR]:", err)
		return nil
	}

	return db
}

func collectData(url string) []*Balls {
	var b *Balls
	BallList := make([]*Balls, 0, 400)

	doc, _ := goquery.NewDocument(url)
	node := doc.Find(".historylist")
	historyDoc := goquery.NewDocumentFromNode(node.Nodes[0])
	historyDoc.Find("td").Each(func(i int, contentSelection *goquery.Selection) {
		if i < 6 {
			return
		}

		if i%11 == 6 {
			Phase := strings.TrimSpace(contentSelection.Text())
			b = new(Balls)
			b.phase = Phase
		}
		if i%11 == 7 {
			content := strings.TrimSpace(contentSelection.Text())
			date := strings.SplitN(content, "（", 2)[0]
			b.date = date
		}
		if i%11 == 8 {
			elist := goquery.NewDocumentFromNode(contentSelection.Nodes[0])
			elist.Find("em").Each(func(j int, eSelection *goquery.Selection) {
				em := strings.TrimSpace(eSelection.Text())
				if j < 7 {
					b.redBalls = append(b.redBalls, em)
				} else {
					b.blueBall = em
				}

			})
			BallList = append(BallList, b)
		}
	})

	return BallList
}

func nilTest() {
	url := "http://baidu.lecai.com/lottery/draw/list/50?type=range_date&start=2015-01-01&end=2017-05-24"
	ballList := collectData(url)
	logger.Println("[INFO] len:", len(ballList))
	db := openDataBase()
	for _, b := range ballList {
		stmt, _ := db.Prepare("insert into golang_color_balls_table(phase,red_balls,blue_ball,phase_date) values(?,?,?,?)")
		stmt.Exec(b.phase, b.redBalls, b.blueBall, b.date)
		stmt.Close()
	}
}
