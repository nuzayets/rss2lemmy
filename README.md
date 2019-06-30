# rss2reddit

## Installation
```
go get github.com/nuzayets/rss2reddit
```

## Set up application
Create an [an agent file for graw](https://turnage.gitbooks.io/graw/content/chapter1.html). Save it somewhere rss2reddit can read.

## Usage
```
Usage of ./rss2reddit:
  -agent string
        full path to agent file (default "/home/username/rss2reddit.agent")
  -feed string
        the feed URL (default "https://blog.golang.org/feed.atom?format=xml")
  -posted string
        writable file to store already-posted links (default "/home/username/rss2reddit.posted")
  -scope int
        posts published more than scope seconds ago will not be posted (default 3600)
  -subreddit string
        subreddit to post in (default "testingground4bots")
  -suffix string
        string to append to post title
```


