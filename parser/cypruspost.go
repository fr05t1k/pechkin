package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fr05t1k/pechkin/storage"
	"net/http"
	"strings"
	"time"
)

const defaultPageUrl = "https://ips.cypruspost.gov.cy/ipswebtrack/IPSWeb_item_events.aspx?itemid=%s&Submit=Submit"

type cyprusPost struct {
	PageUrl string
	Client  *http.Client
}

func (c *cyprusPost) Parse(track string) (events []storage.Event, err error) {
	response, err := c.Client.Get(fmt.Sprintf(c.PageUrl, track))
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return
	}

	doc.Find("table.table-striped.table-bordered tr").Each(func(i int, tr *goquery.Selection) {
		event := storage.Event{}
		descBuilder := strings.Builder{}
		tds := tr.Find("td")
		if tds.Size() < 6 || (!tr.HasClass("tabl1") && !tr.HasClass("tabl1")) {
			return
		}
		tds.Each(func(i int, selection *goquery.Selection) {
			switch i {
			case 0:
				// "3/14/2019 3:51:00 PM"
				event.EventAt, err = time.Parse("1/2/2006 3:04:05 PM", selection.Text())
			default:
				descBuilder.WriteString(selection.Text())
				descBuilder.WriteString("\n")
			}
		})
		event.Description = descBuilder.String()
		events = append(events, event)
	})

	return
}

func NewCyprusPost() *cyprusPost {
	c := cyprusPost{
		PageUrl: defaultPageUrl,
		Client:  http.DefaultClient,
	}

	return &c
}
