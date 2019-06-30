package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/turnage/graw/reddit"
)

var agentFile, postedFile, titleSuffix, subreddit, feedURL string
var scopeSecs int

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return home
}

func getFeed(url string) *gofeed.Feed {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Fatal(err)
	}
	return feed
}

func getBot() reddit.Bot {
	bot, err := reddit.NewBotFromAgentFile(agentFile, 0)
	if err != nil {
		log.Fatal(err)
	}
	return bot
}

func isPosted(subreddit string, url string) bool {
	f, err := os.OpenFile(postedFile, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	match := fmt.Sprintf("%s,%s", subreddit, url)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if line == match {
			return true
		}
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	return false
}

func logPosted(subreddit string, url string) {
	f, err := os.OpenFile(postedFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	line := fmt.Sprintf("%s,%s\n", subreddit, url)
	if _, err = f.WriteString(line); err != nil {
		log.Fatal(err)
	}

}

func main() {
	flag.StringVar(&agentFile, "agent", homeDir()+"/rss2reddit.agent", "full path to agent file")
	flag.StringVar(&postedFile, "posted", homeDir()+"/rss2reddit.posted", "writable file to store already-posted links")
	flag.StringVar(&titleSuffix, "suffix", "", "string to append to post title")
	flag.StringVar(&subreddit, "subreddit", "testingground4bots", "subreddit to post in")
	flag.StringVar(&feedURL, "feed", "https://blog.golang.org/feed.atom?format=xml", "the feed URL")
	flag.IntVar(&scopeSecs, "scope", 3600, "posts published more than scope seconds ago will not be posted")
	flag.Parse()

	var bot reddit.Bot
	minPublishTime := time.Now().Add(time.Duration(-scopeSecs) * time.Second)
	feed := getFeed(feedURL)

	for _, item := range feed.Items {
		// only items within past {scope}
		if item.PublishedParsed.After(minPublishTime) {
			// post
			if bot == nil {
				bot = getBot() // only get bot when we encounter an item so we don't login every cron
			}

			title := fmt.Sprintf("%s%s", item.Title, titleSuffix)
			fmt.Println("Title: " + title)
			fmt.Println("Link: " + item.Link)
			if !isPosted(subreddit, item.Link) {
				fmt.Println(fmt.Sprintf("Posting to %s...", subreddit))
				err := bot.PostLink(subreddit, title, item.Link)
				if err != nil {
					log.Println(err)
				} else {
					logPosted(subreddit, item.Link)
				}
			} else {
				fmt.Println("Already posted!")
			}

		}
	}

}
