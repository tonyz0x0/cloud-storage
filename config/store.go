package config

import (
	cmn "cloud-storage/src/common"
)

const (
	// TempLocalRootDir: Temporal Local Storage Root
	TempLocalRootDir = "/data/fileserver/"
	// CurrentStoreType: Current Store Type, Ceph or OSS
	CurrentStoreType = cmn.StoreOSS
)
