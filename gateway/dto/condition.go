package dto

type CreateConditionDto struct {
	PaymentAmount  int    `json:"payment_amount"`
	PaymentRateId  int    `json:"payment_rate_id"`
	ContractTypeId int    `json:"contract_type_id"`
	LocationTypeId int    `json:"location_type_id"`
	Loc            string `json:"loc"`
}

type ConditionResponse struct {
	Id            int    `json:"id"`
	PaymentAmount int    `json:"payment_amount"`
	PaymentRate   string `json:"payment_rate"`
	ContractType  string `json:"contract_type"`
	Location      string `json:"location"`
}
