package enum

type ErrorType string

const (
	ErrInsufficientMoney        ErrorType = "недостаточно монет"
	ErrReceiverNotFound         ErrorType = "пользователь к переводу не нашелся"
	ErrUserNotAuthorized        ErrorType = "пользователь не авторизован"
	ErrWrongReqFormat           ErrorType = "неверный формат запроса"
	ErrCoinsInappropriateAmount ErrorType = "количество монет должно быть больше нуля"
	ErrReceivingCoinsInfo       ErrorType = "ошибка получения информации о монетах"
	ErrInvalidToken             ErrorType = "неверный или просроченный токен"
	ErrNoUsernameAndPassword    ErrorType = "имя пользователя и пароль обязательны"
	ErrReceivingTransferHistory ErrorType = "ошибка получения истории переводов"
	ErrReceivingPurchaseHistory ErrorType = "ошибка получения информации о покупках"
	ErrBuyWithInsufficientMoney ErrorType = "недостаточно монет для покупки"
	ErrItemNotFound             ErrorType = "товар не найден"
	ErrNotProvidedItem          ErrorType = "не указан товар для покупки"
	ErrGeneratingToken          ErrorType = "ошибка генерации токена"
	ErrWrongCredentials         ErrorType = "неверное имя пользователя или пароль"
	ErrInternalServer           ErrorType = "ошибка сервера"
	ErrCreatingUser             ErrorType = "ошибка создания пользователя"
	ErrNoAuthToken              ErrorType = "нет токена авторизации"
	ErrWrongTokenFormat         ErrorType = "неверный формат токена"
	ErrEqualReceivers           ErrorType = "получатели должны отличаться друг от друга"
)

func (et ErrorType) Error() string {
	return string(et)
}
