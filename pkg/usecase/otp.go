package usecase

import (
	"context"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/config"
	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
	utils "github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type otpusecase struct {
	OtpRepo interfaces.OtpRepository
}

func NewOtpUseCase(otprepo interfaces.OtpRepository) services.OtpUseCase {
	return &otpusecase{
		OtpRepo: otprepo,
	}
}
func (c *otpusecase) TwilioSendOTP(ctx context.Context, phoneNumber string) (string, error) {
	password := config.GetCofig().AUTHTOKEN
	userName := config.GetCofig().ACCOUNTSID
	seviceSid := config.GetCofig().SERVICESID

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Password: password,
		Username: userName,
	})
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo("+91" + phoneNumber)
	params.SetChannel("sms")
	resp, err := client.VerifyV2.CreateVerification(seviceSid, params)
	if err != nil {
		return *resp.Sid, err
	}
	otpsession := domain.OtpSession{
		OtpID:    *resp.Sid,
		PhoneNum: phoneNumber,
	}
	err1 := c.OtpRepo.SaveOtp(ctx, otpsession)
	if err1 != nil {
		return *resp.Sid, err1
	}
	return *resp.Sid, nil
}

func (c *otpusecase) TwilioVerifyOTP(ctx context.Context, code utils.OtpVerify) (domain.OtpSession, error) {
	//create a twilio client with twilio details
	password := config.GetCofig().AUTHTOKEN
	userName := config.GetCofig().ACCOUNTSID
	seviceSid := config.GetCofig().SERVICESID
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Password: password,
		Username: userName,
	})
	session, err := c.OtpRepo.RetrieveSession(ctx, code.OtpID)
	if err != nil {
		return session, err
	}
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo("+91" + session.PhoneNum)
	params.SetCode(code.Otp)
	resp, err := client.VerifyV2.CreateVerificationCheck(seviceSid, params)
	if err != nil {
		return session, err
	} else if *resp.Status == "approved" {
		return session, nil
	}

	return session, nil
}
