package token

type Token struct {
	token string
	// по-хорошему за хранение, проверку и выдачу токенов отвечает отдельный сервис
	// сделал заглушку с интерфейсом поверх, чтобы на остаьной сервис не влияло
	// будет время прикручу авторизацию
}

const (
	JohnSmith = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjAsImlzQWRtaW4iOmZhbHNlLCJuYW1lIjoiSm9obiBTbWl0aCJ9.l5PzO3w-1HuD3zoq85oujufTSzFFTaV9zmSU3zQpxNo"
	JaneDoe   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjAsImlzQWRtaW4iOnRydWUsIm5hbWUiOiJKYW5lIERvZSJ9.-A6bF348vryjwST2vccaW2sgGO6bh7AzmmABdiGKhz0"
)

func (t *Token) Exist() bool {
	switch t.token {
	case JohnSmith:
		return true
	case JaneDoe:
		return true
	default:
		return false
	}
}

func (t *Token) IsAdmin() bool {
	switch t.token {
	case JohnSmith:
		return false
	case JaneDoe:
		return true
	default:
		return false
	}
}

func (t *Token) Make(token string) {
	t.token = token
}
