package main

// Задание 21
// Реализовать паттерн «адаптер» на любом примере.

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
)

// Cервис A выдаёт заказы в формате XML.
type (
	OrderServiceA struct {
		XMLName         xml.Name          `xml:"order"`
		OrderID         uuid.UUID         `xml:"order_id"`
		OrderedAt       time.Time         `xml:"ordered_at"`
		ProcessedAt     time.Time         `xml:"processed_at,omitempty"`
		ShippedAt       time.Time         `xml:"shipped_at,omitempty"`
		Product         []ProductServiceA `xml:"product"`
		CustomerID      uuid.UUID         `xml:"customer_id"`
		ShippingAddress string            `xml:"shipping_address"`
	}
	ProductServiceA struct {
		ProductID   uuid.UUID `xml:"product_id"`
		ProductName string    `xml:"product_name"`
		Quantity    int       `xml:"quantity"`
		Price       float64   `xml:"price"`
		Currency    string    `xml:"currency"`
	}
)

// Сервис B принимает информацию в виде JSON.
type OrderServiceB struct {
	OrderID         uuid.UUID `json:"order_id"`
	ProcessedAt     time.Time `json:"processed_at"`
	ShippedAt       time.Time `json:"shipped_at"`
	CustomerID      uuid.UUID `json:"customer_id"`
	ShippingAddress string    `json:"shipping_address"`
	TotalProducts   int       `json:"total_products"`
	TotalAmountRUB  string    `json:"total_amount_rub"`
}

// AdapterAB выполняет конвертацию данных, полученных от сервиса A в формат, который читает сервис B.
type AdapterAB struct {
	order  OrderServiceA
	writer io.Writer
}

// Adapt конвертирует данные заказа из структуры формата сервиса A в формат сервиса B.
// Эта функция используется, если нужен экземпляр структуры сервиса B.
// Для прямого конвертирования из XML в JSON используется интерфейс Writer.
func (ad *AdapterAB) Adapt(xmlOrder []byte) (*OrderServiceB, error) {
	if err := ad.importXML(xmlOrder); err != nil {
		return nil, fmt.Errorf("could not decode the order: %w", err)
	}
	return ad.convert(), nil
}

// NewAdapter создаёт адаптер, реализующий интерфейс io.Writer.
func NewAdapter(w io.Writer) *AdapterAB {
	return &AdapterAB{writer: w}
}

// Write реализует интерфейс io.Writer. Осуществляет конвертацию заказов из формата сервиса A в формат сервиса B.
func (ad *AdapterAB) Write(p []byte) (n int, err error) {
	if err := ad.importXML(p); err != nil {
		return 0, err
	}
	orderB := ad.convert()
	jsonOrder, err := json.Marshal(orderB)
	if err != nil {
		return 0, err
	}
	return ad.writer.Write(jsonOrder)
}

// importXML раскодирует XML-сообщение в формат структтуры сервиса A.
func (ad *AdapterAB) importXML(p []byte) error {
	if err := xml.Unmarshal(p, &ad.order); err != nil {
		return fmt.Errorf("could not decode the order: %w", err)
	}
	return nil
}

// convert преобразует структуру заказа сервиса A в струтуру сервиса B.
// Валюта конвертируется в рубли и подсчитывается общая стоимость заказов.
func (ad *AdapterAB) convert() *OrderServiceB {
	totalProducts := 0
	totalAmount := 0.0
	for _, product := range ad.order.Product {
		totalProducts += product.Quantity
		priceRUB, err := convertPriceToRub(product.Price, product.Currency)
		if err != nil {
			panic(err)
		}
		totalAmount += float64(product.Quantity) * priceRUB
	}
	return &OrderServiceB{
		OrderID:         ad.order.OrderID,
		ProcessedAt:     ad.order.ProcessedAt,
		ShippedAt:       ad.order.ShippedAt,
		CustomerID:      ad.order.CustomerID,
		ShippingAddress: ad.order.ShippingAddress,
		TotalProducts:   totalProducts,
		TotalAmountRUB:  fmt.Sprintf("%.2fр.", totalAmount),
	}
}

// convertPriceToRub конвертирует валюту по выгодному курсу))
func convertPriceToRub(amount float64, currency string) (float64, error) {
	const (
		eurToRub = 63.83
		usdToRub = 62.90
	)
	var rate float64
	switch currency {
	case "EUR":
		rate = eurToRub
	case "USD":
		rate = usdToRub
	case "RUB":
		rate = 1
	default:
		return -1, errors.New("incorrect currency type")
	}
	return amount * rate, nil
}

// данные заказов для тестирования
var (
	order1 = []byte(`<order>
<order_id>349d828e-16a8-4db9-8461-02fcdc5a47f6</order_id>
<ordered_at>2021-04-01T12:03:11Z</ordered_at>
<processed_at>2021-04-02T15:13:12Z</processed_at>
<shipped_at>2021-04-10T09:54:32Z</shipped_at>
<product>
	<product_id>f6788e2f-4646-4ab8-920c-ee29ab858ecc</product_id>
	<product_name>Apple iPhone 18</product_name>
	<quantity>2</quantity>
	<price>1999</price>
	<currency>USD</currency>
</product>
<product>
	<product_id>8adb2d72-e7e9-4d06-92b4-efa31e29b235</product_id>
	<product_name>Apple Watch 7</product_name>
	<quantity>1</quantity>
	<price>999</price>
	<currency>USD</currency>
</product>
<customer_id>a3ca6f8c-fe97-44ab-ba4f-251de72b4c59</customer_id>
<shipping_address>SPb, Nevsky pr. 128 fl.99</shipping_address>
</order>
`)

	order2 = []byte(`<order>
<order_id>cd509f1a-5c51-4e85-b328-7afb5ac66d03</order_id>
<ordered_at>2021-04-09T19:03:11Z</ordered_at>
<processed_at>2021-04-10T21:13:12Z</processed_at>
<product>
	<product_id>26783ad6-6b71-42ec-9ebb-9ebd40f85e87</product_id>
	<product_name>Doshirak</product_name>
	<quantity>33</quantity>
	<price>1.55</price>
	<currency>EUR</currency>
</product>
<customer_id>06da9344-6f6e-4305-af88-f3654f7707ed</customer_id>
<shipping_address>SPb, Nevsky pr. 7 fl.9</shipping_address>
</order>
`)
)

func main() {
	// пример использования адаптера как Writer'а
	a := NewAdapter(os.Stdout)
	a.Write(order1)
	fmt.Println()

	// для конвертации в структуру конструктор не нужен
	order, err := new(AdapterAB).Adapt(order2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *order)
}
