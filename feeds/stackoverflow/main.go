package stackoverflow

import (
	"../../feeds"
	"../../models"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"strings"
)

type PublicFeedConfig struct {
	url  string
	host string
}

type PublicFeed struct {
	*feeds.BaseFeed
	config *PublicFeedConfig
}

func NewPublicFeed(name string) *PublicFeed {
	config := &PublicFeedConfig{}
	config.host = "https://stackoverflow.com"
	return &PublicFeed{
		config:   config,
		BaseFeed: feeds.NewBaseFeed(name),
	}
}

func (f *PublicFeed) Connect() {
	counter := 0
	doc := f.GetDocument(fmt.Sprintf("%s/jobs?r=true&tl=javascript+react+react.js+node.js+angular.js+php+golang&sort=p", f.config.host))
	doc.Find(".listResults").Children().Each(func(i int, s *goquery.Selection) {
		if counter < f.Limit() {
			id, exists := s.Attr("data-jobid")
			if exists {
				href := fmt.Sprintf("%s/jobs/%s", f.config.host, id)
				job := f.GetDocument(href)

				title := job.Find(".fs-headline1 a").Text()
				salary, _ := job.Find(".-salary").Attr("title")
				company := job.Find(".employer").Text()
				position := job.Find(".-remote").Text()
				apply, exists := job.Find("._apply").Attr("href")

				if exists {
					u, err := url.Parse(href)
					if err != nil {
						log.Fatal(err)
					}

					post := &models.Post{
						Path:     u.Path,
						Name:     f.Name(),
						Host:     f.config.host,
						Title:    strings.TrimSpace(title),
						Salary:   strings.TrimSpace(salary),
						Position: strings.TrimSpace(position),
						Apply:    strings.TrimSpace(apply),
						Company:  strings.TrimSpace(company),
					}
					saved, err := f.SavePost(post)
					if err != nil {
						log.Fatal(err)
					}
					if saved {
						log.Println(fmt.Sprintf("Post:%v saved successfully ", post))
						counter++
					}
				}
			}
		}
	})
}