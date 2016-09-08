# redis-recommend
## A Simple Redis recommendation engine written in [Go](https://golang.org/).

###

## About
This is a simple recommendation engine written in [Go](https://golang.org/) using [Redis](http://redis.io). The Redis client Go library used is [Redigo](https://github.com/garyburd/redigo).

## Usage

Rate an item:

```
redis-recommend rate <user> <item> <score>  
```

Find (n) similar users for all users:

```
redis-recommend batch-update [--results=<n>]
```

Get (n) suggested items for a user:
```
redis-recommend suggest <user> [--results=<n>]
```

Get the probable score a user would give to an item:
```
redis-recommend get-probability <user> <item>
```
