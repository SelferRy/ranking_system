package bandit

import (
	"testing"

	"github.com/SelferRy/ranking_system/internal/domain/entity"
)

func TestUCB1Selector_SelectBanner(t *testing.T) {
	tests := []struct {
		name     string
		stats    []entity.BannerStat
		expected entity.BannerID
		wantErr  bool
	}{
		{
			name:     "no banners",
			stats:    []entity.BannerStat{},
			expected: 0,
			wantErr:  true,
		},
		{
			name: "banner with zero impressions should be selected",
			stats: []entity.BannerStat{
				{BannerID: 1, SlotID: 1, GroupID: 1, Impressions: 0, Clicks: 0},
				{BannerID: 2, SlotID: 1, GroupID: 1, Impressions: 100, Clicks: 10},
			},
			expected: 1,
			wantErr:  false,
		},
		{
			name: "banner with highest UCB score should be selected",
			stats: []entity.BannerStat{
				{BannerID: 1, SlotID: 1, GroupID: 1, Impressions: 100, Clicks: 5},  // CTR = 0.05
				{BannerID: 2, SlotID: 1, GroupID: 1, Impressions: 100, Clicks: 10}, // CTR = 0.10
				{BannerID: 3, SlotID: 1, GroupID: 1, Impressions: 100, Clicks: 15}, // CTR = 0.15
			},
			expected: 3,
			wantErr:  false,
		},
		{
			name: "exploration vs exploitation balance",
			stats: []entity.BannerStat{
				{BannerID: 1, SlotID: 1, GroupID: 1, Impressions: 1000, Clicks: 100}, // CTR = 0.10, well-explored
				{BannerID: 2, SlotID: 1, GroupID: 1, Impressions: 100, Clicks: 20},   // CTR = 0.20, less explored
			},
			// Banner 2 should be selected due to exploration bonus
			expected: 2,
			wantErr:  false,
		},
	}

	selector := NewUCB1Service()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selector.SelectBanner(tt.stats)

			if (err != nil) != tt.wantErr {
				t.Errorf("SelectBanner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.expected {
				t.Errorf("SelectBanner() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCalculateTotalShows(t *testing.T) {
	tests := []struct {
		name     string
		stats    []entity.BannerStat
		expected int64
	}{
		{
			name:     "empty stats",
			stats:    []entity.BannerStat{},
			expected: 0,
		},
		{
			name: "multiple banners",
			stats: []entity.BannerStat{
				{Impressions: 100},
				{Impressions: 200},
				{Impressions: 300},
			},
			expected: 600,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateTotalShows(tt.stats); got != tt.expected {
				t.Errorf("calculateTotalShows() = %v, want %v", got, tt.expected)
			}
		})
	}
}
