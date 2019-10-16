// Command testblinktd runs a grpc server for the blinktd that can be used for debugging and testing.
package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/mweb/console-blinktd/internal/pkg/blinktv1"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:7023"))
	if err != nil {
		panic(err)
	}

	grpcs := grpc.NewServer()
	blinktv1.RegisterBlinktAPIServer(grpcs, NewBlinktService())

	if err := grpcs.Serve(lis); err != nil {
		panic(err)
	}

}

type BlinktService struct {
}

func NewBlinktService() BlinktService {
	return BlinktService{}
}

func (b BlinktService) SetPixel(ctx context.Context, req *blinktv1.SetPixelRequest) (*blinktv1.SetPixelResponse, error) {
	fmt.Printf("Set: LED[%d] b: '%s' c: 'r%d:g%d:b%d'\n", req.Index, req.Brightness, req.Color.R, req.Color.G, req.Color.B)
	return &blinktv1.SetPixelResponse{}, nil
}
func (b BlinktService) Clear(context.Context, *blinktv1.ClearRequest) (*blinktv1.ClearResponse, error) {
	fmt.Printf("Clear All\n")
	return &blinktv1.ClearResponse{}, nil
}
func (b BlinktService) Show(context.Context, *blinktv1.ShowRequest) (*blinktv1.ShowResponse, error) {
	fmt.Printf("Show All\n")
	return &blinktv1.ShowResponse{}, nil
}
