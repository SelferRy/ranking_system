package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"github.com/SelferRy/ranking_system/internal/domain/interfaces/repository"
	"github.com/jackc/pgx/v5"
	"time"
)

type managementRepository struct {
	pool Pool
}

func NewManagementRepository(pool Pool) repository.ManagementRepository {
	return &managementRepository{pool: pool}
}

func (r *managementRepository) AddBannerToSlot(
	ctx context.Context,
	banner entity.Banner,
	slotID entity.SlotID) error {
	const queryBanners = `
		INSERT INTO ranking_system.banners (id, description, banned_at, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO NOTHING;
	`
	_, err := r.pool.Exec(ctx, queryBanners, banner.ID, banner.Description, nil, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to add banner %d: %w", int64(banner.ID), err)
	}

	const queryBannerSlot = `
		INSERT INTO ranking_system.banner_slot (banner_id, slot_id)
		VALUES ($1, $2)
		ON CONFLICT (banner_id, slot_id) DO NOTHING;
	`
	_, err = r.pool.Exec(ctx, queryBannerSlot, banner.ID, slotID)
	if err != nil {
		return fmt.Errorf("failed to add banner %d to slot %d: %w", int64(banner.ID), int64(slotID), err)
	}
	return nil
}

func (r *managementRepository) BannerExistsInSlot(
	ctx context.Context,
	bannerID entity.BannerID,
	slotID entity.SlotID) (bool, error) {
	const query = `
		SELECT 1
		FROM ranking_system.banner_slot
		WHERE banner_id = $1 AND slot_id = $2
		LIMIT 1;
	`
	var exists int
	err := r.pool.QueryRow(ctx, query, bannerID, slotID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// when no rows in db
			return false, nil
		}
		// when problem with request to db
		return false, fmt.Errorf("failed to check existence of banner %d in slot %d: %w",
			int64(bannerID), int64(slotID), err)
	}
	// when the id exists in db
	return true, nil
}

func (r *managementRepository) RemoveBannerFromSlot(
	ctx context.Context,
	banner entity.Banner,
	slotID entity.SlotID) error {
	const queryBannerSlot = `
		DELETE FROM ranking_system.banner_slot
		WHERE banner_id = $1 AND slot_id = $2; 
	`
	_, err := r.pool.Exec(ctx, queryBannerSlot, banner.ID, slotID)
	if err != nil {
		return fmt.Errorf("failed to remove banner %d to slot %d: %w", int64(banner.ID), int64(slotID), err)
	}

	const queryBanners = `
		DELETE FROM ranking_system.banners
		WHERE id = $1;
	`
	_, err = r.pool.Exec(ctx, queryBanners, banner.ID)
	if err != nil {
		return fmt.Errorf("failed to remove banner %d from banners table: %w", int64(banner.ID), err)
	}
	return nil
}
