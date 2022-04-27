package tools

import (
	"China_Telecom_Monitor/configs"
	"China_Telecom_Monitor/models"
	"github.com/golang-module/carbon/v2"
	"strconv"
)

func ToSummary(dr *models.DetailRequest, b *models.BalanceNew, username string, time carbon.Carbon) models.Summary {
	var ds models.Summary
	if dr == nil {
		return ds
	}
	if dr.Items == nil || len(dr.Items) == 0 {
		return ds
	}

	var generalUse int64 = 0   // 通用流量使用量
	var generalTotal int64 = 0 // 通用流量总量
	var specialUse int64 = 0   // 专用流量使用量
	var specialTotal int64 = 0 // 专用流量总量

	for _, di := range dr.Items {
		if di == nil || di.Items == nil || len(di.Items) == 0 {
			continue
		}
		for _, div := range di.Items {
			if div == nil {
				continue
			}
			if di.OfferType != 23 && div.UnitTypeId == "3" {
				gu, err := strconv.ParseInt(div.UsageAmount, 10, 64)
				if err != nil {
					configs.Logger.Error(err)
				}
				gt, err := strconv.ParseInt(div.RatableAmount, 10, 64)
				if err != nil {
					configs.Logger.Error(err)
				}
				generalUse += gu
				generalTotal += gt
			} else if div.UnitTypeId == "3" {
				su, err := strconv.ParseInt(div.UsageAmount, 10, 64)
				if err != nil {
					configs.Logger.Error(err)
				}
				st, err := strconv.ParseInt(div.RatableAmount, 10, 64)
				if err != nil {
					configs.Logger.Error(err)
				}
				specialUse += su
				specialTotal += st
			}

		}
	}

	voiceUsage, err := strconv.ParseInt(dr.VoiceUsage, 10, 64)
	if err != nil {
		configs.Logger.Error(err)
	}
	voiceAmount, err := strconv.ParseInt(dr.VoiceAmount, 10, 64)
	if err != nil {
		configs.Logger.Error(err)
	}
	balance, err := strconv.ParseInt(b.TotalBalanceAvailable, 10, 64)
	if err != nil {
		configs.Logger.Error(err)
	}
	return models.Summary{
		Username:     username,
		Use:          dr.Used,
		Total:        dr.Total,
		Balance:      balance,
		VoiceUsage:   voiceUsage,
		VoiceAmount:  voiceAmount,
		GeneralUse:   generalUse,
		GeneralTotal: generalTotal,
		SpecialUse:   specialUse,
		SpecialTotal: specialTotal,
		CreateTime:   carbon.DateTime{time},
	}
}
