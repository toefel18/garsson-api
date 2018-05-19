package migration

// When adding new definitions, create a new const with the next number as its ID. Then add it as the last element
// of the completeDatabaseDefinition array at the bottom of this file.
// WARNING: Do not change any the statements in this file once deployed (formatting can change)!
//          The migration script will create a hash of each query and fail if it doesn't match with the
//          previous version
const (
	v1CreateProductTable = `CREATE TABLE product (
							id SERIAL PRIMARY KEY, 
							name VARCHAR(256),
							price REAL,
							date_added VARCHAR(64))`
)

var completeDatabaseDefinition = []string {
    v1CreateProductTable,
}