package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/daniele/gestione-caselo/backend/eventform/internal/graphql"
	"github.com/daniele/gestione-caselo/backend/eventform/internal/queue"
	appconfig "github.com/daniele/gestione-caselo/backend/internal/config"
)

func corsMiddlewareWithConfig(origins string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origins)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	ctx := context.Background()
	origins := appconfig.GetEnvVariable("ALLOWED_ORIGINS")
	queueURL := appconfig.GetEnvVariable("SQS_QUEUE_URL")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Configure SQS client with custom endpoint if AWS_ENDPOINT_URL is set
	var sqsClient *sqs.Client
	if endpointURL := os.Getenv("AWS_ENDPOINT_URL"); endpointURL != "" {
		sqsClient = sqs.NewFromConfig(cfg, func(o *sqs.Options) {
			o.BaseEndpoint = &endpointURL
		})
	} else {
		sqsClient = sqs.NewFromConfig(cfg)
	}

	queueClient := queue.New(sqsClient, queueURL)

	resolver := graphql.NewResolver(queueClient)
	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	corsMiddlewareWithOrigins := corsMiddlewareWithConfig(origins)
	mux := http.NewServeMux()

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Production: only serve GraphQL endpoint
		mux.Handle("/graphql", srv)
		adapter := httpadapter.NewV2(corsMiddlewareWithOrigins(mux))
		lambda.Start(adapter.ProxyWithContext)
	} else {
		// Local development: serve both GraphQL and playground
		mux.Handle("/graphql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				playground.Handler("GraphQL playground", "/graphql").ServeHTTP(w, r)
			} else {
				srv.ServeHTTP(w, r)
			}
		}))
		fmt.Println("Server starting on :8080")
		fmt.Println("GraphQL playground available at http://localhost:8080/graphql")
		if err := http.ListenAndServe(":8080", corsMiddlewareWithOrigins(mux)); err != nil {
			log.Fatal(err)
		}
	}
}
