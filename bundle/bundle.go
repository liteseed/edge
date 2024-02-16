package bundle

import (
	"fmt"
	"log"
	"math"
	"strings"
	"sync"

	"github.com/everFinance/go-everpay/account"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/liteseed/bungo/cache"
	"github.com/liteseed/bungo/database"
	"github.com/liteseed/bungo/schema"
	"github.com/liteseed/bungo/store"
	"github.com/shopspring/decimal"
)

type Bundle struct {
	cache           *cache.Cache
	db              *database.Database
	store           *store.Store
	bundlePerFeeMap map[string]schema.Fee
	locker          sync.RWMutex
}

func (bundle *Bundle) ProcessSubmitItem(item types.BundleItem, currency string, size int64) (schema.Order, error) {
	if err := utils.VerifyBundleItem(item); err != nil {
		return schema.Order{}, err
	}
	if item.DataReader != nil { // reset io stream to origin of the file
		if _, err := item.DataReader.Seek(0, 0); err != nil {
			return schema.Order{}, err
		}
	}
	if err := bundle.store.AtomicSaveItem(item); err != nil {
		return schema.Order{}, err
	}

	signerAddr, err := utils.ItemSignerAddr(item)
	if err != nil {
		return schema.Order{}, err
	}
	_, accId, err := account.IDCheck(signerAddr)
	if err != nil {
		return schema.Order{}, err
	}
	order := schema.Order{
		ItemId:        item.Id,
		Signer:        accId,
		SignType:      item.SignatureType,
		Size:          size,
		ExpectedBlock: bundle.cache.GetInfo().Height + 10,
		Status:        schema.WaitOnChain,
	}
	// calc fee
	// respFee, err := s.CalcItemFee(currency, order.Size)
	// if err != nil {
	// 	return schema.Order{}, err
	// }
	// order.Decimals = respFee.Decimals
	// order.Fee = respFee.FinalFee
	// order.Currency = strings.ToUpper(currency)

	// if isNoFeeMode {
	// 	order.PaymentStatus = schema.SuccPayment
	// } else {
	// 	order.PaymentExpiredTime = time.Now().Unix()
	// 	order.PaymentStatus = schema.UnPayment
	// }

	// // insert to mysql
	// if err = s.db.InsertOrder(order); err != nil {
	// 	return schema.Order{}, err
	// }
	return order, nil
}

// func (s *Bungo) CalcItemFee(currency string, itemSize int64) (*schema.RespFee, error) {
// 	perFee :=
// 	if perFee == nil {
// 		return nil, fmt.Errorf("not support currency: %s", currency)
// 	}

// 	count := int64(0)
// 	if itemSize > 0 {
// 		count = (itemSize-1)/types.MAX_CHUNK_SIZE + 1
// 	}

// 	chunkFees := decimal.NewFromInt(count).Mul(perFee.PerChunk)
// 	finalFee := perFee.Base.Add(chunkFees)

// 	return &schema.RespFee{
// 		Currency: perFee.Currency,
// 		Decimals: perFee.Decimals,
// 		FinalFee: finalFee.String(),
// 	}, nil
// }

func (b *Bundle) GetBundlePerFees() (map[string]schema.Fee, error) {
	arPrice, err := b.db.GetArPrice()
	if err != nil {
		return nil, err
	}
	tps, err := b.db.GetPrices()
	if err != nil {
		return nil, err
	}
	arFee := b.cache.GetFee()
	res := make(map[string]schema.Fee)
	for _, tp := range tps {
		if tp.Price <= 0.0 {
			continue
		}

		// fee = 1e(tpDecimals) * arPrice * arBaseFee / 1e(arDeciamls) / tpPrice
		baseFee := decimal.NewFromFloat(math.Pow10(tp.Decimals)).Mul(decimal.NewFromFloat(arPrice)).Mul(decimal.NewFromInt(arFee.Base)).
			Div(decimal.NewFromFloat(math.Pow10(12))).Div(decimal.NewFromFloat(tp.Price)).Round(0)

		perChunkFee := decimal.NewFromFloat(math.Pow10(tp.Decimals)).Mul(decimal.NewFromFloat(arPrice)).Mul(decimal.NewFromInt(arFee.PerChunk)).
			Div(decimal.NewFromFloat(math.Pow10(12))).Div(decimal.NewFromFloat(tp.Price)).Round(0)

		res[strings.ToUpper(tp.Symbol)] = schema.Fee{
			Currency: tp.Symbol,
			Decimals: tp.Decimals,
			Base:     baseFee,
			PerChunk: perChunkFee,
		}
	}
	return res, nil
}

func (s *Bundle) ParseAndSaveBundleItems(arId string, data []byte) error {
	if s.store.ExistArIdToItemIds(arId) {
		return nil
	}

	bundle, err := utils.DecodeBundle(data)
	if err != nil {
		return err
	}
	itemIds := make([]string, 0, len(bundle.Items))
	// save items
	for _, item := range bundle.Items {
		if err = s.store.AtomicSaveItem(item); err != nil {
			log.Fatal("s.saveItem(item)", "err", err, "arId", arId)
			return err
		}

		itemIds = append(itemIds, item.Id)
	}

	// save arId to itemIds
	return s.store.SaveArIdToItemIds(arId, itemIds)
}

func (s *Bundle) CalcItemFee(currency string, itemSize int64) (*schema.RespFee, error) {
	perFee := s.GetPerFee(currency)
	if perFee == nil {
		return nil, fmt.Errorf("not support currency: %s", currency)
	}

	count := int64(0)
	if itemSize > 0 {
		count = (itemSize-1)/types.MAX_CHUNK_SIZE + 1
	}

	chunkFees := decimal.NewFromInt(count).Mul(perFee.PerChunk)
	finalFee := perFee.Base.Add(chunkFees)

	return &schema.RespFee{
		Currency: perFee.Currency,
		Decimals: perFee.Decimals,
		FinalFee: finalFee.String(),
	}, nil
}

func (s *Bundle) GetPerFee(tokenSymbol string) *schema.Fee {
	s.locker.RLock()
	defer s.locker.RUnlock()
	perFee, ok := s.bundlePerFeeMap[strings.ToUpper(tokenSymbol)]
	if !ok {
		return nil
	}
	return &perFee
}

func (s *Bundle) SetPerFee(feeMap map[string]schema.Fee) {
	s.locker.Lock()
	s.bundlePerFeeMap = feeMap
	s.locker.Unlock()
}
