package bandit

import (
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"math"
)

type BannerSelector interface {
	SelectBanner(stats []entity.BannerStat) (entity.BannerID, error)
}

type ucb1Selector struct{}

func NewUCB1Service() BannerSelector {
	return &ucb1Selector{}
}

func (s ucb1Selector) SelectBanner(stats []entity.BannerStat) (entity.BannerID, error) {
	totalShows := calculateTotalShows(stats)

	var (
		maxScore float64 = -1
		selectId entity.BannerID
	)

	for _, stat := range stats {
		score := s.calculateScore(stat, totalShows)
		if score > maxScore {
			maxScore = score
			selectId = stat.BannerID
		}
	}

	return selectId, nil
}

func (s ucb1Selector) calculateScore(stat entity.BannerStat, totalShows int64) float64 {
	if stat.Impressions == 0 {
		return math.Inf(1) // choose it 100% in the case
	}
	//ctr := stat.CTR.Value()
	ctr, err := entity.NewCTR(stat.Impressions, stat.Clicks)
	if err != nil {
		return math.Inf(1)
	}
	return ctr.Value() + math.Sqrt(2*math.Log(float64(totalShows))/float64(stat.Impressions))
}

func calculateTotalShows(stats []entity.BannerStat) int64 {
	var total int64 = 0
	for _, stat := range stats {
		total += stat.Impressions
	}
	return total
}
