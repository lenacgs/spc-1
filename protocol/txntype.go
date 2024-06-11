package protocol

type TxnType = int

const (
	UnknownType TxnType = 0
	UpsertType  TxnType = 1
	SelectType  TxnType = 2
	DeleteType  TxnType = 3
)
