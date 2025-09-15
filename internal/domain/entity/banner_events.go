package entity

import "time"

// DomainEvent — general interface for all events.
type DomainEvent interface {
	EventType() string
	EventTime() time.Time
}

// BannerImpressionRecorded - records the fact that the banner was displayed.
type BannerImpressionRecorded struct {
	BannerID BannerID
	SlotID   SlotID
	GroupID  GroupID
	Time     time.Time
}

func (e BannerImpressionRecorded) EventType() string    { return "BannerImpressionRecorded" }
func (e BannerImpressionRecorded) EventTime() time.Time { return e.Time }

// BannerClickRecorded records the fact that the banner was clicked.
type BannerClickRecorded struct {
	BannerID BannerID
	SlotID   SlotID
	GroupID  GroupID
	Occurred time.Time
}

func (e BannerClickRecorded) EventType() string    { return "BannerClickRecorded" }
func (e BannerClickRecorded) EventTime() time.Time { return e.Occurred }

// BannerAssignedToSlot records the addition of a banner to a slot.
type BannerAssignedToSlot struct {
	BannerID BannerID
	SlotID   SlotID
	Occurred time.Time
}

func (e BannerAssignedToSlot) EventType() string    { return "BannerAssignedToSlot" }
func (e BannerAssignedToSlot) EventTime() time.Time { return e.Occurred }

// BannerRemovedFromSlot records the removal of a banner from a slot.
type BannerRemovedFromSlot struct {
	BannerID BannerID
	SlotID   SlotID
	Occurred time.Time
}

func (e BannerRemovedFromSlot) EventType() string    { return "BannerRemovedFromSlot" }
func (e BannerRemovedFromSlot) EventTime() time.Time { return e.Occurred }
