package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/flum1025/tweam/internal/aws"
	"github.com/flum1025/tweam/internal/config"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
)

type server struct {
	config    *config.Config
	sqsClient *aws.SQSClient
}

func NewServer(
	config *config.Config,
) (*server, error) {
	sqsClient, err := aws.NewSQSClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQSClient: %w", err)
	}

	return &server{
		config:    config,
		sqsClient: sqsClient,
	}, nil
}

func (s *server) Run(port int) error {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	s.registerRoutes(router)

	log.Println(fmt.Sprintf("[INFO] listening :%d", port))

	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		router,
	); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *server) registerRoutes(router chi.Router) {
	router.Post("/webhook/worker", s.worker)
	router.Post("/webhook/twistributer", s.twistributer)
}

func parse(body io.Reader, params interface{}) ([]byte, error) {
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body: %w", err)
	}

	if err = json.Unmarshal(buf, params); err != nil {
		return nil, fmt.Errorf("failed to parse request body: %w", err)
	}

	return buf, validator.New().Struct(params)
}
