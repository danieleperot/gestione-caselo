package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/daniele/gestione-caselo/backend/eventform/internal/auth"
	"github.com/daniele/gestione-caselo/backend/eventform/internal/graphql"
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
	endpoint := appconfig.GetEnvVariable("COGNITO_ENDPOINT")
	poolID := appconfig.GetEnvVariable("COGNITO_POOL_ID")
	origins := appconfig.GetEnvVariable("ALLOWED_ORIGINS")

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}}))

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	jwksURL := fmt.Sprintf("%s/%s/.well-known/jwks.json", endpoint, poolID)
	mux.Handle("/graphql", auth.Middleware(auth.Config{JWKSUrl: jwksURL})(srv))

	corsMiddlewareWithOrigins := corsMiddlewareWithConfig(origins)

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		lambda.Start(httpadapter.New(corsMiddlewareWithOrigins(mux)).ProxyWithContext)
	} else {
		fmt.Println("Server starting on :8080")
		fmt.Println("GraphQL playground available at http://localhost:8080")
		if err := http.ListenAndServe(":8080", corsMiddlewareWithOrigins(mux)); err != nil {
			log.Fatal(err)
		}
	}
}
