package kademlia

type KadOptions struct {
	BucketCapacity int
	Alpha          int
	NodeRefresh    int
	NodeExpiration int
	MaxIterations  int
}

var Options = KadOptions{
	BucketCapacity: 20,
	Alpha:          3,
	NodeRefresh:    60 * 60,
	NodeExpiration: 60 * 60,
	MaxIterations:  20,
}
