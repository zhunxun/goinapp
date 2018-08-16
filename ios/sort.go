package ios

import (
	"sort"
)

// SortType represent enumeration of sorting types for Sorted function
type SortType int

const (
	ByPurchaseDate SortType = iota
	ByOriginalPurchaseDate
)

// Sorted sort InApps array of InApp by different values like:
// - ByPurchaseDate
// - ByOriginalPurchaseDate
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
	bi := convertToTime(b[i].PurchaseDateMS)
	bj := convertToTime(b[j].PurchaseDateMS)
	return bi.Before(bj)
}

// byOriginalPurchaseDate type implements sort.Interface and used to sort an array of in-apps by original purchase date
type byOriginalPurchaseDate InApps

func (b byOriginalPurchaseDate) Len() int      { return len(b) }
func (b byOriginalPurchaseDate) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byOriginalPurchaseDate) Less(i, j int) bool {
	bi := convertToTime(b[i].PurchaseDateMS)
	bj := convertToTime(b[j].PurchaseDateMS)
	return bi.Before(bj)
}
