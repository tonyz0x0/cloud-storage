package common

type StoreType int

const (
	_ StoreType = iota
	StoreLocal
	// StoreCeph : Ceph Cluster
	StoreCeph
	// StoreOSS : AliCloud OSS
	StoreOSS
	// StoreMix : Mixed(Ceph and Ali-OSS)
	StoreMix
	// StoreAll : Store in all sources
	StoreAll
)
