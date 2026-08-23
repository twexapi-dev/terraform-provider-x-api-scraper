package provider

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/twexapi-dev/x-api-scraper-go/optionalnullable"
)

func optString(value optionalnullable.OptionalNullable[string]) types.String {
	v, ok := value.GetOrZero()
	if !ok || value.IsNull() {
		return types.StringNull()
	}
	return types.StringValue(v)
}

func optInt64(value optionalnullable.OptionalNullable[int64]) types.Int64 {
	v, ok := value.GetOrZero()
	if !ok || value.IsNull() {
		return types.Int64Null()
	}
	return types.Int64Value(v)
}

func optBool(value optionalnullable.OptionalNullable[bool]) types.Bool {
	v, ok := value.GetOrZero()
	if !ok || value.IsNull() {
		return types.BoolNull()
	}
	return types.BoolValue(v)
}

func asJSON(value any) (types.String, error) {
	if value == nil {
		return types.StringNull(), nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return types.StringNull(), err
	}
	return types.StringValue(string(raw)), nil
}
