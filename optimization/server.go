/*
 * Copyright (c) New Cloud Technologies, Ltd. 2013-2022.
 * Author: Vitaly Isaev <vitaly.isaev@myoffice.team>
 * License: https://github.com/newcloudtechnologies/rbfopt-go/blob/master/LICENSE
 */

package optimization

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

type server struct {
	httpServer *http.Server
	estimator  *costEstimator
	lastError  error
	logger     logr.Logger
}

// Estimate Cost
func (s *server) estimateCostHandler(w http.ResponseWriter, r *http.Request) {
	s.middleware(w, r, s.estimateCost)
}

//nolint:dupl // TODO: need to use more advanced web framework
func (s *server) estimateCost(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, errors.New("invalid method")
	}

	decoder := json.NewDecoder(r.Body)
	request := &estimateCostRequest{}

	err := decoder.Decode(request)
	if err != nil {
		return http.StatusBadRequest, errors.Wrap(err, "json decode")
	}

	response, err := s.estimator.estimateCost(ctx, request)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "estimate cost")
	}

	encoder := json.NewEncoder(w)
	if err = encoder.Encode(response); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "json encode")
	}

	return http.StatusOK, nil
}

// Register report
func (s *server) registerReportHandler(w http.ResponseWriter, r *http.Request) {
	s.middleware(w, r, s.registerReport)
}

//nolint:dupl // TODO: need to use more advanced web framework
func (s *server) registerReport(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) (int, error) {
	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, errors.New("invalid method")
	}

	decoder := json.NewDecoder(r.Body)
	request := &registerReportRequest{}

	err := decoder.Decode(request)
	if err != nil {
		return http.StatusBadRequest, errors.Wrap(err, "json decode")
	}

	response, err := s.estimator.registerReport(ctx, request)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "json encode")
	}

	encoder := json.NewEncoder(w)
	if err = encoder.Encode(response); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "json encode")
	}

	return http.StatusOK, nil
}

type handlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error)

func (s *server) middleware(w http.ResponseWriter, r *http.Request, handler handlerFunc) {
	logger := s.annotateLogger(r)
	ctx := logr.NewContext(r.Context(), logger)

	logger.V(0).Info("request handling started")

	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Error(err, "request body close")
		}
	}()

	statusCode, err := handler(ctx, w, r)
	w.WriteHeader(statusCode)

	if err != nil {
		// cache errors
		s.lastError = err

		logger.Error(err, "request handling finished")
	} else {
		logger.V(0).Info("request handling finished")
	}
}

func (s *server) annotateLogger(r *http.Request) logr.Logger {
	return s.logger.WithValues(
		"url", r.URL,
		"method", r.Method,
		"remote_addr", r.RemoteAddr,
	)
}

func (s *server) quit(ctx context.Context) {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error(err, "http server shutdown")
	}
}

func newServer(logger logr.Logger, endpoint string, estimator *costEstimator) *server {
	handler := http.NewServeMux()

	srv := &server{
		httpServer: &http.Server{
			Addr:    endpoint,
			Handler: handler,
		},
		estimator: estimator,
		logger:    logger,
	}

	handler.HandleFunc("/estimate_cost", srv.estimateCostHandler)
	handler.HandleFunc("/register_report", srv.registerReportHandler)

	go func() {
		if err := srv.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			srv.logger.Error(err, "http server listen and serve")
		}
	}()

	return srv
}
