package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serveREST(port string) error {
	srv := http.Server{
		//!add TLS config
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	shutdown_signal := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		slog.Info(fmt.Sprintf("shutdown signal: %s", s.String()))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdown_signal <- err
		}

		slog.Info("waiting for background tasks to finish")

		app.wait_group.Wait()
		shutdown_signal <- nil
	}()

	slog.Info(fmt.Sprintf("Starting REST server on port %s", port))
	err := srv.ListenAndServe()
	if err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			return nil
		default:
			return err
		}
	}

	err = <-shutdown_signal
	if err != nil {
		return err
	}

	return nil
}
