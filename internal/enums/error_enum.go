package enums

type ErrorType int

const (
	ErrInsufficientMoney ErrorType = iota
	ErrReceiverNotFound
	ErrUserNotAuthorized
	ErrWrongReqFormat
	ErrCoinsInappropriateAmount
	ErrReceivingCoinsInfo
	ErrInvalidToken
	ErrNoUsernameAndPassword
	ErrReceivingTransferHistory
	ErrReceivingPurchaseHistory
	ErrBuyWithInsufficientMoney
	ErrItemNotFound
	ErrNotProvidedItem
	ErrGeneratingToken
	ErrWrongCredentials
	ErrInternalServer
	ErrCreatingUser
	ErrNoAuthToken
	ErrWrongTokenFormat
)

var errorMessages = map[ErrorType]string{
	ErrInsufficientMoney:        "недостаточно монет",
	ErrReceiverNotFound:         "пользователь к переводу не нашелся",
	ErrUserNotAuthorized:        "пользователь не авторизован",
	ErrWrongReqFormat:           "неверный формат запроса",
	ErrCoinsInappropriateAmount: "количество монет должно быть больше нуля",
	ErrReceivingCoinsInfo:       "ошибка получения информации о монетах",
	ErrInvalidToken:             "неверный или просроченный токен",
	ErrNoUsernameAndPassword:    "имя пользователя и пароль обязательны",
	ErrReceivingTransferHistory: "ошибка получения истории переводов",
	ErrReceivingPurchaseHistory: "ошибка получения информации о покупках",
	ErrBuyWithInsufficientMoney: "недостаточно монет для покупки",
	ErrItemNotFound:             "товар не найден",
	ErrNotProvidedItem:          "не указан товар для покупки",
	ErrGeneratingToken:          "ошибка генерации токена",
	ErrWrongCredentials:         "неверное имя пользователя или пароль",
	ErrInternalServer:           "ошибка сервера",
	ErrCreatingUser:             "ошибка создания пользователя",
	ErrNoAuthToken:              "нет токена авторизации",
	ErrWrongTokenFormat:         "неверный формат токена",
}

func (et ErrorType) Error() string {
	return errorMessages[et]
}
