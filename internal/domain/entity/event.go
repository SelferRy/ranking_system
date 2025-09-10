package entity

import "time"

type RotationEvent struct {
	Type      string
	BannerID  BannerID
	SlotID    SlotID
	GroupID   GroupID
	EventTime time.Time
}

//func (bs BannerStat) CTR() float64 {
//	if bs.Impressions == 0 {
//		return 0
//	}
//	return float64(bs.Impressions) / float64(bs.Clicks)
//}
