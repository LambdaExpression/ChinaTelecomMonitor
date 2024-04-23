package models

type Request[C any] struct {
	Content     C                  `json:"content"`
	HeaderInfos RequestHeaderInfos `json:"headerInfos"`
}

type LoginRequestContent struct {
	FieldData LoginRequestFieldData `json:"fieldData"`
	Attach    string                `json:"attach"`
}
type LoginRequestFieldData struct {
	AccountType                string `json:"accountType"`
	Authentication             string `json:"authentication"`
	DeviceUid                  string `json:"deviceUid"`
	IsChinatelecom             string `json:"isChinatelecom"`
	LoginAuthCipherAsymmertric string `json:"loginAuthCipherAsymmertric"`
	LoginType                  string `json:"loginType"`
	PhoneNum                   string `json:"phoneNum"`
	SystemVersion              string `json:"systemVersion"`
}

type QryImportantDataRequestContent struct {
	FieldData QryImportantDataRequestFieldData `json:"fieldData"`
	Attach    string                           `json:"attach"`
}
type QryImportantDataRequestFieldData struct {
	ProvinceCode   string `json:"provinceCode"`
	CityCode       string `json:"cityCode"`
	ShopId         string `json:"shopId"`
	IsChinatelecom string `json:"isChinatelecom"`
	Account        string `json:"account"`
}

type UserFluxPackageRequestContent struct {
	FieldData UserFluxPackageRequestFieldData `json:"fieldData"`
	Attach    string                          `json:"attach"`
}
type UserFluxPackageRequestFieldData struct {
	QueryFlag  string `json:"queryFlag"`
	AccessAuth string `json:"accessAuth"`
	Account    string `json:"account"`
}

type RequestHeaderInfos struct {
	ClientType     string `json:"clientType,omitempty"`
	Code           string `json:"code,omitempty"`
	ShopId         string `json:"shopId,omitempty"`
	Source         string `json:"source,omitempty"`
	SourcePassword string `json:"sourcePassword,omitempty"`
	Timestamp      string `json:"timestamp,omitempty"`
	UserLoginName  string `json:"userLoginName,omitempty"`
	Token          string `json:"token,omitempty"`
}

type Result[D any] struct {
	HeaderInfos  HeaderInfos     `json:"headerInfos"`
	ResponseData ResponseData[D] `json:"responseData"`
}
type HeaderInfos struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
}
type ResponseData[D any] struct {
	ResultCode string `json:"resultCode"`
	ResultDesc string `json:"resultDesc"`
	Attach     string `json:"attach"`
	Data       D      `json:"data"`
}

type LoginData struct {
	// https://appgologin.189.cn:9031/login/client/userLoginNormal
	LoginSuccessResult LoginSuccessResult `json:"loginSuccessResult"`
	LoginFailResult    interface{}        `json:"loginFailResult"`
}
type ImportantData struct {
	// https://appfuwu.189.cn:9021/query/qryImportantData
	BalanceInfo    BalanceInfo      `json:"balanceInfo,omitempty"`
	FlowInfo       FlowInfo         `json:"flowInfo,omitempty"`
	VoiceInfo      VoiceInfo        `json:"voiceInfo,omitempty"`
	IntegralInfo   IntegralInfo     `json:"integralInfo,omitempty"`
	StorageInfo    StorageInfo      `json:"storageInfo,omitempty"`
	ThresholdTypes []ThresholdTypes `json:"thresholdTypes,omitempty"`
}
type UserFluxPackageData struct {
	// https://appfuwu.189.cn:9021/query/userFluxPackage
	MainProductOFFInfo MainProductOFFInfo `json:"mainProductOFFInfo,omitempty"`
	ProductOFFRatable  ProductOFFRatable  `json:"productOFFRatable,omitempty"`
	WarnInfo           interface{}        `json:"warnInfo,omitempty"`
	ButtonInfo         ButtonInfo         `json:"buttonInfo,omitempty"`
	Tips               string             `json:"tips,omitempty"`
	QueryFailInfo      interface{}        `json:"queryFailInfo,omitempty"`
	VoiceMessage       string             `json:"voiceMessage,omitempty"`
}

type LoginSuccessResult struct {
	PhoneNbr     string      `json:"phoneNbr"`
	Token        string      `json:"token"`
	UserID       string      `json:"userId"`
	UserType     string      `json:"userType"`
	IsDirectCon  string      `json:"isDirectCon"`
	PhoneType    string      `json:"phoneType"`
	ProvinceCode string      `json:"provinceCode"`
	CityCode     string      `json:"cityCode"`
	ProvinceName string      `json:"provinceName"`
	CityName     string      `json:"cityName"`
	AreaCode     string      `json:"areaCode"`
	NativeNet    string      `json:"nativeNet"`
	NetType      string      `json:"netType"`
	AccessToken  interface{} `json:"accessToken"`
	MemberType   string      `json:"memberType"`
	Operator     string      `json:"operator"`
	IsNewUser    string      `json:"isNewUser"`
}
type IndexBalanceDataInfo struct {
	Arrear    string `json:"arrear"`
	Balance   string `json:"balance"`
	Title     string `json:"title"`
	ShowFlag  string `json:"showFlag"`
	IsShowRed string `json:"isShowRed"`
	Link      string `json:"link"`
	LinkType  string `json:"linkType"`
}
type PhoneBillRegion struct {
	Title      string `json:"title"`
	SubTitle   string `json:"subTitle"`
	SubTitleHh string `json:"subTitleHh"`
	IconURL    string `json:"iconUrl"`
}
type PhoneBillBars struct {
	Title            string `json:"title"`
	SubTilte         string `json:"subTilte"`
	BarPercent       string `json:"barPercent"`
	BarRightSubTitle string `json:"barRightSubTitle"`
}
type ImportantDataBtns struct {
	Title        string `json:"title"`
	IconURL      string `json:"iconUrl"`
	Link         string `json:"link"`
	LinkType     string `json:"linkType"`
	ProvinceCode string `json:"provinceCode"`
	SceneID      string `json:"sceneId"`
	Recommender  string `json:"recommender"`
}
type AdConfigs struct {
	Title        string `json:"title"`
	IconURL      string `json:"iconUrl"`
	Link         string `json:"link"`
	LinkType     string `json:"linkType"`
	ProvinceCode string `json:"provinceCode"`
	SceneID      string `json:"sceneId"`
	Recommender  string `json:"recommender"`
}
type BalanceInfo struct {
	IndexBalanceDataInfo IndexBalanceDataInfo `json:"indexBalanceDataInfo"`
	PhoneBillRegion      PhoneBillRegion      `json:"phoneBillRegion"`
	LoopTips             []string             `json:"loopTips"`
	PhoneBillBars        []PhoneBillBars      `json:"phoneBillBars"`
	ImportantDataBtns    []ImportantDataBtns  `json:"importantDataBtns"`
	AdConfigs            []AdConfigs          `json:"adConfigs"`
	ErrorMes             string               `json:"errorMes"`
}
type Amount struct {
	Total     string `json:"total"`
	Balance   string `json:"balance"`
	Used      string `json:"used"`
	Over      string `json:"over"`
	Title     string `json:"title"`
	ShowField string `json:"showField"`
	Link      string `json:"link"`
	LinkType  string `json:"linkType"`
}
type FlowRegion struct {
	Title      string `json:"title"`
	SubTitle   string `json:"subTitle"`
	SubTitleHh string `json:"subTitleHh"`
	IconURL    string `json:"iconUrl"`
}
type FlowList struct {
	Title         string `json:"title"`
	SubTitle      string `json:"subTitle"`
	BarPercent    string `json:"barPercent"`
	BarRightCount string `json:"barRightCount"`
	LeftTitle     string `json:"leftTitle"`
	LeftTitleHh   string `json:"leftTitleHh"`
	RightTitle    string `json:"rightTitle"`
	RightTitleHh  string `json:"rightTitleHh"`
	RightTitleEnd string `json:"rightTitleEnd"`
}
type FlowInfo struct {
	CommonFlow        Amount              `json:"commonFlow"`
	SpecialAmount     Amount              `json:"specialAmount"`
	TotalAmount       Amount              `json:"totalAmount"`
	FlowRegion        FlowRegion          `json:"flowRegion"`
	LoopTips          []interface{}       `json:"loopTips"`
	ImportantDataBtns []ImportantDataBtns `json:"importantDataBtns"`
	AdConfigs         []AdConfigs         `json:"adConfigs"`
	FlowList          []FlowList          `json:"flowList"`
	ErrorMes          string              `json:"errorMes"`
}
type VoiceDataInfo struct {
	Total     string `json:"total"`
	Balance   string `json:"balance"`
	Used      string `json:"used"`
	Title     string `json:"title"`
	ShowField string `json:"showField"`
	Link      string `json:"link"`
	LinkType  string `json:"linkType"`
}
type VoiceRegion struct {
	Title      string `json:"title"`
	SubTitle   string `json:"subTitle"`
	SubTitleHh string `json:"subTitleHh"`
	IconURL    string `json:"iconUrl"`
}
type VoiceBars struct {
	Title         string `json:"title"`
	SubTitle      string `json:"subTitle"`
	BarPercent    string `json:"barPercent"`
	BarRightCount string `json:"barRightCount"`
	LeftTitle     string `json:"leftTitle"`
	LeftTitleHh   string `json:"leftTitleHh"`
	RightTitle    string `json:"rightTitle"`
	RightTitleHh  string `json:"rightTitleHh"`
	RightTitleEnd string `json:"rightTitleEnd"`
}
type VoiceInfo struct {
	VoiceDataInfo     VoiceDataInfo       `json:"voiceDataInfo"`
	VoiceRegion       VoiceRegion         `json:"voiceRegion"`
	LoopTips          []interface{}       `json:"loopTips"`
	VoiceBars         []VoiceBars         `json:"voiceBars"`
	ImportantDataBtns []ImportantDataBtns `json:"importantDataBtns"`
	AdConfigs         []AdConfigs         `json:"adConfigs"`
	ErrorMes          string              `json:"errorMes"`
}
type IntegralInfo struct {
	Title    string `json:"title"`
	Integral string `json:"integral"`
	Link     string `json:"link"`
	LinkType string `json:"linkType"`
}
type StorageDataInfo struct {
	Balance  string `json:"balance"`
	Title    string `json:"title"`
	Link     string `json:"link"`
	LinkType string `json:"linkType"`
}
type StorageInfo struct {
	StorageDataInfo   StorageDataInfo     `json:"storageDataInfo"`
	FlowRegion        FlowRegion          `json:"flowRegion"`
	LoopTips          []interface{}       `json:"loopTips"`
	ImportantDataBtns []ImportantDataBtns `json:"importantDataBtns"`
	AdConfigs         []AdConfigs         `json:"adConfigs"`
	FlowList          []FlowList          `json:"flowList"`
	ErrorMes          string              `json:"errorMes"`
}
type ThresholdMesList struct {
	Title        string `json:"title"`
	ButtonTitle  string `json:"buttonTitle"`
	Link         string `json:"link"`
	LinkType     string `json:"linkType"`
	SceneID      string `json:"sceneId"`
	ProvinceCode string `json:"provinceCode"`
	Recommender  string `json:"recommender"`
}
type ThresholdTypes struct {
	Type             string             `json:"type"`
	ThresholdMesList []ThresholdMesList `json:"thresholdMesList"`
}

type MainProductOFFInfo struct {
	ProdOFFNameLink     string `json:"prodOFFNameLink"`
	ProdOFFNameLinkType string `json:"prodOFFNameLinkType"`
	ProductOFFName      string `json:"productOFFName"`
	ShareLink           string `json:"shareLink"`
	ShareLinkType       string `json:"shareLinkType"`
	ShareTipDesc        string `json:"shareTipDesc"`
	ShareTitle          string `json:"shareTitle"`
}
type LeftStructure struct {
	RedFlag         string      `json:"redFlag"`
	Num             string      `json:"num"`
	Unit            string      `json:"unit"`
	TitleCornerMark interface{} `json:"titleCornerMark"`
	Title           string      `json:"title"`
}
type RightStructure struct {
	RedFlag         string      `json:"redFlag"`
	Num             string      `json:"num"`
	Unit            string      `json:"unit"`
	TitleCornerMark interface{} `json:"titleCornerMark"`
	Title           string      `json:"title"`
}
type ProductInfos struct {
	LinkType         string      `json:"linkType"`
	Link             string      `json:"link"`
	Title            string      `json:"title"`
	ProgressBar      string      `json:"progressBar"`
	LeftTitle        string      `json:"leftTitle"`
	LeftHighlight    string      `json:"leftHighlight"`
	LeftIcon         string      `json:"leftIcon"`
	TitleIcon        string      `json:"titleIcon"`
	RightTitle       string      `json:"rightTitle"`
	RightHighlight   string      `json:"rightHighlight"`
	RightCommon      string      `json:"rightCommon"`
	OutOfServiceTime string      `json:"outOfServiceTime"`
	IsInvalid        string      `json:"isInvalid"`
	IsInfiniteAmount string      `json:"isInfiniteAmount"`
	InfiniteTitle    interface{} `json:"infiniteTitle"`
	InfiniteValue    interface{} `json:"infiniteValue"`
	InfiniteUnit     interface{} `json:"infiniteUnit"`
	OrderLevel       int         `json:"orderLevel"`
}
type RatableResourcePackages struct {
	LinkType       string         `json:"linkType"`
	Link           string         `json:"link"`
	Title          string         `json:"title"`
	LeftStructure  LeftStructure  `json:"leftStructure"`
	RightStructure RightStructure `json:"rightStructure"`
	ProductInfos   []ProductInfos `json:"productInfos"`
}
type ProductOFFRatable struct {
	ExceedingUsages         interface{}               `json:"exceedingUsages"`
	RatableResourcePackages []RatableResourcePackages `json:"ratableResourcePackages"`
	ExceedResourcePackages  []interface{}             `json:"exceedResourcePackages"`
	LinkType                interface{}               `json:"linkType"`
	Link                    interface{}               `json:"link"`
	ExceedLinkType          interface{}               `json:"exceedLinkType"`
	ExceedLink              interface{}               `json:"exceedLink"`
	ExceedTitle             interface{}               `json:"exceedTitle"`
}
type ButtonInfo struct {
	ButtonImageURL string `json:"buttonImageUrl"`
	ButtonLink     string `json:"buttonLink"`
	ButtonLinkType string `json:"buttonLinkType"`
	ButtonTitle    string `json:"buttonTitle"`
}
