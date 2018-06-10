package migration

// When adding new definitions, create a new const with the next number as its ID. Then add it as the last element
// of the completeDatabaseDefinition array at the bottom of this file.
// WARNING: Do not change any the statements in this file once deployed (formatting can change)!
//          The migration script will create a hash of each query and fail if it doesn't match with the
//          previous version
const (
    V1UsersTable = `CREATE TABLE user_account (
                      email         VARCHAR(128) PRIMARY KEY,
                      password_hash VARCHAR(128),
                      roles         VARCHAR(256),
                      last_sign_in  VARCHAR(64)
                    )`

    V2ProductTable = `CREATE TABLE product (
                        id             BIGSERIAL PRIMARY KEY,
                        name           VARCHAR(256) NOT NULL,
                        brand          VARCHAR(256) NOT NULL,
                        price_in_cents BIGINT NOT NULL,
                        time_added     VARCHAR(64) NOT NULL
                      )`

    V3CustomerOrderTable = `CREATE TABLE customer_order (
                              id                   BIGSERIAL PRIMARY KEY,
                              status               VARCHAR(128) NOT NULL,
                              time_created         VARCHAR(64),
                              time_prepared        VARCHAR(64),
                              time_paid            VARCHAR(64),
                              waiter_id            VARCHAR(128) NOT NULL REFERENCES user_account (email),
                              bar_handler_id       VARCHAR(128) REFERENCES user_account (email),
                              customer_name        VARCHAR(128),
                              amount_paid_in_cents BIGINT,
                              remark               TEXT
                            )`

    V4CustomerOrderLineTable = `CREATE TABLE customer_order_line (
                                  order_id               BIGINT NOT NULL REFERENCES customer_order (id),
                                  product_id             BIGINT,
                                  product_name           VARCHAR(256),
                                  product_brand          VARCHAR(256),
                                  product_price_in_cents BIGINT,
                                  quantity               BIGINT NOT NULL,
                                  remark                 TEXT NULL,
                                  PRIMARY KEY (order_id, product_id)
                                )`
)


var completeDatabaseDefinition = []string{
    V1UsersTable,
    V2ProductTable,
    V3CustomerOrderTable,
    V4CustomerOrderLineTable,
}
