package riot

import (
	"testing"
)

func TestCompareItems(t *testing.T) {
	tests := []struct {
		name     string
		oldItem  ItemData
		newItem  ItemData
		expected int // number of diffs
	}{
		{
			name: "No changes",
			oldItem: ItemData{
				Name: "Item 1",
				Gold: ItemGold{Total: 1000},
				Stats: map[string]float64{
					"FlatHPPoolMod": 100,
				},
			},
			newItem: ItemData{
				Name: "Item 1",
				Gold: ItemGold{Total: 1000},
				Stats: map[string]float64{
					"FlatHPPoolMod": 100,
				},
			},
			expected: 0,
		},
		{
			name: "Cost increase (nerf)",
			oldItem: ItemData{
				Gold: ItemGold{Total: 1000},
			},
			newItem: ItemData{
				Gold: ItemGold{Total: 1100},
			},
			expected: 1,
		},
		{
			name: "Stat increase (buff)",
			oldItem: ItemData{
				Stats: map[string]float64{
					"FlatHPPoolMod": 100,
				},
			},
			newItem: ItemData{
				Stats: map[string]float64{
					"FlatHPPoolMod": 150,
				},
			},
			expected: 1,
		},
		{
			name: "Multiple changes",
			oldItem: ItemData{
				Gold: ItemGold{Total: 1000},
				Stats: map[string]float64{
					"FlatHPPoolMod": 100,
				},
			},
			newItem: ItemData{
				Gold: ItemGold{Total: 900}, // buff
				Stats: map[string]float64{
					"FlatHPPoolMod": 80, // nerf
				},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diffs := compareItems(tt.oldItem, tt.newItem)
			if len(diffs) != tt.expected {
				t.Errorf("expected %d diffs, got %d", tt.expected, len(diffs))
			}
		})
	}
}
