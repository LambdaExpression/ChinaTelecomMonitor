package tools

import (
	"China_Telecom_Monitor/models"
	"github.com/golang-module/carbon/v2"
	"strconv"
	"strings"
)

func ToSummary(qryImportantData *models.Result[models.ImportantData], username string, time carbon.Carbon) models.Summary {
	var ds models.Summary
	if qryImportantData == nil || qryImportantData.HeaderInfos.Code != "0000" || qryImportantData.ResponseData.ResultCode != "0000" {
		return ds
	}
	data := qryImportantData.ResponseData.Data

	useFlow, _ := strconv.ParseInt(data.FlowInfo.TotalAmount.Used, 10, 64)
	balanceFlow, _ := strconv.ParseInt(data.FlowInfo.TotalAmount.Balance, 10, 64)
	totalFlow := useFlow + balanceFlow

	generalUse, _ := strconv.ParseInt(data.FlowInfo.TotalAmount.Used, 10, 64)
	generalBalance, _ := strconv.ParseInt(data.FlowInfo.TotalAmount.Balance, 10, 64)
	generalTotal := generalUse + generalBalance

	specialUse, _ := strconv.ParseInt(data.FlowInfo.SpecialAmount.Used, 10, 64)
	specialBalance, _ := strconv.ParseInt(data.FlowInfo.SpecialAmount.Balance, 10, 64)
	specialTotal := specialUse + specialBalance

	voiceUsage, _ := strconv.ParseInt(data.VoiceInfo.VoiceDataInfo.Used, 10, 64)
	voiceAmount, _ := strconv.ParseInt(data.VoiceInfo.VoiceDataInfo.Total, 10, 64)

	balanceFloat, _ := strconv.ParseFloat(data.BalanceInfo.IndexBalanceDataInfo.Balance, 64)
	balance := int64(balanceFloat * 100)

	var items []models.SummaryItems
	flowLists := data.FlowInfo.FlowList
	if flowLists != nil && len(flowLists) > 0 {
		items = make([]models.SummaryItems, len(flowLists))
		for i, flowList := range flowLists {
			if !strings.Contains(flowList.Title, "流量") {
				continue
			}
			var use, balanceF int64
			if strings.Contains(flowList.LeftTitle, "已用") {
				use, _ = ToInt64(flowList.LeftTitle)
			}
			if strings.Contains(flowList.RightTitle, "剩余") {
				balanceF, _ = ToInt64(flowList.RightTitleHh)
			}
			items[i] = models.SummaryItems{
				Name:  flowList.Title,
				Use:   use,
				Total: use + balanceF,
			}
		}
	}

	return models.Summary{
		Username:     username,
		Use:          useFlow,
		Total:        totalFlow,
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
