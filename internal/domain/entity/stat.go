package entity

type BannerStat struct {
	SlotID      SlotID
	BannerID    BannerID
	GroupID     GroupID
	Shows       int64
	Clicks      int64
	Description string
}
