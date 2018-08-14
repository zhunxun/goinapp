package ios

import (
	"sort"
	"strconv"
	"time"
)

type SortType int

const (
	ByPurchaseDate SortType = iota
	ByOriginalPurchaseDate
)

func (i InApps) Sorted(by SortType) InApps {
	inapps := i
	switch by {
	case ByPurchaseDate:
		sort.Sort(byPurchaseDate(inapps))
	case ByOriginalPurchaseDate:
		sort.Sort(byOriginalPurchaseDate(inapps))
	default:
		return inapps
	}
	return inapps
}

// byPurchaseDate type implements sort.Interface and used to sort an array pf in-apps by purchase date
type byPurchaseDate InApps

func (b byPurchaseDate) Len() int      { return len(b) }
func (b byPurchaseDate) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byPurchaseDate) Less(i, j int) bool {
	bi, _ := parseDateMS(b[i].PurchaseDateMS)
	bj, _ := parseDateMS(b[j].PurchaseDateMS)
	return bi.Before(*bj)
}

// byOriginalPurchaseDate type implements sort.Interface and used to sort an array pf in-apps by original purchase date
type byOriginalPurchaseDate InApps

func (b byOriginalPurchaseDate) Len() int      { return len(b) }
func (b byOriginalPurchaseDate) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byOriginalPurchaseDate) Less(i, j int) bool {
	bi, _ := parseDateMS(b[i].PurchaseDateMS)
	bj, _ := parseDateMS(b[j].PurchaseDateMS)
	return bi.Before(*bj)
}

func parseDateMS(timeMS string) (*time.Time, error) {
	parseInt, err := strconv.ParseInt(timeMS, 10, 64)
	if err != nil {
		return nil, err
	}
	parseTime := time.Unix(0, parseInt*int64(time.Millisecond))
	return &parseTime, nil
}
