package p

import (
	"db"
	"fmt"
	"log"
	"math"
)

func test() {
	log.Println("foo", math.Abs(-123))                // want "log2zap"
	log.Printf("set some var=%d err=%v", 215, nil)    // want "log2zap"
	log.Printf("set some var=%d err=%v", 215+10, nil) // want "log2zap"

	d := db.Cluster{
		Master:  db.DB{},
		Replica: db.DB{},
	}
	defer func() {
		d.CleanUp()
	}()

	foo := "foo"
	data := map[string]func() string{
		"test": func() string {
			return "map-test"
		},
	}

	fmt.Print("foo", "bar", d.Master.String(), []byte(foo), data["test"]()) // want "log2zap"
}

func formatMessage(msg string, args ...interface{}) string {
	x := fmt.Sprintf(msg, args...)
	if x == "" {
		x = "empty"
	}
	return x
}
