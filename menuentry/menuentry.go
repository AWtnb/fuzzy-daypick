package menuentry

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
)

func toJpDate(s string) string {
	wj := strings.NewReplacer(
		"Sunday", "日",
		"Monday", "月",
		"Tuesday", "火",
		"Wednesday", "水",
		"Thursday", "木",
		"Friday", "金",
		"Saturday", "土",
	)
	return wj.Replace(s)
}

type MenuEntry struct {
	Dates []time.Time
}

func (m MenuEntry) getTable() map[string]string {
	return map[string]string{
		"MM月 d日（aaa）": "01月_2日（Monday）",
		"MM月dd日（aaa）": "01月02日（Monday）",
		"MM月d日（aaa）":  "01月2日（Monday）",
		"M月 d日（aaa）":  "1月_2日（Monday）",
		"M月dd日（aaa）":  "1月02日（Monday）",
		"M月d日（aaa）":   "1月2日（Monday）",
	}
}

func (m MenuEntry) applyFormat(fmt string) []string {
	var ss []string
	for _, d := range m.Dates {
		j := toJpDate(d.Format(fmt))
		ss = append(ss, j)
	}
	return ss
}

func (m MenuEntry) getMenuKeys() []string {
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

func (m MenuEntry) Preview() string {
	fmts := m.getMenuKeys()
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

func ToLines(dates []time.Time, fmt string) []string {
	var ss []string
	for _, d := range dates {
		j := toJpDate(d.Format(fmt))
		ss = append(ss, j)
	}
	return ss
}
