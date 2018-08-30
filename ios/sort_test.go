package ios

import (
	"testing"
)

func TestInApps_Sorted(t *testing.T) {
	inapps := InApps{
		InApp{
			OriginalPurchaseDateMS: 1527811200002,
			CancellationDateMS:     1527811200002,
			ExpiresDateMS:          1527811200002,
			PurchaseDateMS:         1527811200002,
		},
		InApp{
			OriginalPurchaseDateMS: 1527811200001,
			CancellationDateMS:     1527811200001,
			ExpiresDateMS:          1527811200001,
			PurchaseDateMS:         1527811200001,
		},
		InApp{
			OriginalPurchaseDateMS: 1527811200000,
			CancellationDateMS:     1527811200000,
			ExpiresDateMS:          1527811200000,
			PurchaseDateMS:         1527811200000,
		},
		InApp{
			OriginalPurchaseDateMS: 1527811200002,
			CancellationDateMS:     1527811200002,
			ExpiresDateMS:          1527811200002,
			PurchaseDateMS:         1527811200002,
		},
		InApp{
			OriginalPurchaseDateMS: 1527811200005,
			CancellationDateMS:     1527811200005,
			ExpiresDateMS:          1527811200005,
			PurchaseDateMS:         1527811200005,
		},
	}

	max := int64(1527811200005)
	min := int64(1527811200000)

	t.Run("By Original Purchase Date", func(t *testing.T) {
		got := inapps.Sorted(ByOriginalPurchaseDate)
		if got[0].OriginalPurchaseDateMS != 1527811200005 {
			t.Fail()
		}
		if got[4].OriginalPurchaseDateMS != min {
			t.Fail()
		}
	})

	t.Run("By Purchase Date", func(t *testing.T) {
		got := inapps.Sorted(ByPurchaseDate)
		if got[0].OriginalPurchaseDateMS != max {
			t.Fail()
		}
		if got[4].OriginalPurchaseDateMS != min {
			t.Fail()
		}
	})

}
