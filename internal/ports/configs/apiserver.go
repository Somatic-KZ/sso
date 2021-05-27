package configs

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

const (
	defaultTokenTTLInMin         = 120
	defaultRefreshTokenTTLInDays = 30
)

type APIServer struct {
	DSName string `short:"n" long:"ds" env:"DATASTORE" description:"DataStore name (format: mongo/null)" required:"false" default:"mongo"`
	DSDB   string `short:"d" long:"ds-db" env:"DATASTORE_DB" description:"DataStore database name (format: sso)" required:"false" default:"sso"`
	DSURL  string `short:"u" long:"ds-url" env:"DATASTORE_URL" description:"DataStore URL (format: mongodb://localhost:27017)" required:"false" default:"mongodb+srv://root:somatic123@clusterdiploma.yytdj.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"`

	ListenAddr string `short:"l" long:"listen" env:"LISTEN" description:"Listen Address (format: :8080|127.0.0.1:8080)" required:"false" default:":8080"`
	BasePath   string `long:"base-path" env:"BASE_PATH" description:"base path of the host" required:"false" default:""`
	FilesDir   string `long:"files-directory" env:"FILES_DIR" description:"Directory where all static files are located" required:"false" default:"./api"`
	CertFile   string `short:"c" long:"cert" env:"CERT_FILE" description:"Location of the SSL/TLS cert file" required:"false" default:""`
	KeyFile    string `short:"k" long:"key" env:"KEY_FILE" description:"Location of the SSL/TLS key file" required:"false" default:""`

	GrpcListenAddr          string `long:"grpc-listen" env:"GRPC_LISTEN" description:"Grpc Listen Address (format: :4000|127.0.0.1:4000)" required:"false" default:":4000"`

	JWTKey                string `long:"jwt-key" env:"JWT_KEY" description:"JWT secret key" required:"false" default:"technodom-secret"`
	TokenTTLInMin         int64  `long:"token-ttl" env:"TOKEN_TTL" description:"Auth token lifetime duration (in min)" required:"false" default:"120"`
	RefreshTokenTTLInDays int64  `long:"refresh-token-ttl" env:"REFRESH_TOKEN_TTL" description:"Refresh token lifetime duration (in days)" required:"false" default:"30"`

	DelegateTokenTTLInMin int64 `long:"delegate-token-ttl" env:"DEL_TOKEN_TTL" description:"Delegate token lifetime duration (in min)" required:"false" default:"10"`

	PromListenAddr string `long:"prom-listen" env:"PROM_LISTEN" description:"Listen Address (format: :9090|127.0.0.1:9090)" required:"false" default:":9090"`
	Prometheus     bool   `long:"with-prom" env:"WITH_PROM" description:"Enable Prometheus metrics" required:"false"`

	RecoverySpamPenalty   int64 `env:"RECOVERY_SPAM_PENALTY" description:"recovery spam penalty (in sec)" required:"false"`
	VerifySpamPenalty     int64 `env:"VERIFY_SPAM_PENALTY" description:"verify spam penalty (in sec)" required:"false"`
	VerifyLongSpamPenalty int64 `env:"VERIFY_LONG_SPAM_PENALTY" description:"verify long spam penalty (in sec)" required:"false"`

	Dbg       bool `long:"dbg" env:"DEBUG" description:"debug mode"`
	IsTesting bool `long:"testing" env:"APP_TESTING" description:"testing mode"`
}

func (srv *APIServer) Parse()  {
	p := flags.NewParser(srv, flags.Default)

	if _, err := p.Parse(); err != nil {
		log.Println("[ERROR] parsing config")
		os.Exit(1)
	}
}

func (srv *APIServer) EnsureDefaults() {
	if srv.TokenTTLInMin < 1 {
		srv.TokenTTLInMin = defaultTokenTTLInMin
	}

	if srv.RefreshTokenTTLInDays < 1 {
		srv.TokenTTLInMin = defaultRefreshTokenTTLInDays
	}
}
