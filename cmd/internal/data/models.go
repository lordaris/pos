package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryMovement struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ProductID    primitive.ObjectID `bson:"product_id"`
	MovementType string             `bson:"movement_type"`
	Quantity     int                `bson:"quantity"`
	MovementDate time.Time          `bson:"movement_date"`
	UserID       primitive.ObjectID `bson:"user_id,omitempty"`
	Reason       string             `bson:"reason,omitempty"`
}

type Customer struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Barcode string             `bson:"barcode"`
	Name    string             `bson:"name"`
	Email   string             `bson:"email"`
	Phone   string             `bson:"phone"`
	Points  int                `bson:"points"`
}

type Invoice struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	TotalAmount float64            `bson:"total_amount"`
	UserID      primitive.ObjectID `bson:"user_id"`
	CustomerID  primitive.ObjectID `bson:"customer_id,omitempty"`
	SaleDate    time.Time          `bson:"sale_date"`
	ChangeGiven float64            `bson:"change_given"`
	Discount    float64            `bson:"discount,omitempty"`
	Items       []InvoiceItem      `bson:"items"`
	Payments    []InvoicePayment   `bson:"payments"`
}

type InvoiceItem struct {
	ProductID primitive.ObjectID `bson:"product_id"`
	Quantity  int                `bson:"quantity"`
	Price     float64            `bson:"price"`
}

type PaymentMethod struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

type InvoicePayment struct {
	PaymentMethodID primitive.ObjectID `bson:"payment_method_id"`
	Amount          float64            `bson:"amount"`
}
