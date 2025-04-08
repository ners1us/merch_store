package enum

type MessageType string

const (
	SuccessfulTransfer MessageType = "перевод выполнен успешно"
	SuccessfulPurchase MessageType = "покупка прошла успешно"
)

func (mt MessageType) String() string {
	return string(mt)
}
