package module

import (
	"Nicepay-CC/Helper"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type ResponsePayment struct {
	AcquBankCd     string `json:"acquBankCd"`
	AcquBankNm     string `json:"acquBankNm"`
	Amt            string `json:"amt"`
	AuthNo         string `json:"authNo"`
	BillingNm      string `json:"billingNm"`
	CardBrand      string `json:"cardBrand"`
	CardExpYymm    string `json:"cardExpYymm"`
	CardNo         string `json:"cardNo"`
	CcTransType    string `json:"ccTransType"`
	Currency       string `json:"currency"`
	Description    string `json:"description"`
	GoodsNm        string `json:"goodsNm"`
	InstmntMon     string `json:"instmntMon"`
	InstmntType    string `json:"instmntType"`
	IssuBankCd     string `json:"issuBankCd"`
	IssuBankNm     string `json:"issuBankNm"`
	MRefNo         string `json:"mRefNo"`
	MerchantToken  string `json:"merchantToken"`
	MitraCd        string `json:"mitraCd"`
	PayMethod      string `json:"payMethod"`
	PreauthToken   string `json:"preauthToken"`
	ReceiptCode    string `json:"receiptCode"`
	RecurringToken string `json:"recurringToken"`
	ReferenceNo    string `json:"referenceNo"`
	ResultCd       string `json:"resultCd"`
	ResultMsg      string `json:"resultMsg"`
	TXid           string `json:"tXid"`
	TimeStamp      string `json:"timeStamp"`
	TransDt        string `json:"transDt"`
	TransTm        string `json:"transTm"`
}

func Payment(iMid, merchantKey, payment_endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			TXid         string `json:"trxID"`
			ReferenceNo  string `json:"noRef"`
			Amt          string `json:"amount"`
			CardNo       string `json:"CardNo"`
			CardExp      string `json:"CardExp"`
			CardCvv      string `json:"CardCvv"`
			CardHolderNm string `json:"CardHolderNm"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		timestampTrx := time.Now().Format("20060102150405")
		merchantData := timestampTrx + iMid + body.ReferenceNo + body.Amt + merchantKey
		merchantToken := Helper.GenerateMerchantToken(merchantData)

		data := map[string]string{
			"cardNo":        body.CardNo,
			"cardCvv":       body.CardCvv,
			"cardHolderNm":  body.CardHolderNm,
			"callBackUrl":   "http://merchant.com/callbackUrl",
			"cardExpYymm":   body.CardExp,
			"timeStamp":     timestampTrx,
			"merchantToken": merchantToken,
			"tXid":          body.TXid,
		}
		url := "https://dev.nicepay.co.id/" + payment_endpoint
		a, _ := Helper.Request().
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetFormData(data).
			Post(url)
		_ = a.String()

		html := a.String()

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			fmt.Println("Error parsing HTML:", err)
			return
		}

		var response ResponsePayment

		doc.Find("input[type=hidden]").Each(func(i int, s *goquery.Selection) {
			name, exists := s.Attr("name")
			if exists {
				value, _ := s.Attr("value")
				switch name {
				case "resultCd":
					response.ResultCd = value
				case "resultMsg":
					response.ResultMsg = value
				case "tXid":
					response.TXid = value
				case "referenceNo":
					response.ReferenceNo = value
				case "amt":
					response.Amt = value
				case "transDt":
					response.TransDt = value
				case "transTm":
					response.TransTm = value
				case "description":
					response.Description = value
				case "issuBankCd":
					response.IssuBankCd = value
				case "acquBankCd":
					response.AcquBankCd = value
				case "cardNo":
					response.CardNo = value
				case "currency":
					response.Currency = value
				}
			}
		})

		if response.ResultCd == "0000" && response.ResultMsg == "SUCCESS" {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": response.ResultCd,
				"message":    response.ResultMsg,
				"details": gin.H{
					"cardNo":     response.CardNo,
					"currency":   response.Currency,
					"amount":     response.Amt,
					"trxID":      response.TXid,
					"noRef":      response.ReferenceNo,
					"transDate":  response.TransDt,
					"transTime":  response.TransTm,
					"issuBankCd": response.IssuBankCd,
					"acquBankCd": response.AcquBankCd,
				},
			})
		} else if response.ResultCd == "9127" {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": response.ResultCd,
				"message":    response.ResultMsg,
				"details": gin.H{
					"cardNo":    response.CardNo,
					"currency":  response.Currency,
					"amount":    response.Amt,
					"noRef":     response.ReferenceNo,
					"transDate": response.TransDt,
					"transTime": response.TransTm,
				},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": response.ResultCd,
				"message":    response.ResultMsg,
			})
		}

	}
}
