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

func toJpDate(s string) string {
	wj := strings.NewReplacer(
		"Sun", "日",
		"Mon", "月",
		"Tue", "火",
		"Wed", "水",
		"Thu", "木",
		"Fri", "金",
		"Sat", "土",
	)
	return wj.Replace(s)
}

func run(year int, month int, day int, span int) int {
	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	dm := DateMenu{start: start, span: span}
	selected, err := dm.selectDate()
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			return 1
		}
		return 0
	}
	me := MenuEntry{selected}
	df := me.preview()
	if 0 < len(df) {
		for _, d := range selected {
			j := toJpDate(d.Format(df))
			fmt.Printf("%v\n", j)
		}
	}
	return 0
}

type DateMenu struct {
	start time.Time
	span  int
}

func (d DateMenu) getMenu() []time.Time {
	var ts []time.Time
	for i := 0; i < d.span; i++ {
		ts = append(ts, d.start.AddDate(0, 0, i))
	}
	return ts
}

func (d DateMenu) selectDate() ([]time.Time, error) {
	menu := d.getMenu()
	idxs, err := fuzzyfinder.FindMulti(menu, func(i int) string {
		return menu[i].Format("01/02 (Mon)")
	})
	if err != nil {
		return []time.Time{}, err
	}
	sort.Ints(idxs)
	var selected []time.Time
	for _, i := range idxs {
		selected = append(selected, menu[i])
	}
	return selected, nil
}

type MenuEntry struct {
	dates []time.Time
}

func (m MenuEntry) getTable() map[string]string {
	return map[string]string{
		"MM-dd":            "01-02",
		"MM/dd":            "01/02",
		"MM月dd日（aaa）":      "01月02日（Mon）",
		"M月d日（aaa）":        "1月2日（Mon）",
		"YYYY-MM-dd":       "2006-01-02",
		"YYYY/MM/dd":       "2006/01/02",
		"YYYY年MM月dd日（aaa）": "2006年01月02日（Mon）",
		"YYYY年M月d日（aaa）":   "2006年1月2日（Mon）",
	}
}

func (m MenuEntry) applyFormat(fmt string) []string {
	var ss []string
	for _, d := range m.dates {
		j := toJpDate(d.Format(fmt))
		ss = append(ss, j)
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
