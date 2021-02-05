package indeed

import (
	"../../feeds"
	"../../models"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

type PublicFeedConfig struct{
	url string
	host string
}

type PublicFeed struct {
	*feeds.BaseFeed
	config *PublicFeedConfig
}

func NewPublicFeed(name string) *PublicFeed{
	log.Println(fmt.Sprintf("feed %s connected",name))
	config:= &PublicFeedConfig{}
	config.host = "https://www.indeed.co.uk"
	return &PublicFeed{
		config: config,
		BaseFeed : feeds.NewBaseFeed(name),
	}
}

func (f *PublicFeed) Connect(){
	counter := 0
	url := fmt.Sprintf("%s/jobs?q=golang&sort=date&fromage=1&start=100", f.config.host)
	log.Println(fmt.Sprintf("making a request to url: %s",url))
	doc := f.GetDocument(url)
	doc.Find("#resultsCol div.jobsearch-SerpJobCard").Each(func(i int, s *goquery.Selection){
		if counter < f.Limit(){
			id,exists := s.Attr("data-jk")
			if exists{
				path := fmt.Sprintf("/viewjob?jk=%s",id)
				href := fmt.Sprintf("%s%s",f.config.host,path)

				log.Println(fmt.Sprintf("making a request to page: %s",href))

				jobPost:= f.GetDocument(href)
				title := jobPost.Find(".jobsearch-JobInfoHeader-title-job").Text()
				salary := jobPost.Find(".jobsearch-JobMetadataHeader-item").Text()
				position := jobPost.Find(".jobsearch-DesktopStickyContainer-subtitle").Children().Last().Text()
				company := jobPost.Find(".jobsearch-DesktopStickyContainer-subtitle").Children().First().Children().First().Text()

				apply, exists := jobPost.Find("#applyButtonLinkContainer a").Attr("href")
				if exists {
					post := &models.Post{
						Path:     path,
						Name:     f.Name(),
						Host:     f.config.host,
						Title:    strings.TrimSpace(title),
						Apply:    strings.TrimSpace(apply),
						Company:  strings.TrimSpace(company),
						Salary:   strings.TrimSpace(salary),
						Position: strings.TrimSpace(position),
					}
					log.Println(fmt.Sprintf("trying to post: %v",post))
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
	fmt.Println("finished indeed feed")
}