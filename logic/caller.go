package logic

import (
	"errors"

	"github.com/boltdb/bolt"
	core "github.com/jtremback/usc-core/judge"
	"github.com/jtremback/usc-core/wire"
	"github.com/jtremback/usc-judge/access"
)

type Caller struct {
	DB *bolt.DB
}

func (a *Caller) ConfirmChannel(chID string) error {
	var err error
	ch := &core.Channel{}
	err = a.DB.Update(func(tx *bolt.Tx) error {
		ch, err = access.GetChannel(tx, chID)
		if err != nil {
			return err
		}

		ch.OpeningTxEnvelope = ch.Judge.SignEnvelope(ch.OpeningTxEnvelope)

		access.SetChannel(tx, ch)
		if err != nil {
			return errors.New("database error")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *PeerAPI) FinishClosingChannel(ev *wire.Envelope) error {

}
