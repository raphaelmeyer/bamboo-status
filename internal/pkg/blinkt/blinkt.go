package blinkt

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"time"
)

const (
	address = "localhost:7023"
)

type Blinkt struct {
	client BlinktClient
}

type LedState uint

const (
	Red LedState = iota
	Green
	Blue
	Yellow
	Purple
	Cyan
	Off
)

func (state LedState) toPixel(i uint32) *Pixel {
	switch state {
	case Red:
		return &Pixel{Index: i, Brightness: Pixel_Medium, Color: &Color{R: 255, G: 0, B: 0}}
	case Green:
		return &Pixel{Index: i, Brightness: Pixel_Medium, Color: &Color{R: 0, G: 127, B: 0}}
	case Blue:
		return &Pixel{Index: i, Brightness: Pixel_Medium, Color: &Color{R: 0, G: 0, B: 127}}
	case Yellow:
		return &Pixel{Index: i, Brightness: Pixel_Medium, Color: &Color{R: 127, G: 127, B: 0}}
	case Purple:
		return &Pixel{Index: i, Brightness: Pixel_Medium, Color: &Color{R: 127, G: 0, B: 127}}
	case Cyan:
		return &Pixel{Index: i, Brightness: Pixel_Medium, Color: &Color{R: 0, G: 127, B: 127}}
	case Off:
		return &Pixel{Index: i, Brightness: Pixel_Off, Color: &Color{R: 0, G: 0, B: 0}}
	}
	return nil
}

func NewBlinkt() (*Blinkt, func(), error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	client := NewBlinktClient(conn)
	close := func() {
		conn.Close()
	}

	return &Blinkt{client}, close, nil
}

func (bt *Blinkt) Clear() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := bt.client.Clear(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	_, err = bt.client.Show(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	return nil
}

func (bt *Blinkt) SetLed(index uint32, state LedState) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := bt.client.SetPixel(ctx, state.toPixel(index))
	if err != nil {
		return err
	}

	_, err = bt.client.Show(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	return nil
}
