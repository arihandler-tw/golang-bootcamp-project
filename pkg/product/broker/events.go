package broker

const (
	ProductsTopic = "products"
)

type ProductCreationRequest struct {
	Id          *string
	Price       float32
	Description string
}
