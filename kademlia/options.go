package kademlia

type KadOptions struct {
	BucketCapacity int
	Alpha          int
	NodeRefresh    int
	NodeExpiration int
}

var Options = KadOptions{
	BucketCapacity: 3,
	Alpha:          3,
	NodeRefresh:    60 * 60,
	NodeExpiration: 60 * 60,
}
