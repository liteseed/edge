package schema

var (
	// bucket
	Chunks              = "chunks"               // key: chunkStartOffset, val: chunks
	TransactionOffset   = "transaction_offset"   // key: dataRoot+dataSize; val: transactionOffset
	TransactionMetadata = "transaction_metadata" // key: id (Transaction ID), val: arTx; not include data
	Constants           = "constants"

	// tasks
	PendingTasks = "pending_tasks" // key: taskId(taskType+"_"+arId), val: "0x01"
	Tasks        = "tasks"         // key: taskId(taskType+"_"+arId), val: task

	// bundle bucketName
	BundleItemBinary = "bundle_item_binary"
	BundleItemMeta   = "bundle_item_metadata"

	// parse arTx data to bundle items
	BundleWaitParseArIdBucket = "bundle_wait_parse_arId_bucket" // key: arId, val: "0x01"
	BundleArIdToItemIdsBucket = "bundle_arId_to_itemIds_bucket" // key: arId, val: json.marshal(itemIds)

	//statistic
	OrderStatistics = "order_statistics"
)
