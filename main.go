package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/AWtnb/fuzzy-daypick/datemenu"
)

func main() {
	var (
		start   string
		year    int
		month   int
		day     int
		span    int
		weekday bool
	)
	flag.StringVar(&start, "start", "", "start day in yyyyMMdd format")
	flag.IntVar(&year, "year", 0, "year to start menu")
	flag.IntVar(&month, "month", 0, "month to start menu")
	flag.IntVar(&day, "day", 0, "day to start menu")
	flag.IntVar(&span, "span", 30, "span of date menu")
	flag.BoolVar(&weekday, "weekday", false, "skip saturday and sunday")
	flag.Parse()
	ln := len(start)
	if 4 <= ln {
		switch ln {
		case 4:
			start = fmt.Sprintf("%s0101", start)
		case 5:
			start = fmt.Sprintf("%s0%s01", start[0:4], start[4:5])
		case 6:
			start = fmt.Sprintf("%s01", start)
		case 7:
			start = fmt.Sprintf("%s0%s", start[0:6], start[6:7])
		}
		year, err := strconv.Atoi(start[0:4])
		if err != nil {
			fmt.Println(err)
			return
		}
		month, err = strconv.Atoi(start[4:6])
		if err != nil {
			fmt.Println(err)
			return
		}
		day, err = strconv.Atoi(start[6:8])
		if err != nil {
			fmt.Println(err)
			return
		}
		os.Exit(run(year, month, day, span, weekday))
		return
	}
	now := time.Now()
	if year < 1 && month < 1 && day < 1 {
		os.Exit(run(now.Year(), int(now.Month()), now.Day(), span, weekday))
		return
	}
	if year < 1 {
		year = now.Year()
	}
	if month < 1 {
		month = 1
	}
	if day < 1 {
		day = 1
	}
	os.Exit(run(year, month, day, span, weekday))
}

func run(year int, month int, day int, span int, weekday bool) int {
	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	dm := datemenu.DateMenu{Start: start, Span: span, Weekday: weekday}
	selected, err := dm.SelectDate()
	if err != nil {
		return 1
	}
	me := datemenu.MenuEntry{Dates: selected}
	df := me.Preview()
	if 0 < len(df) {
		ss := datemenu.ToLines(selected, df)
		dl := datemenu.DateLines{Lines: ss}
		f := dl.SelectPrefix()
		for _, s := range ss {
			fmt.Printf("%s\n", f+s)
		}
	}
	return 0
}
