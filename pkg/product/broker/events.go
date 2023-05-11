package broker

const (
	ProductsTopic = "products"
)

type ProductCreationRequest struct {
	id          *string
	price       float32
	description string
}
