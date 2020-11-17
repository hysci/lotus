// +build !debug
// +build !2k
// +build !testground

package build

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/actors/builtin/power"
)

var DrandSchedule = map[abi.ChainEpoch]DrandEnum{
	0:                  DrandIncentinet,
}

const UpgradeCreeperHeight = 54720
const UpgradeBreezeHeight = 51910
const BreezeGasTampingDuration = 120
const RcPos = -2640

const UpgradeSmokeHeight = 72070

func init() {
	miner.UpgradeRcHeight = UpgradeBreezeHeight + RcPos
	power.ConsensusMinerMinPower = big.NewInt(20 << 30)
	miner.SupportedProofTypes = map[abi.RegisteredSealProof]struct{}{
		abi.RegisteredSealProof_StackedDrg4GiBV1:  {},
		abi.RegisteredSealProof_StackedDrg16GiBV1: {},
	}
}

const BlockDelaySecs = uint64(builtin.EpochDurationSeconds)

const PropagationDelaySecs = uint64(6)
