package token

type IToken interface {
	Exist() bool
	IsAdmin() bool
	Make(token string)
}

func New(token string) IToken {
	t := Token{}
	t.Make(token)
	return &t
}
