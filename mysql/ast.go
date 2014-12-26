package mysql

type (
	Statement interface {
		statement()
	}

	Identifier interface {
		identifier()
	}
)

type (
	DropTableStatement struct {
		TableNames []TableNameIdentifier
	}
)

func (x *DropTableStatement) statement() {}

type (
	TableNameIdentifier struct {
		Name string
		Database string
	}
)

func (x *TableNameIdentifier) identifier() {}
