package models

type DetailRequest struct {
	UsageCommon       int64                 `json:"usageCommon"`
	VoiceUsage        string                `json:"voiceUsage"`
	Used              int64                 `json:"used"`
	Xvalue            int                   `json:"xvalue"`
	ParaFieldResult   string                `json:"paraFieldResult"`
	UsedCommon        int                   `json:"usedCommon"`
	Result            int                   `json:"result"`
	Total             int64                 `json:"total"`
	VoiceBalance      string                `json:"voiceBalance"`
	Balance           int64                 `json:"balance"`
	ServiceResultCode int                   `json:"serviceResultCode"`
	VoiceAmount       string                `json:"voiceAmount"`
	IsUnlimit         string                `json:"isUnlimit"`
	BalanceCommon     int                   `json:"balanceCommon"`
	TotalCommon       int                   `json:"totalCommon"`
	Items             []*DetailItemsRequest `json:"items"`
}

type DetailItemsRequest struct {
	OfferType      int                        `json:"offerType"`
	ProductOFFName string                     `json:"productOFFName"`
	ProductOfferId string                     `json:"productOfferId"`
	Items          []*DetailItemsItemsRequest `json:"items"`
}

type DetailItemsItemsRequest struct {
	NameType            string `json:"nameType"`
	OwnerType           string `json:"ownerType"`
	UnitTypeId          string `json:"unitTypeId"`
	RatableAmount       string `json:"ratableAmount"`
	UsageAmount         string `json:"usageAmount"`
	RatableResourceID   string `json:"ratableResourceID"`
	BeginTime           string `json:"beginTime"`
	EndTime             string `json:"endTime"`
	RatableResourcename string `json:"ratableResourcename"`
	BalanceAmount       string `json:"balanceAmount"`
	OwnerID             string `json:"ownerID"`
}

type BalanceNew struct {
	Result                int    `json:"result"`
	ServiceResultCode     string `json:"serviceResultCode"`
	TotalBalanceAvailable string `json:"totalBalanceAvailable"`
	PaymentFlag           string `json:"paymentFlag"`
	FeeType               int    `json:"feeType"`
	Items                 []struct {
		BalanceAvailable string `json:"balanceAvailable"`
		BalanceTypeFlag  string `json:"balanceTypeFlag"`
	} `json:"items"`
	ParaFieldResult string `json:"paraFieldResult"`
}
