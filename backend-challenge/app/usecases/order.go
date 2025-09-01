package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/thanhfphan/kart-challenge/app/dto"
	"github.com/thanhfphan/kart-challenge/app/models"
	"github.com/thanhfphan/kart-challenge/app/repos"
	"github.com/thanhfphan/kart-challenge/config"
	"github.com/thanhfphan/kart-challenge/env"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
)

var _ Order = (*order)(nil)

type Order interface {
	PlaceOrder(ctx context.Context, req *dto.OrderRequest) (*dto.OrderResponse, error)
}

type order struct {
	cfg *config.Config
	env *env.Env

	baseRepo      repos.Repo
	orderRepo     repos.Order
	orderItemRepo repos.OrderItem
	productRepo   repos.Product
	promoCodeRepo repos.PromoCode
}

func newOrder(cfg *config.Config, env *env.Env, repos repos.Repo) (Order, error) {
	return &order{
		cfg:           cfg,
		env:           env,
		baseRepo:      repos,
		orderRepo:     repos.Order(),
		orderItemRepo: repos.OrderItem(),
		productRepo:   repos.Product(),
		promoCodeRepo: repos.PromoCode(),
	}, nil
}

func (u *order) PlaceOrder(ctx context.Context, req *dto.OrderRequest) (*dto.OrderResponse, error) {
	log := logging.FromContext(ctx)
	log.Infof("Placing order with %d items", len(req.Items))

	productIDs := make([]int64, 0, len(req.Items))
	itemMap := make(map[int64]*dto.OrderItem)

	for _, item := range req.Items {
		productID, err := strconv.ParseInt(item.ProductID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID: %s", item.ProductID)
		}
		productIDs = append(productIDs, productID)
		itemMap[productID] = &item
	}

	products, err := u.productRepo.GetByIDList(ctx, productIDs)
	if err != nil {
		log.Errorf("Failed to get products: %v", err)
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	if len(products) != len(productIDs) {
		return nil, fmt.Errorf("some products not found")
	}

	// Calculate total and create order items
	var total float64
	orderItems := make([]*models.OrderItem, 0, len(products))
	responseItems := make([]dto.OrderItemResponse, 0, len(products))
	responseProducts := make([]dto.ProductResponse, 0, len(products))

	for _, product := range products {
		item := itemMap[product.ID]
		itemTotal := product.Price * float64(item.Quantity)
		total += itemTotal

		orderItems = append(orderItems, &models.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})

		responseItems = append(responseItems, dto.OrderItemResponse{
			ProductID: fmt.Sprintf("%d", product.ID),
			Quantity:  item.Quantity,
		})

		responseProducts = append(responseProducts, dto.ProductResponse{
			ID:       fmt.Sprintf("%d", product.ID),
			Name:     product.Name,
			Price:    product.Price,
			Category: product.Category,
			Image: &dto.ProductImage{
				Thumbnail: product.ThumbnailURL,
				Mobile:    product.MobileURL,
				Tablet:    product.TabletURL,
				Desktop:   product.DesktopURL,
			},
		})
	}

	var (
		finalTotal   float64
		discount     float64
		createdOrder *models.Order
		coupon       *models.PromoCode
	)

	if req.CouponCode != "" {
		coupon, err = u.promoCodeRepo.GetCode(ctx, req.CouponCode)
		if err != nil {
			log.Errorf("Failed to get promo code: %v", err)
			return nil, fmt.Errorf("failed to get promo code: %w", err)
		}

		if !coupon.IsActive {
			return nil, fmt.Errorf("promo code is not active")
		}

		discount = total * coupon.DiscountPct / 100
	}

	finalTotal = total - discount

	// Create order in transaction
	err = u.baseRepo.WithTransaction(ctx, func(tx repos.Repo) error {
		createdOrder, err = tx.Order().Create(ctx, &models.Order{
			ID:         uuid.New().String(),
			Total:      finalTotal,
			Discounts:  discount,
			CouponCode: req.CouponCode,
			Status:     "completed",
		})
		if err != nil {
			return err
		}

		for _, item := range orderItems {
			item.OrderID = createdOrder.ID
		}

		err = tx.OrderItem().CreateMany(ctx, orderItems)
		if err != nil {
			return err
		}

		if coupon != nil {
			// If we want to deactivate the coupon after use, we can do it here
			// err = tx.PromoCode().UpdateWithMap(ctx, coupon, map[string]interface{}{
			// 	// TODO: Use pessimistic lock to avoid concurrent use of the same coupon
			// 	"is_active": false, // Deactivate the coupon after use
			// })
			// if err != nil {
			// 	return err
			// }
		}

		return nil
	})
	if err != nil {
		log.Errorf("Failed to create order: %v", err)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	log.Infof("Order created successfully: %s, total: %.2f", createdOrder.ID, createdOrder.Total)

	return &dto.OrderResponse{
		ID:        createdOrder.ID,
		Total:     createdOrder.Total,
		Discounts: createdOrder.Discounts,
		Items:     responseItems,
		Products:  responseProducts,
	}, nil
}
