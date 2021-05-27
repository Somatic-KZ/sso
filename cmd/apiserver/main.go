package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/JetBrainer/sso/internal/adapters/database"
	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain/manager/auth"
	monitoring2 "github.com/JetBrainer/sso/internal/domain/manager/monitoring"
	"github.com/JetBrainer/sso/internal/domain/models"
	"github.com/JetBrainer/sso/internal/domain/validation"
	"github.com/JetBrainer/sso/internal/ports/configs"
	"github.com/JetBrainer/sso/internal/ports/grpc"
	"github.com/JetBrainer/sso/internal/ports/grpc/resources"
	"github.com/JetBrainer/sso/internal/ports/http"
	"github.com/JetBrainer/sso/internal/ports/monitoring"
	"github.com/JetBrainer/sso/pkg/logger"
	"github.com/JetBrainer/sso/pkg/signal"
	"golang.org/x/sync/errgroup"
)

var (
	version = "unknown"
)

func main() {
	fmt.Printf("sso %s\n", version)

	opts := new(configs.APIServer)
	opts.EnsureDefaults()
	opts.Parse()

	logger.Setup("DEBUG")

	appCtx, appCtxCancel := context.WithCancel(context.Background())
	defer appCtxCancel()

	go signal.CatchTermination(appCtxCancel)

	ds, err := setupDatabase(opts)
	if err != nil {
		log.Println(err)
		return
	}
	defer ds.Close()

	metrics := setupMonitoring(appCtx, opts)
	monitoringManager := monitoring2.New(ds, metrics)
	verifyMan := resources.NewVerify(ds)

	authManager := auth.New(
		ds,
		[]byte(opts.JWTKey),
		time.Duration(opts.TokenTTLInMin)*time.Minute,
		time.Duration(opts.RefreshTokenTTLInDays)*time.Hour*24,
		time.Duration(opts.DelegateTokenTTLInMin)*time.Minute,
	)
	verificationManager := authManager.VerificationManager()
	if opts.VerifySpamPenalty > 0 {
		verificationManager.WithSpamPenalty(time.Duration(opts.VerifySpamPenalty) * time.Second)
	}

	recoveryManager := authManager.RecoveryManager()
	if opts.RecoverySpamPenalty > 0 {
		recoveryManager.WithSpamPenalty(time.Duration(opts.RecoverySpamPenalty) * time.Second)
	}

	servers, serversCtx := errgroup.WithContext(appCtx)

	if opts.IsTesting {
		log.Printf("[INFO] the service is running in test mode")
		authManager.Testing()
	}

	httpSrv := http.NewAPIServer(
		serversCtx,
		opts,
		http.WithAuthManager(authManager),
		http.WithMonitoringManager(monitoringManager),
		http.WithValidator(validation.New(authManager.Users())),
		http.WithVersion(version),
	)
	servers.Go(func() error {
		if err := httpSrv.Run(); err != nil {
			return errors.New(fmt.Sprintf("HTTP server: %v", err))
		}

		httpSrv.Wait()
		return nil
	})

	grpcSrv := grpc.NewAPIServer(
		serversCtx,
		opts,
		grpc.WithAuthManager(authManager),
		grpc.WithVersion(version),
		grpc.WithVerifyManager(verifyMan),
		grpc.WithDatastore(ds),
	)
	servers.Go(func() error {
		if err := grpcSrv.Run(); err != nil {
			return errors.New(fmt.Sprintf("GRPC server: %v", err))
		}

		grpcSrv.Wait()
		return nil
	})

	if err := servers.Wait(); err != nil {
		log.Printf("[INFO] process terminated, %s", err)
		return
	}
}

func setupDatabase(opts *configs.APIServer) (drivers.DataStore, error) {
	ds, err := database.New(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataBaseName:  opts.DSDB,
		DataStoreName: opts.DSName,
	})
	if err != nil {
		return nil, err
	}

	if err := ds.Connect(); err != nil {
		errText := fmt.Sprintf("[ERROR] cannot connect to datastore %s: %v", opts.DSName, err)
		return nil, errors.New(errText)
	}

	fmt.Println("Connected to", ds.Name())

	return ds, nil
}

func setupMonitoring(ctx context.Context, opts *configs.APIServer) *models.Metrics {
	promSrv := monitoring.NewPrometheusSrv(ctx, opts.PromListenAddr)

	if opts.Prometheus {
		go func() {
			if err := promSrv.Run(); err != nil {
				log.Printf("[ERROR] Prometheus is not working: %s", err)
			}
		}()
	}

	return promSrv.Metrics()
}
