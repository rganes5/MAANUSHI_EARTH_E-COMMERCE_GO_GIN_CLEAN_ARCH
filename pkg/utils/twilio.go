package utils

import (
	"errors"
	"fmt"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

func TwilioSendOTP(phoneNumber string) (string, error) {
	//create a twilio client with twilio details
	password := config.GetCofig().AUTHTOKEN
	// fmt.Println("password", password)
	fmt.Println(config.GetCofig().AUTHTOKEN)
	userName := config.GetCofig().ACCOUNTSID
	seviceSid := config.GetCofig().SERVICESID
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Password: password,
		Username: userName,
	})
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")
	resp, err := client.VerifyV2.CreateVerification(seviceSid, params)
	// fmt.Println(resp)
	if err != nil {
		return "", err
	}
	return *resp.Sid, nil
}

func TwilioVerifyOTP(phoneNumber string, code string) error {
	//create a twilio client with twilio details
	password := config.GetCofig().AUTHTOKEN
	userName := config.GetCofig().ACCOUNTSID
	seviceSid := config.GetCofig().SERVICESID
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Password: password,
		Username: userName,
	})

	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(seviceSid, params)
	if err != nil {
		// fmt.Println("1response:", resp.Status)
		return errors.New("invalid otp")
	} else if *resp.Status == "approved" {
		// fmt.Println("2response:", resp.Status)
		return nil
	} else {
		// fmt.Println("3response:", resp.Status)
		return errors.New("invalid otp")
	}
}
