package middleware

import (
	"github.com/thanhfphan/kart-challenge/pkg/contxt"
	"github.com/thanhfphan/kart-challenge/pkg/logging"

	"github.com/gin-gonic/gin"
)

func SetupAppContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := contxt.ContextWithAppWrapper(c.Request.Context(), &contxt.AppContext{})

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func SetCommonData() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := logging.FromContext(ctx)

		ctxWrapper, err := contxt.GetAppWrapper(ctx)
		if err != nil {
			log.Warnf("Get context wrapper failed with err=%v", err)
			c.Next()
			return
		}

		md := c.Request.Header

		if val := md.Get("accept-language"); len(val) > 0 {
			ctxWrapper.Set("accept-language", val)
		}

		if val := md.Get("token"); len(val) > 0 {
			ctxWrapper.Set("token", val)
		}

		if val := md.Get("x-request-id"); len(val) > 0 {
			ctxWrapper.Set("x-request-id", val)
		}

		ctx = contxt.ContextWithAppWrapper(ctx, ctxWrapper)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
