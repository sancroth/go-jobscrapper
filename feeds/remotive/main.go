package remotive

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
	config.host = "https://remotive.io"
	return &PublicFeed{
		config:   config,
		BaseFeed: feeds.NewBaseFeed(name),
	}
}

func (f *PublicFeed) Connect() {
	counter := 0
	doc := f.GetDocument(fmt.Sprintf("%s/remote-jobs/software-dev?live_jobs[toggle]&live_jobs[sortBy]=live_jobs_sort_by_date&live_jobs[menu][category]=Software Development", f.config.host))
	doc.Find("div#hits ul").Children().Each(func(i int, s *goquery.Selection) {
		if counter < f.Limit() {
			href, exists := s.Find("a.job-tile-title").Attr("href")
			if exists {

				job := f.GetDocument(href)
				title := job.Find(".content .h1").Text()
				company := job.Find(".content .company").Text()
				apply, exists := job.Find(".apply-wrapper a").Attr("apply-url")

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