package kademlia

type KadOptions struct {
	BucketCapacity int
	Alpha          int
	IdLength       uint
	NodeRefresh    int
	NodeExpiration int
}

var Options = KadOptions{
	BucketCapacity: 8,
	Alpha:          3,
	IdLength:       160,
	NodeRefresh:    60 * 60,
	NodeExpiration: 60 * 60,
}
