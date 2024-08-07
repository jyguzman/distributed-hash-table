package kademlia

type KadOptions struct {
	BucketCapacity int
	Alpha          int
	TRefresh       int
	TExpiration    int
	MaxIterations  int
}

var Options = KadOptions{
	BucketCapacity: 3,
	Alpha:          3,
	TRefresh:       60 * 60,
	TExpiration:    60 * 60,
	MaxIterations:  20,
}
