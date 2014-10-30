package main

import "github.com/bradfitz/gomemcache/memcache"
import "strings"
import "strconv"
import "bufio"
import "time"
import "fmt"
import "os"

var client *memcache.Client

func main() {
	host := os.Args[1]

	fmt.Println("connect to " + host)
	client = memcache.New(host)

	// test connection
	err := Set("test", "connection", 1)
	if err != nil {
		fmt.Println("can not connect to memcached")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	log("welcome")

	for true {
		input, _ := reader.ReadString('\n')

		if len(input) < 3 {
			continue
		}

		args := parseInput(input)

		switch strings.ToLower(args[0]) {
		case "set":
			if len(args) < 3 {
				log("set: key value required")
				break
			}

			key, value := args[1], args[2]

			expire := 60
			if len(args) > 3 {
				var err error
				expire, err = strconv.Atoi(args[3])
				if err != nil {
					log("invalid expire")
					break
				}
			}

			err := Set(key, value, int32(expire))

			if err != nil {
				fmt.Println("set error")
				fmt.Println(err.Error())
			} else {
				log("set success")
			}
		case "get":
			if len(args) < 2 {
				log("get: key required")
				break
			}

			key := args[1]
			value := Get(key)

			if value == "" {
				log(key + " not found")
			} else {
				log(value)
			}
		case "q":
			os.Exit(0)
		}

		time.Sleep(time.Second * 1)
	}
}

func Set(key, value string, expire int32) error {
	return client.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: expire,
	})
}

func Get(key string) string {
	item, err := client.Get(key)

	if err != nil {
		return ""
	}

	return string(item.Value)
}

func parseInput(input string) []string {
	var args []string

	ss := strings.Split(strings.TrimSpace(input), " ")

	for _, v := range ss {
		if v != "" {
			args = append(args, v)
		}
	}

	return args
}

func log(msg string) {
	fmt.Println(msg)
	fmt.Println()
}
