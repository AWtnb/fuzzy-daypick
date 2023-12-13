package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		year  int
		month int
		day   int
		span  int
	)
	flag.IntVar(&year, "year", 0, "year to start menu")
	flag.IntVar(&month, "month", 0, "month to start menu")
	flag.IntVar(&day, "day", 0, "day to start menu")
	flag.IntVar(&span, "span", 30, "span of date menu")
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
	os.Exit(run(year, month, day, span))
}

func run(year int, month int, day int, span int) int {
	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	menu := getDayMenu(start, span)
	idxs, err := fuzzyfinder.FindMulti(menu, func(i int) string {
		return menu[i].Format("01/02 (Mon)")
	})
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			fmt.Printf("ERROR: %s\n", err)
			return 1
		}
		return 0
	}
	sort.Ints(idxs)
	var selected []time.Time
	for _, i := range idxs {
		selected = append(selected, menu[i])
	}
	df := selectFormat(selected)
	if 0 < len(df) {
		for _, d := range selected {
			fmt.Printf("%v\n", d.Format(df))
		}
	}
	return 0
}

func getDayMenu(start time.Time, span int) []time.Time {
	var ds []time.Time
	for i := 0; i < span; i++ {
		ds = append(ds, start.AddDate(0, 0, i))
	}
	return ds
}

var FormatMap = map[string]string{
	"MM-dd":       "01-02",
	"MM/dd":       "01/02",
	"MM年dd月（ddd）": "01月02日（Mon）",
}

// type DateEntry struct {
// 	Dates []time.Time
// }

// func (d DateEntry) Format(fmt string) []string {
// 	var ss []string
// 	f := FormatMap[fmt]
// 	for _, d := range d.Dates {
// 		ss = append(ss, d.Format(f))
// 	}
// 	return ss
// }

func selectFormat(dates []time.Time) string {
	var fmts []string
	for k := range FormatMap {
		fmts = append(fmts, k)
	}
	sort.Strings(fmts)
	idx, err := fuzzyfinder.Find(fmts, func(i int) string {
		return fmts[i]
	}, fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
		if i == -1 {
			return ""
		}
		f := FormatMap[fmts[i]]
		var ds []string
		for _, d := range dates {
			ds = append(ds, d.Format(f))
		}
		return strings.Join(ds, "\n")
	}))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return FormatMap[fmts[idx]]
}
