package buyer

import (
	"context"
	"vinted/otel-workshop/internal/random"
	"vinted/otel-workshop/pb/genproto/otelworkshop"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	resp, err := b.client.ListProducts(ctx, &otelworkshop.Empty{})
	if err != nil {
		return err
	}

	b.logger.WithField("count", len(resp.Products)).Info("listed products")

	if len(resp.Products) == 0 {
		return nil
	}

	product := random.Item(resp.Products)
	person := randomPerson()
	quantity := random.Int64(product.Quantity)

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

	b.logger.WithFields(logrus.Fields{
		"name":     person.name,
		"surname":  person.surname,
		"quantity": quantity,
		"color":    product.Color,
		"product":  product.Name,
	}).Info("bought product")

	return nil
}
