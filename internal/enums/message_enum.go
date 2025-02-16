package enums

type MessageType int

const (
	SuccessfulTransfer MessageType = iota
	SuccessfulPurchase
)

var messages = map[MessageType]string{
	SuccessfulTransfer: "перевод выполнен успешно",
	SuccessfulPurchase: "покупка прошла успешно",
}

func (mt MessageType) String() string {
	return messages[mt]
}
