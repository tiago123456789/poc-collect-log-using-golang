// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"net"

// 	"github.com/hibiken/asynq"
// 	pb "github.com/tiago123456789/poc-grpc/proto"
// 	"google.golang.org/grpc"
// )

// type server struct {
// 	pb.UnimplementedLogServer
// }

// const redisAddr = "127.0.0.1:6379"

// var client *asynq.Client

// var (
// 	port = flag.Int("port", 50051, "The server port")
// )

// func (s *server) NewLine(ctx context.Context, in *pb.Line) (*pb.OKReply, error) {

// 	jsonEncode, err := json.Marshal(in)
// 	fmt.Println(string(jsonEncode))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// _, err = client.Enqueue(
// 	// 	asynq.NewTask("logs", jsonEncode),
// 	// 	asynq.MaxRetry(1),
// 	// )

// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	return &pb.OKReply{Message: "OK"}, nil
// }

// func main() {
// 	client = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
// 	defer client.Close()

// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}

// 	s := grpc.NewServer()
// 	pb.RegisterLogServer(s, &server{})

//		log.Printf("server listening at %v", lis.Addr())
//		if err := s.Serve(lis); err != nil {
//			log.Fatalf("failed to serve: %v", err)
//		}
//	}
package main

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/tiago123456789/poc-collect-log-golang/model"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		logData := new(model.Log)
		if err := c.BodyParser(logData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		logInBytes, _ := json.Marshal(logData)
		_, err := client.Enqueue(
			asynq.NewTask("log", logInBytes),
			asynq.MaxRetry(2),
		)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(202)
	})

	app.Listen(":3000")
}
