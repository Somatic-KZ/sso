package monitoring

import (
	"context"
	"log"
	"net/http"

	"github.com/JetBrainer/sso/internal/domain/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusSrv struct {
	masterCtx context.Context
	Address   string
	metrics   *models.Metrics
}

func NewPrometheusSrv(masterCtx context.Context, addr string) *PrometheusSrv {
	return &PrometheusSrv{
		masterCtx: masterCtx,
		Address:   addr,
		metrics:   metrics(),
	}
}

func metrics() *models.Metrics {
	signInByPhoneRequests := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_in_by_phone_requests_total",
		Help: "Total number of signin by phone requests.",
	})
	signInByPhoneErrorInvalidLoginOrPasswordRequests := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_in_by_phone_error_invalid_login_or_password_total",
		Help: "Total number of signin errors 'Invalid login or password'.",
	})
	signInByPhoneErrorRecoveryRequiredRequests := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_in_by_phone_error_recovery_required_total",
		Help: "Total number of signin errors 'Recovery required'.",
	})
	signInByPhoneUserNotFoundErrors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_in_by_phone_error_user_not_found_total",
		Help: "Total number of signin errors 'User not found'.",
	})
	signInByPhoneVerificationRequiredErrors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_in_by_phone_error_verification_required_total",
		Help: "Total number of signin errors 'Verification required'.",
	})
	signInByPhoneInternalServerErrors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_in_by_phone_error_internal_server_error_total",
		Help: "Total number of signin errors 'Internal server error'.",
	})
	signInByPhoneInvalidNumberFormatErrors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_in_by_phone_error_invalid_number_format_error_total",
		Help: "Total number of signin errors 'Invalid Number Format'.",
	})
	signInByPhoneSuccessfulOperation := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_in_by_phone_successful_operation_total",
		Help: "Total number of successful operations.",
	})
	signInValidationErrors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_in_by_phone_validation_error_total",
		Help: "Total number of signin errors 'Validation error'.",
	})

	signUPInternalServerErrors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_up_error_internal_server_error_total",
		Help: "Total number of signup errors 'Internal server error'.",
	})
	signUpValidationErrors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_up_error_validation_error_total",
		Help: "Total number of signup errors 'Validation error'.",
	})
	signUpByPhoneRequests := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_up_requests_total",
		Help: "Total number of signup by phone requests.",
	})
	signUpSuccessfulOperation := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_sign_up_successful_operation_total",
		Help: "Total number of successful operations.",
	})

	refreshRequests := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_refresh_requests_total",
		Help: "Total number of refresh requests.",
	})
	refreshInternalServerErrors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_refresh_internal_server_error_total",
		Help: "Total number of refresh errors 'Internal server error'.",
	})
	refreshSuccessfulOperation := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_refresh_successful_operation_total",
		Help: "Total number of successful operations.",
	})
	refreshInvalidTokenErrors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sso_refresh_invalid_token_error_total",
		Help: "Total number of refresh errors 'Invalid token error'.",
	})

	return &models.Metrics{
		SignInByPhoneRequests:                     &signInByPhoneRequests,
		SignInByPhoneInvalidLoginOrPasswordErrors: &signInByPhoneErrorInvalidLoginOrPasswordRequests,
		SignInByPhoneRecoveryRequiredErrors:       &signInByPhoneErrorRecoveryRequiredRequests,
		SignInByPhoneUserNotFoundErrors:           &signInByPhoneUserNotFoundErrors,
		SignInByPhoneVerificationRequiredErrors:   &signInByPhoneVerificationRequiredErrors,
		SignInByPhoneInternalServerErrors:         &signInByPhoneInternalServerErrors,
		SignInByPhoneInvalidNumberFormatErrors:    &signInByPhoneInvalidNumberFormatErrors,
		SignInByPhoneSuccessfulOperation:          &signInByPhoneSuccessfulOperation,
		SignInValidationErrors:                    &signInValidationErrors,

		SignUpByPhoneRequests:      &signUpByPhoneRequests,
		SignUpValidationErrors:     &signUpValidationErrors,
		SignUpInternalServerErrors: &signUPInternalServerErrors,
		SignUpSuccessfulOperation:  &signUpSuccessfulOperation,

		RefreshRequests:             &refreshRequests,
		RefreshInternalServerErrors: &refreshInternalServerErrors,
		RefreshSuccessfulOperation:  &refreshSuccessfulOperation,
		RefreshInvalidTokenErrors:   &refreshInvalidTokenErrors,
	}
}

func (p *PrometheusSrv) Metrics() *models.Metrics {
	return p.metrics
}

func (p *PrometheusSrv) Run() error {
	prometheus.MustRegister(*p.metrics.SignInByPhoneRequests)
	prometheus.MustRegister(*p.metrics.SignInByPhoneInvalidLoginOrPasswordErrors)
	prometheus.MustRegister(*p.metrics.SignInByPhoneRecoveryRequiredErrors)
	prometheus.MustRegister(*p.metrics.SignInByPhoneUserNotFoundErrors)
	prometheus.MustRegister(*p.metrics.SignInByPhoneVerificationRequiredErrors)
	prometheus.MustRegister(*p.metrics.SignInByPhoneInternalServerErrors)
	prometheus.MustRegister(*p.metrics.SignInByPhoneInvalidNumberFormatErrors)
	prometheus.MustRegister(*p.metrics.SignInValidationErrors)
	prometheus.MustRegister(*p.metrics.SignInByPhoneSuccessfulOperation)

	prometheus.MustRegister(*p.metrics.SignUpInternalServerErrors)
	prometheus.MustRegister(*p.metrics.SignUpByPhoneRequests)
	prometheus.MustRegister(*p.metrics.SignUpValidationErrors)
	prometheus.MustRegister(*p.metrics.SignUpSuccessfulOperation)

	prometheus.MustRegister(*p.metrics.RefreshRequests)
	prometheus.MustRegister(*p.metrics.RefreshInternalServerErrors)
	prometheus.MustRegister(*p.metrics.RefreshSuccessfulOperation)
	prometheus.MustRegister(*p.metrics.RefreshInvalidTokenErrors)

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("[INFO] serving Prometheus HTTP on %s", p.Address)
	if err := http.ListenAndServe(p.Address, nil); err != nil {
		log.Printf("[ERROR] ошибка обслуживания метрик Prometheus на %s: %v", p.Address, err)
		return err
	}

	return nil
}
