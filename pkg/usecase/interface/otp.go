package interfaces

import (
	"context"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type OtpUseCase interface {
	TwilioSendOTP(context.Context, string) (string, error)
	TwilioVerifyOTP(context.Context, utils.OtpVerify) (domain.OtpSession, error)
}
