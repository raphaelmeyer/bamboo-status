package blinkt

import (
	"bamboo/internal/pkg/blinktv1"
	"context"
	"google.golang.org/grpc"
	"time"
)

const (
	address = "localhost:7023"
)

type Blinkt struct {
	client blinktv1.BlinktAPIClient
}

type LedState uint

const (
	Red LedState = iota
	Green
	Blue
	Yellow
	Purple
	Cyan
	White
	Off
)

func (state LedState) toPixel(i uint32) *blinktv1.SetPixelRequest {
	switch state {
	case Red:
		return &blinktv1.SetPixelRequest{Index: i, Brightness: blinktv1.Brightness_BRIGHTNESS_LOW, Color: &blinktv1.Color{R: 191, G: 0, B: 0}}
	case Green:
		return &blinktv1.SetPixelRequest{Index: i, Brightness: blinktv1.Brightness_BRIGHTNESS_LOW, Color: &blinktv1.Color{R: 0, G: 127, B: 0}}
	case Blue:
		return &blinktv1.SetPixelRequest{Index: i, Brightness: blinktv1.Brightness_BRIGHTNESS_LOW, Color: &blinktv1.Color{R: 0, G: 0, B: 127}}
	case Yellow:
		return &blinktv1.SetPixelRequest{Index: i, Brightness: blinktv1.Brightness_BRIGHTNESS_LOW, Color: &blinktv1.Color{R: 191, G: 127, B: 0}}
	case Purple:
		return &blinktv1.SetPixelRequest{Index: i, Brightness: blinktv1.Brightness_BRIGHTNESS_LOW, Color: &blinktv1.Color{R: 127, G: 0, B: 127}}
	case Cyan:
		return &blinktv1.SetPixelRequest{Index: i, Brightness: blinktv1.Brightness_BRIGHTNESS_LOW, Color: &blinktv1.Color{R: 0, G: 127, B: 127}}
	case White:
		return &blinktv1.SetPixelRequest{Index: i, Brightness: blinktv1.Brightness_BRIGHTNESS_LOW, Color: &blinktv1.Color{R: 127, G: 127, B: 127}}
	case Off:
		return &blinktv1.SetPixelRequest{Index: i, Brightness: blinktv1.Brightness_BRIGHTNESS_LOW, Color: &blinktv1.Color{R: 0, G: 0, B: 0}}
	}
	return nil
}

func NewBlinkt() (*Blinkt, func(), error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	client := blinktv1.NewBlinktAPIClient(conn)
	close := func() {
		conn.Close()
	}

	return &Blinkt{client}, close, nil
}

func (bt *Blinkt) Clear() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := bt.client.Clear(ctx, &blinktv1.ClearRequest{})
	if err != nil {
		return err
	}

	_, err = bt.client.Show(ctx, &blinktv1.ShowRequest{})
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

	_, err = bt.client.Show(ctx, &blinktv1.ShowRequest{})
	if err != nil {
		return err
	}

	return nil
}
