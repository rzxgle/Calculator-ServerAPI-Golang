package main

import (
	"apicalc/src/pb/calc"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	calc.CalcServiceServer
}

func (s *server) Calc(stream calc.CalcService_CalcServer) error {
	var quantity int32 = 0
	var total int32 = 0
	for {
		input, err := stream.Recv()
		if err == io.EOF {
			avg := float64(float64(total) / float64(quantity))
			return stream.SendAndClose(&calc.Output{
				Quantity: quantity,
				Average:  avg,
				Total:    total,
			})
		}
		if err != nil {
			return err
		}
		quantity++
		total += input.GetValue()

		fmt.Printf("input; %+v\n", input)
	}
}

func main() {
	log.Println("Starting grpc server")
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln("error on get listener. error: ", err)
	}

	s := grpc.NewServer()
	calc.RegisterCalcServiceServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalln("error on serve. error: ", err)
	}
}
