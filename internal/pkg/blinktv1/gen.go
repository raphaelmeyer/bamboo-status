package blinktv1

//go:generate go get github.com/raphaelmeyer/blinktd@v1.0.0
//go:generate sh -c "protoc -I$(go list -m -f '{{ .Dir }}/proto' github.com/raphaelmeyer/blinktd) --go_out=plugins=grpc:. blinkt_api.proto blinkt_types.proto"
