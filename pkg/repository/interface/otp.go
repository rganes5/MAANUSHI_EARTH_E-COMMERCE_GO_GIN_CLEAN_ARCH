package interfaces

import (
	"context"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
)

type OtpRepository interface {
	SaveOtp(ctx context.Context, otpsession domain.OtpSession) error
	RetrieveSession(context.Context, string) (domain.OtpSession, error)
}
