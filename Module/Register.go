package module

import (
	"eska/Helper"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Response struct {
	ResultCd     string      `json:"resultCd"`
	ResultMsg    string      `json:"resultMsg"`
	TXid         string      `json:"tXid"`
	ReferenceNo  string      `json:"referenceNo"`
	PayMethod    string      `json:"payMethod"`
	Amt          string      `json:"amt"`
	TransDt      string      `json:"transDt"`
	TransTm      string      `json:"transTm"`
	Description  string      `json:"description"`
	BankCd       interface{} `json:"bankCd"`
	VacctNo      interface{} `json:"vacctNo"`
	MitraCd      interface{} `json:"mitraCd"`
	PayNo        interface{} `json:"payNo"`
	Currency     interface{} `json:"currency"`
	GoodsNm      interface{} `json:"goodsNm"`
	BillingNm    interface{} `json:"billingNm"`
	VacctValidDt interface{} `json:"vacctValidDt"`
	VacctValidTm interface{} `json:"vacctValidTm"`
	PayValidDt   interface{} `json:"payValidDt"`
	PayValidTm   interface{} `json:"payValidTm"`
	RequestURL   interface{} `json:"requestURL"`
	PaymentExpDt interface{} `json:"paymentExpDt"`
	PaymentExpTm interface{} `json:"paymentExpTm"`
	QrContent    interface{} `json:"qrContent"`
	QrURL        interface{} `json:"qrUrl"`
}

func Request(timetrx, iMid, referenceNo, merTok, amt, goodsNm, billingNm, billingPhone, billingEmail, billingAddr, billingCity, billingState, billingPostCd, billingCountry, deliveryNm, deliveryPhone, deliveryAddr, deliveryCity, deliveryState, deliveryPostCd, deliveryCountry, description string) string {
	var req = `{
	"timeStamp":"` + timetrx + `",
	"iMid":"` + iMid + `",
	"payMethod":"01",
	"currency":"IDR",
	"amt":"` + amt + `",
	"referenceNo":"` + referenceNo + `",
	"goodsNm":"` + goodsNm + `",
	"billingNm":"` + billingNm + `",
	"billingPhone":"` + billingPhone + `",
	"billingEmail":"` + billingEmail + `",
	"billingAddr":"` + billingAddr + `",
	"billingCity":"` + billingCity + `",
	"billingState":"` + billingState + `",
	"billingPostCd":"` + billingPostCd + `",
	"billingCountry":"` + billingCountry + `",
	"deliveryNm":"` + deliveryNm + `",
	"deliveryPhone":"` + deliveryPhone + `",
	"deliveryAddr":"` + deliveryAddr + `",
	"deliveryCity":"` + deliveryCity + `",
	"deliveryState":"` + deliveryState + `",
	"deliveryPostCd":"` + deliveryPostCd + `",
	"deliveryCountry":"` + deliveryCountry + `",
	"dbProcessUrl":"https://ptsv2.com/t/test-nicepay-v2",
	"description":"` + description + `",
	"merchantToken":"` + merTok + `",
	"reqDomain":"eskaseptian.com",
	"cartData":"{\"count\":1,\"item\":[{\"img_url\":\"https:\/\/images.ctfassets.net\/od02wyo8cgm5\/14Njym0dRLAHaVJywF8WFL\/1910357dd0da0ae38b61bc8ad3db86e4\/cloudflyer_2-fw19-grey_lime-m-g1.png\",\"goods_name\":\"Shoe\",\"goods_detail\":\"Shoe\",\"goods_amt\":` + amt + `]}",
	"instmntType":"2",
	"instmntMon":"1",
	"recurrOpt":"0"
}`
	return req
}

func Register(iMid, merchantKey, reg_endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Amt             string `json:"amt"`
			GoodsNm         string `json:"goodsNm"`
			BillingNm       string `json:"billingNm"`
			BillingPhone    string `json:"billingPhone"`
			BillingEmail    string `json:"billingEmail"`
			BillingAddr     string `json:"billingAddr"`
			BillingCity     string `json:"billingCity"`
			BillingState    string `json:"billingState"`
			BillingPostCd   string `json:"billingPostCd"`
			BillingCountry  string `json:"billingCountry"`
			DeliveryNm      string `json:"deliveryNm"`
			DeliveryPhone   string `json:"deliveryPhone"`
			DeliveryAddr    string `json:"deliveryAddr"`
			DeliveryCity    string `json:"deliveryCity"`
			DeliveryState   string `json:"deliveryState"`
			DeliveryPostCd  string `json:"deliveryPostCd"`
			DeliveryCountry string `json:"deliveryCountry"`
			Description     string `json:"description"`
			CartData        string `json:"cartData"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		timestampTrx := time.Now().Format("20060102150405")
		refNo := "ORDITM" + timestampTrx

		merchantData := timestampTrx + iMid + refNo + body.Amt + merchantKey
		merchantToken := Helper.GenerateMerchantToken(merchantData)

		url := "https://dev.nicepay.co.id/" + reg_endpoint
		register_nicepay := &Response{}
		a, _ := Helper.Request().
			SetHeader("Content-Type", "application/json").
			SetBody(Request(
				timestampTrx,
				iMid,
				refNo,
				merchantToken,
				body.Amt,
				body.GoodsNm,
				body.BillingNm,
				body.BillingPhone,
				body.BillingEmail,
				body.BillingAddr,
				body.BillingCity,
				body.BillingState,
				body.BillingPostCd,
				body.BillingCountry,
				body.DeliveryNm,
				body.DeliveryPhone,
				body.DeliveryAddr,
				body.DeliveryCity,
				body.DeliveryState,
				body.DeliveryPostCd,
				body.DeliveryCountry,
				body.Description)).
			SetResult(register_nicepay).
			Post(url)
		_ = a.String()

		resultCd := fmt.Sprint(register_nicepay.ResultCd)

		if resultCd == "0000" {
			c.JSON(http.StatusCreated, gin.H{
				"resultCd":    register_nicepay.ResultCd,
				"resultMsg":   register_nicepay.ResultMsg,
				"trxID":       register_nicepay.TXid,
				"noRef":       register_nicepay.ReferenceNo,
				"payMethod":   register_nicepay.PayMethod,
				"amount":      register_nicepay.Amt,
				"transDt":     register_nicepay.TransDt,
				"transTm":     register_nicepay.TransTm,
				"description": register_nicepay.Description,
			})
		} else {
			c.JSON(401, gin.H{"resultMsg": register_nicepay.ResultMsg, "resultCd": register_nicepay.ResultCd})
		}
	}
}
