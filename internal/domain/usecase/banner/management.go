package banner

//func (uc DeliveryUseCase) AddBannerToSlot(ctx context.Context, slotID, bannerID int64) error {
//	exists, err := uc.bannerRepo.ExistsInSlot(ctx, slotID, bannerID)
//	if err != nil {
//		return fmt.Errorf("check banner existence: %w", err)
//	}
//	if exists {
//		uc.log.Info(
//			"banner already exists in the slot",
//			logger.Int64("slotID", slotID),
//			logger.Int64("bannerID", bannerID),
//		)
//		return nil
//	}
//	if err := uc.bannerRepo.AddToSlot(ctx, slotID, bannerID); err != nil {
//		return fmt.Errorf("add banner to slot: %w", err)
//	}
//	uc.log.Info(
//		"banner added to slot",
//		logger.Int64("slotID", slotID),
//		logger.Int64("bannerID", bannerID),
//	)
//	return nil
//}
//
//func (uc DeliveryUseCase) RemoveBannerFromSlot(ctx context.Context, slotID, bannerID int64) error {
//	if err := uc.bannerRepo.RemoveFromSlot(ctx, slotID, bannerID); err != nil {
//		return fmt.Errorf("remove banner from slot: %w", err)
//	}
//	uc.log.Info(
//		"banner removed from slot",
//		logger.Int64("slotID", slotID),
//		logger.Int64("bannerID", bannerID),
//	)
//	return nil
//}
