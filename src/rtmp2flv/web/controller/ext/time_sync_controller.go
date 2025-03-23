package ext

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
)

func FetchIso8601Time(ctx *gin.Context) {
	now := time.Now()
	timeRfcStr := now.Format(time.RFC3339)
	ctx.String(http.StatusOK, timeRfcStr)
}

func FetchUnixMillisecondTime(ctx *gin.Context) {
	now := time.Now()
	unixMillisecondTime := now.UnixNano() / int64(time.Millisecond)
	ctx.String(http.StatusOK, fmt.Sprintf("%d", unixMillisecondTime))
}

// 判断给定的时间是否在范围内，需要使用url编码，如“+”是%2B
func CheckTimeRang(ctx *gin.Context) {
	clientIso8601Time := ctx.Query("clientIso8601Time")
	if clientIso8601Time == "" {
		logs.Error("Query param clientIso8601Time is rquired")
		http.Error(ctx.Writer, "Query param clientIso8601Time is rquired", http.StatusBadRequest)
		return
	}
	absMillisecondStr := ctx.Query("absMillisecond")
	if absMillisecondStr == "" {
		logs.Error("Query param absMillisecond is rquired")
		http.Error(ctx.Writer, "Query param absMillisecond is rquired", http.StatusBadRequest)
		return
	}
	absMillisecond, err := strconv.ParseInt(absMillisecondStr, 10, 64)
	if err != nil {
		logs.Error("absMillisecond is not a number")
		http.Error(ctx.Writer, "absMillisecond is not a number", http.StatusBadRequest)
		return
	}
	clientTime, err := time.Parse(time.RFC3339, clientIso8601Time)
	if err != nil {
		logs.Error("client time parse error: %v", err)
		http.Error(ctx.Writer, "client time parse error", http.StatusBadRequest)
		return
	}
	sinceTime := time.Since(clientTime)
	if math.Abs(float64(sinceTime)) > float64(time.Duration(absMillisecond)*time.Millisecond) {
		now := time.Now()
		timeRfcStr := now.Format(time.RFC3339)
		ctx.String(http.StatusOK, timeRfcStr)
	}
	ctx.String(http.StatusOK, "")
}
