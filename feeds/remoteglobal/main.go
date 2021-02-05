package remoteglobal

import (
	"../../feeds"
	"../../models"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
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

type Data struct{
	Html string `json:"html"`
}

func NewPublicFeed(name string) *PublicFeed{
	log.Println(fmt.Sprintf("feed %s connected",name))
	config:= &PublicFeedConfig{}
	config.host = "remoteglobal.com"
	config.url = fmt.Sprintf("https://%s/jm-ajax/get_listings",config.host)
	return &PublicFeed{
		config: config,
		BaseFeed : feeds.NewBaseFeed(name),
	}
}

func (f *PublicFeed) Connect(){
	counter := 0
	log.Println(fmt.Sprintf("connect : %s fetching",f.config.url))
	resp, err := http.Get(f.config.url)
	var data Data
	if err:= json.NewDecoder(resp.Body).Decode(&data); err!=nil{
		log.Fatal(err)
	}
	log.Println("decoding data.HTML")
	doc,err := goquery.NewDocumentFromReader(strings.NewReader(data.Html))
	if err!=nil{
		log.Fatal(err)
	}

	doc.Find("li").Each(func(i int , s *goquery.Selection){
		log.Println("for each li")
		if counter < f.Limit(){
			log.Println("for each li next")
			v,_ :=s.Html()
			log.Println(fmt.Sprintf("reading : %v",v))
			href, exists := s.Find("a").Attr("href")
			log.Println(href)
			log.Println(exists)
			if exists{
				log.Println(fmt.Sprintf("parsing : %s",href))
				u,err := url.Parse(href)
				if err!=nil{
					log.Fatal(err)
				}
				log.Println(fmt.Sprintf("reading : %s",href))
				job := f.GetDocument(href)
				title := job.Find(".container .entry-title").Text()
				company := job.Find(".website").Text()
				apply, exists := job.Find(".application_details a").Attr("href")
				log.Println("checking if href exists")
				if exists {
					log.Println(fmt.Sprintf("creating post for : %s",href))
					post := &models.Post{
						Path:     u.Path,
						Name:     f.Name(),
						Host:     f.config.host,
						Title:    strings.TrimSpace(title),
						Apply:    strings.TrimSpace(apply),
						Company:  strings.TrimSpace(company),
					}
					log.Println(fmt.Sprintf("save : %v",post))
					saved, err := f.SavePost(post)
					if err != nil {
						log.Fatal(err)
					}
					if saved {
						log.Println(fmt.Sprintf("Post:%v saved successfully ", post))
						counter++
						log.Println(counter)
					}
				}

			}
		}
	})
}