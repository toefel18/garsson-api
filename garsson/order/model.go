package order

import "github.com/gocraft/dbr"

type customerOrderEntity struct {
    ID                int64
    Status            string
    TimeCreated       string
    TimePrepared      dbr.NullString
    TimePaid          dbr.NullString
    WaiterID          string
    BarHandlerID      dbr.NullString
    CustomerName      dbr.NullString
    AmountPaidInCents dbr.NullInt64
    Remark            dbr.NullString
}

type customerOrderLineEntity struct {
    OrderID             int64
    ProductID           int64
    ProductName         string
    ProductBrand        dbr.NullString
    ProductPriceInCents int64
    Quantity            int64
    Remark              dbr.NullString
}

// ProductEntity is the same as the data
type ProductEntity struct {
    ID           int64  `json:"id"`
    Name         string `json:"name"`
    PriceInCents int64  `json:"priceInCents"`
    TimeAdded    string `json:"timeAdded,omitempty"`
}

// CustomerOrder is the public interface, requires multiple queries to run
type CustomerOrder struct {
    ID                int64                `json:"id"`
    Status            string               `json:"status"`
    TimeCreated       string               `json:"timeCreated"`
    TimePrepared      string               `json:"timePrepared,omitempty"`
    Waiter            string               `json:"waiter"`
    BarHandler        string               `json:"barHandler,omitempty"`
    CustomerName      string               `json:"customerName,omitempty"`
    TimePaid          string               `json:"timePaid,omitempty"`
    AmountPaidInCents int64                `json:"amountPaidInCents,omitempty"`
    Remark            string               `json:"remark,omitempty"`
    OrderLines        []*CustomerOrderLine `json:"orderLines"`
}

type CustomerOrderLine struct {
    ProductName         string `json:"productName"`
    ProductBrand        string `json:"productBrand,omitempty"`
    ProductPriceInCents int64 `json:"productPriceInCents"`
    Quantity            int64  `json:"quantity"`
    Remark              string `json:"remark,omitempty"`
}
