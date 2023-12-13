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
	me := MenuEntry{selected}
	df := me.preview()
	if 0 < len(df) {
		for _, d := range selected {
			fmt.Printf("%v\n", d.Format(df))
		}
	}
	return 0
}

func getDayMenu(start time.Time, span int) []time.Time {
	var ts []time.Time
	for i := 0; i < span; i++ {
		ts = append(ts, start.AddDate(0, 0, i))
	}
	return ts
}

type MenuEntry struct {
	dates []time.Time
}

func (m MenuEntry) getTable() map[string]string {
	return map[string]string{
		"MM-dd":            "01-02",
		"MM/dd":            "01/02",
		"MM月dd日（ddd）":      "01月02日（Mon）",
		"M月d日（ddd）":        "1月2日（Mon）",
		"YYYY-MM-dd":       "2006-01-02",
		"YYYY/MM/dd":       "2006/01/02",
		"YYYY年MM月dd日（ddd）": "2006年01月02日（Mon）",
		"YYYY年M月d日（ddd）":   "2006年1月2日（Mon）",
	}
}

func (m MenuEntry) applyFormat(fmt string) []string {
	var ss []string
	for _, d := range m.dates {
		ss = append(ss, d.Format(fmt))
	}
	return ss
}

func (m MenuEntry) getCommonFormats() []string {
	var ss []string
	for f := range m.getTable() {
		ss = append(ss, f)
	}
	sort.Strings(ss)
	return ss
}

func (m MenuEntry) toGoFormat(commonFmt string) string {
	return m.getTable()[commonFmt]
}

func (m MenuEntry) preview() string {
	fmts := m.getCommonFormats()
	idx, err := fuzzyfinder.Find(fmts, func(i int) string {
		return fmts[i]
	}, fuzzyfinder.WithPreviewWindow(func(i, _, _ int) string {
		if i == -1 {
			return ""
		}
		f := m.toGoFormat(fmts[i])
		return strings.Join(m.applyFormat(f), "\n")
	}))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return m.toGoFormat(fmts[idx])
}
