// dataCenter.go
package colorballs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
	"util"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbHostIP   = "127.9.26.130:3306" //IP地址
	dbUserName = "admin2g2HTwv"      //用户名
	dbPassword = "fP_DVT_QPiUq"      //密码
	dbName     = "golang"            //表名
	url        = "http://baidu.lecai.com/lottery/draw/list/50?type=range_date&start=%s&end=%s"
)

type Balls struct {
	Phase     int
	RedBalls  []string
	BlueBall  string
	PhaseDate string
}

func parseUrlData(url string) []*Balls {
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
			b.Phase = util.Atoi(Phase)
		}
		if i%11 == 7 {
			content := strings.TrimSpace(contentSelection.Text())
			date := strings.SplitN(content, "（", 2)[0]
			b.PhaseDate = date
		}
		if i%11 == 8 {
			elist := goquery.NewDocumentFromNode(contentSelection.Nodes[0])
			elist.Find("em").Each(func(j int, eSelection *goquery.Selection) {
				em := strings.TrimSpace(eSelection.Text())
				if j < 6 {
					b.RedBalls = append(b.RedBalls, em)
				} else {
					b.BlueBall = em
				}

			})
			BallList = append(BallList, b)
		}
	})

	return BallList
}

func openDataBase() *sql.DB {
	info := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUserName, dbPassword, dbHostIP, dbName)
	db, err := sql.Open("mysql", info)
	if err != nil {
		log.Println("[ERROR]:", err)
		return nil
	}

	return db
}

func InitPhaseDate() string {
	b := queryNewData()
	var phaseDate string = "2015-01-27"
	if b != nil {
		phaseDate = b.PhaseDate
	}

	return phaseDate
}
func Update(phaseDate string) {
	nowDate := time.Now().Format("2006-01-02")
	mUrl := fmt.Sprintf(url, phaseDate, nowDate)
	ballList := parseUrlData(mUrl)
	log.Println("[INFO] len:", len(ballList))
	db := openDataBase()
	for _, b := range ballList {
		stmt, _ := db.Prepare("insert into golang_color_balls_table(phase,red_balls,blue_ball,phase_date) values(?,?,?,?)")
		ball := util.Encoder(b.RedBalls)
		stmt.Exec(b.Phase, string(ball), b.BlueBall, b.PhaseDate)
		stmt.Close()
	}
	db.Close()
	phaseDate = nowDate
}

func queryNewData() *Balls {
	b := new(Balls)
	db := openDataBase()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM golang_color_balls_table order by phase DESC")
	defer rows.Close()
	for rows.Next() {
		var red string
		rows.Columns()
		err = rows.Scan(&b.Phase, &red, &b.BlueBall, &b.PhaseDate)
		if err != nil {
			return nil
		}
		err = json.Unmarshal([]byte(red), b.RedBalls)
		if err != nil {
			return nil
		}
		break
	}
	return b
}

func GetAllBalls() []*Balls {
	db := openDataBase()
	defer db.Close()

	mballs := make([]*Balls, 0, 200)
	rows, err := db.Query("SELECT * FROM golang_color_balls_table order by phase DESC")
	if err != nil {
		log.Println("[ERROR] Query:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var red string
		b := new(Balls)
		rows.Columns()
		err = rows.Scan(&b.Phase, &red, &b.BlueBall, &b.PhaseDate)
		if err != nil {
			log.Println("[ERROR] Scan:", err)
			break
		}
		b.RedBalls = make([]string, 0, 6)
		err = json.Unmarshal([]byte(red), &b.RedBalls)
		if err != nil {
			log.Println("[ERROR] Unmarshal:", err)
			break
		}
		mballs = append(mballs, b)
	}
	return mballs
}

type ballTimes struct {
	name  string
	count int
}

type bTimesSlice []*ballTimes

func (b bTimesSlice) Len() int {
	return len(b)
}
func (b bTimesSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b bTimesSlice) Less(i, j int) bool {
	return b[i].count < b[j].count
}

func GetNewBalls() *Balls {
	mballs := GetAllBalls()

	redList := make([]*ballTimes, 0, 33)
	blueList := make([]*ballTimes, 0, 16)
	redMap := make(map[string]*ballTimes)
	blueMap := make(map[string]*ballTimes)
	for _, b := range mballs {
		for _, v := range b.RedBalls {
			t, ok := redMap[v]
			if !ok {
				t = new(ballTimes)
				t.name = v
				t.count = 1
				redMap[v] = t
				continue
			} else {
				t.count++
			}
		}
		t, ok := blueMap[b.BlueBall]
		if !ok {
			t = new(ballTimes)
			t.name = b.BlueBall
			t.count = 1
			blueMap[b.BlueBall] = t

		} else {
			t.count++
		}
	}

	for _, v := range redMap {
		redList = append(redList, v)
	}
	for _, v := range blueMap {
		blueList = append(blueList, v)
	}
	sort.Sort(bTimesSlice(redList))
	sort.Sort(bTimesSlice(blueList))

	nb := new(Balls)
	nb.BlueBall = blueList[0].name
	for i, v := range redList {
		nb.RedBalls = append(nb.RedBalls, v.name)
		if i == 5 {
			break
		}
	}
	nb.Phase = mballs[0].Phase + 1

	return nb
}
