package bandit

import (
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"math"
)

type BannerSelector interface {
	SelectBanner(stats []entity.BannerStat) (int64, error)
}

type ucb1Selector struct{}

func NewUCB1Service() BannerSelector {
	return &ucb1Selector{}
}

func (s ucb1Selector) SelectBanner(stats []entity.BannerStat) (int64, error) {
	totalShows := calculateTotalShows(stats)

	var (
		maxScore float64 = -1
		selectId int64
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
	if stat.Shows == 0 {
		return math.Inf(1) // choose it 100% in the case
	}
	ctr := float64(stat.Clicks) / float64(stat.Shows)
	return ctr + math.Sqrt(2*math.Log(float64(totalShows))/float64(stat.Shows))
}

func calculateTotalShows(stats []entity.BannerStat) int64 {
	var total int64 = 0
	for _, stat := range stats {
		total += stat.Shows
	}
	return total
}

// TODO: пробежаться по типам в entity и изменить с int64 на int где это избыточно - чтобы показать,
//  что понятно зачем использую int64. Это демо проект. Нагрузка не будет большой. Важна именно демонстрация навыков.
