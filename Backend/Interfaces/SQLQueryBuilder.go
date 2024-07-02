package Interfaces

// SQLQueryBuilder defines methods for building a SQL query dynamically
type SQLQueryBuilder interface {
	Where(field, value string) SQLQueryBuilder
	Order(field, direction string) SQLQueryBuilder
	Offset(offset int) SQLQueryBuilder
	Limit(limit int) SQLQueryBuilder
	Build() string
	BuildFullQueryOn(tableName string) string
	SelectFields(fields []string) SQLQueryBuilder
}
