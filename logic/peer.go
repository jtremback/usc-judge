package logic

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/jtremback/usc-core/wire"
	"github.com/jtremback/usc-judge/access"
)

type PeerAPI struct {
	db *bolt.DB
}

func (a *PeerAPI) AddChannel(ev *wire.Envelope) error {
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

		acct0, err := access.GetAccount(tx, otx.Pubkeys[0])
		if err != nil {
			return err
		}

		acct1, err := access.GetAccount(tx, otx.Pubkeys[1])
		if err != nil {
			return err
		}

		judge, err := access.GetJudge(tx, acct0.Judge.Pubkey)
		if err != nil {
			return err
		}

		ch, err := judge.AddChannel(ev, otx, acct0, acct1)

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

func (a *PeerAPI) StartClosingChannel(ev *wire.Envelope) error {

}
