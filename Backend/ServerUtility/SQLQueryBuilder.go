package ServerUtility

import (
	"Fetch-Rewards-API/Backend/Interfaces"
	_ "Fetch-Rewards-API/Backend/Interfaces"
	"fmt"
	"strings"
)

// MySQLQueryBuilder implements SQLQueryBuilder interface for MySQL queries
type MySQLQueryBuilder struct {
	whereClauses []string
	orderClause  string
	offset       int
	limit        int
	fields       []string // new field to hold selected fields
}

// Where adds a WHERE clause to the query builder
func (qb *MySQLQueryBuilder) Where(field, value string) Interfaces.SQLQueryBuilder {
	qb.whereClauses = append(qb.whereClauses, fmt.Sprintf("%s LIKE '%s'", field, value))
	return qb
}

// Order adds an ORDER BY clause to the query builder
func (qb *MySQLQueryBuilder) Order(field, direction string) Interfaces.SQLQueryBuilder {
	qb.orderClause = fmt.Sprintf("ORDER BY %s %s", field, direction)
	return qb
}

// Offset sets the OFFSET clause for pagination
func (qb *MySQLQueryBuilder) Offset(offset int) Interfaces.SQLQueryBuilder {
	qb.offset = offset
	return qb
}

// Limit sets the LIMIT clause for pagination
func (qb *MySQLQueryBuilder) Limit(limit int) Interfaces.SQLQueryBuilder {
	qb.limit = limit
	return qb
}

// SelectFields sets the fields to be selected in the query
func (qb *MySQLQueryBuilder) SelectFields(fields []string) Interfaces.SQLQueryBuilder {
	qb.fields = fields
	return qb
}

// Build constructs the WHERE clause of the SQL query
func (qb *MySQLQueryBuilder) Build() string {
	whereClause := strings.Join(qb.whereClauses, " AND ")
	return whereClause
}

// BuildFullQuery constructs the full SQL query with WHERE, ORDER BY, LIMIT, and OFFSET clauses
func (qb *MySQLQueryBuilder) BuildFullQueryOn(tableName string) string {
	whereClause := qb.Build()

	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT ")
	if len(qb.fields) > 0 {
		queryBuilder.WriteString(strings.Join(qb.fields, ", "))
	} else {
		queryBuilder.WriteString("*") // Select all fields if no specific fields are provided
	}
	queryBuilder.WriteString(fmt.Sprintf(" FROM %s", tableName))
	if whereClause != "" {
		queryBuilder.WriteString(fmt.Sprintf(" WHERE %s", whereClause))
	}
	if qb.orderClause != "" {
		queryBuilder.WriteString(fmt.Sprintf(" %s", qb.orderClause))
	}
	if qb.limit > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))
	}
	if qb.offset > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" OFFSET %d", qb.offset))
	}

	return queryBuilder.String()
}

// NewMySQLQueryBuilder creates a new instance of MySQLQueryBuilder
func NewMySQLQueryBuilder() Interfaces.SQLQueryBuilder {
	return &MySQLQueryBuilder{
		whereClauses: []string{"1 = 1"},
		orderClause:  "",
		offset:       0,
		limit:        0,
		fields:       nil,
	}
}
