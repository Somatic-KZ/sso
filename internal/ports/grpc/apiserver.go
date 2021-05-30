package grpc

import (
	"context"
	"log"
	"net"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain/manager/auth"
	"github.com/JetBrainer/sso/internal/ports/configs"
	"github.com/JetBrainer/sso/internal/ports/grpc/resources"
	"github.com/Somatic-KZ/sso-client/protobuf"
	"google.golang.org/grpc"
)

type APIServer struct {
	Address     string
	IsTesting   bool
	authManager *auth.Manager

	server    protobuf.SSOServer
	verifyMan *resources.Verify

	db drivers.DataStore

	idleConnsClosed chan struct{}
	masterCtx       context.Context
	version         string
}

func NewAPIServer(ctx context.Context, config *configs.APIServer, opts ...APIServerOption) *APIServer {
	srv := &APIServer{
		Address:         config.GrpcListenAddr,
		IsTesting:       config.IsTesting,
		idleConnsClosed: make(chan struct{}), // способ определить незавершенные соединения
		masterCtx:       ctx,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (srv *APIServer) Run() error {
	listener, err := net.Listen("tcp", srv.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// Сначала монтируем ресуры
	// (всегда должны использоваться сгенерированные в *.pb.go файлах интерфейсы)
	srv.setupResources()
	// Затем регистрируем их
	srv.registerServices(grpcServer)

	go srv.GracefulShutdown(grpcServer)
	log.Printf("[INFO] serving GRPC on \"%s\"", srv.Address)
	if err := grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

// setupResources монтирует необходимые grpc-ресурсы для
// обработки клиентских запросов
func (srv *APIServer) setupResources() {
	srv.verifyMan = resources.NewVerify(srv.db)
}

// registerServices регистрирует все необходимые сервисы для работы grpc сервера
func (srv *APIServer) registerServices(grpcServer *grpc.Server) {
	protobuf.RegisterSSOServer(grpcServer, srv.verifyMan)
}

// GracefulShutdown обрабатывает все оставшиеся соединения до остановки
func (srv *APIServer) GracefulShutdown(grpcServer *grpc.Server) {
	<-srv.masterCtx.Done()
	log.Printf("[INFO] shutting down gRPC server")
	grpcServer.GracefulStop()
	close(srv.idleConnsClosed)
}

// Wait ожидает момента завершения обработки всех соединений.
func (srv *APIServer) Wait() {
	<-srv.idleConnsClosed
	log.Println("[INFO] gRPC server has processed all idle connections")
}
