package mysql

type (
	Statement interface {
		statement()
	}

	Identifier interface {
		identifier()
	}

	AlterSpecification interface {
		alterspecification()
	}

	ColumnDefinition struct {
		DataTypeDefinition DataTypeDefinition
	}

	DataTypeDefinition interface {
		data_type_definition()
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
		TableName TableNameIdentifier
		AlterSpecifications []AlterSpecification
	}
)

func (x *DropTableStatement) statement() {}
func (x *DropDatabaseStatement) statement() {}
func (x *CreateDatabaseStatement) statement() {}
func (x *AlterTableStatement) statement() {}

type (
	TableNameIdentifier struct {
		Name string
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
)

func (x *TableNameIdentifier) identifier() {}
func (x *DatabaseNameIdentifier) identifier() {}
func (x *ColumnNameIdentifier) identifier() {}
func (x *IndexNameIdentifier) identifier() {}

type (
	AlterSpecificationDropColumn struct {
		ColumnName ColumnNameIdentifier
	}
	AlterSpecificationDropIndex struct {
		IndexName IndexNameIdentifier
	}
	AlterSpecificationAddColumn struct {
		ColumnName ColumnNameIdentifier
		ColumnDefinition ColumnDefinition
	}
)

func (x *AlterSpecificationDropColumn) alterspecification() {}
func (x *AlterSpecificationDropIndex) alterspecification() {}
func (x *AlterSpecificationAddColumn) alterspecification() {}

type (
	DataTypeDefinitionSimple struct {
		Type DataType
	}
	DataTypeDefinitionNumber struct {
		Type DataType
		Length uint
		Unsigned bool
		Zerofill bool
	}
	DataTypeDefinitionFraction struct {
		Type DataType
		Length uint
		Decimals uint
		Unsigned bool
		Zerofill bool
	}
	DataTypeDefinitionString struct {
		Type DataType
		Length uint
		CharsetName string
		CollationName string
	}
	DataTypeDefinitionTextBlob struct {
		Type DataType
		Binary bool
		CharsetName string
		CollationName string
	}
)

func (x *DataTypeDefinitionSimple) data_type_definition() {}
func (x *DataTypeDefinitionNumber) data_type_definition() {}
func (x *DataTypeDefinitionFraction) data_type_definition() {}
func (x *DataTypeDefinitionString) data_type_definition() {}
func (x *DataTypeDefinitionTextBlob) data_type_definition() {}
