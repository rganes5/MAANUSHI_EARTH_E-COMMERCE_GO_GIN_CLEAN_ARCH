package utils

import (
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
	// fmt.Println(err)
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
		return err
	} else if *resp.Status == "approved" {
		return nil
	}

	return nil
}
