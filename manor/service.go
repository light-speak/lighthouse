package manor

import "github.com/light-speak/lighthouse/manor/kitex_gen/manor/rpc/manor"

func Start() error {
	svr := manor.NewServer(new(ManorImpl))
	return svr.Run()
}
