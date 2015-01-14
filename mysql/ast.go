package mysql

import (
	"fmt"
	"strings"
)

type (
	Statement interface {
		statement()
		ToQuery() string
	}

	Identifier interface {
		identifier()
	}

	AlterSpecification interface {
		alterspecification()
		ToQuery() string
	}

	ColumnDefinition struct {
		DataTypeDefinition DataTypeDefinition
		Nullable           bool
		AutoIncrement      bool
		Default            DefaultDefinition
	}

	DataTypeDefinition interface {
		data_type_definition()
		ToQuery() string
	}

	CreateDefinition interface {
		create_definition()
		ToQuery() string
	}

	DefaultDefinition interface {
		default_definition()
		ToQuery() string
	}
)

type (
	DropTableStatement struct {
		TableNames []TableNameIdentifier
	}
	DropDatabaseStatement struct {
		DatabaseName DatabaseNameIdentifier
	}
	CreateDatabaseStatement struct {
		DatabaseName DatabaseNameIdentifier
	}
	AlterTableStatement struct {
		TableName           TableNameIdentifier
		AlterSpecifications []AlterSpecification
	}

	CreateTableStatement struct {
		TableName         TableNameIdentifier
		CreateDefinitions []CreateDefinition
		TableOptions      []TableOption
	}

	CommentStatement struct {
		Content string
	}
)

func (x *DropTableStatement) statement() {}
func (x *DropTableStatement) ToQuery() string {
	var tableNames []string
	for _, table := range x.TableNames {
		tableNames = append(tableNames, table.ToQuery())
	}
	return "DROP TABLE " + strings.Join(tableNames, ", ") + ";"
}

func (x *DropDatabaseStatement) statement() {}
func (x *DropDatabaseStatement) ToQuery() string {
	return "DROP DATABASE " + x.DatabaseName.ToQuery() + ";"
}
func (x *CreateDatabaseStatement) statement() {}
func (x *CreateDatabaseStatement) ToQuery() string {
	return "CREATE DATABASE " + x.DatabaseName.ToQuery() + ";"
}

func (x *AlterTableStatement) statement() {}
func (x *AlterTableStatement) ToQuery() string {
	var specQueries []string
	for _, spec := range x.AlterSpecifications {
		specQueries = append(specQueries, spec.ToQuery())
	}
	return "ALTER TABLE " + x.TableName.ToQuery() + " " + strings.Join(specQueries, ", ") + ";"
}
func (x *CreateTableStatement) statement() {}
func (x *CreateTableStatement) ToQuery() string {
	var options []string
	for _, option := range x.TableOptions {
		options = append(options, option.ToQuery())
	}
	var defs []string
	for _, def := range x.CreateDefinitions {
		defs = append(defs, def.ToQuery())
	}
	return "CREATE TABLE " + x.TableName.ToQuery() + " (\n\t"  +  strings.Join(defs, ",\n\t") + "\n) " + strings.Join(options, " ") + ";"
}
func (x *CommentStatement) statement() {}
func (x *CommentStatement) ToQuery() string {
	return "TODO"
}

type (
	TableNameIdentifier struct {
		Name     string
		Database string
	}
	DatabaseNameIdentifier struct {
		Name string
	}
	ColumnNameIdentifier struct {
		Name string
	}
	IndexNameIdentifier struct {
		Name string
	}

	EngineNameIdentifier struct {
		Name string
	}
)

func (x *TableNameIdentifier) identifier() {}

func (x *TableNameIdentifier) ToQuery() string {
	if x.Database == "" {
		return "`" + x.Name + "`"
	} else {
		return fmt.Sprintf("`%s`.`%s`", x.Database, x.Name)
	}
}

func (x *DatabaseNameIdentifier) identifier() {}
func (x *DatabaseNameIdentifier) ToQuery() string {
	return "`" + x.Name + "`"
}
func (x *ColumnNameIdentifier) identifier() {}
func (x *ColumnNameIdentifier) ToQuery() string {
	return "`" + x.Name + "`"
}

func (x *IndexNameIdentifier) identifier() {}
func (x *IndexNameIdentifier) ToQuery() string {
	return "`" + x.Name + "`"
}

type (
	AlterSpecificationDropColumn struct {
		ColumnName ColumnNameIdentifier
	}
	AlterSpecificationDropIndex struct {
		IndexName IndexNameIdentifier
	}
	AlterSpecificationAddColumn struct {
		ColumnName       ColumnNameIdentifier
		ColumnDefinition ColumnDefinition
	}
)

func (x *AlterSpecificationDropColumn) alterspecification() {}
func (x *AlterSpecificationDropColumn) ToQuery() string {
	return "DROP " + x.ColumnName.ToQuery()
}

func (x *AlterSpecificationDropIndex) alterspecification() {}
func (x *AlterSpecificationDropIndex) ToQuery() string {
	return "DROP INDEX " + x.IndexName.ToQuery()
}
func (x *AlterSpecificationAddColumn) alterspecification() {}
func (x *AlterSpecificationAddColumn) ToQuery() string {
	return "ADD " + x.ColumnName.ToQuery() + " " + x.ColumnDefinition.ToQuery()
}

func (x ColumnDefinition) ToQuery() string {
	result := ""
	result += x.DataTypeDefinition.ToQuery()
	if !x.Nullable {
		result += " NOT NULL"
	}
	if x.AutoIncrement {
		result += " AUTO_INCREMENT"
	}
	result += " " + x.Default.ToQuery()

	return result
}

type (
	DataTypeDefinitionSimple struct {
		Type DataType
	}
	DataTypeDefinitionNumber struct {
		Type     DataType
		Length   uint
		Unsigned bool
		Zerofill bool
	}
	DataTypeDefinitionFraction struct {
		Type     DataType
		Length   uint
		Decimals uint
		Unsigned bool
		Zerofill bool
	}
	DataTypeDefinitionString struct {
		Type          DataType
		Length        uint
		CharsetName   string
		CollationName string
	}
	DataTypeDefinitionTextBlob struct {
		Type          DataType
		Binary        bool
		CharsetName   string
		CollationName string
	}
)

func (x *DataTypeDefinitionSimple) data_type_definition() {}
func (x *DataTypeDefinitionSimple) ToQuery() string {
	return x.Type.String()
}
func (x *DataTypeDefinitionNumber) data_type_definition() {}
func (x *DataTypeDefinitionNumber) ToQuery() string {
	result := x.Type.String()
	if x.Length != 0 {
		result += fmt.Sprintf("(%d)", x.Length)
	}
	if x.Unsigned {
		result += " UNSIGNED"
	}
	if x.Zerofill {
		result += " ZEROFILL"
	}
	return result
}

func (x *DataTypeDefinitionFraction) data_type_definition() {}
func (x *DataTypeDefinitionFraction) ToQuery() string {
	result := x.Type.String()
	if x.Decimals == 0 {
		if x.Length != 0 {
			result += fmt.Sprintf("(%d)", x.Length)
		}
	} else {
		result += fmt.Sprintf("(%d, %d)", x.Length, x.Decimals)
	}
	if x.Unsigned {
		result += " UNSIGNED"
	}
	if x.Zerofill {
		result += " ZEROFILL"
	}
	return result
}
func (x *DataTypeDefinitionString) data_type_definition() {}
func (x *DataTypeDefinitionString) ToQuery() string {
	result := x.Type.String()
	if x.Length > 0 {
		result += fmt.Sprintf("(%d)", x.Length)
	}
	if x.CharsetName != "" {
		result += fmt.Sprintf(" CHARACTER SET %s", x.CharsetName)
	}
	if x.CollationName != "" {
		result += fmt.Sprintf(" COLLATE %s", x.CollationName)
	}
	return result
}

func (x *DataTypeDefinitionTextBlob) data_type_definition() {}
func (x *DataTypeDefinitionTextBlob) ToQuery() string {
	result := x.Type.String()
	if x.Binary {
		result += " BINARY"
	}
	if x.CharsetName != "" {
		result += fmt.Sprintf(" CHARACTER SET %s", x.CharsetName)
	}
	if x.CollationName != "" {
		result += fmt.Sprintf(" COLLATE %s", x.CollationName)
	}
	return result
}

type (
	CreateDefinitionColumn struct {
		ColumnName       ColumnNameIdentifier
		ColumnDefinition ColumnDefinition
	}
	CreateDefinitionPrimaryIndex struct {
		Columns []ColumnNameIdentifier
	}

	CreateDefinitionUniqueIndex struct {
		Name    IndexNameIdentifier
		Columns []ColumnNameIdentifier
	}
	CreateDefinitionIndex struct {
		Name    IndexNameIdentifier
		Columns []ColumnNameIdentifier
	}
)

func (x *CreateDefinitionColumn) create_definition()       {}
func (x *CreateDefinitionColumn) ToQuery() string      {
	return x.ColumnName.ToQuery() + " " + x.ColumnDefinition.ToQuery()
}
func (x *CreateDefinitionPrimaryIndex) create_definition() {}
func (x *CreateDefinitionPrimaryIndex) ToQuery() string {
	var columns []string
	for _, column := range x.Columns {
		columns = append(columns, column.ToQuery())
	}
	return "PRIMARY KEY ( " + strings.Join(columns, ",") +  " )"
}
func (x *CreateDefinitionUniqueIndex) create_definition()  {}
func (x *CreateDefinitionUniqueIndex) ToQuery() string {
	var columns []string
	for _, column := range x.Columns {
		columns = append(columns, column.ToQuery())
	}
	name := ""
	if x.Name.Name != "" {
		name = x.Name.ToQuery()
	}
	return "UNIQUE KEY " + name + " ( " + strings.Join(columns, ",") +  " )"
}
func (x *CreateDefinitionIndex) create_definition()        {}
func (x *CreateDefinitionIndex) ToQuery() string {
	var columns []string
	for _, column := range x.Columns {
		columns = append(columns, column.ToQuery())
	}
	name := ""
	if x.Name.Name != "" {
		name = x.Name.ToQuery()
	}
	return "INDEX " + name + " ( " + strings.Join(columns, ",") +  " )"
}

type (
	DefaultDefinitionString struct {
		Value string
	}

	DefaultDefinitionEmpty struct {
	}

	DefaultDefinitionNull struct {
	}

	DefaultDefinitionCurrentTimestamp struct {
		OnUpdate bool
	}
)

func (x *DefaultDefinitionEmpty) default_definition() {}
func (x *DefaultDefinitionEmpty) ToQuery() string {
	return ""
}
func (x *DefaultDefinitionNull) default_definition() {}
func (x *DefaultDefinitionNull) ToQuery() string {
	return "DEFAULT NULL"
}
func (x *DefaultDefinitionString) default_definition() {}
func (x *DefaultDefinitionString) ToQuery() string {
	return "DEFAULT \"" + x.Value + "\""
}
func (x *DefaultDefinitionCurrentTimestamp) default_definition() {}
func (x *DefaultDefinitionCurrentTimestamp) ToQuery() string {
	if x.OnUpdate {
		return "DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"
	} else {
		return "DEFAULT CURRENT_TIMESTAMP"
	}
}

type TableOption struct {
	Key   string
	Value string
}

func (x *TableOption) ToQuery() string {
	return x.Key + "=" + x.Value
}
