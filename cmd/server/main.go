package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/luuan11/clean-architecture/internal/infra/database"
	"github.com/luuan11/clean-architecture/internal/infra/graphql/graph"
	"github.com/luuan11/clean-architecture/internal/infra/grpc/pb"
	"github.com/luuan11/clean-architecture/internal/infra/grpc/service"
	"github.com/luuan11/clean-architecture/internal/infra/web/webserver"
	"github.com/luuan11/clean-architecture/internal/usecase"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "orders")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database")

	orderRepository := database.NewOrderRepository(db)

	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	go startRESTServer(createOrderUseCase, listOrdersUseCase)
	go startGRPCServer(createOrderUseCase, listOrdersUseCase)
	startGraphQLServer(createOrderUseCase, listOrdersUseCase)
}

func startRESTServer(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) {
	port := getEnv("REST_PORT", "8080")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	orderHandler := webserver.NewWebOrderHandler(createOrderUseCase, listOrdersUseCase)

	r.Post("/order", orderHandler.CreateOrder)
	r.Get("/order", orderHandler.ListOrders)

	log.Printf("REST server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start REST server: %v", err)
	}
}

func startGRPCServer(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) {
	port := getEnv("GRPC_PORT", "50051")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)

	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}

func startGraphQLServer(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) {
	port := getEnv("GRAPHQL_PORT", "8081")

	resolver := graph.NewOrderResolver(createOrderUseCase, listOrdersUseCase)
	schema, err := resolver.BuildSchema()
	if err != nil {
		log.Fatalf("Failed to build GraphQL schema: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)

	log.Printf("GraphQL server listening on port %s", port)
	log.Printf("GraphiQL available at http://localhost:%s/graphql", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start GraphQL server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
