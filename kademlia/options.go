package kademlia

import "fmt"

type KadOptions struct {
	Protocol       string
	BucketCapacity int
}

var Options = KadOptions{
	Protocol:       "udp",
	BucketCapacity: 8,
}

func SetOptions(protocol string, k int) error {
	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf("network must be \"tcp\" or \"udp\"")
	}
	Options.Protocol = protocol
	Options.BucketCapacity = k
	return nil
}
