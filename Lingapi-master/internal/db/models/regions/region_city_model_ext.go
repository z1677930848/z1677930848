package regions

import (
	"encoding/json"

	"github.com/TeaOSLab/EdgeAPI/internal/remotelogs"
)

// DecodeCodes returns decoded codes slice.
func (this *RegionCity) DecodeCodes() []string {
	if len(this.Codes) == 0 {
		return []string{}
	}
	var result []string
	if err := json.Unmarshal(this.Codes, &result); err != nil {
		remotelogs.Error("RegionCity.DecodeCodes", err.Error())
		return []string{}
	}
	return result
}

// DecodeCustomCodes returns decoded custom codes slice.
func (this *RegionCity) DecodeCustomCodes() []string {
	if len(this.CustomCodes) == 0 {
		return []string{}
	}
	var result []string
	if err := json.Unmarshal(this.CustomCodes, &result); err != nil {
		remotelogs.Error("RegionCity.DecodeCustomCodes", err.Error())
		return []string{}
	}
	return result
}

// DisplayName returns custom name if set, otherwise name.
func (this *RegionCity) DisplayName() string {
	if len(this.CustomName) > 0 {
		return this.CustomName
	}
	return this.Name
}

// AllCodes merges codes and custom codes.
func (this *RegionCity) AllCodes() []string {
	codes := this.DecodeCodes()
	custom := this.DecodeCustomCodes()
	return append(codes, custom...)
}
