package store

import (
	"encoding/json"
	"log"
	"os"

	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/liteseed/bungo/schema"
)

type IStore interface {
	Put(bucket, key string, value interface{}) (err error)

	Get(bucket, key string) (data []byte, err error)

	GetStream(bucket, key string) (data *os.File, err error)

	GetAllKey(bucket string) (keys []string, err error)

	Delete(bucket, key string) (err error)

	Close() (err error)

	Type() string

	Exist(bucket, key string) bool
}

type Store struct {
	KVDb IStore
}

func (s *Store) AtomicSync(oldOffset, newOffset uint64, root, size string) error {

	if err := s.SaveAllDataEndOffset(newOffset); err != nil {
		log.Fatal("s.store.SaveAllDataEndOffset(newEndOffset)", "err", err)
		return err
	}
	
	if err := s.SaveTransactionOffSet(root, size, newOffset); err != nil {
		_ = s.RollbackAllDataEndOffset(oldOffset)
		return err
	}
	return nil
}

func (s *Store) Close() error {
	return s.KVDb.Close()
}

func (s *Store) LoadAllDataEndOffset() (offset uint64) {
	key := "allDataEndOffset"
	data, err := s.KVDb.Get(schema.Constants, key)
	if err != nil || data == nil {
		offset = 0
		return
	}
	offset = base64ToInt(string(data))
	return
}

func (s *Store) LoadChunk(chunkStartOffset uint64) (chunk *types.GetChunk, err error) {
	chunk = &types.GetChunk{}
	data, err := s.KVDb.Get(schema.Chunks, intToBase64(chunkStartOffset))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, chunk)
	return
}

func (s *Store) LoadTransactionMetadata(id string) (transaction *types.Transaction, err error) {
	transaction = &types.Transaction{}
	data, err := s.KVDb.Get(schema.TransactionMetadata, id)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, transaction)
	return
}

func (s *Store) RollbackAllDataEndOffset(preDataEndOffset uint64) (err error) {
	key := "allDataEndOffset"
	val := []byte(intToBase64(preDataEndOffset))

	return s.KVDb.Put(schema.Constants, key, val)
}

func (s *Store) SaveAllDataEndOffset(allDataEndOffset uint64) (err error) {
	key := "allDataEndOffset"
	val := []byte(intToBase64(allDataEndOffset))

	return s.KVDb.Put(schema.Constants, key, val)
}

func (s *Store) SaveTransactionMetadata(transaction types.Transaction) error {
	transaction.Data = "" // only store transaction metadata, exclude data
	key := transaction.ID
	val, err := json.Marshal(&transaction)
	if err != nil {
		return err
	}
	return s.KVDb.Put(schema.TransactionMetadata, key, val)
}

func (s *Store) DoesMetadataExist(id string) bool {
	_, err := s.LoadTransactionMetadata(id)
	return err != schema.ErrNotExist
}

func (s *Store) SaveTransactionOffSet(root, size string, offset uint64) (err error) {
	return s.KVDb.Put(schema.TransactionOffset, generateOffsetKey(root, size), []byte(intToBase64(offset)))
}

func (s *Store) LoadTxDataEndOffSet(dataRoot, dataSize string) (txDataEndOffset uint64, err error) {
	data, err := s.KVDb.Get(schema.TransactionOffset, generateOffsetKey(dataRoot, dataSize))
	if err != nil {
		return
	}
	txDataEndOffset = base64ToInt(string(data))
	return
}

func (s *Store) IsExistTxDataEndOffset(dataRoot, dataSize string) bool {
	_, err := s.LoadTxDataEndOffSet(dataRoot, dataSize)
	return err != schema.ErrNotExist
}

func (s *Store) SaveChunk(chunkStartOffset uint64, chunk types.GetChunk) error {
	chunkJs, err := chunk.Marshal()
	if err != nil {
		return err
	}
	err = s.KVDb.Put(schema.Chunks, intToBase64(chunkStartOffset), chunkJs)

	return err
}

func (s *Store) IsExistChunk(chunkStartOffset uint64) bool {
	_, err := s.LoadChunk(chunkStartOffset)
	return err != schema.ErrNotExist
}

func (s *Store) SavePeers(peers map[string]int64) error {
	peersB, err := json.Marshal(peers)
	key := "peer-list"
	if err != nil {
		return err
	}
	return s.KVDb.Put(schema.Constants, key, peersB)
}

func (s *Store) LoadPeers() (peers map[string]int64, err error) {
	key := "peer-list"
	peers = make(map[string]int64, 0)
	data, err := s.KVDb.Get(schema.Constants, key)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &peers)
	return
}

func (s *Store) IsExistPeers() bool {
	_, err := s.LoadPeers()
	return err != schema.ErrNotExist
}

// about tasks

func (s *Store) PutTaskPendingPool(taskId string) error {
	return s.KVDb.Put(schema.PendingTasks, taskId, []byte("0x01"))
}

func (s *Store) LoadAllPendingTaskIds() ([]string, error) {
	taskIds, err := s.KVDb.GetAllKey(schema.PendingTasks)
	if err != nil {
		if err == schema.ErrNotExist {
			return nil, nil
		}
		return nil, err
	}
	return taskIds, err
}

func (s *Store) DelPendingPoolTaskId(taskId string) error {
	return s.KVDb.Delete(schema.PendingTasks, taskId)
}

func (s *Store) SaveTask(taskId string, tk schema.Task) error {
	val, err := json.Marshal(&tk)
	if err != nil {
		return err
	}
	return s.KVDb.Put(schema.Tasks, taskId, val)
}

func (s *Store) LoadTask(taskId string) (tk *schema.Task, err error) {
	tk = &schema.Task{}
	data, err := s.KVDb.Get(schema.Tasks, taskId)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, tk)
	return
}

// about bundle
func (s *Store) AtomicSaveItem(item types.BundleItem) (err error) {
	if err = s.SaveItemBinary(item); err != nil {
		return
	}
	if err = s.SaveItemMeta(item); err != nil {
		_ = s.DelItemBinary(item.Id)
	}
	return
}

func (s *Store) AtomicDelItem(itemId string) (err error) {
	err = s.DelItemMeta(itemId)
	if err != nil {
		return
	}
	return s.DelItemBinary(itemId)
}

func (s *Store) SaveItemBinary(item types.BundleItem) (err error) {
	if item.DataReader != nil {
		binaryReader, err := utils.GenerateItemBinaryStream(&item)
		if err != nil {
			return err
		}
		return s.KVDb.Put(schema.BundleItemBinary, item.Id, binaryReader)
	} else {
		return s.KVDb.Put(schema.BundleItemBinary, item.Id, item.ItemBinary)
	}
}

func (s *Store) LoadItemBinary(itemId string) (binaryReader *os.File, itemBinary []byte, err error) {
	itemBinary, err = s.KVDb.Get(schema.BundleItemBinary, itemId)
	return
}

func (s *Store) IsExistItemBinary(itemId string) bool {
	return s.KVDb.Exist(schema.BundleItemBinary, itemId)
}

func (s *Store) DelItemBinary(itemId string) (err error) {
	return s.KVDb.Delete(schema.BundleItemBinary, itemId)
}

func (s *Store) SaveItemMeta(item types.BundleItem) (err error) {
	item.Data = "" // without data
	meta, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return s.KVDb.Put(schema.BundleItemMeta, item.Id, meta)
}

func (s *Store) LoadItemMeta(itemId string) (meta types.BundleItem, err error) {
	meta = types.BundleItem{}
	data, err := s.KVDb.Get(schema.BundleItemMeta, itemId)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &meta)
	return
}

func (s *Store) DelItemMeta(itemId string) (err error) {
	return s.KVDb.Delete(schema.BundleItemMeta, itemId)
}

// bundle items to arTx

func (s *Store) SaveWaitParseBundleArId(arId string) error {
	return s.KVDb.Put(schema.BundleWaitParseArIdBucket, arId, []byte("0x01"))
}

func (s *Store) LoadWaitParseBundleArIds() (arIds []string, err error) {
	arIds, err = s.KVDb.GetAllKey(schema.BundleWaitParseArIdBucket)
	return
}

func (s *Store) DelParsedBundleArId(arId string) error {
	return s.KVDb.Delete(schema.BundleWaitParseArIdBucket, arId)
}

func (s *Store) SaveArIdToItemIds(arId string, itemIds []string) error {
	itemIdsJs, err := json.Marshal(itemIds)
	if err != nil {
		return err
	}
	return s.KVDb.Put(schema.BundleArIdToItemIdsBucket, arId, itemIdsJs)
}

func (s *Store) LoadArIdToItemIds(arId string) (itemIds []string, err error) {
	itemIds = make([]string, 0)
	data, err := s.KVDb.Get(schema.BundleArIdToItemIdsBucket, arId)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &itemIds)
	return
}

func (s *Store) ExistArIdToItemIds(arId string) bool {
	_, err := s.LoadArIdToItemIds(arId)
	return err != schema.ErrNotExist
}

func (s *Store) UpdateRealTimeStatistic(data []byte) error {
	key := "RealTimeOrderStatistic"
	return s.KVDb.Put(schema.OrderStatistics, key, data)
}

func (s *Store) GetRealTimeStatistic() ([]byte, error) {
	key := "RealTimeOrderStatistic"
	return s.KVDb.Get(schema.OrderStatistics, key)
}
