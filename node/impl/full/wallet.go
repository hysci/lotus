package full

import (
	"context"

	"github.com/filecoin-project/go-address"
	"go.uber.org/fx"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/crypto"

	"github.com/filecoin-project/lotus/chain/stmgr"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
	"github.com/filecoin-project/lotus/lib/sigs"
)

type WalletAPI struct {
	fx.In

	StateManager *stmgr.StateManager
	Wallet       *wallet.Wallet
}

func (a *WalletAPI) WalletNew(ctx context.Context, typ crypto.SigType) (address.Address, error) {
	return a.Wallet.GenerateKey(typ)
}

func (a *WalletAPI) WalletHas(ctx context.Context, addr address.Address) (bool, error) {
	return a.Wallet.HasKey(addr)
}

func (a *WalletAPI) WalletList(ctx context.Context) ([]address.Address, error) {
	return a.Wallet.ListAddrs()
}

func (a *WalletAPI) WalletBalance(ctx context.Context, addr address.Address) (types.BigInt, error) {
	act, err := a.StateManager.LoadActorTsk(ctx, addr, types.EmptyTSK)
	if xerrors.Is(err, types.ErrActorNotFound) {
		return big.Zero(), nil
	} else if err != nil {
		return big.Zero(), err
	}
	return act.Balance, nil
}

func (a *WalletAPI) WalletSign(ctx context.Context, k address.Address, msg []byte) (*crypto.Signature, error) {
	keyAddr, err := a.StateManager.ResolveToKeyAddress(ctx, k, nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to resolve ID address: %w", keyAddr)
	}
	return a.Wallet.Sign(ctx, keyAddr, msg)
}

func (a *WalletAPI) WalletSignMessage(ctx context.Context, k address.Address, msg *types.Message) (*types.SignedMessage, error) {
	mcid := msg.Cid()

	if msg.Method == 0 {
		return nil, xerrors.Errorf("the method is lock, pelease input passwd")
	}

	sig, err := a.WalletSign(ctx, k, mcid.Bytes())
	if err != nil {
		return nil, xerrors.Errorf("failed to sign message: %w", err)
	}

	return &types.SignedMessage{
		Message:   *msg,
		Signature: *sig,
	}, nil
}

func (a *WalletAPI) WalletSignMessage2(ctx context.Context, k address.Address, msg *types.Message, passwd string) (*types.SignedMessage, error) {
	mcid := msg.Cid()

	oldPasswd := wallet.WalletPasswd
	wallet.WalletPasswd = passwd
	defer func() {
		wallet.WalletPasswd = oldPasswd
	}()

	sig, err := a.WalletSign(ctx, k, mcid.Bytes())
	if err != nil {
		return nil, xerrors.Errorf("failed to sign message: %w", err)
	}

	return &types.SignedMessage{
		Message:   *msg,
		Signature: *sig,
	}, nil
}

func (a *WalletAPI) WalletVerify(ctx context.Context, k address.Address, msg []byte, sig *crypto.Signature) (bool, error) {
	return sigs.Verify(sig, k, msg) == nil, nil
}

func (a *WalletAPI) WalletDefaultAddress(ctx context.Context) (address.Address, error) {
	return a.Wallet.GetDefault()
}

func (a *WalletAPI) WalletSetDefault(ctx context.Context, addr address.Address) error {
	return a.Wallet.SetDefault(addr)
}

func (a *WalletAPI) WalletExport(ctx context.Context, addr address.Address, passwd string) (*types.KeyInfo, error) {
	oldPasswd := wallet.WalletPasswd
	wallet.WalletPasswd = passwd
	defer func() {
		wallet.WalletPasswd = oldPasswd
	}()
	return a.Wallet.Export(addr)
}

func (a *WalletAPI) WalletImport(ctx context.Context, ki *types.KeyInfo) (address.Address, error) {
	return a.Wallet.Import(ki)
}

func (a *WalletAPI) WalletDelete(ctx context.Context, addr address.Address) error {
	return a.Wallet.DeleteKey(addr)
}

func (a *WalletAPI) WalletValidateAddress(ctx context.Context, str string) (address.Address, error) {
	return address.NewFromString(str)
}

func (a *WalletAPI) WalletLock(ctx context.Context) error {
	if wallet.IsSetup() {
		if wallet.WalletPasswd != "" {
			wallet.WalletPasswd = ""
			return nil
		} else {
			return xerrors.Errorf("Wallet is locked")
		}
	}

	return xerrors.Errorf("Passwd is not setup")
}

func (a *WalletAPI) WalletUnlock(ctx context.Context, passwd string) error {
	if wallet.IsSetup() {
		if wallet.WalletPasswd == "" {
			err := wallet.CheckPasswd([]byte(passwd))
			if err != nil {
				return err
			}
			wallet.WalletPasswd = passwd
			return nil
		} else {
			return xerrors.Errorf("Wallet is unlocked")
		}
	}

	return xerrors.Errorf("Passwd is not setup")
}
func (a *WalletAPI) WalletIsLock(ctx context.Context) (bool, error) {
	if wallet.IsSetup() {
		if wallet.WalletPasswd == "" {
			return true, nil
		}
		return false, nil
	}
	return false, xerrors.Errorf("Passwd is not setup")
}

func (a *WalletAPI) WalletChangePasswd(ctx context.Context, newPasswd string) (bool, error) {
	if wallet.IsSetup() {
		if wallet.WalletPasswd != "" {
			addr_list, err := a.Wallet.ListAddrs()
			if err != nil {
				return false, err
			}
			addr_all := make(map[address.Address]*types.KeyInfo)
			for _, v := range addr_list {
				addr_all[v], err = a.Wallet.Export(v)
				if err != nil {
					return false, err
				}
				err = a.Wallet.DeleteKey2(v)
				if err != nil {
					return false, err
				}
			}

			setDefault := true
			defalutAddr, err := a.Wallet.GetDefault()
			if err != nil {
				setDefault = false
			}

			if len(newPasswd) != 16 {
				return false, xerrors.Errorf("passwd must 16 character")
			}

			err = wallet.ResetPasswd([]byte(newPasswd))
			if err != nil {
				return false, err
			}

			for k, v := range addr_all {
				addr, err := a.Wallet.Import(v)
				if err != nil {
					return false, nil
				} else if addr != k {
					return false, xerrors.Errorf("import error")
				}
			}

			if setDefault {
				err = a.Wallet.SetDefault(defalutAddr)
				if err != nil {
					return false, err
				}
			}

			a.Wallet.ClearCache()

			return true, nil
		}
		return false, xerrors.Errorf("Wallet is locked")
	}
	return false, xerrors.Errorf("Passwd is not setup")
}

func (a *WalletAPI) WalletClearPasswd(ctx context.Context) (bool, error) {
	if wallet.IsSetup() {
		if wallet.WalletPasswd != "" {
			addr_list, err := a.Wallet.ListAddrs()
			if err != nil {
				return false, err
			}
			addr_all := make(map[address.Address]*types.KeyInfo)
			for _, v := range addr_list {
				addr_all[v], err = a.Wallet.Export(v)
				if err != nil {
					return false, err
				}
				err = a.Wallet.DeleteKey2(v)
				if err != nil {
					return false, err
				}
			}

			setDefault := true
			defalutAddr, err := a.Wallet.GetDefault()
			if err != nil {
				setDefault = false
			}

			err = wallet.ClearPasswd()
			if err != nil {
				return false, err
			}

			for k, v := range addr_all {
				addr, err := a.Wallet.Import(v)
				if err != nil {
					return false, nil
				} else if addr != k {
					return false, xerrors.Errorf("import error")
				}
			}

			if setDefault {
				err = a.Wallet.SetDefault(defalutAddr)
				if err != nil {
					return false, err
				}
			}

			a.Wallet.ClearCache()

			return true, nil
		}
		return false, xerrors.Errorf("Wallet is locked")
	}
	return false, xerrors.Errorf("Passwd is not setup")
}
