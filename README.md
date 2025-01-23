# Blog Aggregator

This is a simple blog aggregator that fetches blog posts from different sources and displays them in your CLI.

## Pre-requisites

- Go 1.20 or higher
- PostgreSQL 13.2 or higher (you can use Docker to run it)

## Installation

```bash
go install github.com/ricardolopes86/blog-aggregator
```

In your home directory, create a `~/.gatorconfig.json` file with the following content:

```json
{"db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable","current_user_name":"<your_name_here>"}
```

Gator doesn't use any authentication mechanism, so you can put any name you want in the `current_user_name` field.

After its installation, you can run the following command to see the available options:

```bash
gator 
```

## Usage

Gator offers a few commands to interact with the blog aggregator:
- `login`: logs you in the aggregator
- `register`: registers a new user
- `reset`: resets the aggregator (it will delete all the data)
- `users`: lists all the users
- `agg`: aggregates the blog posts. This command accepts a time parameter in the format `1h`, `1d`, `1w`, `1m`, `1y` to specify the time between every aggregation. If no time is provided, it will return an error.
- `add`: adds a new feed from where it will collect posts to add to the database. This command accepts a URL parameter to specify the feed URL.
- `feeds`: lists all the feeds for the logged user
- `follow`: follows a feed. This command accepts a feed URL parameter to specify the feed to follow.
- `unfollow`: unfollows a feed. This command accepts a feed URL parameter to specify the feed to unfollow.
- `following`: lists all the feeds the logged user is following
- `browse`: browses the blog posts. This command accepts a number parameter in the format to specify the `limit` of posts to fetch from the database. If no `limit` is provided, it will be limited to 2.