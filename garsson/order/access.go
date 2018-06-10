package order

import (
    "github.com/gocraft/dbr"
    "github.com/toefel18/garsson-api/garsson/db"
)

// TODO modify to ptr
func QueryProducts(sess dbr.SessionRunner) ([]ProductEntity, error) {
    var products = []ProductEntity{} // do not replace will nil slice declaration
    var err error
    _, err = sess.Select("*").From(db.ProductTable).Load(&products)
    return products, err
}

// TODO modify to ptr
func QueryProductByID(sess dbr.SessionRunner, id int64) (ProductEntity, error) {
    var products = ProductEntity{} // do not replace will nil slice declaration
    err := sess.Select("*").From(db.ProductTable).Where("id = ? ", id).LoadOne(&products)
    return products, err
}

func queryOrderEntityByID(sess dbr.SessionRunner, id int64) (*customerOrderEntity, error) {
    var order *customerOrderEntity
    if err := sess.Select("*").From(db.CustomerOrderTable).Where("id = ?", id).LoadOne(&order); err != nil {
        return nil, err
    } else {
        return order, nil
    }
}

func queryOrderLinesByOrderID(sess dbr.SessionRunner, orderId int64) ([]*customerOrderLineEntity, error) {
    var orderLines []*customerOrderLineEntity
    if _, err := sess.Select("*").From(db.CustomerOrderLineTable).Where("order_id = ?", orderId).Load(&orderLines); err != nil {
        return nil, err
    } else {
        return orderLines, nil
    }
}
