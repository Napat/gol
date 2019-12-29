package main

import (
	"os"

	"github.com/Napat/gol"
)

func main() {

	// Create log dir
	os.MkdirAll("log", os.ModePerm)

	// Log to log/ingress_YYYYMMDD.log, Rotate file by Bangkok date time.
	x, err := gol.NewWithHookAutoFileDate("log/ingress.log", "Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	// Log to log/egress_YYYYMMDD.log, Rotate file by Bangkok date time.
	y, err := gol.NewJSONWithHookAutoFileDate("log/egress.log", "Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	x.GDebugf("GInfof debug")

	x.WithFields(gol.Fields{
		"isAwesome": true,
		"star":      9999,
	}).GInfof("Thailand has many awesome places to explore")

	x.WithFields(gol.Fields{
		"isBeauty": true,
		"star":     9999,
	}).GErrorf("Music is beautiful in a way that nothing else can be")

	y.WithFields(gol.Fields{
		"service":     "golapp",
		"status code": 0,
	}).GDebugf("Healthy")

}
