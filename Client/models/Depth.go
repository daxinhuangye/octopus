package models

import (
	"crypto/hmac"
	"crypto/sha256"
	_ "encoding/hex"
	"fmt"
	"net/url"
	_ "reflect"
	"sort"
	"strings"
	"time"
	"tsEngine/tsCrypto"
	"tsEngine/tsTime"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/tidwall/gjson"
)

//行情数据 btcusdt
func (this *Binance) Binance(symbol string, num float64) (float64, float64, int64) {

	api := "https://" + BA_HOST + "/api/v1/depth?symbol=" + symbol + "&limit=5"

	content := this.request(api)
	if content == "" {
		return 0, 0, 0
	}

	//卖盘
	tick := gjson.Get(content, "asks").Array()
	temp := tick[0].Array()
	asks := temp[0].Float()

	//买盘
	tick = gjson.Get(content, "bids").Array()
	temp = tick[0].Array()
	bids := temp[0].Float()

	//时间
	ts := gjson.Get(content, "lastUpdateId").Int()

	return bids, asks, ts

}
