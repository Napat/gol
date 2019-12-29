package gol

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// HookPrinter ...
type HookPrinter struct {
}

func (h *HookPrinter) Fire(entry *logrus.Entry) error {
	s, _ := entry.String()
	fmt.Printf("%s", s)

	// WARNING
	// Using go-routine with entry.String() will cause race-condition: https://github.com/sirupsen/logrus/issues/1012
	// go func() {
	// 	s, _ := entry.String()
	// 	fmt.Printf("here is the msg: <%s>\n", s)
	//  fmt.Printf("%s", s)
	// }()
	return nil
}

func (h *HookPrinter) Levels() []logrus.Level {
	return logrus.AllLevels
}

type HookAutoFileDate struct {
	Logger       *Gol
	TimeLocation *time.Location
	previousDate string
	Logpath      string
}

// Fire HookAutoFileDate Auto print log to stdout
func (h *HookAutoFileDate) Fire(entry *logrus.Entry) error {
	s, _ := entry.String()
	fmt.Fprintf(os.Stdout, "%s", s)

	h.FileDateHandler()
	return nil
}

func (h *HookAutoFileDate) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *HookAutoFileDate) FileDateHandler() {

	ctime := time.Now()
	// fmt.Println("Location:", ctime.Location(), ":Time:", ctime)
	// fmt.Println(ctime.In(h.TimeLocation).Format("2006-01-02 15:04:05"))
	currentDate := ctime.In(h.TimeLocation).Format("20060102")

	if currentDate == h.previousDate {
		// do nothing
		// fmt.Printf("NOT Change!\n")
		return
	}
	// fmt.Printf("Change!\n")

	dir := filepath.Dir(h.Logger.FullPath)
	base := filepath.Base(h.Logger.FullPath)
	ext := filepath.Ext(h.Logger.FullPath)
	filename := strings.TrimSuffix(base, ext)
	// fmt.Printf("fullpath: %v, dir: %v, base: %v, filename: %v ext: %v\n", h.Logger.FullPath, dir, base, filename, ext)

	logpath := fmt.Sprintf("%s/%s_%s%s", dir, filename, currentDate, ext)
	// fmt.Printf("logpath: %s\n", logpath)

	h.Logger.setFileOutput(logpath, false, false)
	h.Logpath = logpath
	h.previousDate = currentDate

}
