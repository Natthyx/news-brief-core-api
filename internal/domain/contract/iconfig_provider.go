package contract

import "time"

type IConfigProvider interface {
	GetSendActivationEmail() bool
	GetAppBaseURL() string
	GetFrontendBaseURL() string
	GetRefreshTokenExpiry() time.Duration
	GetPasswordResetTokenExpiry() time.Duration
	GetEmailVerificationTokenExpiry() time.Duration
	GetAIServiceAPIKey() string
}
