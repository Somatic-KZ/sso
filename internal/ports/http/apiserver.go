package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/JetBrainer/sso/internal/domain/manager/auth"
	"github.com/JetBrainer/sso/internal/domain/manager/monitoring"
	"github.com/JetBrainer/sso/internal/ports/configs"
	"github.com/JetBrainer/sso/internal/ports/http/resources"
	v1 "github.com/JetBrainer/sso/internal/ports/http/resources/auth/v1"
	v12 "github.com/JetBrainer/sso/internal/ports/http/resources/swagger/v1"
	"github.com/JetBrainer/sso/pkg/validation"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

const compressLevel = 5

type APIServer struct {
	Address           string
	CertFile, KeyFile *string
	BasePath          string
	FilesDir          string
	authManager       *auth.Manager
	monitManager      *monitoring.Manager
	validator         *validation.Validator
	idleConnsClosed   chan struct{}
	IsTesting         bool
	version           string // версия приложения
	masterCtx         context.Context
}

func NewAPIServer(ctx context.Context, config *configs.APIServer, opts ...APIServerOption) *APIServer {
	srv := &APIServer{
		masterCtx:       ctx,
		Address:         config.ListenAddr,
		BasePath:        config.BasePath,
		FilesDir:        config.FilesDir,
		idleConnsClosed: make(chan struct{}), // способ определить незавершенные соединения
		IsTesting:       config.IsTesting,
	}

	for _, opt := range opts {
		opt(srv)
	}

	if config.CertFile != "" {
		srv.CertFile = &config.CertFile
	}

	if config.KeyFile != "" {
		srv.KeyFile = &config.KeyFile
	}

	return srv
}

// setupRouter инициализирует HTTP роутер.
// Функция используется для подключения middleware и маппинга ресурсов.
func (srv *APIServer) setupRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NoCache) // no-cache
	//r.Use(middleware.RequestID) // вставляет request ID в контекст каждого запроса
	r.Use(middleware.Logger)    // логирует начало и окончание каждого запроса с указанием времени обработки
	r.Use(middleware.Recoverer) // управляемо обрабатывает паники и выдает stack trace при их возникновении
	r.Use(middleware.RealIP)    // устанавливает RemoteAddr для каждого запроса с заголовками X-Forwarded-For или X-Real-IP
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins(srv.IsTesting),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Mount("/api/v1/auth", v1.NewAuth(srv.authManager,srv.monitManager.Metrics(), srv.validator).Routes())

	// монтируем дополнительные ресурсы
	r.Mount("/version", resources.VersionResource{Version: srv.version}.Routes())
	r.Mount("/health", resources.NewHealth(srv.monitManager).Routes())
	if srv.IsTesting {
		r.Mount("/files", resources.FilesResource{FilesDir: srv.FilesDir}.Routes())
		r.Mount("/swagger", v12.SwaggerResource{FilesPath: "/files", BasePath: srv.BasePath}.Routes())
	}

	return r
}

// getAllowedOrigins возвращает список хостов для C.O.R.S.
func allowedOrigins(testing bool) []string {
	if testing {
		return []string{"*"}
	}

	return []string{}
}

// Run запускает HTTP или HTTPS листенер в зависимости от того как заполнена
// структура HTTPServer{}.
func (srv *APIServer) Run() error {
	const (
		readTimeout  = 5 * time.Second
		writeTimeout = 30 * time.Second
	)

	s := &http.Server{
		Addr:         srv.Address,
		Handler:      chi.ServerBaseContext(srv.masterCtx, srv.setupRouter()),
		ReadTimeout:  readTimeout,  // wait() + tls handshake + req.headers + req.body
		WriteTimeout: writeTimeout, // все что выше + response
	}

	go srv.GracefulShutdown(s)
	log.Printf("[INFO] serving HTTP on \"%s\"", srv.Address)

	var err error
	if srv.CertFile == nil && srv.KeyFile == nil {
		if err = s.ListenAndServe(); err != nil {
			return err
		}
	} else {
		if err = s.ListenAndServeTLS(*srv.CertFile, *srv.KeyFile); err != nil {
			return err
		}
	}

	return nil
}

func (srv *APIServer) GracefulShutdown(httpSrv *http.Server) {
	<-srv.masterCtx.Done()

	if err := httpSrv.Shutdown(context.Background()); err != nil {
		log.Printf("[ERROR] HTTP server Shutdown: %v", err)
	}

	log.Println("[INFO] HTTP server has processed all idle connections")
	close(srv.idleConnsClosed)
}

// Wait ожидает момента завершения обработки всех соединений.
func (srv *APIServer) Wait() {
	<-srv.idleConnsClosed
}
