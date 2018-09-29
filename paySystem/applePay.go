package paySystem

const (
	sandUrl = "https://sandbox.itunes.apple.com/verifyReceipt"
	buyUrl  = "https://buy.itunes.apple.com/verifyReceipt"
)

type iapPayDataReq struct {
	Receipt_data string `json:"receipt-data"`
}

func AppleIAPPay(receipt string, productId string) {

}
