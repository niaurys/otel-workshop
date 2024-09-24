package buyer

import (
	"context"
	"vinted/otel-workshop/internal/random"
	"vinted/otel-workshop/pb/genproto/otelworkshop"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	name = "vinted/otel-workshop/internal/buyer"
)

var (
	tracer = otel.Tracer(name)
)

type Buyer interface {
	Buy(context.Context) error
}

type RandomBuyer struct {
	client otelworkshop.ShopServiceClient
	logger *logrus.Logger
}

func NewRandomBuyer(logger *logrus.Logger, shopAddress string) (*RandomBuyer, error) {
	conn, err := grpc.NewClient(shopAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	client := otelworkshop.NewShopServiceClient(conn)

	return &RandomBuyer{
		client: client,
		logger: logger,
	}, nil
}

var names = []string{
	"John",
	"Jane",
	"Jack",
	"Jill",
	"James",
}

var surnames = []string{
	"Doe",
	"Smith",
	"Johnson",
	"Brown",
	"Williams",
}

type person struct {
	name    string
	surname string
}

func randomPerson() *person {
	return &person{
		name:    random.Item(names),
		surname: random.Item(surnames),
	}
}

func (b *RandomBuyer) Buy(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "buyer.ListProducts")
	defer span.End()

	resp, err := b.client.ListProducts(ctx, &otelworkshop.Empty{})
	if err != nil {
		return err
	}

	count := len(resp.Products)

	span.SetAttributes(
		attribute.Int("count", count),
	)

	span.AddEvent("listed products")

	b.logger.WithField("count", count).Info("listed products")

	if len(resp.Products) == 0 {
		return nil
	}

	product := random.Item(resp.Products)
	person := randomPerson()
	quantity := random.Int64(product.Quantity)

	ctx, span = tracer.Start(ctx, "buyer.BuyProduct")
	defer span.End()

	_, err = b.client.BuyProduct(ctx, &otelworkshop.BuyProductRequest{
		Name:    person.name,
		Surname: person.surname,
		Product: &otelworkshop.Product{
			Name:     product.Name,
			Color:    product.Color,
			Quantity: quantity,
		},
	})
	if err != nil {
		return err
	}

	span.SetAttributes(
		attribute.String("name", person.name),
		attribute.String("surname", person.surname),
		attribute.Int64("quantity", quantity),
		attribute.String("color", product.Color),
		attribute.String("product", product.Name),
	)

	span.AddEvent("bought product")

	b.logger.WithFields(logrus.Fields{
		"name":     person.name,
		"surname":  person.surname,
		"quantity": quantity,
		"color":    product.Color,
		"product":  product.Name,
	}).Info("bought product")

	return nil
}
