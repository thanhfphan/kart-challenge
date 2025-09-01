package usecases

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/thanhfphan/kart-challenge/app/models"
	"github.com/thanhfphan/kart-challenge/app/repos"
	"github.com/thanhfphan/kart-challenge/config"
	"github.com/thanhfphan/kart-challenge/env"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
)

var _ PromoCode = (*promoCode)(nil)

type PromoCode interface {
	ProcessCouponFile(ctx context.Context, filePath string) error
}

type promoCode struct {
	cfg *config.Config
	env *env.Env

	promoCodeRepo repos.PromoCode
	repo          repos.Repo
}

func newPromoCode(cfg *config.Config, env *env.Env, repos repos.Repo) (PromoCode, error) {
	return &promoCode{
		cfg:           cfg,
		env:           env,
		promoCodeRepo: repos.PromoCode(),
		repo:          repos,
	}, nil
}

func (p *promoCode) ProcessCouponFile(ctx context.Context, filePath string) error {
	log := logging.FromContext(ctx)
	log.Infof("Starting to process coupon file: %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open coupon file %s: %w", filePath, err)
	}
	defer file.Close()

	// Read file line by line using streaming
	scanner := bufio.NewScanner(file)
	var promoCodes []*models.PromoCode
	batchSize := 1000 // Process in batches for memory efficiency
	currentTime := time.Now().Unix()

	lineCount := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		lineCount++

		promoCode := &models.PromoCode{
			Code:        line,
			Description: fmt.Sprintf("Promo code %s", line),
			DiscountPct: 10.0, // Default 10% discount
			IsActive:    true,
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
		}

		promoCodes = append(promoCodes, promoCode)

		if len(promoCodes) >= batchSize {
			if err := p.processBatch(ctx, promoCodes); err != nil {
				return fmt.Errorf("failed to process batch at line %d: %w", lineCount, err)
			}
			promoCodes = promoCodes[:0] // Reset slice but keep capacity
		}
	}

	// Process remaining codes
	if len(promoCodes) > 0 {
		if err := p.processBatch(ctx, promoCodes); err != nil {
			return fmt.Errorf("failed to process final batch: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading coupon file: %w", err)
	}

	log.Infof("Successfully processed %d coupon codes from file: %s", lineCount, filePath)
	return nil
}

func (p *promoCode) processBatch(ctx context.Context, promoCodes []*models.PromoCode) error {
	log := logging.FromContext(ctx)

	err := p.promoCodeRepo.BulkUpsert(ctx, promoCodes)
	if err != nil {
		log.Errorf("Failed to bulk upsert promo codes: %v", err)
		return err
	}

	return nil
}
