package logic

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/jtremback/usc-core/wire"
	"github.com/jtremback/usc-judge/access"
)

type Peer struct {
	db *bolt.DB
}

func (a *Peer) AddChannel(ev *wire.Envelope) error {
	var err error

	otx := &wire.OpeningTx{}
	err = proto.Unmarshal(ev.Payload, otx)
	if err != nil {
		return err
	}

	err = a.db.Update(func(tx *bolt.Tx) error {
		_, err = access.GetChannel(tx, otx.ChannelId)
		if err != nil {
			return errors.New("channel already exists")
		}

		cpt, err = access.GetCounterparty(tx, otx.Pubkeys[0])
		if err != nil {
			return err
		}

		acct, err = access.GetAccount(tx, otx.Pubkeys[1])
		if err != nil {
			return err
		}

		ch, err := core.NewChannel(ev, acct, cpt)

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
