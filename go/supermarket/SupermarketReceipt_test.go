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
	amount []float64
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
	var notABundle = NewBundle([]Product{}, false)//replace testcases with this
	// ARRANGE
	var toothbrush = Product{name: "toothbrush", unit: Each}
	var apples = Product{name: "apples", unit: Kilo}
	var catalog = NewFakeCatalog()
	catalog.addProduct(toothbrush, 0.99)
	catalog.addProduct(apples, 1.99)

	var teller = NewTeller(catalog)
	teller.addSpecialOffer(TenPercentDiscount, toothbrush, 10.0, *notABundle)

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
	var testBundle = *NewBundle([]Product{toothbrush, apples}, true)
	//var notABundle = NewBundle([]Product{}, false)//replace testcases with this
	test_cases := []TestCase{
		{
			name:     "two products no ofer",
			products: []Product{toothbrush, apples},
            expected_total: 2.98,
			offers:   []SpecialOffer{},
			amount:   []float64{1.0,1.0},
		},
		{
			name:     "10 percent discount",
			products: []Product{toothbrush},
            expected_total: 0.891,
			offers:   []SpecialOffer{SpecialOffer{TenPercentDiscount, toothbrush, 10.0, Bundle{[]Product{}, false}}},
			amount:   []float64{1.0},
		},
		{
			name:     "ThreeForTwo",
			products: []Product{apples},
            expected_total: 3.98,
			offers:   []SpecialOffer{SpecialOffer{ThreeForTwo, apples, 0.0, Bundle{[]Product{}, false}}},
			amount:   []float64{3.0},
		},
		{
			name:     "TwoForAmount",
			products: []Product{apples},
            expected_total: 3.5,
			offers:   []SpecialOffer{SpecialOffer{TwoForAmount, apples, 3.5, Bundle{[]Product{}, false}}},
			amount:   []float64{2.0},
		},
		{
			name:     "FiveForAmount",
			products: []Product{apples},
            expected_total: 3.5,
			offers:   []SpecialOffer{SpecialOffer{FiveForAmount, apples, 3.5, Bundle{[]Product{}, false}}},
			amount:   []float64{5.0},
		},
		{
			name:     "BasicBundleProductApple",//asumes bundles should have one product by defnition, a bundle cant be 3oranges and 1 deodorant, thats why specialofer just got an array with one product each
			products: []Product{toothbrush, apples},
            expected_total: 2.68,//(0.99+1.99)*0.9
			offers:   []SpecialOffer{SpecialOffer{FiveForAmount, apples, 10.0, testBundle}},//when there is a bundle the speccialoffer product is not used dont matters?
			amount:   []float64{1.0,1.0},
		},
		{
			name:     "BasicBundleProductToothbrush",//asumes bundles should have one product by defnition, a bundle cant be 3oranges and 1 deodorant, thats why specialofer just got an array with one product each
			products: []Product{toothbrush, apples},
            expected_total: 2.68,//(0.99+1.99)*0.9
			offers:   []SpecialOffer{SpecialOffer{FiveForAmount, toothbrush, 10.0, testBundle}},
			amount:   []float64{1.0,1.0},
		},
		{
			name:     "BrokenBasicBundleApple",
			products: []Product{toothbrush, apples},
            expected_total: 3.97,
			offers:   []SpecialOffer{SpecialOffer{FiveForAmount, apples, 10.0, testBundle}},
			amount:   []float64{2.0,1.0},
		},
		{
			name:     "BrokenBasicBundleToothbrush",
			products: []Product{toothbrush, apples},
            expected_total: 3.97,
			offers:   []SpecialOffer{SpecialOffer{FiveForAmount, toothbrush, 10.0, testBundle}},
			amount:   []float64{2.0,1.0},
		},
    }
	for _, testcase := range test_cases {
		t.Run(testcase.name, func(t *testing.T) {
			var cart = NewShoppingCart()
			for i, product := range testcase.products {
				cart.addItemQuantity(product, testcase.amount[i])
			}
			for _, offer := range testcase.offers {
				teller.addSpecialOffer(offer.offerType, offer.product, offer.argument, offer.bundle)
			}
			var receipt = teller.checksOutArticlesFrom(cart)
			totalPrice := receipt.totalPrice()
			assert.Equal(t, testcase.expected_total, totalPrice)
		})
	}
}
