package usecases

import (
	"context"
	"fmt"
	"strconv"

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
	GetByID(ctx context.Context, id string) (*dto.OrderResponse, error)
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

	// Validate and get products
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

	// Get products from database
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

	// Apply promo code discount if provided
	// var discount float64
	// if req.CouponCode != "" {
	// 	valid, discountPct, err := u.promoCodeRepo.ValidateCode(ctx, req.CouponCode)
	// 	if err != nil {
	// 		log.Errorf("Failed to validate promo code: %v", err)
	// 		return nil, fmt.Errorf("failed to validate promo code: %w", err)
	// 	}

	// 	if valid {
	// 		discount = total * (discountPct / 100)
	// 		log.Infof("Applied promo code %s: %.2f%% discount (%.2f)", req.CouponCode, discountPct, discount)
	// 	} else {
	// 		log.Warnf("Invalid promo code: %s", req.CouponCode)
	// 		return nil, fmt.Errorf("invalid promo code: %s", req.CouponCode)
	// 	}
	// }

	// finalTotal := total - discount

	// // Create order
	// order := &models.Order{
	// 	ID:         uuid.New().String(),
	// 	Total:      finalTotal,
	// 	Discounts:  discount,
	// 	CouponCode: req.CouponCode,
	// 	Status:     "completed",
	// }

	// // Save order with items
	// createdOrder, err := u.orderRepo.CreateWithItems(ctx, order, orderItems)
	// if err != nil {
	// 	log.Errorf("Failed to create order: %v", err)
	// 	return nil, fmt.Errorf("failed to create order: %w", err)
	// }

	// log.Infof("Order created successfully: %s, total: %.2f", createdOrder.ID, createdOrder.Total)

	return &dto.OrderResponse{
		// ID:        createdOrder.ID,
		// Total:     createdOrder.Total,
		// Discounts: createdOrder.Discounts,
		Items:    responseItems,
		Products: responseProducts,
	}, nil
}

func (u *order) GetByID(ctx context.Context, id string) (*dto.OrderResponse, error) {
	log := logging.FromContext(ctx)
	log.Infof("Getting order by id=%s", id)

	order, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		log.Errorf("Get order by id=%s err=%v", id, err)
		return nil, err
	}

	orderItems, err := u.orderItemRepo.GetByOrderID(ctx, id)
	if err != nil {
		log.Errorf("Get order items by order_id=%s err=%v", id, err)
		return nil, err
	}

	responseItems := make([]dto.OrderItemResponse, 0, len(orderItems))
	responseProducts := make([]dto.ProductResponse, 0, len(orderItems))

	for _, item := range orderItems {
		responseItems = append(responseItems, dto.OrderItemResponse{
			ProductID: fmt.Sprintf("%d", item.ProductID),
			Quantity:  item.Quantity,
		})

		itemProduct, err := u.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			log.Errorf("Get product by id=%d err=%v", item.ProductID, err)
			return nil, err
		}

		if itemProduct != nil {
			responseProducts = append(responseProducts, dto.ProductResponse{
				ID:       fmt.Sprintf("%d", itemProduct.ID),
				Name:     itemProduct.Name,
				Price:    itemProduct.Price,
				Category: itemProduct.Category,
				Image: &dto.ProductImage{
					Thumbnail: itemProduct.ThumbnailURL,
					Mobile:    itemProduct.MobileURL,
					Tablet:    itemProduct.TabletURL,
					Desktop:   itemProduct.DesktopURL,
				},
			})
		}
	}

	return &dto.OrderResponse{
		ID:        order.ID,
		Total:     order.Total,
		Discounts: order.Discounts,
		Items:     responseItems,
		Products:  responseProducts,
	}, nil
}
