package api

import (
	"net/http"
	"strings"
	"time"

	"login-meta-jatis/model/response"
	"login-meta-jatis/provider"
	"login-meta-jatis/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func authorization(logger provider.ILogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqID := util.GenerateReqID()
		ctx.Set("req-id", reqID)
		ctx.Header("X-REQUEST-ID", reqID)

		bearerToken := strings.Split(ctx.GetHeader("Authorization"), " ")

		unauthorizedResp := response.ErrorResponse{
			Error: response.Error{
				Code:    http.StatusUnauthorized,
				Title:   "Unauthorized",
				Details: "invalid token",
			},
		}

		invalidTokenMsg := "Invalid bearer token"

		if len(bearerToken) != 2 {
			logger.WithFields(provider.AppLog, logrus.Fields{"REQUEST_ID": reqID}).Error(invalidTokenMsg)
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				unauthorizedResp,
			)
			return
		}

		if bearerToken[0] != "Bearer" {
			logger.WithFields(provider.AppLog, logrus.Fields{"REQUEST_ID": reqID}).Error(invalidTokenMsg)
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				unauthorizedResp,
			)
			return
		}

		token := bearerToken[1]

		ctx.Set("token", token)
		ctx.Next()
	}
}

func setReqID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqID := util.GenerateReqID()
		ctx.Set("req-id", reqID)
		ctx.Header("X-REQUEST-ID", reqID)
		ctx.Next()
	}
}

func loggingMiddleware(logger provider.ILogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Starting time
		startTime := time.Now()

		// Processing request
		ctx.Next()

		// End Time
		endTime := time.Now()

		// execution time
		latencyTime := endTime.Sub(startTime)

		// Request method
		reqMethod := ctx.Request.Method

		// Request route
		reqUri := ctx.Request.RequestURI

		// status code
		statusCode := ctx.Writer.Status()

		// Request IP
		clientIP := ctx.ClientIP()

		var reqID string
		val, ok := ctx.Get("req-id")
		if ok {
			reqID = val.(string)
		}

		logger.WithFields(
			provider.AppLog,
			logrus.Fields{
				"METHOD":     reqMethod,
				"URI":        reqUri,
				"STATUS":     statusCode,
				"LATENCY":    latencyTime,
				"CLIENT_IP":  clientIP,
				"REQUEST_ID": reqID,
			},
		).Info("HTTP REQUEST")

		ctx.Next()
	}
}
