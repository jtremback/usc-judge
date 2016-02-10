package access

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	core "github.com/jtremback/usc-core/judge"
	"github.com/tv42/compound"
)

// compound index types
type ssb struct {
	A string
	B string
	C []byte
}

func MakeBuckets(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Indexes"))
		_, err = tx.CreateBucketIfNotExists([]byte("Channels"))
		_, err = tx.CreateBucketIfNotExists([]byte("Judges"))
		_, err = tx.CreateBucketIfNotExists([]byte("Accounts"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func SetJudge(tx *bolt.Tx, jd *core.Judge) error {
	b, err := json.Marshal(jd)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("Judges")).Put(jd.Pubkey, b)
	if err != nil {
		return err
	}

	return nil
}

func SetAccount(tx *bolt.Tx, acct *core.Account) error {
	b, err := json.Marshal(acct)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("Accounts")).Put([]byte(acct.Pubkey), b)
	if err != nil {
		return err
	}

	// Relations

	b, err = json.Marshal(acct.Judge)
	if err != nil {
		return err
	}

	err = tx.Bucket([]byte("Judges")).Put(acct.Judge.Pubkey, b)
	if err != nil {
		return err
	}

	return nil
}

func PopulateAccount(tx *bolt.Tx, acct *core.Account) error {
	jd := &core.Judge{}
	err := json.Unmarshal(tx.Bucket([]byte("Judges")).Get([]byte(acct.Judge.Pubkey)), jd)
	if err != nil {
		return err
	}
	acct.Judge = jd

	return nil
}

func SetChannel(tx *bolt.Tx, ch *core.Channel) error {
	b, err := json.Marshal(ch)
	if err != nil {
		return err
	}
	err = tx.Bucket([]byte("Channels")).Put([]byte(ch.ChannelId), b)
	if err != nil {
		return err
	}

	// Relations

	// Judge

	b, err = json.Marshal(ch.Judge)
	if err != nil {
		return err
	}

	tx.Bucket([]byte("Judges")).Put(ch.Judge.Pubkey, b)

	// Accounts

	for _, acct := range ch.Accounts {
		b, err := json.Marshal(acct)
		if err != nil {
			return err
		}

		tx.Bucket([]byte("Accounts")).Put(acct.Pubkey, b)
	}

	// Indexes

	// Judge Pubkey

	err = tx.Bucket([]byte("Indexes")).Put(compound.Key(ssb{
		"Judge",
		"Pubkey",
		ch.Judge.Pubkey}), []byte(ch.ChannelId))
	if err != nil {
		return err
	}

	return nil
}

func PopulateChannel(tx *bolt.Tx, ch *core.Channel) error {
	accts := []*core.Account{}
	for _, acct := range ch.Accounts {
		err := json.Unmarshal(tx.Bucket([]byte("Accounts")).Get([]byte(acct.Pubkey)), acct)
		if err != nil {
			return err
		}
		err = PopulateAccount(tx, acct)
		if err != nil {
			return err
		}

		accts = append(accts, acct)
	}

	jd := &core.Judge{}
	err := json.Unmarshal(tx.Bucket([]byte("Judges")).Get([]byte(ch.Judge.Pubkey)), jd)
	if err != nil {
		return err
	}

	ch.Accounts = accts
	ch.Judge = jd

	return nil
}
