package failoverGrpctest

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"goworkwebook/syntax/grpc001/myGrpc"
	"log"
)

type FailedServer struct {
	myGrpc.UnimplementedUserServiceServer
	Name string
}

func (s *FailedServer) GetByID(ctx context.Context, request *myGrpc.GetByIDRequest) (*myGrpc.GetByIDResponse, error) {
	log.Println("进来了 failover")
	return nil, status.Errorf(codes.Unavailable, "假装我被熔断了")
}
