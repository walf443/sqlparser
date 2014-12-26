package mysql

type (
	Statement interface {
		statement()
	}

	Expression interface {
		expression()
	}
)

type (
	DropTableStatement struct {
		TableNames []TableNameExpression
	}
)

func (x *DropTableStatement) statement() {}

type (
	TableNameExpression struct {
		Lit string
	}
)

func (x *TableNameExpression) expression() {}
