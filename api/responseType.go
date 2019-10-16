package api

type rightContentResponse struct {
	Content string `json:"content"`
}

func newRightContentResponse(content string) *rightContentResponse {
	return &rightContentResponse{Content: content}
}

// account的基础信息
type accountStation struct {
	Account           string `json:"account"`
	Address           string `json:"account_address"`
	Name              string `json:"account_name"`
	Phone             string `json:"account_phone"`
	Charge            string `json:"charge"`
	CurrentMeter      string `json:"current_meter"`
	Meter             string `json:"meter"`
	Paid              bool   `json:"paid"`
	Default_          bool   `json:"default"`
	UnpaidPeriodCount int    `json:"unpaid_period_count"`
}

type accountsResponse struct {
	Account           string `json:"account"`
	Address           string `json:"account_address"`
	Name              string `json:"account_name"`
	Charge            string `json:"charge"`
	CurrentMeter      string `json:"current_meter"`
	Meter             string `json:"meter"`
	Paid              bool   `json:"paid"`
	Default_          bool   `json:"default"`
	UnpaidPeriodCount int    `json:"unpaid_period_count"`
}

func newAccountsResponse(as *accountStation) *accountsResponse {
	return &accountsResponse{
		Account:           as.Account,
		Address:           as.Address,
		Name:              as.Name,
		Charge:            as.Charge,
		CurrentMeter:      as.CurrentMeter,
		Meter:             as.Meter,
		Paid:              as.Paid,
		Default_:          as.Default_,
		UnpaidPeriodCount: as.UnpaidPeriodCount,
	}
}

type defaultAccountResponse struct {
	Account      string `json:"account"`
	Address      string `json:"account_address"`
	Name         string `json:"account_name"`
	Charge       string `json:"charge"`
	CurrentMeter string `json:"current_meter"`
	Meter        string `json:"meter"`
	Paid         bool   `json:"paid"`
}

func newDefaultAccountResponse(as *accountStation) *defaultAccountResponse {
	return &defaultAccountResponse{
		Account:      as.Account,
		Address:      as.Address,
		Name:         as.Name,
		Charge:       as.Charge,
		CurrentMeter: as.CurrentMeter,
		Meter:        as.Meter,
		Paid:         as.Paid,
	}
}
