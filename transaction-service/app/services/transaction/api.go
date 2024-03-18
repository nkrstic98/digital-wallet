package transaction

type Service interface {
	AddMoney(input AddMoneyInput) (float64, error)
	TransferMoney(input TransferMoneyInput) error
}
