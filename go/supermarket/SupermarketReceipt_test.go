package supermarket

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type FakeCatalog struct {
    _products map[string]Product
    _prices map[string]float64
}

type TestCase struct {
	name       string
	products   []Product
	expected_total float64
	offers     []SpecialOffer
	amount float64
}


func (c FakeCatalog) unitPrice(product Product) float64 {
	return c._prices[product.name]
}

func (c FakeCatalog) addProduct(product Product, price float64) {
	c._products[product.name] = product
	c._prices[product.name] = price
}

func NewFakeCatalog() *FakeCatalog {
	var c FakeCatalog
	c._products = make(map[string]Product)
	c._prices = make(map[string]float64)
	return &c
}

func TestTenPercentDiscount(t *testing.T) {
	// ARRANGE
	var toothbrush = Product{name: "toothbrush", unit: Each}
	var apples = Product{name: "apples", unit: Kilo}
	var catalog = NewFakeCatalog()
	catalog.addProduct(toothbrush, 0.99)
	catalog.addProduct(apples, 1.99)

	var teller = NewTeller(catalog)
	teller.addSpecialOffer(TenPercentDiscount, toothbrush, 10.0)

	var cart = NewShoppingCart()
	cart.addItemQuantity(apples, 2.5)

	// ACT
	var receipt = teller.checksOutArticlesFrom(cart)

	// ASSERT
	assert.Equal(t, 4.975, receipt.totalPrice())
	assert.Equal(t, 0, len(receipt.discounts))
	require.Equal(t, 1, len(receipt.items))
	var receiptItem = receipt.items[0]
    assert.Equal(t, 1.99, receiptItem.price)
	assert.Equal(t, 2.5*1.99, receiptItem.totalPrice)
	assert.Equal(t, 2.5, receiptItem.quantity)
}

func TestApp(t *testing.T) {

	var toothbrush = Product{name: "toothbrush", unit: Each}
	var apples = Product{name: "apples", unit: Kilo}
	var catalog = NewFakeCatalog()
	catalog.addProduct(toothbrush, 0.99)
	catalog.addProduct(apples, 1.99)

	var teller = NewTeller(catalog)

	test_cases := []TestCase{
		{
			name:     "two products no ofer",
			products: []Product{toothbrush},
            expected_total: 1.98,
			offers:   []SpecialOffer{},
			amount:   2.0,
		},
		{
			name:     "10 percent discount",
			products: []Product{toothbrush},
            expected_total: 0.891,
			offers:   []SpecialOffer{SpecialOffer{TenPercentDiscount, toothbrush, 10.0}},
			amount:   1.0,
		},
		{
			name:     "ThreeForTwo",
			products: []Product{apples},
            expected_total: 3.98,
			offers:   []SpecialOffer{SpecialOffer{ThreeForTwo, apples, 0.0}},
			amount:   3.0,
		},
		{
			name:     "TwoForAmount",
			products: []Product{apples},
            expected_total: 3.5,
			offers:   []SpecialOffer{SpecialOffer{TwoForAmount, apples, 3.5}},
			amount:   2.0,
		},
		{
			name:     "FiveForAmount",
			products: []Product{apples},
            expected_total: 3.5,
			offers:   []SpecialOffer{SpecialOffer{FiveForAmount, apples, 3.5}},
			amount:   5.0,
		},
    }
	for _, testcase := range test_cases {
		t.Run(testcase.name, func(t *testing.T) {
			var cart = NewShoppingCart()
			for _, product := range testcase.products {
				cart.addItemQuantity(product, testcase.amount)
			}
			for _, offer := range testcase.offers {
				teller.addSpecialOffer(offer.offerType, offer.product, offer.argument)
			}
			var receipt = teller.checksOutArticlesFrom(cart)
			totalPrice := receipt.totalPrice()
			assert.Equal(t, testcase.expected_total, totalPrice)
		})
	}
}
