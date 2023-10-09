package server

import (
	"crypto/md5"
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/buraksezer/consistent"
)

type ConsistentHashMember string

// String this method comply with consistent.Member interface
func (m ConsistentHashMember) String() string {
	return string(m)
}

type ConsistentHasher struct {
}

// Sum64 hash data bytes and return sum
func (ch ConsistentHasher) Sum64(data []byte) uint64 {
	// Use md5 to hash the data
	// Use math/big and hex to process and return the hash as uint64
	bi := big.NewInt(0)
	h := md5.New()
	h.Write(data)
	hexstr := hex.EncodeToString(h.Sum(nil))
	bi.SetString(hexstr, 16)

	result := bi.Int64()
	if result < 0 {
		result *= -1
	}
	return uint64(result)
}

// ServerCache define cache to be used for various data associated with the server
type ServerCache[T ICacheable] struct {
	hashRing         map[string]*ServerCachePartition[T]
	consistentHasher *consistent.Consistent
}

// NewServerCache return new server cache instance
func NewServerCache[T ICacheable](numberOfCachePartitions int) *ServerCache[T] {
	hashRing := make(map[string]*ServerCachePartition[T])
	consistentConfig := consistent.Config{
		Hasher:            &ConsistentHasher{},
		PartitionCount:    numberOfCachePartitions,
		Load:              1.25,
		ReplicationFactor: 20, // Virtual nodes
	}
	consistentHasher := consistent.New(nil, consistentConfig)

	// Initialize the map that contained the consistent hashing ring
	for i := 0; i < numberOfCachePartitions; i++ {
		partitionID := strconv.Itoa(i)
		consistentHasher.Add(ConsistentHashMember(partitionID))
		hashRing[partitionID] = NewServerCachePartition[T](partitionID)
	}

	return &ServerCache[T]{
		hashRing:         hashRing,
		consistentHasher: consistentHasher,
	}
}

// Get get cache value by key
func (sc *ServerCache[T]) Get(key string) *T {
	// Find the partitionID associate with the map we need to look for the key
	partitionID := sc.consistentHasher.LocateKey([]byte(key))
	if partitionID != nil {
		// Use the partition to get the cache data
		cachePartition := sc.hashRing[partitionID.String()]
		return cachePartition.Get(key)
	}
	return nil
}

// Set set cache value by key
func (sc *ServerCache[T]) Set(key string, value *T) {
	// Find the partitionID associate with the map we need to look for the key
	partitionID := sc.consistentHasher.LocateKey([]byte(key))
	if partitionID != nil {
		// Use the partition to set the cache data
		cachePartition := sc.hashRing[partitionID.String()]
		cachePartition.Set(key, value)
	}
}
