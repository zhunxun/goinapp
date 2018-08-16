package ios

import (
	"sort"
	"time"
)

type SortType int

const (
	ByPurchaseDate SortType = iota
	ByOriginalPurchaseDate
)

func (i InApps) Sorted(by SortType) InApps {
	switch by {
	case ByPurchaseDate:
		sort.Sort(byPurchaseDate(i))
	case ByOriginalPurchaseDate:
		sort.Sort(byOriginalPurchaseDate(i))
	}
	return i
}

// byPurchaseDate type implements sort.Interface and used to sort an array ff in-apps by purchase date
type byPurchaseDate InApps

func (b byPurchaseDate) Len() int      { return len(b) }
func (b byPurchaseDate) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byPurchaseDate) Less(i, j int) bool {
	bi, _ := parseDateMS(b[i].PurchaseDateMS)
	bj, _ := parseDateMS(b[j].PurchaseDateMS)
	return bi.Before(*bj)
}

// byOriginalPurchaseDate type implements sort.Interface and used to sort an array of in-apps by original purchase date
type byOriginalPurchaseDate InApps

func (b byOriginalPurchaseDate) Len() int      { return len(b) }
func (b byOriginalPurchaseDate) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byOriginalPurchaseDate) Less(i, j int) bool {
	bi, _ := parseDateMS(b[i].PurchaseDateMS)
	bj, _ := parseDateMS(b[j].PurchaseDateMS)
	return bi.Before(*bj)
}

func parseDateMS(timeMS int64 /*string*/) (*time.Time, error) {
	//parseInt, err := strconv.Atoi(timeMS)
	//if err != nil {
	//	return nil, err
	//}
	//parseTime := time.Unix(0, int64(parseInt)*int64(time.Millisecond))
	parseTime := time.Unix(0, timeMS*int64(time.Millisecond))
	return &parseTime, nil
}
