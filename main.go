package main

import (
	"arisu.land/tsubaki/graphql"
	"arisu.land/tsubaki/pkg"
	"arisu.land/tsubaki/pkg/infra"
	"arisu.land/tsubaki/pkg/is"
	"arisu.land/tsubaki/pkg/middleware"
	"arisu.land/tsubaki/pkg/util"
	"arisu.land/tsubaki/routers"
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"golang.org/x/xerrors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	version    string
	commitHash string
	buildDate  string
	log        = slog.Make(sloghuman.Sink(os.Stdout))
)

func init() {
	// bit of warnings for now x3
	if is.Root() {
		log.Warn(context.Background(), "It is not recommended to run Tsubaki with `sudo` or under root.")
	}

	if is.Docker() {
		log.Warn(context.Background(), "Make sure to backup your projects (if using filesystem storage) under a volume!")
	}

	if is.Kubernetes() {
		log.Warn(context.Background(), "Make sure to backup your projects under a PVC!")
	}
}

func main() {
	pkg.SetVersion(version, commitHash, buildDate)
	util.PrintBanner(version, commitHash, buildDate)

	container, err := infra.NewContainer()
	if err != nil {
		log.Error(
			context.Background(),
			"An error occurred while initializing container:",
			slog.Error(xerrors.Errorf("%w", err)))

		os.Exit(1)
	}

	log.Info(context.Background(), "Parsing GraphQL schema from ./schema.gql")
	gql := graphql.NewGraphQLManager(container)

	if err := gql.GenerateSchema(); err != nil {
		log.Fatal(context.Background(), "Unable to parse GraphQL schema:", slog.Error(xerrors.Errorf("%v", err)))
		os.Exit(1)
	}

	log.Info(context.Background(), "Starting up GraphQL server...")

	router := chi.NewRouter()
	router.Use(middleware.LogMiddleware)
	router.Use(middleware.NewErrorHandler(container).Serve)
	router.Mount("/", routers.NewMainRouter(container))
	router.Mount("/health", routers.NewHealthRouter())
	router.Mount("/graphql", routers.NewGraphQLRouter(container, gql))
	router.Mount("/metrics", routers.NewMetricsRouter())
	router.Mount("/version", routers.NewVersionRouter())

	addr := fmt.Sprintf("%s:%d", container.Config.Host, container.Config.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Listen for syscall signals so Docker can properly destroy the server
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Run the server
		log.Info(context.Background(), fmt.Sprintf("Now listening on address: %s", addr))

		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(context.Background(), "Error has occurred while trying to listen to server:", slog.Error(xerrors.Errorf("%v", err)))
		}
	}()

	<-sigint

	log.Warn(context.Background(), "Closing off GraphQL server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// Wait for connections to die off
	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Fatal(context.Background(), "Graceful shutdown timed out! Forcing exit!")
		}
	}()

	defer func() {
		// Now shutdown the container
		err = container.Close()
		if err != nil {
			log.Fatal(context.Background(), "Unable to cleanup container:", slog.Error(xerrors.Errorf("%v", err)))
		}

		// Cancel the shutdown hook
		cancel()
	}()

	// Now we kill the server ^w^
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(context.Background(), "Unable to shutdown server:", slog.Error(xerrors.Errorf("%v", err)))
		os.Exit(1)
	}
}
