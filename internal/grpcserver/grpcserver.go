package grpcserver

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/BiteLikeASnake/kiddy_line_processor/internal/model"
	grpc "google.golang.org/grpc"
)

type grpcServer struct{}

//Реализует метод SubscribeOnSportsLines интерфейса LinesServer
func (s *grpcServer) SubscribeOnSportsLines(stream Lines_SubscribeOnSportsLinesServer) error {
	ctx := stream.Context()
	//
	funcFirstTimeInvoked := true //маркер первого запуска метода
	exitSender := make(chan int)
	for {
		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return fmt.Errorf("SubscribeOnSportsLines: %v", ctx.Err())
		default:
		}

		// receive data from stream
		query, err := stream.Recv()
		if funcFirstTimeInvoked {
			funcFirstTimeInvoked = false
		} else {
			exitSender <- 1
		}

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("SubscribeOnSportsLines: %v", err)
		}
		fmt.Println("Получены данные")

		//запускаем отправитель линий
		go func() {
			senderFirstTimeInvoked := true //маркер первого запуска отправителя
			lines := query.GetLine()

			frequency := query.GetFrequency()
			//Создаем возвращаемый тип
			resp := LinesResponse{}
			for {

				select {
				case <-exitSender:
					fmt.Println("Выход из отправителя")
					return
				default:
				}

				if senderFirstTimeInvoked {
					senderFirstTimeInvoked = false
					//return lines with values
					for _, line := range lines {
						lineval, err := model.Storage.ReturnLineCurrentVal(line)
						if err != nil {
							fmt.Println(err)
							return
						}
						resp.Resp = append(resp.Resp, &LineDelta{Line: line, Delta: lineval})
					}
				} else {
					//return lines with deltas
					for _, line := range resp.Resp {
						linedelta, err := model.Storage.ReturnLineDelta(line.GetLine())
						if err != nil {
							fmt.Println(err)
							return
						}
						line.Delta = linedelta
					}
				}
				if err := stream.Send(&resp); err != nil {
					fmt.Printf("send error %v", err)
				}
				time.Sleep(time.Second * time.Duration(frequency))
			}
		}()

	}
}

//StartServer запускает gRPC сервер
func StartServer(port string) error {
	lis, err := net.Listen("tcp", port) //port = ":9000"
	if err != nil {
		return fmt.Errorf("StartServer: %v", err)
	}
	s := grpcServer{}
	grpcServer := grpc.NewServer()
	RegisterLinesServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("StartServer: %v", err)
	}
	return nil
}
