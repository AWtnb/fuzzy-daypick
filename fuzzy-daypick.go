package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/AWtnb/fuzzy-daypick/datelines"
	"github.com/AWtnb/fuzzy-daypick/datemenu"
	"github.com/AWtnb/fuzzy-daypick/menuentry"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		year    int
		month   int
		day     int
		span    int
		weekday bool
	)
	flag.IntVar(&year, "year", 0, "year to start menu")
	flag.IntVar(&month, "month", 0, "month to start menu")
	flag.IntVar(&day, "day", 0, "day to start menu")
	flag.IntVar(&span, "span", 30, "span of date menu")
	flag.BoolVar(&weekday, "weekday", false, "skip saturday and sunday")
	flag.Parse()
	now := time.Now()
	if year < 1 {
		year = now.Year()
	}
	if month < 1 {
		month = int(now.Month())
	}
	if day < 1 {
		day = now.Day()
	}
	os.Exit(run(year, month, day, span, weekday))
}

func run(year int, month int, day int, span int, weekday bool) int {
	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	dm := datemenu.DateMenu{Start: start, Span: span, Weekday: weekday}
	selected, err := dm.SelectDate()
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			return 1
		}
		return 0
	}
	me := menuentry.MenuEntry{Dates: selected}
	df := me.Preview()
	if 0 < len(df) {
		ss := menuentry.ToLines(selected, df)
		dl := datelines.DateLines{Lines: ss}
		f := dl.SelectPrefix()
		for _, s := range ss {
			fmt.Printf("%s\n", f+s)
		}
	}
	return 0
}
