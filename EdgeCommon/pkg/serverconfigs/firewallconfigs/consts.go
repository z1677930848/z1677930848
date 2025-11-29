// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package firewallconfigs

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/ipconfigs"
	"github.com/iwind/TeaGo/types"
)

const (
	GlobalBlackListId int64 = 2000000000
	GlobalWhiteListId int64 = 2000000001
	GlobalGreyListId  int64 = 2000000002

	DefaultEventLevel = "critical"
)

func FindGlobalListIdWithType(listType ipconfigs.IPListType) int64 {
	switch listType {
	case ipconfigs.IPListTypeBlack:
		return GlobalBlackListId
	case ipconfigs.IPListTypeWhite:
		return GlobalWhiteListId
	case ipconfigs.IPListTypeGrey:
		return GlobalGreyListId
	}
	return 0
}

func FindGlobalListNameWithType(listType ipconfigs.IPListType) string {
	switch listType {
	case ipconfigs.IPListTypeBlack:
		return "global blacklist"
	case ipconfigs.IPListTypeWhite:
		return "global whitelist"
	case ipconfigs.IPListTypeGrey:
		return "global greylist"
	}
	return "global blacklist"
}

func IsGlobalListId(listId int64) bool {
	return listId == GlobalBlackListId || listId == GlobalWhiteListId || listId == GlobalGreyListId
}

func FindGlobalListIds() []int64 {
	return []int64{GlobalBlackListId, GlobalWhiteListId, GlobalGreyListId}
}

func FindGlobalListIdStrings() []string {
	return []string{types.String(GlobalBlackListId), types.String(GlobalWhiteListId), types.String(GlobalGreyListId)}
}
