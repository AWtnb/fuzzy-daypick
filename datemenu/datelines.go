package datemenu

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

type DateLines struct {
	Lines []string
}

func (dl DateLines) getTable() map[string]string {
	return map[string]string{
		"(Null)":       "",
		"dash":         "- ",
		"bullet":       "・",
		"indent(2)":    "  ",
		"indent(Full)": "　",
	}
}

func (dl DateLines) getMenuKeys() []string {
	var ss []string
	for f := range dl.getTable() {
		ss = append(ss, f)
	}
	slices.Sort(ss)
	return ss
}

func (dl DateLines) applyPrefix(pre string) []string {
	var ss []string
	for _, d := range dl.Lines {
		ss = append(ss, pre+d)
	}
	return ss
}

func (dl DateLines) toPrefix(menuKey string) string {
	return dl.getTable()[menuKey]
}

func (dl DateLines) SelectPrefix() string {
	keys := dl.getMenuKeys()
	idx, err := fuzzyfinder.Find(keys, func(i int) string {
		return keys[i]
	}, fuzzyfinder.WithPreviewWindow(func(i, _, _ int) string {
		if i == -1 {
			return ""
		}
		p := dl.toPrefix(keys[i])
		return strings.Join(dl.applyPrefix(p), "\n")
	}), fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return dl.toPrefix(keys[idx])
}
