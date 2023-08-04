package Helper

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

func Request() *resty.Request {
	return resty.New().
		RemoveProxy().
		SetDebug(true).
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(15)).
		SetTimeout(180 * time.Second).
		SetDisableWarn(true).
		SetRetryCount(6).
		SetRetryWaitTime(5 * time.Second).
		AddRetryCondition(
			// RetryConditionFunc type is for retry condition function
			// input: non-nil Response OR request execution error
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusServiceUnavailable
			},
		).
		R()
	//EnableTrace()
}

func GenerateMerchantToken(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
