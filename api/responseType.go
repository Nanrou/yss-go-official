package api

type rightContentResponse struct {
	Content string `json:"content"`
}

func newRightContentResponse(content string) *rightContentResponse {
	return &rightContentResponse{Content: content}
}

// account的基础信息
type accountStation struct {
	Account           string
	Address           string
	Name              string
	Phone             string
	Charge            string
	CurrentMeter      string
	Meter             string
	Paid              bool
	Default_          bool
	UnpaidPeriodCount int
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

type feeDetail struct {
	Address          string
	Name             string
	Charge           string
	CurrentMeter     string
	MeterReadingDate string
	Paid             bool
	PreviousMeter    string
	WaterCharge      string
	WaterProperty    string
	Wsf              string
	Xfft             string
	Ljf              string
	Ecjydf           string
	Szyf             string
	Cjhys            string
	Wyj              string
	Wswyj            string
}

type accountWithoutBindingResponse struct {
	Account       string   `json:"account"`
	Name          string   `json:"account_name"`
	Charge        string   `json:"charge"`
	CurrentMeter  string   `json:"current_meter"`
	CurrentPeriod string   `json:"current_period"`
	Paid          bool     `json:"paid"`
	PreviousMeter string   `json:"previous_meter"`
	OtherFee      otherFee `json:"other_fee"`
	WaterCharge   string   `json:"water_charge"`
}

type otherFee struct {
	Wsf    string `json:"污水费"`
	Xfft   string `json:"消防分摊费"`
	Ljf    string `json:"垃圾费"`
	Ecjydf string `json:"二次加压电费"`
	Szyf   string `json:"水资源费"`
	Cjhys  string `json:"超计划用水费"`
	Wyj    string `json:"违约金"`
	Wswyj  string `json:"污水违约金"`
}

func newAccountWithoutBindingResponse(fd *feeDetail, account string, date string) *accountWithoutBindingResponse {
	return &accountWithoutBindingResponse{
		Account:       account,
		Name:          fd.Name,
		Charge:        fd.Charge,
		CurrentMeter:  fd.CurrentMeter,
		CurrentPeriod: date,
		Paid:          fd.Paid,
		PreviousMeter: fd.PreviousMeter,
		OtherFee: otherFee{
			Wsf:    fd.Wsf,
			Xfft:   fd.Xfft,
			Ljf:    fd.Ljf,
			Ecjydf: fd.Ecjydf,
			Szyf:   fd.Szyf,
			Cjhys:  fd.Cjhys,
			Wyj:    fd.Wyj,
			Wswyj:  fd.Wswyj,
		},
		WaterCharge: fd.WaterCharge,
	}
}

type feeHistoryResponse struct {
	Account           string         `json:"account"`
	Address           string         `json:"account_address"`
	Name              string         `json:"account_name"`
	Charge            string         `json:"charge"`
	CurrentMeter      string         `json:"current_meter"`
	Meter             string         `json:"meter"`
	Paid              bool           `json:"paid"`
	PaidHistory       [] paidHistory `json:"paid_history"`
	UnpaidPeriodCount int            `json:"unpaid_period_count"`
}

type paidHistory struct {
	Date   string `json:"date"`
	Charge string `json:"charge"`
	BillId string `json:"bill_id"`
}

func newFeeHistoryResponse(as *accountStation, ph [] paidHistory) *feeHistoryResponse {
	return &feeHistoryResponse{
		Account:           as.Account,
		Address:           as.Address,
		Name:              as.Name,
		Charge:            as.Charge,
		CurrentMeter:      as.CurrentMeter,
		Meter:             as.Meter,
		Paid:              as.Paid,
		PaidHistory:       ph,
		UnpaidPeriodCount: as.UnpaidPeriodCount,
	}
}

type feeHistoryWithoutPaidResponse struct {
	Account       string           `json:"account"`
	Name          string           `json:"account_name"`
	UnpaidHistory [] unpaidHistory `json:"unpaid_history"`
}

type unpaidHistory struct {
	CurrentPeriod string   `json:"current_period"`
	Charge        string   `json:"charge"`
	CurrentMeter  string   `json:"current_meter"`
	PreviousMeter string   `json:"previous_meter"`
	OtherFee      otherFee `json:"other_fee"`
}

func newUnpaidHistory (detail *feeDetail, date string) *unpaidHistory {
	return &unpaidHistory{
		CurrentPeriod: date,
		Charge:        detail.Charge,
		CurrentMeter:  detail.CurrentMeter,
		PreviousMeter: detail.PreviousMeter,
		OtherFee: otherFee{
			Wsf:    detail.Wsf,
			Xfft:   detail.Xfft,
			Ljf:    detail.Ljf,
			Ecjydf: detail.Ecjydf,
			Szyf:   detail.Szyf,
			Cjhys:  detail.Cjhys,
			Wyj:    detail.Wyj,
			Wswyj:  detail.Wswyj,
		},
	}
}

func newFeeHistoryWithoutPaidResponse(as *accountStation, history [] unpaidHistory) *feeHistoryWithoutPaidResponse {
	return &feeHistoryWithoutPaidResponse{
		Account:       as.Account,
		Name:          as.Name,
		UnpaidHistory: history,
	}
}
