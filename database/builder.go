package database

import "bytes"

type JoinType uint8
type Relation uint8
type SortType uint8

const (
	Left  JoinType = 0
	Right JoinType = 1
	Inner JoinType = 2
	Full  JoinType = 3

	AND Relation = 0
	OR 	Relation = 1

	ASC  SortType = 0
	DESC SortType = 1
	Auto SortType = 2
)

type Table struct {
	db string
	name string
	alias string
}

type Field struct {
	table *Table
	name  string
	alias string
}

type Condition struct {
	relation Relation
	conditions []*Condition
}

type Join struct {
	table *Table
	style JoinType
	condition *Condition
}

type Sort struct {
	column 	*Field
	sort 	SortType
	detail  string
}

func T(name string) *Table{
	return TAlias(name, "")
}

func TAlias(name, alias string) *Table{
	return TDBAlias("", name, alias)
}

func TDBAlias(db, name, alias string) *Table {
	return &Table{
		db:    db,
		name:  name,
		alias: alias,
	}
}

func JoinTable(table *Table, condition *Condition, style JoinType) *Join {
	switch style {
	case Left:
		return LeftJoin(table, condition)
	case Right:
		return RightJoin(table, condition)
	case Inner:
		return InnerJoin(table, condition)
	case Full:
		return FullJoin(table, condition)
	default:
		return LeftJoin(table, condition)
	}
}

func LeftJoin(table *Table, condition *Condition) *Join {
	return &Join{
		table:     table,
		style:     Left,
		condition: condition,
	}
}

func RightJoin(table *Table, condition *Condition) *Join {
	return &Join{
		table:     table,
		style:     Right,
		condition: condition,
	}
}

func InnerJoin(table *Table, condition *Condition) *Join {
	return &Join{
		table:     table,
		style:     Inner,
		condition: condition,
	}
}

func FullJoin(table *Table, condition *Condition) *Join {
	return &Join{
		table:     table,
		style:     Full,
		condition: condition,
	}
}

func (table *Table) String(delimiters []byte) string {
	if delimiters == nil || len(delimiters) < 1 {
		delimiters = []byte{0x00, 0x00}
	} else if len(delimiters) == 1{
		delimiters = append(delimiters, delimiters[0])
	}

	var buffer bytes.Buffer
	if table.db != "" {
		buffer.WriteByte(delimiters[0])
		buffer.WriteString(table.db)
		buffer.WriteByte(delimiters[1])
		buffer.WriteByte(0x2e)
	}

	buffer.WriteByte(delimiters[0])
	buffer.WriteString(table.name)
	buffer.WriteByte(delimiters[1])

	if table.alias != "" {
		buffer.WriteString(" AS ")
		buffer.WriteByte(delimiters[0])
		buffer.WriteString(table.alias)
		buffer.WriteByte(delimiters[1])
	}

	return buffer.String()
}

func (field *Field) String(delimiters []byte) string {
	if delimiters == nil || len(delimiters) < 1 {
		delimiters = []byte{0x00, 0x00}
	} else if len(delimiters) == 1{
		delimiters = append(delimiters, delimiters[0])
	}

	var buffer bytes.Buffer
	buffer.WriteString(field.table.String(delimiters))
	buffer.WriteByte(0x2e)
	buffer.WriteByte(delimiters[0])
	buffer.WriteString(field.name)
	buffer.WriteByte(delimiters[1])

	if field.alias != "" {
		buffer.WriteString(" AS ")
		buffer.WriteByte(delimiters[0])
		buffer.WriteString(field.alias)
		buffer.WriteByte(delimiters[1])
	}

	return buffer.String()
}

func (sort *Sort) String(delimiters []byte) string {
	var (
		field = sort.column.String(delimiters)
		buffer bytes.Buffer
	)
	buffer.WriteString(field)
	buffer.WriteByte(0x20)
	if sort.sort == ASC {
		buffer.WriteString("ASC")
	} else if sort.sort == DESC {
		buffer.WriteString("DESC")
	} else {
		buffer.Truncate(0)
		buffer.WriteString("FIELD(")
		buffer.WriteString(field)
		buffer.WriteString(", ")
		buffer.WriteString(sort.detail)
	}

	return buffer.String()
}

type QueryBuilder struct {
	columns 	[]*Field
	tables 		[]*Table
	joins 		[]*Join
	orders 		[]*Sort
	groups 		[]*Field
	having		*Condition
	conditions 	*Condition
	limit 		uint
	offset 		uint
	forUpdate 	bool
	delimiters  []byte
}

func CreateQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		columns:    make([]*Field, 0),
		tables:     make([]*Table, 0),
		joins:      make([]*Join, 0),
		orders:     make([]*Sort, 0),
		groups:     make([]*Field, 0),
		having:     nil,
		conditions: nil,
		limit:      0,
		offset:     0,
		forUpdate:  false,
	}
}

func CreateQueryBuilderForTable(table *Table) *QueryBuilder {
	builder := &QueryBuilder{
		columns:    make([]*Field, 0),
		tables:     make([]*Table, 0),
		joins:      make([]*Join, 0),
		orders:     make([]*Sort, 0),
		groups:     make([]*Field, 0),
		having:     nil,
		conditions: nil,
		limit:      0,
		offset:     0,
		forUpdate:  false,
	}

	if table != nil {
		return builder.From(table)
	}

	return builder
}

func (builder *QueryBuilder) Columns(fields ...*Field) *QueryBuilder{
	builder.columns = fields
	return builder
}

func (builder *QueryBuilder) From(table *Table) *QueryBuilder{
	if len(builder.tables) > 0 {
		builder.tables[0] = table
	}

	builder.tables = append(builder.tables, table)
	return builder
}

func (builder *QueryBuilder) AddTable(table *Table) *QueryBuilder{
	builder.tables = append(builder.tables, table)
	return builder
}

func (builder *QueryBuilder) Where(condition *Condition) *QueryBuilder{
	builder.conditions = condition
	return builder
}

func (builder *QueryBuilder) Join(join *Join) *QueryBuilder{
	builder.joins = append(builder.joins, join)
	return builder
}

func (builder *QueryBuilder) Having(condition *Condition) *QueryBuilder{
	builder.having = condition
	return builder
}

func (builder *QueryBuilder) SortBy(sort *Sort) *QueryBuilder{
	builder.orders = append(builder.orders, sort)
	return builder
}

func (builder *QueryBuilder) GroupBy(fields ...*Field) *QueryBuilder{
	builder.groups = fields
	return builder
}

func (builder *QueryBuilder) Page(page, size uint) *QueryBuilder{
	builder.limit = size
	builder.offset = page * size

	return builder
}

func (builder *QueryBuilder) ForUpdate() *QueryBuilder {
	builder.forUpdate = true
	return builder
}