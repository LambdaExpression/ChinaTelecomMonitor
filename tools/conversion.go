package tools

import (
	"China_Telecom_Monitor/models"
	"strconv"
	"strings"

	"github.com/golang-module/carbon/v2"
)

func ParseSummary(qryImportantData *models.Result[models.ImportantData], username string, time carbon.Carbon) models.Summary {
	var ds models.Summary
	if qryImportantData == nil || qryImportantData.HeaderInfos.Code != "0000" || qryImportantData.ResponseData.ResultCode != "0000" {
		return ds
	}
	data := qryImportantData.ResponseData.Data

	usedData, _ := strconv.ParseInt(data.TrafficInfo.TotalAmount.Used, 10, 64)
	balanceData, _ := strconv.ParseInt(data.TrafficInfo.TotalAmount.Balance, 10, 64)
	totalData := usedData + balanceData

	generalUse, _ := strconv.ParseInt(data.TrafficInfo.CommonTraffic.Used, 10, 64)
	generalBalance, _ := strconv.ParseInt(data.TrafficInfo.CommonTraffic.Balance, 10, 64)
	generalTotal := generalUse + generalBalance

	specialUse, _ := strconv.ParseInt(data.TrafficInfo.SpecialAmount.Used, 10, 64)
	specialBalance, _ := strconv.ParseInt(data.TrafficInfo.SpecialAmount.Balance, 10, 64)
	specialTotal := specialUse + specialBalance

	voiceUsage, _ := strconv.ParseInt(data.VoiceInfo.VoiceDataInfo.Used, 10, 64)
	voiceAmount, _ := strconv.ParseInt(data.VoiceInfo.VoiceDataInfo.Total, 10, 64)

	balanceFloat, _ := strconv.ParseFloat(data.BalanceInfo.IndexBalanceDataInfo.Balance, 64)
	balance := int64(balanceFloat * 100)

	var items []models.SummaryItems
	trafficLists := data.TrafficInfo.TrafficList
	if trafficLists != nil && len(trafficLists) > 0 {
		items = make([]models.SummaryItems, len(trafficLists))
		for i, trafficList := range trafficLists {
			if !strings.Contains(trafficList.Title, "流量") {
				continue
			}
			var use, balanceF int64
			if strings.Contains(trafficList.LeftTitle, "已用") {
				use, _ = ParseTraffic(trafficList.LeftTitleHh)
			}
			if strings.Contains(trafficList.RightTitle, "剩余") {
				balanceF, _ = ParseTraffic(trafficList.RightTitleHh)
			}
			items[i] = models.SummaryItems{
				Name:  trafficList.Title,
				Use:   use,
				Total: use + balanceF,
			}
		}
	}

	return models.Summary{
		Username:     username,
		Use:          usedData,
		Total:        totalData,
		Balance:      balance,
		VoiceUsage:   voiceUsage,
		VoiceAmount:  voiceAmount,
		GeneralUse:   generalUse,
		GeneralTotal: generalTotal,
		SpecialUse:   specialUse,
		SpecialTotal: specialTotal,
		CreateTime:   carbon.DateTime{Carbon: time},
		Items:        items,
	}

}
