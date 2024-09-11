package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
	"go.elara.ws/go-lemmy"
)

type AgentFile struct {
	InstanceBaseUrl string `json:"instance_base_url"`
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

var agentFile, postedFile, community, titleSuffix, feedURL string
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

func getAgentFile() AgentFile {
	agentBytes, err := os.ReadFile(agentFile)
	if err != nil {
		log.Fatal(err)
	}
	var agent AgentFile
	err = json.Unmarshal(agentBytes, &agent)
	if err != nil {
		log.Fatal(err)
	}
	return agent
}

func getClient(ctx context.Context) *lemmy.Client {
	agent := getAgentFile()
	client, err := lemmy.New(agent.InstanceBaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = client.ClientLogin(ctx, lemmy.Login{
		UsernameOrEmail: agent.UsernameOrEmail,
		Password:        agent.Password,
	})
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func getCommunityId(ctx context.Context, client *lemmy.Client, community string) int64 {
	res, err := client.Community(ctx, lemmy.GetCommunity{
		Name: lemmy.NewOptional(community),
	})
	if err != nil {
		log.Fatal(err)
	}

	return res.CommunityView.Community.ID
}

func isPosted(community string, url string) bool {
	f, err := os.OpenFile(postedFile, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	match := fmt.Sprintf("%s,%s", community, url)

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

func logPosted(community string, url string) {
	f, err := os.OpenFile(postedFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	line := fmt.Sprintf("%s,%s\n", community, url)
	if _, err = f.WriteString(line); err != nil {
		log.Fatal(err)
	}

}

func main() {
	flag.StringVar(&agentFile, "agent", homeDir()+"/rss2lemmy.agent", "full path to agent file")
	flag.StringVar(&postedFile, "posted", homeDir()+"/rss2lemmy.posted", "writable file to store already-posted links")
	flag.StringVar(&titleSuffix, "suffix", "", "string to append to post title")
	flag.StringVar(&community, "community", "testingground4bots", "community to post in")
	flag.StringVar(&feedURL, "feed", "https://blog.golang.org/feed.atom?format=xml", "the feed URL")
	flag.IntVar(&scopeSecs, "scope", 3600, "posts published more than scope seconds ago will not be posted")
	flag.Parse()

	var client *lemmy.Client
	var communityId int64
	ctx := context.Background()
	minPublishTime := time.Now().Add(time.Duration(-scopeSecs) * time.Second)
	feed := getFeed(feedURL)

	for _, item := range feed.Items {
		if item.PublishedParsed.After(minPublishTime) {
			if client == nil {
				client = getClient(ctx) // only get client when we encounter an item so we don't login every cron
				communityId = getCommunityId(ctx, client, community)
			}

			title := fmt.Sprintf("%s%s", item.Title, titleSuffix)
			fmt.Println("Title: " + title)
			fmt.Println("Link: " + item.Link)
			if !isPosted(community, item.Link) {
				fmt.Printf("Posting to %s...", community)
				res, err := client.CreatePost(ctx, lemmy.CreatePost{
					Name:        title,
					URL:         lemmy.NewOptional(item.Link),
					CommunityID: communityId,
				})

				if err != nil {
					log.Println(err)
				} else if res.Error.IsValid() {
					log.Print(res.Error.String())
				} else {
					logPosted(community, item.Link)
				}
			} else {
				fmt.Println("Already posted!")
			}

		}
	}

}
