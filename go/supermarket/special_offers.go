package supermarket

type SpecialOfferType int

const (
	TenPercentDiscount SpecialOfferType = iota
	ThreeForTwo
	TwoForAmount
	FiveForAmount
)

type Bundle struct {
	products []Product
    bundled bool
}

func NewBundle(products []Product, bundled bool) *Bundle {
    var b Bundle
	b.products = products
    b.bundled = bundled
    return &b
}

func (b *Bundle) isValid(productQuantities map[Product]float64) bool {
    for _, prod := range b.products {
        if productQuantities[prod] > 1 {
            return false
        }
    }
    return true
}

type SpecialOffer struct {
	offerType SpecialOfferType
	product Product
	argument float64
	bundle Bundle
}

type Discount struct {
	product Product
	description string
	discountAmount float64
}

