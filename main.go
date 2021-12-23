package main

import (
	"arisu.land/tsubaki/graphql"
	"arisu.land/tsubaki/pkg"
	"arisu.land/tsubaki/pkg/infra"
	"arisu.land/tsubaki/pkg/is"
	"arisu.land/tsubaki/pkg/logging"
	"arisu.land/tsubaki/pkg/middleware"
	"arisu.land/tsubaki/pkg/ratelimiter"
	"arisu.land/tsubaki/pkg/sessions"
	"arisu.land/tsubaki/pkg/util"
	"arisu.land/tsubaki/routers"
	"arisu.land/tsubaki/routers/integrations"
	"context"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	version         string
	commitHash      string
	buildDate       string
	profilerEnabled = flag.Bool("p", false, "enables profiling")
)

func init() {
	pkg.SetVersion(version, commitHash, buildDate)
	logrus.SetFormatter(&logging.Formatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	flag.Parse()
}

func main() {
	util.PrintBanner(version, commitHash, buildDate)

	if profilerEnabled != nil && *profilerEnabled {
		logrus.Info("Profiling is now enabled on the server!")
	}

	// bit of warnings for now x3
	if is.Root() {
		logrus.Warn("Make sure to run Tsubaki without using `sudo`, admin privileges on Windows, or under the `root` account.")
	}

	if is.Docker() {
		logrus.Warn("Make sure you use the `-v` flag when running the Tsubaki Docker image if you're using the Filesystem Storage Provider.")
	}

	if is.Kubernetes() {
		logrus.Warn("Make sure you create a persisted volume claim in your deployment or statefulset.")
	}

	container, err := infra.NewContainer()
	if err != nil {
		logrus.WithField("step", "container init").Errorf("Unable to initialize container: %v", err)
		os.Exit(1)
	}

	logrus.WithField("step", "graphql").Info("Parsing GraphQL schema from ./schema.gql!")
	gql := graphql.NewGraphQLManager(container)

	if err := gql.GenerateSchema(); err != nil {
		logrus.WithField("step", "graphql").Errorf("Unable to parse GraphQL schema: %v", err)
		os.Exit(1)
	}

	logrus.WithField("step", "server").Info("Starting up server...")

	rl := ratelimiter.NewRatelimiter(container.Redis)
	router := chi.NewRouter()
	sessions.NewSessionManager(container.Redis, container.Database)

	router.Use(rl.Middleware)
	router.Use(middleware.Headers)
	router.Use(middleware.LogMiddleware)
	router.Use(sessions.Sessions.Middleware)
	router.Use(middleware.NewErrorHandler(container).Serve)
	router.Mount("/", routers.NewMainRouter(container))
	router.Mount("/health", routers.NewHealthRouter())
	router.Mount("/graphql", routers.NewGraphQLRouter(container, gql))
	router.Mount("/metrics", routers.NewMetricsRouter())
	router.Mount("/version", routers.NewVersionRouter())
	router.Mount("/integrations", integrations.NewIntegrationsRouter())

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
		logrus.WithField("step", "server").Infof("Running Tsubaki under address '%s'", addr)

		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Unable to run server: %s", err)
		}
	}()

	<-sigint

	logrus.WithField("step", "shutdown").Warn("Closing off server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// Wait for connections to die off
	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			logrus.WithField("step", "shutdown requests").Warn("Reached deadline to close off incoming requests...")
		}
	}()

	defer func() {
		// Now shutdown the container
		err = container.Close()
		if err != nil {
			logrus.WithField("step", "shutdown container").Errorf("Unable to close container resources: %v", err)
		}

		// Cancel the shutdown hook
		cancel()
	}()

	// Now we kill the server ^w^
	if err := server.Shutdown(shutdownCtx); err != nil {
		logrus.WithField("step", "shutdown server").Errorf("Unable to shutdown server: %v", err)
		os.Exit(1)
	}
}
