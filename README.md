# rss2lemmy

## Installation

```
go install github.com/nuzayets/rss2lemmy@latest
```

## Set up application

Create an agent file and save it somewhere the application can read. See example.agent.

## Usage

```
Usage of ./rss2lemmy:
  -agent string
        full path to agent file (default "/home/coder/rss2lemmy.agent")
  -community string
        community to post in (default "testingground4bots")
  -feed string
        the feed URL (default "https://blog.golang.org/feed.atom?format=xml")
  -posted string
        writable file to store already-posted links (default "/home/coder/rss2lemmy.posted")
  -scope int
        posts published more than scope seconds ago will not be posted (default 3600)
  -suffix string
        string to append to post title
```

You can execute it with a cron. Each time the cron runs, it will post anything new that is within the scope.
