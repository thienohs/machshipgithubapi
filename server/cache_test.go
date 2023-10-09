package server

import "testing"

func TestCacheSetAndRetrieveValue(t *testing.T) {
	tests := map[string]struct {
		NumberOfCachePartitions int
		InputsWithExpected      []struct {
			Key      string
			Value    *TestCacheableStruct
			Expected *TestCacheableStruct
		}
	}{
		"Test normal key, 1 partition": {
			NumberOfCachePartitions: 1,
			InputsWithExpected: []struct {
				Key      string
				Value    *TestCacheableStruct
				Expected *TestCacheableStruct
			}{
				{
					Key: "Key1",
					Value: &TestCacheableStruct{
						data: "Data1",
					},
					Expected: &TestCacheableStruct{
						data: "Data1",
					},
				},
			},
		},
		"Test normal key, 7 partitions": {
			NumberOfCachePartitions: 7,
			InputsWithExpected: []struct {
				Key      string
				Value    *TestCacheableStruct
				Expected *TestCacheableStruct
			}{
				{
					Key: "Key1",
					Value: &TestCacheableStruct{
						data: "Data1",
					},
					Expected: &TestCacheableStruct{
						data: "Data1",
					},
				},
				{
					Key: "Key2",
					Value: &TestCacheableStruct{
						data: "Data2",
					},
					Expected: &TestCacheableStruct{
						data: "Data2",
					},
				},
				{
					Key: "Key3",
					Value: &TestCacheableStruct{
						data: "Data3",
					},
					Expected: &TestCacheableStruct{
						data: "Data3",
					},
				},
			},
		},
		"Test multiple empty keys, 7 partitions": {
			NumberOfCachePartitions: 7,
			InputsWithExpected: []struct {
				Key      string
				Value    *TestCacheableStruct
				Expected *TestCacheableStruct
			}{
				{
					Key: "",
					Value: &TestCacheableStruct{
						data: "Data1",
					},
					Expected: &TestCacheableStruct{
						data: "Data1",
					},
				},
				{
					Key: "",
					Value: &TestCacheableStruct{
						data: "Data2",
					},
					Expected: &TestCacheableStruct{
						data: "Data2",
					},
				},
				{
					Key: "",
					Value: &TestCacheableStruct{
						data: "Data3",
					},
					Expected: &TestCacheableStruct{
						data: "Data3",
					},
				},
			},
		},
		"Test 0 partition": {
			NumberOfCachePartitions: 0,
			InputsWithExpected: []struct {
				Key      string
				Value    *TestCacheableStruct
				Expected *TestCacheableStruct
			}{
				{
					Key: "",
					Value: &TestCacheableStruct{
						data: "Data1",
					},
					Expected: nil,
				},
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			cache := NewServerCache[TestCacheableStruct](test.NumberOfCachePartitions)
			for _, input := range test.InputsWithExpected {
				cache.Set(input.Key, input.Value)
				result := cache.Get(input.Key)
				if result == nil && input.Expected != nil {
					t.Errorf("expected saved result, got %v", result)
				} else if input.Expected == nil && result != nil {
					t.Errorf("expected nil, got %v", result)
				} else if result != nil && input.Expected != nil && result.data != input.Expected.data {
					t.Errorf("expected %v for key %v, got %v", input.Expected, input.Key, result)
				}
			}
		})
	}
}
