package module

import (
	"Nicepay-CC/Helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RequestInq(timetrx, iMid, merTok, tXid, referenceNo, amt string) string {
	var req = `{
	"timeStamp":"` + timetrx + `",
	"tXid":"` + tXid + `",
	"iMid":"` + iMid + `",
	"referenceNo":"` + referenceNo + `",
	"amt":"` + amt + `",
	"merchantToken":"` + merTok + `"
}`
	return req
}

type ResponseInq struct {
	TXid           string      `json:"tXid"`
	IMid           string      `json:"iMid"`
	Currency       string      `json:"currency"`
	Amt            string      `json:"amt"`
	InstmntMon     string      `json:"instmntMon"`
	InstmntType    string      `json:"instmntType"`
	ReferenceNo    string      `json:"referenceNo"`
	GoodsNm        string      `json:"goodsNm"`
	PayMethod      string      `json:"payMethod"`
	BillingNm      string      `json:"billingNm"`
	ReqDt          string      `json:"reqDt"`
	ReqTm          string      `json:"reqTm"`
	Status         string      `json:"status"`
	ResultCd       string      `json:"resultCd"`
	ResultMsg      string      `json:"resultMsg"`
	CardNo         string      `json:"cardNo"`
	PreauthToken   interface{} `json:"preauthToken"`
	AcquBankCd     string      `json:"acquBankCd"`
	IssuBankCd     string      `json:"issuBankCd"`
	VacctValidDt   interface{} `json:"vacctValidDt"`
	VacctValidTm   interface{} `json:"vacctValidTm"`
	VacctNo        interface{} `json:"vacctNo"`
	BankCd         interface{} `json:"bankCd"`
	PayNo          interface{} `json:"payNo"`
	MitraCd        interface{} `json:"mitraCd"`
	ReceiptCode    interface{} `json:"receiptCode"`
	CancelAmt      interface{} `json:"cancelAmt"`
	TransDt        string      `json:"transDt"`
	TransTm        string      `json:"transTm"`
	RecurringToken interface{} `json:"recurringToken"`
	CcTransType    string      `json:"ccTransType"`
	PayValidDt     interface{} `json:"payValidDt"`
	PayValidTm     interface{} `json:"payValidTm"`
	MRefNo         interface{} `json:"mRefNo"`
	AcquStatus     string      `json:"acquStatus"`
	CardExpYymm    string      `json:"cardExpYymm"`
	AcquBankNm     string      `json:"acquBankNm"`
	IssuBankNm     string      `json:"issuBankNm"`
	DepositDt      interface{} `json:"depositDt"`
	DepositTm      interface{} `json:"depositTm"`
	PaymentExpDt   interface{} `json:"paymentExpDt"`
	PaymentExpTm   interface{} `json:"paymentExpTm"`
	PaymentTrxSn   interface{} `json:"paymentTrxSn"`
	CancelTrxSn    interface{} `json:"cancelTrxSn"`
	UserID         interface{} `json:"userId"`
	ShopID         interface{} `json:"shopId"`
	AuthNo         string      `json:"authNo"`
}

func Status(iMid, merchantKey, inq_endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			TXid        string `json:"trxID"`
			ReferenceNo string `json:"noRef"`
			Amt         string `json:"amount"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		timestampTrx := time.Now().Format("20060102150405")
		merchantData := timestampTrx + iMid + body.ReferenceNo + body.Amt + merchantKey
		merchantToken := Helper.GenerateMerchantToken(merchantData)

		url := inq_endpoint
		inq_nicepay := &ResponseInq{}
		a, _ := Helper.Request().
			SetHeader("Content-Type", "application/json").
			SetBody(RequestInq(timestampTrx, iMid, merchantToken, body.TXid, body.ReferenceNo, body.Amt)).
			SetResult(inq_nicepay).
			Post(url)
		_ = a.String()

		if inq_nicepay.ResultCd == "0000" && inq_nicepay.ResultMsg == "paid" {
			c.JSON(http.StatusOK, gin.H{
				"status_bayar": "Y",
				"details": gin.H{
					"trxID":       inq_nicepay.TXid,
					"currency":    inq_nicepay.Currency,
					"amount":      inq_nicepay.Amt,
					"instmntMon":  inq_nicepay.InstmntMon,
					"instmntType": inq_nicepay.InstmntType,
					"noRef":       inq_nicepay.ReferenceNo,
					"goodsNm":     inq_nicepay.GoodsNm,
					"payMethod":   inq_nicepay.PayMethod,
					"billingNm":   inq_nicepay.BillingNm,
					"reqDt":       inq_nicepay.ReqDt,
					"reqTm":       inq_nicepay.ReqTm,
					"status":      inq_nicepay.Status,
					"cardNo":      inq_nicepay.CardNo,
					"acquBankCd":  inq_nicepay.AcquBankCd,
					"issuBankCd":  inq_nicepay.IssuBankCd,
					"transDt":     inq_nicepay.TransDt,
					"transTm":     inq_nicepay.TransTm,
					"ccTransType": inq_nicepay.CcTransType,
					"acquStatus":  inq_nicepay.AcquStatus,
					"cardExpYymm": inq_nicepay.CardExpYymm,
					"acquBankNm":  inq_nicepay.AcquBankNm,
					"issuBankNm":  inq_nicepay.IssuBankNm,
					"authNo":      inq_nicepay.AuthNo,
				},
			})
		} else if inq_nicepay.ResultCd == "0000" && inq_nicepay.ResultMsg == "init" {
			c.JSON(http.StatusOK, gin.H{
				"status_bayar": "N",
				"details": gin.H{
					"trxID":       inq_nicepay.TXid,
					"currency":    inq_nicepay.Currency,
					"amount":      inq_nicepay.Amt,
					"instmntMon":  inq_nicepay.InstmntMon,
					"instmntType": inq_nicepay.InstmntType,
					"noRef":       inq_nicepay.ReferenceNo,
					"goodsNm":     inq_nicepay.GoodsNm,
					"payMethod":   inq_nicepay.PayMethod,
					"billingNm":   inq_nicepay.BillingNm,
					"reqDt":       inq_nicepay.ReqDt,
					"reqTm":       inq_nicepay.ReqTm,
					"status":      inq_nicepay.Status,
				},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"resultCd":  inq_nicepay.ResultCd,
				"resultMsg": inq_nicepay.ResultMsg,
			})
		}
	}
}
