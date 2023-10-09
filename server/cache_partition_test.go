package server

import "testing"

type TestCacheableStruct struct {
	data string
}

func (tcs TestCacheableStruct) String() string {
	return tcs.data
}

func TestCachePartitionSetAndRetrieveValue(t *testing.T) {
	tests := map[string]struct {
		ID       string
		Key      string
		Value    *TestCacheableStruct
		Expected *TestCacheableStruct
	}{
		"Test normal key": {
			ID:  "Partition1",
			Key: "Key1",
			Value: &TestCacheableStruct{
				data: "Data1",
			},
			Expected: &TestCacheableStruct{
				data: "Data1",
			},
		},
		"Test empty key": {
			ID:  "Partition1",
			Key: "",
			Value: &TestCacheableStruct{
				data: "Data2",
			},
			Expected: &TestCacheableStruct{
				data: "Data2",
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			cachePartition := NewServerCachePartition[TestCacheableStruct](test.ID)
			cachePartition.Set(test.Key, test.Value)
			result := cachePartition.Get(test.Key)
			if result == nil {
				t.Errorf("expected saved result, got %v", result)
			} else if result.data != test.Expected.data {
				t.Errorf("expected %v, got %v", test.Expected.data, result.data)
			}
		})
	}
}
