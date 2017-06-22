// novel.go
package novel

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

var (
	url = "http://www.tianzeba.com/shengxu/"
)

func ParseUrlData() {
	doc, _ := goquery.NewDocument(url)
	node := doc.Find(".chapterlist")
	fmt.Println(node)
	historyDoc := goquery.NewDocumentFromNode(node.Nodes[0])
	historyDoc.Find("dd").Each(func(i int, contentSelection *goquery.Selection) {
		title := strings.TrimSpace(contentSelection.Text())
		fmt.Println(title)
		link, _ := contentSelection.Find("a").Attr("href")
		fmt.Println(link)
	})
}
