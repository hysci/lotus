package types

import (
	"bytes"

	"github.com/filecoin-project/go-state-types/exitcode"
	"github.com/filecoin-project/go-state-types/abi"
)

type MessageReceipt struct {
	ExitCode exitcode.ExitCode
	Return   []byte
	GasUsed  int64
	TotalCost abi.TokenAmount
}

func (mr *MessageReceipt) Equals(o *MessageReceipt) bool {
	return mr.ExitCode == o.ExitCode && bytes.Equal(mr.Return, o.Return) && mr.GasUsed == o.GasUsed
}
