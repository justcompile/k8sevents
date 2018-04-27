package exts

import "testing"

func TestInArray(t *testing.T) {
	tables := []struct {
		x int
		y []int
		z bool
	}{
		{1, []int{2, 3, 4}, false},
		{2, []int{2, 3, 4}, true},
		{3, []int{2, 3, 4}, true},
		{4, []int{2, 3, 5}, false},
	}

	for _, table := range tables {
		exists, _ := InArray(table.x, table.y)
		if exists != table.z {
			t.Errorf("Error in asserting value in array. Should have been %v but was %v", table.z, exists)
		}
	}
}
