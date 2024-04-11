package model

import "errors"

var (
	ErrorNotFound            = errors.New("Баннер не найден")
	ErrorInternalServerError = errors.New("Внутренняя ошибка сервера")
	ErrorBadRequest          = errors.New("Некорректные данные")
)
