package support

import (
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
