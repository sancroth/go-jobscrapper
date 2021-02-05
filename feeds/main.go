package feeds

import (
	"../models"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
)

// PublicFeed interface
type PublicFeed interface {
	Connect()                                 // connect to a feed method
	Name() string                             // name of the connected feed
	Limit() int                               // batch size per cycle
	SavePost(post *models.Post) (bool, error) // save to db
	GetDocument(url string) *goquery.Document // parse a new post
}

type BaseFeed struct {
	limit int
	name  string
	posts *models.PostsHandler
}

func NewBaseFeed(name string) *BaseFeed {
	fmt.Println(fmt.Sprintf("creating %s feed",name))
	return &BaseFeed{
		limit: 10,
		name:  name,
		posts: models.NewPostsHandler(),
	}
}

func (f *BaseFeed) Name() string {
	return f.name
}

func (f *BaseFeed) Limit() int {
	return f.limit
}

func (f *BaseFeed) SavePost(post *models.Post) (bool, error) {
	c, err := f.posts.GetPostsCount(post.Name, post.Path)
	if err != nil {
		return false, err
	}
	if c == 1 {
		return false, nil
	}
	return true, f.posts.SavePost(post)
}

func (f *BaseFeed) GetDocument(url string) *goquery.Document {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var doc *goquery.Document
	tasks:= chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error{
			node,err := dom.GetDocument().Do(ctx)
			if err!=nil{
				return err
			}
			res,err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			if err!=nil{
				return err
			}
			doc,err = goquery.NewDocumentFromReader(strings.NewReader(res))
			if err!=nil{
				return err
			}
			return nil
		}),
	}
	if err:= chromedp.Run(ctx,tasks); err!=nil{
		log.Fatal(err)
	}

	return doc
}
