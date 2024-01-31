package datemenu

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
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
	j := wj.Replace(s)
	if strings.HasPrefix(j, "_") {
		pad := ""
		if strings.Index(j, "月") == 2 {
			pad = " "
		}
		j = strings.Replace(j, "_", pad, 1)
	}
	return j
}

func ToLines(dates []time.Time, fmt string) []string {
	var ss []string
	for _, d := range dates {
		j := toJpDate(d.Format(fmt))
		ss = append(ss, j)
	}
	return ss
}

func isPaddable(ns []int) bool {
	slices.Sort(ns)
	return ns[0] < 10 && 10 <= ns[len(ns)-1]
}

func isSingleMonth(ts []time.Time) bool {
	sort.Slice(ts, func(i int, j int) bool {
		return ts[i].Before(ts[j])
	})
	return ts[0].Month() == ts[len(ts)-1].Month()
}

type MenuEntry struct {
	Dates []time.Time
}

func (m MenuEntry) months() (ms []int) {
	for _, m := range m.Dates {
		i, err := strconv.Atoi(m.Format("1"))
		if err == nil {
			ms = append(ms, i)
		}
	}
	return
}

func (m MenuEntry) days() (ds []int) {
	for _, m := range m.Dates {
		i, err := strconv.Atoi(m.Format("2"))
		if err == nil {
			ds = append(ds, i)
		}
	}
	return
}

func (m MenuEntry) getTable() map[string]string {
	ds := map[string]string{
		"d日": "2日（Monday）",
	}
	if isPaddable(m.days()) {
		ds["dd日"] = "02日（Monday）"
		ds["_d日"] = "_2日（Monday）"
	}
	if isSingleMonth(m.Dates) {
		return ds
	}
	ms := map[string]string{
		"M月": "1月",
	}
	if isPaddable(m.months()) {
		ms["MM月"] = "01月"
		ms["_M月"] = "_1月"
	}
	mp := map[string]string{}
	for km, vm := range ms {
		for kd, vd := range ds {
			mp[km+kd] = vm + vd
		}
	}
	return mp
}

func (m MenuEntry) applyFormat(fmt string) []string {
	return ToLines(m.Dates, fmt)
}

func (m MenuEntry) getMenuKeys() []string {
	var ss []string
	for f := range m.getTable() {
		ss = append(ss, f)
	}
	slices.Sort(ss)
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
