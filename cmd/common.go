package cmd

import (
	"co2/carbon"
	"co2/database"
	"co2/types"
)

var fs FsWrapper = &impl{}

type FsWrapper interface {
	Services() types.CarbonConfig
}

type impl struct{}

// Looks through all the registered stores and returns all
// the carbon services that are defined within those stores.
//
// This will never to too deep into the stores when looking
// for services since we want it to be fast. Usually a depth of 2
// is enough.
//
// Each of the returned configurations will have the store
// they belong to injected as well so they can retrieve
// the required data if ever needed.
func (i *impl) Services() types.CarbonConfig {
	stores := database.Stores()
	configs := types.CarbonConfig{}

	for _, store := range stores {
		files := carbon.Configurations(store.Path, 2)

		for k, v := range files {
			v.Store = &store
			configs[k] = v
		}
	}

	return configs
}

// Replaces the default Fs instance with a custom
// implementation.
//
// Note that this exists for the sole purpose of unit testing.
// It makes it easy to replace how we access the carbon methods
// during tests.
func WrapFs(custom FsWrapper) {
	fs = custom
}
