package product

import (
	"vinted/otel-workshop/internal/random"
	"vinted/otel-workshop/pb/genproto/otelworkshop"
)

//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate protoc --go_out=../../pb --go-grpc_out=../../pb --proto_path=../../pb ../../pb/workshop.proto

const (
	NameShoes = "shoes"
	NameHat   = "hat"
	NameSocks = "socks"
	NamePants = "pants"
	NameShirt = "shirt"

	ColorRed    = "red"
	ColorBlue   = "blue"
	ColorGreen  = "green"
	ColorYellow = "yellow"
	ColorBlack  = "black"
)

var names = []string{
	NameShoes,
	NameHat,
	NameSocks,
	NamePants,
	NameShirt,
}

var colors = []string{
	ColorRed,
	ColorBlue,
	ColorGreen,
	ColorYellow,
	ColorBlack,
}

func New() *otelworkshop.Product {
	return &otelworkshop.Product{
		Name:  random.Item(names),
		Color: random.Item(colors),
	}
}

func Names() []string {
	return names
}

func Colors() []string {
	return colors
}
