// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package userconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"

type UserTicketStatus = string

const (
	UserTicketStatusNone   UserTicketStatus = "none"
	UserTicketStatusSolved UserTicketStatus = "solved"
	UserTicketStatusClosed UserTicketStatus = "closed"
)

func UserTicketStatusName(status UserTicketStatus) string {
	switch status {
	case UserTicketStatusNone:
		return "pending"
	case UserTicketStatusSolved:
		return "solved"
	case UserTicketStatusClosed:
		return "closed"
	}
	return ""
}

func FindAllUserTicketStatusList() []*shared.Definition {
	return []*shared.Definition{
		{Name: "pending", Code: UserTicketStatusNone},
		{Name: "solved", Code: UserTicketStatusSolved},
		{Name: "closed", Code: UserTicketStatusClosed},
	}
}
