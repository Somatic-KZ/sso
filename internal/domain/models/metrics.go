package models

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	SignInByPhoneRequests                     *prometheus.Counter
	SignInByPhoneInvalidLoginOrPasswordErrors *prometheus.Counter
	SignInByPhoneInvalidNumberFormatErrors    *prometheus.Counter
	SignInByPhoneRecoveryRequiredErrors       *prometheus.Counter
	SignInByPhoneUserNotFoundErrors           *prometheus.Counter
	SignInByPhoneVerificationRequiredErrors   *prometheus.Counter
	SignInByPhoneInternalServerErrors         *prometheus.Counter
	SignInByPhoneSuccessfulOperation          *prometheus.Counter
	SignInValidationErrors                    *prometheus.Counter

	SignUpByPhoneRequests      *prometheus.Counter
	SignUpInternalServerErrors *prometheus.Counter
	SignUpValidationErrors     *prometheus.Counter
	SignUpSuccessfulOperation  *prometheus.Counter

	RefreshRequests             *prometheus.Counter
	RefreshInternalServerErrors *prometheus.Counter
	RefreshSuccessfulOperation  *prometheus.Counter
	RefreshInvalidTokenErrors   *prometheus.Counter
}
