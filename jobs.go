package bungo

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/everFinance/go-everpay/account"
	"github.com/everFinance/go-everpay/config"
	sdkSchema "github.com/everFinance/go-everpay/sdk/schema"
	paySchema "github.com/everFinance/go-everpay/token/schema"
	tokUtils "github.com/everFinance/go-everpay/token/utils"
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/liteseed/bungo/database"
	"github.com/liteseed/bungo/schema"
	"github.com/panjf2000/ants/v2"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

const (
	ItemPaymentAction   = "payment"
	ApikeyPaymentAction = "apikeyPayment"
)

func (s *Bungo) runJobs(bundleInterval int) {
	// update cache
	s.scheduler.Every(2).Minute().SingletonMode().Do(s.updateAnchor)
	s.scheduler.Every(2).Minute().SingletonMode().Do(s.updateArFee)
	s.scheduler.Every(30).Seconds().SingletonMode().Do(s.updateInfo)
	s.scheduler.Every(5).Minute().SingletonMode().Do(s.updatePeerMap)
	s.scheduler.Every(5).Minute().SingletonMode().Do(s.updateTokenPrice)
	s.scheduler.Every(1).Minute().SingletonMode().Do(s.updateBundlePerFee)

	s.scheduler.Every(bundleInterval).Seconds().SingletonMode().Do(s.onChainBundleItems) // can set a longer time, if the items are less. such as 2m
	// onChainBundleItems by upload order
	s.scheduler.Every(bundleInterval).Seconds().SingletonMode().Do(s.onChainItemsBySeq)
	s.scheduler.Every(3).Minute().SingletonMode().Do(s.watchArTx)
	s.scheduler.Every(2).Minute().SingletonMode().Do(s.retryOnChainArTx)

	// s.scheduler.Every(10).Seconds().SingletonMode().Do(s.parseAndSaveBundleTx) // todo stream

	// manager taskStatus
	s.scheduler.Every(5).Seconds().SingletonMode().Do(s.watcherAndCloseTasks)

	s.scheduler.Every(1).Minute().SingletonMode().Do(s.updateBundler)

	// delete tmp file, one may be repeat request same data,tmp file can be reserve with short time
	s.scheduler.Every(2).Minute().SingletonMode().Do(s.deleteTmpFile)

	//statistic
	s.scheduler.Every(1).Minute().SingletonMode().Do(s.UpdateRealTime)
	go s.ProduceDailyStatistic()
	s.scheduler.Every(1).Day().At("00:01").SingletonMode().Do(s.ProduceDailyStatistic)

	// kafka
	if len(s.KWriters) > 0 {
		s.scheduler.Every(10).Seconds().SingletonMode().Do(s.broadcastItemToKafka)
		s.scheduler.Every(1).Minute().SingletonMode().Do(s.broadcastBlockToKafka)
	}

	s.scheduler.StartAsync()
}

func (s *Bungo) updateAnchor() {
	anchor, err := fetchAnchor(s.arCli, s.cache.GetPeers())
	if err == nil {
		s.cache.UpdateAnchor(anchor)
	}
}

func (s *Bungo) updateInfo() {
	info, err := fetchArInfo(s.arCli, s.cache.GetPeers())
	if err == nil && info != nil {
		s.cache.UpdateInfo(*info)
	}
}

func (s *Bungo) updateArFee() {
	txPrice, err := fetchArFee(s.arCli, s.cache.GetPeers())
	if err == nil {
		s.cache.UpdateFee(txPrice)
	}
}

// update peer list concurrency, check peer available, save in db
func (s *Bungo) updatePeerMap() {
	peers, err := s.arCli.GetPeers()
	if err != nil {
		return
	}

	availablePeers := filterPeers(peers, s.cache.GetConstTx())
	if len(availablePeers) == 0 {
		return
	}

	peerMap := updatePeerMap(s.cache.GetPeerMap(), availablePeers)

	s.cache.UpdatePeers(peerMap)
	if err = s.store.SavePeers(s.cache.GetPeerMap()); err != nil {
		log.Warn("save new peer list fail")
	}
}

// bundle

func (s *Bungo) updateTokenPrice() {
	// update symbol
	tps := make([]schema.TokenPrice, 0)
	for _, tok := range s.everpaySdk.GetTokens() {
		tps = append(tps, schema.TokenPrice{
			Symbol:    strings.ToUpper(tok.Symbol),
			Decimals:  tok.Decimals,
			ManualSet: false,
			UpdatedAt: time.Time{},
		})
	}
	if err := s.database.InsertPrices(tps); err != nil {
		log.Error("s.query.InsertPrices(tps)", "err", err)
		return
	}

	// update fee
	tps, err := s.database.GetPrices()
	if err != nil {
		log.Error("s.query.GetPrices()", "err", err)
		return
	}
	for _, tp := range tps {
		if tp.ManualSet {
			continue
		}
		price, err := config.GetTokenPriceByRedstone(tp.Symbol, "USDC", "")
		if err != nil {
			log.Error("config.GetTokenPriceByRedstone(tp.Symbol,\"USDC\")", "err", err, "symbol", tp.Symbol)
			continue
		}
		if price <= 0.0 {
			// log.Error("GetTokenPriceByRedstone return price less than 0.0", "token", tp.Symbol)
			continue
		}
		// update tokenPrice
		if err := s.database.UpdatePrice(tp.Symbol, price); err != nil {
			log.Error("s.query.UpdateFee(tp.Symbol,fee)", "err", err, "symbol", tp.Symbol, "fee", price)
		}
	}
}

func (s *Bungo) updateBundlePerFee() {
	feeMap, err := s.GetBundlePerFees()
	if err != nil {
		log.Error("s.GetBundlePerFees()", "err", err)
		return
	}
	print(feeMap)
}

func (s *Bungo) watcherAndCloseTasks() {
	tasks := s.taskMg.GetTasks()
	now := time.Now().Unix()
	for _, tk := range tasks {
		if tk.Close || tk.Timestamp == 0 { // timestamp == 0  means do not start
			continue
		}
		// spend time not more than 30 minutes
		if now-tk.Timestamp > 30*60 {
			if err := s.taskMg.CloseTask(tk.ArId, tk.TaskType); err != nil {
				log.Error("watcherAndCloseTasks closeJob", "err", err, "jobId", assembleTaskId(tk.ArId, tk.TaskType))
				continue
			}
		}
	}
}

func (s *Bungo) watchEverReceiptTxs() {
	startCursor, err := s.database.GetLastEverRawId()
	if err != nil {
		panic(err)
	}
	subTx := s.everpaySdk.Cli.SubscribeTxs(sdkSchema.FilterQuery{
		StartCursor: int64(startCursor),
		Address:     s.bundler.Signer.Address,
		Action:      paySchema.TxActionTransfer,
	})
	defer subTx.Unsubscribe()

	for {
		select {
		case tt := <-subTx.Subscribe():
			if tt.To != s.bundler.Signer.Address {
				continue
			}
			_, from, err := account.IDCheck(tt.From)
			if err != nil {
				log.Error("account.IDCheck(tt.From)", "err", err, "from", tt.From)
				continue
			}

			res := schema.ReceiptEverTx{
				RawId:    uint64(tt.RawId),
				EverHash: tt.EverHash,
				Nonce:    tt.Nonce,
				Symbol:   tt.TokenSymbol,
				TokenTag: tokUtils.Tag(tt.ChainType, tt.TokenSymbol, tt.TokenID),
				From:     from,
				Amount:   tt.Amount,
				Data:     tt.Data,
				Sig:      tt.Sig,
				Status:   schema.UnSpent,
			}

			if err = s.database.InsertReceiptTx(res); err != nil {
				log.Error("s.query.InsertReceiptTx(res)", "err", err)
			}
		}
	}
}

func processPayItems(query *database.Database, itemIds []string, urtx schema.ReceiptEverTx) error {
	// get orders by itemIds
	ordArr, err := getUnPaidOrdersByItemIds(query, itemIds)
	if err != nil {
		log.Error("s.query.GetUnPaidOrder", "err", err, "id", urtx.RawId)
		if err == gorm.ErrRecordNotFound {
			log.Warn("need refund about not find order", "id", urtx.RawId)
			// update receipt status is unrefund and waiting refund
			if err = query.UpdateReceiptStatus(urtx.RawId, schema.UnRefund, nil); err != nil {
				log.Error("s.query.UpdateReceiptStatus2", "err", err, "id", urtx.RawId)
			}
		}
		return err
	}
	// check currency, orders currency must == paymentTxSymbol
	if err = checkOrdersCurrency(ordArr, urtx.Symbol); err != nil {
		log.Error("checkOrdersCurrency(ordArr, urtx.Symbol)", "err", err, "urtx", urtx.EverHash)
		if err = query.UpdateReceiptStatus(urtx.RawId, schema.UnRefund, nil); err != nil {
			log.Error("s.query.UpdateReceiptStatus3", "err", err, "id", urtx.RawId)
		}
		return err
	}

	// check amount
	if err = checkOrdersAmount(ordArr, urtx.Amount); err != nil {
		log.Error("checkOrdersAmount(ordArr, urtx.Amount)", "err", err, "urtx", urtx.EverHash)
		if err = query.UpdateReceiptStatus(urtx.RawId, schema.UnRefund, nil); err != nil {
			log.Error("s.query.UpdateReceiptStatus4", "err", err, "id", urtx.RawId)
		}
		return err
	}
	// update order payment status
	dbTx := query.Db.Begin()
	for _, ord := range ordArr {
		if err = query.UpdateOrderPay(ord.ID, urtx.EverHash, schema.SuccPayment, dbTx); err != nil {
			log.Error("s.query.UpdateOrderPay(ord.ID,schema.SuccPayment,dbTx)", "err", err)
			dbTx.Rollback()
			break
		}
	}

	if err = query.UpdateReceiptStatus(urtx.RawId, schema.Spent, dbTx); err != nil {
		log.Error("s.query.UpdateReceiptStatus(urtx.ID,schema.Spent,dbTx)", "err", err)
		dbTx.Rollback()
		return err
	}
	dbTx.Commit()

	return nil
}

func (s *Bungo) mergeReceiptEverTxs() {
	unspentRpts, err := s.database.GetReceiptsByStatus(schema.UnSpent)
	if err != nil {
		log.Error("s.query.GetUnSpentReceipts()", "err", err)
		return
	}
	for _, urtx := range unspentRpts {
		action, itemIds, err := parseTxData(urtx.Data)
		if err != nil {
			log.Error("parseItemIds(urtx.Data)", "err", err, "urtx", urtx.EverHash)
			if err = s.database.UpdateReceiptStatus(urtx.RawId, schema.UnRefund, nil); err != nil {
				log.Error("s.query.UpdateReceiptStatus1", "err", err, "id", urtx.RawId)
			}
			continue
		}

		switch action {
		case ItemPaymentAction:
			if err := processPayItems(s.database, itemIds, urtx); err != nil {
				log.Error("processPayItemOrder", "err", err)
				continue
			}
		default:
			log.Error(fmt.Sprintf("not support the action: %s", action))
			if err = s.database.UpdateReceiptStatus(urtx.RawId, schema.UnRefund, nil); err != nil {
				log.Error("s.query.UpdateReceiptStatus7", "err", err, "id", urtx.RawId)
			}
			continue
		}
	}
}

func ecrecoverPubkey(everHash, everSig string) (pubkey string, err error) {
	signature := common.FromHex(everSig)
	sig := make([]byte, len(signature))
	copy(sig, signature)
	if len(sig) != 65 {
		err = fmt.Errorf("invalid length of signture: %d", len(sig))
		return
	}

	if sig[64] != 27 && sig[64] != 28 && sig[64] != 1 && sig[64] != 0 {
		err = fmt.Errorf("invalid signature type")
		return
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	recoverPub, err := crypto.Ecrecover(common.FromHex(everHash), sig)
	if err != nil {
		err = fmt.Errorf("can not ecrecover: %v", err)
		return
	}

	pubkey = common.Bytes2Hex(recoverPub)
	return
}

func checkOrdersAmount(ordArr []schema.Order, txAmount string) error {
	txAmountInt, ok := new(big.Int).SetString(txAmount, 10)
	if !ok {
		return errors.New("txAmount incorrect")
	}
	totalFee := big.NewInt(0)
	for _, ord := range ordArr {
		fee, ok := new(big.Int).SetString(ord.Fee, 10)
		if !ok {
			return errors.New("order fee incorrect")
		}
		totalFee = new(big.Int).Add(totalFee, fee)
	}
	if txAmountInt.Cmp(totalFee) < 0 {
		return errors.New("payAmount fee not enough")
	}
	return nil
}

func checkOrdersCurrency(ordArr []schema.Order, txSymbol string) error {
	for _, ord := range ordArr {
		if strings.ToUpper(ord.Currency) != strings.ToUpper(txSymbol) {
			return errors.New("currency incorrect")
		}
	}
	return nil
}

func parseTxData(txData string) (action string, itemIds []string, err error) {
	res := gjson.Parse(txData)
	// appName must be bungo
	if res.Get("appName").String() != "bungo" {
		return "", nil, errors.New("txData.appName not be bungo")
	}

	// action
	act := res.Get("action").String()
	switch act {
	case ItemPaymentAction: // bundle item order payment
		itemIds = make([]string, 0)
		for _, it := range res.Get("itemIds").Array() {
			itemIds = append(itemIds, it.String())
		}
		if len(itemIds) == 0 {
			return "", nil, errors.New("itemIds is empty")
		}
		return ItemPaymentAction, itemIds, nil
	case ApikeyPaymentAction:
		return ApikeyPaymentAction, nil, nil
	default:
		return "", nil, errors.New(fmt.Sprintf("not support action: %s", act))
	}
}

func getUnPaidOrdersByItemIds(query *database.Database, itemIds []string) ([]schema.Order, error) {
	ordArr := make([]schema.Order, 0, len(itemIds))
	for _, itemId := range itemIds {
		ord, err := query.GetUnPaidOrder(itemId)
		if err != nil {
			log.Error("s.query.GetUnPaidOrder(itemId)", "err", err, "itemId", itemId)
			return nil, err
		}
		ordArr = append(ordArr, ord)
	}
	return ordArr, nil
}

func (s *Bungo) collectFee() {
	collectAddr := s.config.FeeCollectAddress()
	if collectAddr == "" {
		log.Warn("s.config.FeeCollectAddress()", "collectAddr", "null")
		return
	}

	// check bundler address token balance
	tokBals, err := s.everpaySdk.Cli.Balances(s.bundler.Signer.Address)
	if err != nil {
		log.Error("s.everpaySdk.Cli.Balances(s.bundler.Signer.Address)", "err", err, "bundler", s.bundler.Signer.Address)
		return
	}

	for _, tokBal := range tokBals.Balances {
		amt, ok := new(big.Int).SetString(tokBal.Amount, 10)
		if !ok {
			continue
		}
		if amt.Cmp(big.NewInt(0)) <= 0 {
			continue
		}
		mmap := map[string]string{
			"appName": "bungo",
			"action":  "feeCollection",
			"bundler": s.bundler.Signer.Address,
		}
		data, _ := json.Marshal(mmap)
		_, err = s.everpaySdk.Transfer(tokBal.Tag, amt, collectAddr, string(data))
		if err != nil {
			log.Error("s.everpaySdk.Transfer(tokBal.Tag,amt,collectAddr,\"\")", "err", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *Bungo) refundReceipt() {
	recpts, err := s.database.GetReceiptsByStatus(schema.UnRefund)
	if err != nil {
		log.Error("s.query.GetReceiptsByStatus(schema.UnRefund)", "err", err)
		return
	}

	for _, rpt := range recpts {
		// update rpt status is refund
		if err := s.database.UpdateReceiptStatus(rpt.RawId, schema.Refund, nil); err != nil {
			log.Error("s.query.UpdateReceiptStatus(rpt.ID,schema.Refund,nil)", "err", err, "id", rpt.RawId)
			continue
		}
		// send everTx transfer for refund
		amount, ok := new(big.Int).SetString(rpt.Amount, 10)
		if !ok {
			log.Error("new(big.Int).SetString(rpt.Amount,10) failed", "amt", rpt.Amount)
			continue
		}
		// everTx data
		mmap := map[string]string{
			"appName":        "bungo",
			"action":         "refund",
			"refundEverHash": rpt.EverHash,
		}
		data, _ := json.Marshal(mmap)
		everTx, err := s.everpaySdk.Transfer(rpt.TokenTag, amount, rpt.From, string(data))
		if err != nil { // notice: if refund failed, then need manual check and refund
			log.Error("s.everpaySdk.Transfer", "err", err)
			// update receipt status is unrefund
			if err := s.database.UpdateRefundErr(rpt.RawId, err.Error()); err != nil {
				log.Error("s.query.UpdateRefundErr(rpt.RawId, err.Error())", "err", err, "id", rpt.RawId)
			}
			continue
		}
		log.Info("refund receipt success...", "receipt everHash", rpt.EverHash, "refund everHash", everTx.HexHash())
	}
}

func (s *Bungo) onChainBundleItems() {
	ords, err := s.database.GetNeedOnChainOrders()
	if err != nil {
		log.Error("s.query.GetNeedOnChainOrders()", "err", err)
		return
	}
	if len(ords) == 0 {
		return
	}
	// send arTx to arweave
	arTx, onChainItemIds, err := s.onChainOrds(ords)
	if err != nil {
		log.Error("s.onChainOrds()", "err", err)
		return
	}

	s.updateOnChainInfo(onChainItemIds, arTx, schema.PendingOnChain)

}

func (s *Bungo) onChainItemsBySeq() {
	ords, err := s.database.GetNeedOnChainOrdersSorted()
	if err != nil {
		log.Error("s.query.GetNeedOnChainOrders()", "err", err)
		return
	}

	if len(ords) == 0 {
		return
	}
	arTx, onChainItemIds, err := s.onChainOrds(ords)
	if err != nil {
		return
	}

	// block until arTx.confirmation > 3
	status := arTxWatcher(s.arCli, arTx.ID)
	if !status {
		log.Error("watch tx status failed", "ar id", arTx.ID, "status", status)
		return
	}

	s.updateOnChainInfo(onChainItemIds, arTx, schema.SuccOnChain)
}

func (s *Bungo) onChainOrds(ords []schema.Order) (arTx types.Transaction, onChainItemIds []string, err error) {
	// once total size limit 2 GB
	itemIds := make([]string, 0, len(ords))
	totalSize := int64(0)
	for _, ord := range ords {
		if totalSize+ord.Size > schema.MaxPerOnChainSize {
			continue
		}
		od, exist := s.database.ExistProcessedOrderItem(ord.ItemId)
		if exist {
			if err = s.database.UpdateOrdOnChainStatus(od.ItemId, od.Status, nil); err != nil {
				log.Error("s.query.UpdateOrdOnChainStatus(od.ItemId,od.OnChainStatus)", "err", err, "itemId", od.ItemId)
			}
			continue
		}
		itemIds = append(itemIds, ord.ItemId)
		totalSize += od.Size
	}

	// send arTx to arweave
	return s.onChainBundleTx(itemIds)
}

func (s *Bungo) updateOnChainInfo(onChainItemIds []string, arTx types.Transaction, onChainStatus string) {
	// insert arTx record
	onChainItemIdsJs, err := json.Marshal(onChainItemIds)
	if err != nil {
		log.Error("json.Marshal(itemIds)", "err", err, "onChainItemIds", onChainItemIds)
		return
	}
	if err = s.database.InsertArTx(schema.OnChainTx{
		ArId:      arTx.ID,
		CurHeight: s.cache.GetInfo().Height,
		DataSize:  arTx.DataSize,
		Reward:    arTx.Reward,
		Status:    schema.PendingOnChain,
		ItemIds:   onChainItemIdsJs,
		ItemNum:   len(onChainItemIds),
	}); err != nil {
		log.Error("s.query.InsertArTx", "err", err)
		return
	}

	// update order onChainStatus to pending
	for _, itemId := range onChainItemIds {
		if err = s.database.UpdateOrdOnChainStatus(itemId, onChainStatus, nil); err != nil {
			log.Error("s.query.UpdateOrdOnChainStatus(item.Id,schema.PendingOnChain)", "err", err, "itemId", itemId)
		}
	}
}

func (s *Bungo) watchArTx() {
	txs, err := s.database.GetArTxByStatus(schema.PendingOnChain)
	if err != nil {
		log.Error("s.query.GetArTxByStatus(schema.PendingOnChain)", "err", err)
		return
	}

	for _, tx := range txs {
		// check onchain status
		arTxStatus, err := s.arCli.GetTransactionStatus(tx.ArId)
		if err != nil {
			if err != goar.ErrPendingTx && s.cache.GetInfo().Height-tx.CurHeight > 50 {
				// arTx has expired
				if err = s.database.UpdateArTxStatus(tx.ArId, schema.FailedOnChain, nil, nil); err != nil {
					log.Error("UpdateArTxStatus(tx.ArId,schema.FailedOnChain)", "err", err)
				}
			}
			continue
		}

		// update status success
		if arTxStatus.NumberOfConfirmations > 3 {
			dbTx := s.database.Db.Begin()
			if err = s.database.UpdateArTxStatus(tx.ArId, schema.SuccOnChain, arTxStatus, dbTx); err != nil {
				log.Error("UpdateArTxStatus(tx.ArId,schema.SuccOnChain)", "err", err)
				dbTx.Rollback()
				continue
			}

			// update order onchain status
			bundleItemIds := make([]string, 0)
			if err = json.Unmarshal(tx.ItemIds, &bundleItemIds); err != nil {
				log.Error("json.Unmarshal(tx.ItemIds,&bundleItemIds)", "err", err, "itemsJs", tx.ItemIds)
				dbTx.Rollback()
				continue
			}

			for _, itemId := range bundleItemIds {
				if err = s.database.UpdateOrdOnChainStatus(itemId, schema.SuccOnChain, dbTx); err != nil {
					log.Error("s.query.UpdateOrdOnChainStatus(itemId,schema.SuccOnChain,dbTx)", "err", err, "item", itemId)
					dbTx.Rollback()
					continue
				}
			}
			dbTx.Commit()
		}
	}
}

func (s *Bungo) retryOnChainArTx() {
	txs, err := s.database.GetArTxByStatus(schema.FailedOnChain)
	if err != nil {
		log.Error("s.query.GetArTxByStatus(schema.PendingOnChain)", "err", err)
		return
	}
	if len(txs) == 0 {
		return
	}
	for _, tx := range txs {
		itemIds := make([]string, 0)
		if err = json.Unmarshal(tx.ItemIds, &itemIds); err != nil {
			return
		}
		arTx, _, err := s.onChainBundleTx(itemIds)
		if err != nil {
			return
		}
		// update onChain
		if err = s.database.UpdateArTx(tx.ID, arTx.ID, s.cache.GetInfo().Height, arTx.DataSize, arTx.Reward, schema.PendingOnChain); err != nil {
			log.Error("s.query.UpdateArTx", "err", err, "id", tx.ID, "arId", arTx.ID)
		}
	}
}

func (s *Bungo) onChainBundleTx(itemIds []string) (arTx types.Transaction, onChainItemIds []string, err error) {
	onChainItems := make([]types.BundleItem, 0)
	bundle := &types.Bundle{}

	onChainItems, err = s.getOnChainBundle(itemIds)
	// assemble and send to arweave
	bundle, err = utils.NewBundle(onChainItems...)
	if err != nil {
		log.Error("utils.NewBundle(onChainItems...)", "err", err)
		return
	}

	// verify bundle, ensure that the bundle is exactly right before sending
	if _, err = utils.DecodeBundle(bundle.BundleBinary); err != nil {
		err = errors.New(fmt.Sprintf("Verify bundle failed; err:%v", err))
		return
	}

	// get onChainItemIds
	for _, item := range onChainItems {
		onChainItemIds = append(onChainItemIds, item.Id)
	}

	arTxtags := []types.Tag{
		{Name: "App-Name", Value: "bungo"},
		{Name: "App-Version", Value: "1.0.0"},
		{Name: "Action", Value: "Bundle"},
		{Name: "Protocol-Name", Value: "U"},
		{Name: "Action", Value: "Burn"},
		{Name: "App-Name", Value: "SmartWeaveAction"},
		{Name: "App-Version", Value: "0.3.0"},
		{Name: "Input", Value: `{"function":"mint"}`},
		{Name: "Contract", Value: "KTzTXT_ANmF84fWEKHzWURD1LWd9QaFR9yfYUwH2Lxw"},
	}

	// speed arTx Fee

	concurrentNum := s.config.Param.ChunkConcurrentNum
	if len(bundle.BundleBinary) > 0 {
		log.Debug("use binary submit bundle arTx", "binary length:", len(bundle.BundleBinary))
		price := calculatePrice(s.cache.GetFee(), int64(len(bundle.BundleBinary)))
		speedFactor := calculateFactor(price, s.config.GetSpeedFee())
		arTx, err = s.bundler.SendBundleTxSpeedUp(context.TODO(), concurrentNum, bundle.BundleBinary, arTxtags, speedFactor)
	} else {
		fileInfo, err1 := bundle.BundleDataReader.Stat()
		if err1 != nil {
			err = err1
			return
		}
		if fileInfo.Size() == 0 {
			err = errors.New("bundle.BundleDataReader is null")
			return
		}
		price := calculatePrice(s.cache.GetFee(), fileInfo.Size())
		speedFactor := calculateFactor(price, s.config.GetSpeedFee())
		arTx, err = s.bundler.SendBundleTxSpeedUp(context.TODO(), concurrentNum, bundle.BundleDataReader, arTxtags, speedFactor)
	}
	if err != nil {
		log.Error("s.bundler.SendBundleTxSpeedUp(bundle.BundleBinary,arTxtags)", "err", err)
		return
	}
	log.Info("Send bundle arTx", "arTx", arTx.ID)

	// bungo broadcast tx data
	if err := s.arseedCli.SubmitTxConcurrent(context.TODO(), concurrentNum, arTx); err != nil {
		log.Error("s.arseedCli.SubmitTxConcurrent(arTx)", "err", err, "arId", arTx.ID)
	}
	return
}

func (s *Bungo) getOnChainBundle(itemIds []string) (onChainItems []types.BundleItem, err error) {
	onChainItems = make([]types.BundleItem, 0, len(itemIds))
	for _, itemId := range itemIds {
		_, itemBinary, err := s.store.LoadItemBinary(itemId)
		if err != nil {
			log.Error("s.store.LoadItemBinary(itemId)", "err", err, "itemId", itemId)
			continue
		}
		item, err := utils.DecodeBundleItem(itemBinary)
		if err != nil {
			log.Error("utils.DecodeBundleItem(itemBinary)", "err", err, "itemId", itemId)
			continue
		}
		onChainItems = append(onChainItems, *item)
	}
	if len(onChainItems) == 0 {
		err = errors.New("onChainItems is null")
		return
	}

	// the end off item.Data not be "", because when the case viewblock decode failed. // todo viewblock used stream function decode item, so this is a bug for them
	endItem := onChainItems[len(onChainItems)-1]
	if endItem.Data == "" {
		// find a data != "" item and push to end off
		idx := -1
		for i, item := range onChainItems {
			if item.Data != "" {
				idx = i
				break
			}
		}
		if idx == -1 {
			err = errors.New("all bundle items data are null")
			return
		}
		newEndItem := onChainItems[idx]
		onChainItems = append(onChainItems[:idx], onChainItems[idx+1:]...)
		onChainItems = append(onChainItems, newEndItem)
	}
	return
}

func (s *Bungo) getOnChainBundleStream(itemIds []string) (onChainItems []types.BundleItem, err error) {
	onChainItems = make([]types.BundleItem, 0, len(itemIds))
	for _, itemId := range itemIds {
		binaryReader, _, err := s.store.LoadItemBinary(itemId)
		if err != nil {
			log.Error("s.store.LoadItemBinary(itemId)", "err", err, "itemId", itemId)
			continue
		}

		item, err2 := utils.DecodeBundleItemStream(binaryReader)
		if err2 != nil {
			log.Error("utils.DecodeBundleItem(itemBinary)", "err", err2, "itemId", itemId)
			continue
		}
		binaryReader.Close()
		os.Remove(binaryReader.Name())
		onChainItems = append(onChainItems, *item)
	}
	if len(onChainItems) == 0 {
		err = errors.New("onChainItems is null")
		return
	}
	return
	// the end off item.Data not be "", because when the case viewblock decode failed. // todo viewblock used stream function decode item, so this is a bug for them
	// todo
}

func (s *Bungo) processExpiredOrd() {
	ords, err := s.database.GetExpiredOrders()
	if err != nil {
		log.Error("GetExpiredOrders()", "err", err)
		return
	}
	for _, ord := range ords {
		if err = s.database.UpdateOrdToExpiredStatus(ord.ID); err != nil {
			log.Error("UpdateOrdToExpiredStatus", "err", err, "id", ord.ID)
			continue
		}
		// can not delete
		// 1. exist paid order
		if s.database.ExistPaidOrd(ord.ItemId) {
			continue
		}
		// 2. this order not the latest unpaid order
		if !s.database.IsLatestUnpaidOrd(ord.ItemId, ord.PaymentExpiredTime) {
			continue
		}
		// delete bundle item from store
		if err = s.DelItem(ord.ItemId); err != nil {
			log.Error("DelItem", "err", err, "itemId", ord.ItemId)
			continue
		}
		// delete manifest table
		if err = s.database.DelManifest(ord.ItemId); err != nil {
			log.Error("s.query.DelManifest", "err", err, "itemId", ord.ItemId)
			continue
		}
	}
}

func (s *Bungo) parseAndSaveBundleTx() {
	arIds, err := s.store.LoadWaitParseBundleArIds()
	if err != nil {
		if err != schema.ErrNotExist {
			log.Error("s.store.LoadWaitParseBundleArIds()", "err", err)
		}
		return
	}
	for _, arId := range arIds {
		// get tx data
		arTxMeta, err := s.store.LoadTransactionMetadata(arId)
		if err != nil {
			log.Error("s.store.LoadTxMeta(arId)", "err", err, "arId", arId)
			continue
		}

		data, err := getArTxData(arTxMeta.DataRoot, arTxMeta.DataSize, s.store)
		if err != nil {
			// log.Error("get data failed, if is not_exist_record,then wait submit chunks fully", "err", err, "arId", arId)
			continue
		}
		if err := s.ParseAndSaveBundleItems(arId, data); err != nil {
			log.Error("ParseAndSaveBundleItems", "err", err, "arId", arId)
		}
		// del wait db
		if err = s.store.DelParsedBundleArId(arId); err != nil {
			log.Error("DelParsedBundleArId", "err", err, "arId", arId)
		}
	}
}

func (s *Bungo) updateBundler() {
	// update bundler balance
	addr := s.bundler.Signer.Address
	bal, err := s.arCli.GetWalletBalance(addr)
	if err != nil {
		return
	}
	metricBundlerBalance(bal, addr)
}

func (s *Bungo) deleteTmpFile() {
	tmpFileMapLock.Lock()
	defer tmpFileMapLock.Unlock()
	for tmpFileName, cnt := range tmpFileMap {
		if cnt <= 0 {
			delete(tmpFileMap, tmpFileName)
			os.Remove(tmpFileName)
		}
	}
}

func filterPeers(peers []string, constTx *types.Transaction) map[string]bool {
	var mu sync.Mutex
	var wg sync.WaitGroup
	availablePeers := make(map[string]bool, 0)
	p, err := ants.NewPoolWithFunc(50, func(peer interface{}) {
		defer wg.Done()
		pStr := peer.(string)
		pNode := goar.NewTempConn()
		pNode.SetTempConnUrl("http://" + pStr)
		pNode.SetTimeout(time.Second * 10)
		_, code, err := pNode.SubmitTransaction(constTx)
		if err != nil {
			return
		}
		// if the resp like this ,we believe this peer is available
		if code/100 == 2 || code/100 == 4 {
			mu.Lock()
			availablePeers[pStr] = true
			mu.Unlock()
		}
	})
	// submit a legal Tx to peers, if the response is timely and acceptable, then the peer is good for submit tx
	for _, peer := range peers {
		wg.Add(1)
		err = p.Invoke(peer)
		if err != nil {
			log.Warn("concurrency err", "err", err)
		}
	}
	wg.Wait()
	p.Release()
	return availablePeers
}

func updatePeerMap(oldPeerMap map[string]int64, availablePeers map[string]bool) map[string]int64 {
	for peer, cnt := range oldPeerMap {
		if _, ok := availablePeers[peer]; !ok {
			if cnt > 0 {
				oldPeerMap[peer] -= 1
			}
		} else {
			oldPeerMap[peer] += 1
		}
	}

	for peer, _ := range availablePeers {
		if _, ok := oldPeerMap[peer]; !ok {
			oldPeerMap[peer] = 1
		}
	}
	return oldPeerMap
}

func calculateFactor(price, speedFee int64) int64 {
	if speedFee == 0 {
		return 0
	}
	val := (price+speedFee)*100/price - 100
	if val == 0 {
		val = 1
	}
	return val
}

func arTxWatcher(arCli *goar.Client, arTxHash string) bool {
	// loop watcher on chain tx status
	// total time 59 minute
	tmp := 0
	for i := 1; i <= 21; i++ {
		// sleep
		num := 60 + i*10
		time.Sleep(time.Second * time.Duration(num))

		tmp += num
		log.Debug("watcher tx sleep time", "wait total time(s)", tmp)
		log.Debug("retry get tx status", "txHash", arTxHash)

		// watcher on-chain tx
		status, err := arCli.GetTransactionStatus(arTxHash)
		if err != nil {
			if err.Error() == goar.ErrPendingTx.Error() {
				log.Debug("tx is pending", "txHash", arTxHash)
			} else {
				log.Error("get tx status", "err", err, "txHash", arTxHash)
			}
			continue
		}

		// when err is nil
		// confirms block height must >= 3
		if status.NumberOfConfirmations < 3 {
			log.Debug("bungo send sequence tx must more than 2 block confirms", "txHash", arTxHash, "currentConfirms", status.NumberOfConfirmations)
			continue
		} else {
			return true
		}
	}
	return false
}

func (s *Bungo) UpdateRealTime() {
	data, err := s.database.GetOrderRealTimeStatistic()
	if err != nil {
		log.Error("s.query.GetOrderRealTimeStatistic()", "err", err)
		return
	}
	if err := s.store.UpdateRealTimeStatistic(data); err != nil {
		log.Error("s.store.KVDb.Put()", "err", err)
	}
}

func (s *Bungo) ProduceDailyStatistic() {
	now := time.Now()
	var start time.Time
	var firstOrder schema.Order
	var osc schema.OrderStatistic
	err := s.database.Db.Model(&schema.Order{}).First(&firstOrder).Error
	//Not found
	if err != nil {
		return
	}
	err = s.database.Db.Model(&schema.OrderStatistic{}).Last(&osc).Error
	if err == nil {
		start = osc.Date.Add(24 * time.Hour)
	} else {
		start = time.Date(firstOrder.CreatedAt.Year(), firstOrder.CreatedAt.Month(), firstOrder.CreatedAt.Day(), 0, 0, 0, 0, now.Location())
	}
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	//If yesterday's record already exists, return
	if !s.database.WhetherExec(schema.TimeRange{Start: end.Add(-24 * time.Hour), End: end}) {
		return
	}
	for start.Before(end) {
		results, err := s.database.GetDailyStatisticByDate(schema.TimeRange{Start: start, End: start.Add(24 * time.Hour)})
		if err != nil {
			log.Error("s.ProduceDailyStatistic()", "err", err)
			start = start.Add(24 * time.Hour)
			continue
		}
		if len(results) == 0 {
			s.database.Db.Model(&schema.OrderStatistic{}).Create(&schema.OrderStatistic{
				Date: start,
			})
		} else {
			s.database.Db.Model(&schema.OrderStatistic{}).Create(&schema.OrderStatistic{Date: start, Totals: results[0].Totals, TotalDataSize: results[0].TotalDataSize})
		}
		start = start.Add(24 * time.Hour)
	}
}

func (s *Bungo) broadcastItemToKafka() {
	kafkaOrdInfos, err := s.database.GetKafkaOrderInfos()
	if err != nil {
		log.Error("s.query.GetKafkaOrderInfos()", "err", err)
		return
	}
	if len(kafkaOrdInfos) == 0 {
		return
	}

	itemTopicKw, ok := s.KWriters[ItemTopic]
	if !ok {
		log.Error("s.KWriters[ItemTopic]", "err", "not exist")
		return
	}
	for _, ord := range kafkaOrdInfos {
		// post to kafka
		meta, err := s.store.LoadItemMeta(ord.ItemId)
		if err != nil {
			log.Error("s.store.LoadItemMeta(ord.ItemId)", "err", err, "itemId", ord.ItemId)
			continue
		}
		// 0x hex to base64
		addr := ord.Signer
		if meta.SignatureType != types.ArweaveSignType {
			addr, err = Base64Address(meta.Owner)
			if err != nil {
				log.Error("Base64Address(meta.Owner)", "err", err, "signer", ord.Signer, "owner", meta.Owner)
				continue
			}
		}

		// data content type
		contentType := getTagValue(meta.Tags, schema.ContentType)

		kItem := schema.KafkaBundleItem{
			SignatureType: meta.SignatureType,
			Signature:     meta.Signature,
			Owner:         meta.Owner,
			Target:        meta.Target,
			Anchor:        meta.Anchor,
			Tags:          meta.Tags,
			Id:            meta.Id,
			Size:          ord.Size,
			Address:       addr,
			Type:          contentType,
		}
		itemBy, err := json.Marshal(kItem)
		if err != nil {
			log.Error("json.Marshal(kItem)", "err", err)
			continue
		}
		if err = itemTopicKw.Write(itemBy); err != nil {
			log.Error("itemTopicKw.Write(itemBy)", "err", err)
			continue
		}

		if err = s.database.KafkaDone(ord.ID); err != nil {
			log.Error("s.query.KafkaDone(ord.ID)", "err", err, "id", ord.ID)
			continue
		}
	}
}

func (s *Bungo) broadcastBlockToKafka() {
	kafkaOnchains, err := s.database.GetKafkaOnChains()
	if err != nil {
		log.Error("s.query.GetKafkaOnChains()", "err", err)
		return
	}

	blockkw, ok := s.KWriters[BlockTopic]
	if !ok {
		log.Error("s.KWriters[BlockTopic]", "err", "not exist")
		return
	}
	for _, onchain := range kafkaOnchains {
		// get block timestamp and previous id
		block, err := s.arCli.GetBlockByID(onchain.BlockId)
		if err != nil {
			log.Error("s.arCli.GetBlockByID(onchain.BlockId)", "err", err, "onchain.BlockId", onchain.BlockId)
			continue
		}
		// post to kafka
		itemIds := make([]string, 0)
		if err := json.Unmarshal(onchain.ItemIds, &itemIds); err != nil {
			log.Error("json.Unmarshal(onchain.ItemIds,&itemIds)", "err", err)
			continue
		}
		onchainInfo := schema.KafkaOnChainInfo{
			BundleIn:  onchain.ArId,
			ItemIds:   itemIds,
			Id:        onchain.BlockId,
			Height:    onchain.BlockHeight,
			Timestamp: block.Timestamp,
			Previous:  block.PreviousBlock,
		}

		body, err := json.Marshal(onchainInfo)
		if err != nil {
			log.Error("json.Marshal(onchainInfo)", "err", err)
			continue
		}
		if err = blockkw.Write(body); err != nil {
			log.Error("blockkw.Write(body)", "err", err)
			continue
		}

		if err = s.database.KafkaOnChainDone(onchain.ID); err != nil {
			log.Error("s.query.KafkaOnChainDone(onchain.ID)", "err", err, "id", onchain.ID)
			continue
		}
	}
}

func Base64Address(pubkey string) (string, error) {
	bby, err := utils.Base64Decode(pubkey)
	if err != nil {
		return "", err
	}
	aa := sha256.Sum256(bby)
	return utils.Base64Encode(aa[:]), nil
}
