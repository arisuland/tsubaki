// ☔ Arisu: Translation made with simplicity, yet robust.
// Copyright (C) 2020-2022 Noelware
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package server

import (
	"arisu.land/tsubaki/graphql"
	"arisu.land/tsubaki/internal"
	"arisu.land/tsubaki/pkg"
	"arisu.land/tsubaki/pkg/ratelimit"
	"arisu.land/tsubaki/pkg/sessions"
	"arisu.land/tsubaki/server/middleware"
	"arisu.land/tsubaki/server/routes"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start(path string) error {
	if internal.Root() {
		logrus.Warn("Make sure you are running Tsubaki using `sudo`, admin privileges on Windows, or under the `root` account!")
	}

	if internal.Docker() {
		logrus.Warn("Make sure you have a volume mounted if you're testing locally for projects!")
	}

	if internal.Kubernetes() {
		logrus.Warn("Make sure your volume is persisted using a PVC!")
	}

	err := pkg.NewContainer(path)
	if err != nil {
		return err
	}

	logrus.WithField("step", "graphql").Info("Parsing GraphQL schema...")
	gql := graphql.NewGraphQLManager()

	if err := gql.GenerateSchema(); err != nil {
		return err
	}

	logrus.WithField("step", "server").Info("Starting up HTTP server!")
	rl := ratelimit.NewRatelimiter(pkg.GlobalContainer.Redis)
	router := chi.NewRouter()
	sesh := sessions.NewSessionManager(pkg.GlobalContainer.Redis, pkg.GlobalContainer.Prisma)

	router.Use(rl.Middleware)
	router.Use(sesh.Middleware)
	router.Use(middleware.Logging)
	router.Use(middleware.ErrorReporter)
	router.Mount("/", routes.NewMainRouter())
	router.Mount("/ping", routes.NewPingRouter())
	router.Mount("/metrics", routes.NewMetricsRouter())
	router.Mount("/version", routes.NewVersionRouter())
	router.Mount("/graphql", routes.NewGraphQLRouter(pkg.GlobalContainer, gql))

	addr := fmt.Sprintf("%s:%d", pkg.GlobalContainer.Config.Host, pkg.GlobalContainer.Config.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		WriteTimeout: 10 * time.Second,
	}

	// Listen for syscall signals so Docker can properly destroy the server
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logrus.WithField("step", "server").Infof("Listening under '%s'!", addr)
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
		// Cache all ratelimits + sessions
		err = rl.Close()
		if err != nil {
			logrus.WithField("step", "cache->ratelimits").Errorf("Unable to cache all ratelimits: %v", err)
		}

		err = sesh.Close()
		if err != nil {
			logrus.WithField("step", "cache->sessions").Errorf("Unable to cache all sessions: %v", err)
		}

		// Shutdown the container
		err = pkg.GlobalContainer.Close()
		if err != nil {
			logrus.WithField("step", "shutdown->container").Errorf("Unable to close resources: %v", err)
		}

		cancel()
	}()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}
