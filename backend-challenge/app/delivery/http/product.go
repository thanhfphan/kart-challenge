package http

import (
	"net/http"

	"github.com/thanhfphan/kart-challenge/app/dto"
	"github.com/thanhfphan/kart-challenge/pkg/logging"

	"github.com/gin-gonic/gin"
)

// TODO: add wrapper context
func (a *app) handleGetProduct() gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.Request.Context()
		log := logging.FromContext(ctx)

		var req dto.IDRequest
		if err := ginctx.BindQuery(&req); err != nil {
			log.Warnf("Failed to bind uri: %v", err)

			ginctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		product, err := a.productUC.Get(ctx, req.ID)
		if err != nil {
			log.Errorf("Failed to get product: %v", err)
			ginctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ginctx.JSON(http.StatusOK, gin.H{
			"data": product,
		})
	}
}
