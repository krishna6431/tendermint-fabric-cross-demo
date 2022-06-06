package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	yaml "gopkg.in/yaml.v2"
)

var (
	KeyAdmin      = []byte("Admin")
	KeyIsSetAdmin = []byte("IsSetAdmin")
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(admin string, isSetAdmin bool) Params {
	return Params{
		Admin:      admin,
		IsSetAdmin: isSetAdmin,
	}
}

func (p Params) Validate() error {
	if err := validateAdmin(p.Admin); err != nil {
		return err
	}
	return validateIsSetAdmin(p.IsSetAdmin)
}

func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAdmin, &p.Admin, validateAdmin),
		paramtypes.NewParamSetPair(KeyIsSetAdmin, &p.IsSetAdmin, validateIsSetAdmin),
	}
}

func validateAdmin(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateIsSetAdmin(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
