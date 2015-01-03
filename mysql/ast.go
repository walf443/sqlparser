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
	}

	ColumnDefinition struct {
		DataTypeDefinition DataTypeDefinition
		Nullable           bool
		AutoIncrement      bool
		Default            DefaultDefinition
	}

	DataTypeDefinition interface {
		data_type_definition()
	}

	CreateDefinition interface {
		create_definition()
	}

	DefaultDefinition interface {
		default_definition()
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
	}

	CommentStatement struct {
		Content string
	}
)

func (x *DropTableStatement) statement()      {}
func (x *DropTableStatement) ToQuery() string {
	var tableNames []string
	for _,table := range x.TableNames {
		tableNames = append(tableNames, table.ToQuery())
	}
	return "DROP TABLE " + strings.Join(tableNames, ", ")
}

func (x *DropDatabaseStatement) statement()   {}
func (x *DropDatabaseStatement) ToQuery() string {
	return "TODO"
}
func (x *CreateDatabaseStatement) statement() {}
func (x *CreateDatabaseStatement) ToQuery() string {
	return "TODO"
}
func (x *AlterTableStatement) statement()     {}
func (x *AlterTableStatement) ToQuery() string {
	return "TODO"
}
func (x *CreateTableStatement) statement()    {}
func (x *CreateTableStatement) ToQuery() string {
	return "TODO"
}
func (x *CommentStatement) statement()        {}
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

func (x *TableNameIdentifier) identifier()    {}

func (x *TableNameIdentifier) ToQuery() string {
	if x.Database == "" {
		return "`" + x.Name + "`"
	} else {
		return fmt.Sprintf("`%s`.`%s`", x.Database, x.Name)
	}
}

func (x *DatabaseNameIdentifier) identifier() {}
func (x *ColumnNameIdentifier) identifier()   {}
func (x *IndexNameIdentifier) identifier()    {}

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
func (x *AlterSpecificationDropIndex) alterspecification()  {}
func (x *AlterSpecificationAddColumn) alterspecification()  {}

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

func (x *DataTypeDefinitionSimple) data_type_definition()   {}
func (x *DataTypeDefinitionNumber) data_type_definition()   {}
func (x *DataTypeDefinitionFraction) data_type_definition() {}
func (x *DataTypeDefinitionString) data_type_definition()   {}
func (x *DataTypeDefinitionTextBlob) data_type_definition() {}

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
func (x *CreateDefinitionPrimaryIndex) create_definition() {}
func (x *CreateDefinitionUniqueIndex) create_definition()  {}
func (x *CreateDefinitionIndex) create_definition()        {}

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

func (x *DefaultDefinitionEmpty) default_definition()            {}
func (x *DefaultDefinitionNull) default_definition()             {}
func (x *DefaultDefinitionString) default_definition()           {}
func (x *DefaultDefinitionCurrentTimestamp) default_definition() {}
