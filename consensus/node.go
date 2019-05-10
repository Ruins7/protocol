package consensus

import (
	"github.com/Oneledger/protocol/config"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmnode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
)

type Node = tmnode.Node

func NewNode(application abcitypes.Application, cfg *config.Server, rootDir string) (*tmnode.Node, error) {
	nodecfg, err := parseConfig(cfg, rootDir)
	if err != nil {
		return nil, err
	}
	// log.Dump("TM Node settings", nodecfg.CFG)
	clientCreator := proxy.NewLocalClientCreator(application)
	return tmnode.NewNode(
		&nodecfg.CFG,
		nodecfg.privValidator,
		nodecfg.nodeKey,
		clientCreator,
		nodecfg.genesisProvider,
		tmnode.DefaultDBProvider,
		nodecfg.metricsProvider,
		nodecfg.logger)
}
