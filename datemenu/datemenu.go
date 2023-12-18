package datemenu

import (
	"sort"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
)

type DateMenu struct {
	Start   time.Time
	Span    int
	Weekday bool
}

func (d DateMenu) getMenu() []time.Time {
	var ts []time.Time
	for i := 0; i < d.Span; i++ {
		t := d.Start.AddDate(0, 0, i)
		if d.Weekday && (t.Weekday() == time.Saturday || t.Weekday() == time.Sunday) {
			continue
		}
		ts = append(ts, t)
	}
	return ts
}

func (d DateMenu) SelectDate() ([]time.Time, error) {
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
