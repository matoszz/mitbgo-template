package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"entgo.io/ent/dialect"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	ent "github.com/datumforge/go-template/internal/ent/generated"
)

const (
	defaultListenAddr = ":17608"
)

var (
	enablePlayground bool
	serveDevMode     bool
	mw               string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the example Graph API",
	RunE: func(cmd *cobra.Command, args []string) error {
		return serve(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Bool("debug", false, "enable server debug")
	viperBindFlag("server.debug", serveCmd.Flags().Lookup("debug"))

	serveCmd.Flags().String("listen", defaultListenAddr, "address to listen on")
	viperBindFlag("server.listen", serveCmd.Flags().Lookup("listen"))

	serveCmd.Flags().Duration("shutdown-grace-period", 5*time.Second, "server shutdown grace period")
	viperBindFlag("server.shutdown-grace-period", serveCmd.Flags().Lookup("shutdown-grace-period"))

	serveCmd.Flags().String("listen", "0.0.0.0:3001", "address to listen on")
	viperBindFlag("api.listen", serveCmd.Flags().Lookup("listen"))

	// only available as a CLI arg because these should only be used in dev environments
	serveCmd.Flags().BoolVar(&serveDevMode, "dev", false, "dev mode: enables playground")
	serveCmd.Flags().BoolVar(&enablePlayground, "playground", false, "enable the graph playground")
}

func serve(ctx context.Context) error {
	if serveDevMode {
		enablePlayground = true
	}

	cOpts := []ent.Option{}

	if viper.GetBool("debug") {
		cOpts = append(cOpts,
			ent.Log(logger.Named("ent").Debugln),
			ent.Debug(),
		)
	}

	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1", cOpts...)
	if err != nil {
		logger.Error("failed opening connection to sqlite", zap.Error(err))
		return err
	}
	defer client.Close()

	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx); err != nil {
		logger.Errorf("failed creating schema resources", zap.Error(err))
		return err
	}

	// TODO (sfunk): auth middleware

	srv := echo.New()
	if err != nil {
		logger.Error("failed to create server", zap.Error(err))
	}

	r := api.NewResolver(client, logger.Named("resolvers"))
	handler := r.Handler(enablePlayground)

	// srv.AddHandler(handler)

	if err := runWithContext(ctx); err != nil {
		logger.Error("failed to run server", zap.Error(err))
	}

	return err
}

// RunWithContext listens and serves the echo server on the configured address.
// See ServeWithContext for more details.
func runWithContext(ctx context.Context) error {
	listener, err := net.Listen("tcp", viper.GetString("server.listen"))
	if err != nil {
		return err
	}

	defer listener.Close() //nolint:errcheck // No need to check error.

	return serveWithContext(ctx, listener)
}

// serveWithContext serves an http server on the provided listener.
// Serve blocks until SIGINT or SIGTERM are signalled,
// or if the http serve fails.
// A graceful shutdown will be attempted
func serveWithContext(ctx context.Context, listener net.Listener) error {
	logger := logger.With(zap.String("address", listener.Addr().String()))

	logger.Info("starting server")

	srv := &http.Server{
		Handler: handler(),
	}

	var (
		exit = make(chan error, 1)
		quit = make(chan os.Signal, 2) //nolint:gomnd
	)

	// Serve in a go routine.
	// If serve returns an error, capture the error to return later.
	go func() {
		if err := srv.Serve(listener); err != nil {
			exit <- err

			return
		}

		exit <- nil
	}()

	// close server to kill active connections.
	defer srv.Close() //nolint:errcheck // server is being closed, we'll ignore this.

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var err error

	select {
	case err = <-exit:
		return err
	case sig := <-quit:
		logger.Warn(fmt.Sprintf("%s received, server shutting down", sig.String()))
	case <-ctx.Done():
		logger.Warn("context done, server shutting down")

		// Since the context has already been canceled, the server would immediately shutdown.
		// We'll reset the context to allow for the proper grace period to be given.
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, viper.GetDuration("server.shutdown-grace-period"))
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown timed out", zap.Error(err))

		return err
	}

	return nil
}

// handler returns a new http.Handler for serving requests.
func handler() http.Handler {
	engine := echo.New()

	engine.Use(middleware.RequestID())
	engine.Use(middleware.Recover())

	engine.HideBanner = true
	engine.HidePort = true

	engine.Debug = viper.GetBool("debug")

	p := prometheus.NewPrometheus("echo", nil)

	p.Use(engine)

	return engine
}
