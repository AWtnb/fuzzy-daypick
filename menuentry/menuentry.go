package menuentry

import (
	"fmt"
	"slices"
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

type MenuEntry struct {
	Dates []time.Time
}

func (m MenuEntry) toIntSlice(fmt string) []int {
	var sl []int
	for _, m := range m.Dates {
		i, err := strconv.Atoi(m.Format(fmt))
		if err == nil {
			sl = append(sl, i)
		}
	}
	return sl
}

func (m MenuEntry) getTable() map[string]string {
	ds := map[string]string{
		"d日": "2日（Monday）",
	}
	if isPaddable(m.toIntSlice("2")) {
		ds["dd日"] = "02日（Monday）"
		ds["_d日"] = "_2日（Monday）"
	}
	ms := map[string]string{
		"M月": "1月",
	}
	if isPaddable(m.toIntSlice("1")) {
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
