package config

import "time"

type AuthConfig struct {
	// callback url to frontend
	CallbackURL string `mapstructure:"AUTH_CALLBACK_URL"`

	// for creating access and refresh tokens
	Secret string `mapstructure:"AUTH_TOKEN_SECRET"`
	// access token time to live
	AccessTokenExpDur time.Duration `mapstructure:"AUTH_ACCESS_TOKEN_EXPIRE_DURATION"`
	// refresh token time to live
	RefreshTokenExpDur time.Duration `mapstructure:"AUTH_REFRESH_TOKEN_EXPIRE_DURATION"`

	// verification token time to live
	VerifyTokenExpDur time.Duration `mapstructure:"AUTH_VERIFY_TOKEN_EXPIRE_DURATION"`
	// verification token request cooldown before new request
	VerifyTokenCooldown time.Duration `mapstructure:"AUTH_VERIFY_TOKEN_COOLDOWN"`

	// reset token time to live
	ResetTokenExpDur time.Duration `mapstructure:"AUTH_RESET_TOKEN_EXPIRE_DURATION"`
	// reset token request cooldown tiem before new request
	ResetTokenCooldown time.Duration `mapstructure:"AUTH_RESET_TOKEN_COOLDOWN"`

	// account deletion token time to live
	DeleteTokenExpDur time.Duration `mapstructure:"AUTH_DELETE_TOKEN_EXPIRE_DURATION"`
	// account deletion token request cooldown before new request
	DeleteTokenCooldown time.Duration `mapstructure:"AUTH_DELETE_TOKEN_COOLDOWN"`
}
