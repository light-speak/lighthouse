package lighthouse

import "github.com/light-speak/lighthouse/env"

type Options struct {
	Env *EnvOptions
}

type EnvOptions struct {
	Path string
}

func Init(options *Options) error {
	if options.Env != nil {
		err := env.Init(&options.Env.Path)
		if err != nil {
			return err
		}
	} else {
		err := env.Init(nil)
		if err != nil {
			return err
		}
	}

	return nil
}
