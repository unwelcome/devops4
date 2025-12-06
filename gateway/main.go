package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	pb "github.com/unwelcome/devops4/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var diceEmojis = map[int32]string{
	1: "⚀", 2: "⚁", 3: "⚂", 4: "⚃", 5: "⚄", 6: "⚅",
}

func main() {
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
	defer conn.Close()

	client := pb.NewDiceServiceClient(conn)

	for {
		resp, err := client.Roll(context.Background(), &pb.RollRequest{OperationId: uuid.New().String()})
		if err != nil {
			log.Printf("Error rolling dice: %v", err)
		} else {
			fmt.Printf("Dice: %s (Value: %d)\n", diceEmojis[resp.Number], resp.Number)
		}
		time.Sleep(3 * time.Second)
	}
}
