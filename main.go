package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/RedisLabs/redis-recommend/redrec"
	"github.com/docopt/docopt-go"
)

var rr *redrec.Redrec
var err error

func main() {
	usage := `

Usage:
  redis-recommend rate <user> <item> <score>
  redis-recommend suggest <user> [--results=<n>]
  redis-recommend get-probability <user> <item>
  redis-recommend batch-update [--results=<n>]
  redis-recommend -h | --help
  redis-recommend --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  --results=<n>  Num of suggestions to get [default: 100]`

	arguments, _ := docopt.Parse(usage, nil, true, "redis-recommend", false)

	rr, err = redrec.New("redis://localhost:6379")
	chekErrorAndExit(err)

	if arguments["rate"].(bool) {
		user := arguments["<user>"].(string)
		item := arguments["<item>"].(string)
		score, err := strconv.ParseFloat(arguments["<score>"].(string), 64)
		chekErrorAndExit(err)
		rate(user, item, score)
	}

	if arguments["get-probability"].(bool) {
		user := arguments["<user>"].(string)
		item := arguments["<item>"].(string)
		getProbability(user, item)
	}

	if arguments["suggest"].(bool) {
		user := arguments["<user>"].(string)
		results, err := strconv.Atoi(arguments["--results"].(string))
		chekErrorAndExit(err)
		suggest(user, results)
	}

	if arguments["batch-update"].(bool) {
		results, err := strconv.Atoi(arguments["--results"].(string))
		chekErrorAndExit(err)
		update(results)
	}
}

func chekErrorAndExit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		rr.CloseConn()
		os.Exit(1)
	}
}

func rate(user string, item string, score float64) {
	fmt.Printf("User %s ranked item %s with %.2f\n", user, item, score)
	err := rr.Rate(item, user, score)
	chekErrorAndExit(err)
}

func getProbability(user string, item string) {
	score, err := rr.CalcItemProbability(item, user)
	chekErrorAndExit(err)
	fmt.Printf("%s %s %.2f\n", user, item, score)
}

func suggest(user string, max int) {
	fmt.Printf("Getting %d results for user %s\n", max, user)
	rr.UpdateSuggestedItems(user, max)
	s, err := rr.GetUserSuggestions(user, max)
	chekErrorAndExit(err)
	fmt.Println("results:")
	fmt.Println(s)
}

func update(max int) {
	fmt.Printf("Updating DB\n")
	err := rr.BatchUpdateSimilarUsers(max)
	chekErrorAndExit(err)
}
