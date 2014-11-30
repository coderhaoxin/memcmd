package main

import "github.com/bradfitz/gomemcache/memcache"
import . "github.com/tj/go-debug"
import "strings"
import "strconv"
import "bufio"
import "time"
import "fmt"
import "os"

var client *memcache.Client
var debug = Debug("memcmd")

func init() {
	var host string

	if len(os.Args) >= 2 && !strings.HasPrefix(os.Args[1], "-") {
		host = os.Args[1]
	} else {
		host = "localhost:11211"
	}

	fmt.Println("connect to " + host)
	client = memcache.New(host)

	// test connection
	err := set("test", "connection", 1)
	if err != nil {
		fmt.Println("can not connect to memcached")
		os.Exit(1)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	log("welcome")

	for true {
		input, _ := reader.ReadString('\n')

		if len(input) < 2 {
			continue
		}

		args := parseInput(input)
		debug("args: %v", args)

		cmd := strings.ToLower(args[0])
		debug("cmd: %s", cmd)

		switch cmd {
		case "set":
			runSet(args)
		case "get":
			runGet(args)
		case "getmulti":
			getMulti(args[1:])
		case "delete":
			runDelete(args)
		case "touch":
			runTouch(args)
		case "deleteall":
			runDeleteAll()
		case "flushall":
			runFlushAll()
		case "q":
			os.Exit(0)
		default:
			log("usage:", "set", "get", "getMulti", "delete", "touch", "deleteAll", "flushAll", "q")
		}

		time.Sleep(time.Second * 1)
	}
}

// run cmd
func runSet(args []string) {
	if len(args) < 3 {
		log("set: key value required")
		return
	}

	key, value := args[1], args[2]

	expire := 60
	if len(args) > 3 {
		var err error
		expire, err = strconv.Atoi(args[3])
		if err != nil {
			log("invalid expire")
			return
		}
	}

	assertError(set(key, value, int32(expire)), "set error:")
}

func runGet(args []string) {
	if len(args) < 2 {
		log("get: key required")
		return
	}

	key := args[1]
	value := get(key)

	if value == "" {
		log(key + " not found")
	} else {
		log(value)
	}
}

func runDelete(args []string) {
	if len(args) < 2 {
		log("delete: key required")
		return
	}

	key := args[1]

	assertError(delete(key), "delete error:")
}

func runTouch(args []string) {
	if len(args) < 3 {
		log("touch: key and seconds required")
		return
	}

	expire, err := strconv.Atoi(args[2])

	if err != nil {
		log("invalid seconds")
		return
	}

	assertError(touch(args[1], int32(expire)), "touch error:")
}

func runDeleteAll() {
	assertError(deleteAll(), "delete all error:")
}

func runFlushAll() {
	assertError(flushAll(), "flush all error:")
}

// apis
func set(key, value string, expire int32) error {
	return client.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: expire,
	})
}

func get(key string) string {
	item, err := client.Get(key)

	if err != nil {
		return ""
	}

	return string(item.Value)
}

func getMulti(keys []string) {
	debug("keys: %v", keys)

	if len(keys) == 0 {
		log("getMulti: keys required")
		return
	}

	items, err := client.GetMulti(keys)

	if err != nil {
		debug("get multi error: %s", err.Error())
		log("not found")
		return
	}

	for k, v := range items {
		log(k+":", string(v.Value))
	}
}

func delete(key string) error {
	return client.Delete(key)
}

func touch(key string, seconds int32) error {
	return client.Touch(key, seconds)
}

func deleteAll() error {
	return client.DeleteAll()
}

func flushAll() error {
	return client.FlushAll()
}

// utils
func parseInput(input string) []string {
	debug("input: %v", input)

	var args []string

	ss := strings.Split(strings.TrimSpace(input), " ")

	for _, v := range ss {
		if v != "" {
			args = append(args, v)
		}
	}

	return args
}

func assertError(e error, message string) {
	if e != nil {
		log(message, e.Error())
	} else {
		log("success")
	}
}

func log(messages ...string) {
	for _, s := range messages {
		fmt.Print(s + " ")
	}

	fmt.Println()
}
