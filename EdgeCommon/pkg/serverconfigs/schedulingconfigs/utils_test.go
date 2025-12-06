// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package schedulingconfigs

import "testing"

func TestFindSchedulingType(t *testing.T) {
	t.Logf("%p", FindSchedulingType("roundRobin"))
	t.Logf("%p", FindSchedulingType("roundRobin"))
	t.Logf("%p", FindSchedulingType("roundRobin"))
}
