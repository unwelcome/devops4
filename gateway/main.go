package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	pb "github.com/unwelcome/devops4/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type DiceResponse struct {
	Value int32  `json:"value"`
	Emoji string `json:"emoji"`
}

var diceEmojis = map[int32]string{
	1: "⚀", 2: "⚁", 3: "⚂", 4: "⚃", 5: "⚄", 6: "⚅",
}

var diceClient pb.DiceServiceClient

func main() {
	// 1. Установка gRPC-соединения
	conn := setupGrpcConnection()
	defer conn.Close()

	// Инициализация клиента
	diceClient = pb.NewDiceServiceClient(conn)

	// Настройка HTTP-сервера
	app := fiber.New()

	// Регистрируем обработчик для пути /dice
	app.Get("/dice", diceHandler)

	port := ":8080"

	// Запускаем HTTP-сервер
	if err := app.Listen(port); err != nil {
		log.Fatalf("Fiber server failed: %v", err)
	}
}

func setupGrpcConnection() *grpc.ClientConn {
	// Загружаем CA, чтобы доверять сертификату Nginx
	caPem, err := os.ReadFile("/certs/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caPem)

	// Config для TLS. ServerName должен совпадать с CN в сертификате (nginx)
	creds := credentials.NewTLS(&tls.Config{
		RootCAs:    certPool,
		ServerName: "nginx",
	})

	conn, err := grpc.NewClient("nginx:443", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func diceHandler(c *fiber.Ctx) error {
	// Устанавливаем таймаут для gRPC-запроса
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	operationID := uuid.New().String()

	// Выполняем gRPC-запрос
	resp, err := diceClient.Roll(ctx, &pb.RollRequest{OperationId: operationID})
	if err != nil {
		log.Printf("Error rolling dice (gRPC call failed): %v", err)

		// Fiber позволяет легко отправлять статус и сообщение об ошибке
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Ошибка gRPC-запроса",
			"message": err.Error(),
		})
	}

	// Формируем ответ
	httpResponse := DiceResponse{
		Value: resp.Number,
		Emoji: diceEmojis[resp.Number],
	}

	// Fiber автоматически устанавливает Content-Type: application/json
	log.Printf("ID: %s; Method: GET; Path: /dice; Status: Success\n", operationID)
	return c.JSON(httpResponse)
}
