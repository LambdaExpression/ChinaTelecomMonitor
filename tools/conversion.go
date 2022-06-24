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

func ToSummary2(dr *models.DetailRequest, fp *models.FlowPackage, b *models.BalanceNew, username string, time carbon.Carbon) models.Summary {
	var ds models.Summary
	if fp == nil || fp.Result != 10000 {
		return ds
	}

	if fp.UserPackageBalance.Items == nil || len(fp.UserPackageBalance.Items) == 0 {
		return ds
	}

	var generalUse int64 = 0   // 通用流量使用量
	var generalTotal int64 = 0 // 通用流量总量
	var specialUse int64 = 0   // 专用流量使用量
	var specialTotal int64 = 0 // 专用流量总量

	items := make([]models.SummaryItems, len(fp.UserPackageBalance.Items))

	for i, fu := range fp.UserPackageBalance.Items {
		if fu.ProductOFFName == "国内流量" {
			t, err := ToInt64(fu.RatableAmount)
			if err != nil {
				configs.Logger.Error(err)
			} else {
				generalTotal += t
			}

			u, err := ToInt64(fu.BalanceAmount)
			if err != nil {
				configs.Logger.Error(err)
			} else {
				u = t - u
				generalUse += u
			}

			items[i] = models.SummaryItems{
				Name:  fu.ProductOFFName,
				Use:   u,
				Total: t,
			}
		} else {

			t, err := ToInt64(fu.RatableAmount)
			if err != nil {
				configs.Logger.Error(err)
			} else {
				specialTotal += t
			}

			u, err := ToInt64(fu.BalanceAmount)
			if err != nil {
				configs.Logger.Error(err)
			} else {
				u = t - u
				specialUse += u
			}

			items[i] = models.SummaryItems{
				Name:  fu.ProductOFFName,
				Use:   u,
				Total: t,
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
		Use:          fp.UserPackageBalance.Used,
		Total:        fp.UserPackageBalance.Total,
		Balance:      balance,
		VoiceUsage:   voiceUsage,
		VoiceAmount:  voiceAmount,
		GeneralUse:   generalUse,
		GeneralTotal: generalTotal,
		SpecialUse:   specialUse,
		SpecialTotal: specialTotal,
		CreateTime:   carbon.DateTime{time},
		Items:        items,
	}
}
