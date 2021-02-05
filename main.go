package main

import (
	"./feeds"
	"./feeds/indeed"
	"./feeds/remoteglobal"
	"./feeds/remotive"
	"./feeds/stackoverflow"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
)

func init(){
	if err:= godotenv.Load();err !=nil{
		log.Fatal(err)
	}
}

func createFeed(feedName string) feeds.PublicFeed{
	switch feedName{
	case "indeed":
		fmt.Println("calling indeed")
		return indeed.NewPublicFeed(feedName)
	case "remoteglobal":
		fmt.Println("calling remote global")
		return remoteglobal.NewPublicFeed(feedName)
	case "remotive":
		fmt.Println("calling remotive global")
		return remotive.NewPublicFeed(feedName)
	case "stackoverlow":
		fmt.Println("calling stackoverlow global")
		return stackoverflow.NewPublicFeed(feedName)
	default:
		fmt.Println("calling default")
		return nil
	}
}

func parseData() {
	fmt.Println("calling feeds")
	go createFeed("indeed").Connect()
	go createFeed("remoteglobal").Connect()
	go createFeed("remotive").Connect()
	go createFeed("stackoverflow").Connect()
}

/* will implement later
func broadcastData() {

}
*/

func main() {
	parseData()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit
}
