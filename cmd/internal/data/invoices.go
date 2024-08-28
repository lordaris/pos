package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionInvoice = "invoices"
)

type Invoice struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	TicketNumber int                `bson:"ticket_number"`
	TicketTotal  float64            `bson:"ticket_total"`
	UserID       primitive.ObjectID `bson:"user_id"`
	Username     string             `bson:"username"`
	CustomerID   primitive.ObjectID `bson:"customer_id,omitempty"`
	SaleDate     time.Time          `bson:"sale_date"`
	ChangeGiven  float64            `bson:"change_given"`
	Discount     float64            `bson:"discount,omitempty"`
	Items        []InvoiceItem      `bson:"items"`
	Payments     []InvoicePayment   `bson:"payments"`
}

type InvoiceItem struct {
	ProductName        string  `bson:"product_name"`
	Barcode            int     `bson:"barcode"`
	Quantity           int     `bson:"quantity"`
	Price              float64 `bson:"price"`
	PromotionType      string  `bson:"promotion_type"`
	PaidQuantity       int     `bson:"paid_quantity"`
	FreeQuantity       int     `bson:"free_quantity"`
	DiscountPercentage int     `bson:"discount_percentage"`
	PriceWithDiscount  float64 `bswon:"price_with_discount"`
	DiscountPrice      float64 `bson:"discount_price"`

	TotalAmount float64 `bson:"total_amount"`
}

type PaymentMethod struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

type InvoicePayment struct {
	PaymentMethodID primitive.ObjectID `bson:"payment_method_id"`
	Amount          float64            `bson:"amount"`
}
