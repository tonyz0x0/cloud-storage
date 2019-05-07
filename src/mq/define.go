package mq

import (
	cmn "cloud-storage/src/common"
)

type TransferData struct {
	FileHash      string
	CurLocation   string // tmp local location
	DestLocation  string
	DestStoreType cmn.StoreType
}
