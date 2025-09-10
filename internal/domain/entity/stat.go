package entity

type BannerStat struct {
	BannerID    BannerID
	SlotID      SlotID
	GroupID     GroupID
	Impressions int64
	Clicks      int64
	Description string
}

func NewBannerStat(bannerID BannerID, slotID SlotID, groupID GroupID, impressions, clicks int64) (BannerStat, error) {
	return BannerStat{
		BannerID:    bannerID,
		SlotID:      slotID,
		GroupID:     groupID,
		Impressions: impressions,
		Clicks:      clicks,
		Description: "",
	}, nil
}
