package entity

import "time"

type (
	Banner struct {
		ID          int64
		Description string
	}

	Slot struct {
		ID          int64
		Description string
	}

	SocialDemographicGroup struct {
		ID          int64
		Description string
	}

	BannerStat struct {
		SlotID      int64
		BannerID    int64
		GroupID     int64
		Shows       int64
		Clicks      int64
		Description string
	}

	RotationEvent struct {
		Type      string
		BannerID  int64
		SlotID    int64
		GroupID   int64
		EventTime time.Time
	}
)

//func (bs BannerStat) CTR() float64 {
//	if bs.Shows == 0 {
//		return 0
//	}
//	return float64(bs.Shows) / float64(bs.Clicks)
//}
