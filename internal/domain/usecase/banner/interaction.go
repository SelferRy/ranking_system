package banner

//func (uc DeliveryUseCase) RegisterClick(ctx context.Context, slotID, bannerID, groupID int64) error {
//	if err := uc.statsRepo.IncrementClick(ctx, slotID, bannerID, groupID); err != nil {
//		return fmt.Errorf("increment banner click: %w", err)
//	}
//	event := entity.RotationEvent{
//		Type:      "click",
//		BannerID:  bannerID,
//		SlotID:    slotID,
//		GroupID:   groupID,
//		EventTime: time.Now().UTC(),
//	}
//	if err := uc.producer.Send(ctx, event); err != nil {
//		uc.log.Error("failed to send click event", logger.Error(err))
//	}
//	return nil
//}
