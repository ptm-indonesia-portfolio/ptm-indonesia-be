package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"ptm-indonesia/bootstrap"

	"github.com/gofiber/fiber/v3"
)

func main() {
	application, cleanup, err := bootstrap.InitializeHTTPApplication()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "bootstrap application: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	signalContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	listenAddress := application.Config.Address()
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "listen on %s: %v\n", listenAddress, err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(
		os.Stdout,
		"%s running at %s (bind: %s)\n",
		application.Config.App.Name,
		startupURL(application.Config.App.Host, application.Config.App.Port),
		listenAddress,
	)

	if err := application.App.Listener(listener, fiber.ListenConfig{
		GracefulContext:       signalContext,
		ShutdownTimeout:       application.Config.App.ShutdownTimeout,
		DisableStartupMessage: true,
	}); err != nil && signalContext.Err() == nil {
		application.Logger.WithError(err).Error("http server failed")
		os.Exit(1)
	}
}

func startupURL(host, port string) string {
	displayHost := strings.TrimSpace(host)
	if displayHost == "" || displayHost == "0.0.0.0" || displayHost == "::" || displayHost == "[::]" {
		displayHost = "localhost"
	}

	displayHost = strings.TrimPrefix(displayHost, "[")
	displayHost = strings.TrimSuffix(displayHost, "]")

	return fmt.Sprintf("http://%s", net.JoinHostPort(displayHost, port))
}
