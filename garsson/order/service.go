package order

import (
    "github.com/gocraft/dbr"
)

func FindOrderByID(sess dbr.SessionRunner, id int64) (*CustomerOrder, error) {
    if order, err := queryOrderEntityByID(sess, id); err != nil {
        return nil, err
    } else if lines, err := queryOrderLinesByOrderID(sess, order.ID); err != nil {
        return nil, err
    } else {
        return mapOrderToPublicAPI(order, lines)
    }
}

func FindOrdersWithStatus(sess dbr.SessionRunner, status []string) ([]*CustomerOrder, error) {
    orders, err := queryOrdersWithStatus(sess, status)
    if err != nil {
        return nil, err
    }

    customerOrders := make([]*CustomerOrder, 0)
    for _, v := range orders {
        if orderLines, err := queryOrderLinesByOrderID(sess, v.ID); err != nil {
            return nil, err
        } else if customerOrder, err := mapOrderToPublicAPI(v, orderLines); err != nil{
            return nil, err
        } else {
            customerOrders = append(customerOrders, customerOrder)
        }
    }
    return customerOrders, nil
}


func mapOrderToPublicAPI(order *customerOrderEntity, lines []*customerOrderLineEntity) (*CustomerOrder, error) {
    publicOrder := CustomerOrder{
        ID:                order.ID,
        Status:            order.Status,
        TimeCreated:       order.TimeCreated,
        TimePrepared:      order.TimePrepared.String,
        TimePaid:          order.TimePaid.String,
        Waiter:            order.WaiterID,
        BarHandler:        order.BarHandlerID.String,
        CustomerName:      order.CustomerName.String,
        AmountPaidInCents: order.AmountPaidInCents.Int64,
        Remark:            order.Remark.String,
        OrderLines:        mapOrderLinesToPublicAPI(lines),
    }

    return &publicOrder, nil
}

func mapOrderLinesToPublicAPI(lines []*customerOrderLineEntity) []*CustomerOrderLine {
    orderLines := []*CustomerOrderLine{} // provide empty array if none found
    for _, line := range lines {
        orderLine := CustomerOrderLine{
            ProductName:         line.ProductName,
            ProductBrand:        line.ProductBrand.String,
            ProductPriceInCents: line.ProductPriceInCents,
            Quantity:            line.Quantity,
            Remark:              line.Remark.String,
        }
        orderLines = append(orderLines, &orderLine)
    }
    return orderLines
}
