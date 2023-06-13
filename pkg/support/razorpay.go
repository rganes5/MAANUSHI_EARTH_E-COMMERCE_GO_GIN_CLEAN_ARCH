package support

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"

	razorpay "github.com/razorpay/razorpay-go"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/config"
)

func GenerateNewRazorPayOrder(razorPayAmount int, razorPayReceipt string) (razorPayOrderId string, err error) {
	// get razor pay key and secret
	razorPayKey := config.GetCofig().RAZORPAYKEY
	razorPaySecret := config.GetCofig().RAZORPAYSECRET

	//create a razorpay client
	client := razorpay.NewClient(razorPayKey, razorPaySecret)

	data := map[string]interface{}{
		"amount":   razorPayAmount,
		"currency": "INR",
		"receipt":  razorPayReceipt,
	}
	// create an order on razor pay
	body, err := client.Order.Create(data, nil)
	if err != nil {
		return razorPayOrderId, fmt.Errorf("failed to create a razorpay order from support for amount %v", razorPayAmount)
	}
	//The value retrieved from a map is of type interface{} which is a generic type that can hold values of any type.
	//To assign the retrieved value to the razorPayOrderId variable, a type assertion (string) is used.
	razorPayOrderId = body["id"].(string)
	return razorPayOrderId, nil
}

func VerifyRazorPayment(razorpay_payment_id, razorpay_order_id, razorpay_signature string) error {
	// get razor pay key and secret
	razorPayKey := config.GetCofig().RAZORPAYKEY
	razorPaySecret := config.GetCofig().RAZORPAYSECRET

	//verify signature
	data := razorpay_order_id + "|" + razorpay_payment_id
	h := hmac.New(sha256.New, []byte(razorPaySecret))
	_, err := h.Write([]byte(data))
	if err != nil {
		return errors.New("failed to verify signature")
	}

	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(razorpay_signature)) != 1 {
		return errors.New("razorpay signature does not match")
	}

	//Verify Payment Signature
	client := razorpay.NewClient(razorPayKey, razorPaySecret)

	// fetch payment and verify
	payment, err := client.Payment.Fetch(razorpay_payment_id, nil, nil)

	if err != nil {
		return err
	}

	// check payment status
	if payment["status"] != "captured" {
		return fmt.Errorf("failed to verify payment \n razorpay payment with payment_id %v", razorpay_payment_id)
	}

	return nil

}
