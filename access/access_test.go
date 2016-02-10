package access

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/boltdb/bolt"
	core "github.com/jtremback/usc-core/judge"
	"github.com/jtremback/usc-core/wire"
)

func TestSetJudge(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	jd := &core.Judge{
		Name:    "joe",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
	}
	ju2 := &core.Judge{}

	db.Update(func(tx *bolt.Tx) error {
		err := SetJudge(tx, jd)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		err := json.Unmarshal(tx.Bucket([]byte("Judges")).Get(jd.Pubkey), ju2)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	if !reflect.DeepEqual(jd, ju2) {
		t.Fatal("structs not equal :(")
	}
}

func TestSetAccount(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	acct := &core.Account{
		Name:   "boogie",
		Pubkey: []byte{40, 40, 40},
		Judge: &core.Judge{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	db.Update(func(tx *bolt.Tx) error {
		err := SetAccount(tx, acct)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		ma2 := &core.Account{}
		json.Unmarshal(tx.Bucket([]byte("Accounts")).Get(acct.Pubkey), ma2)
		fmt.Println("shibbbb")
		if !reflect.DeepEqual(acct, ma2) {
			t.Fatal("Account incorrect")
		}

		fromDB := tx.Bucket([]byte("Judges")).Get(acct.Judge.Pubkey)
		jd := &core.Judge{}
		json.Unmarshal(fromDB, jd)

		if !reflect.DeepEqual(acct.Judge, jd) {
			t.Fatal("Judge incorrect", acct.Judge, jd, string(tx.Bucket([]byte("Judges")).Get(acct.Judge.Pubkey)))
		}
		return nil
	})
}

func TestPopulateAccount(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	acct := &core.Account{
		Name:   "boogie",
		Pubkey: []byte{40, 40, 40},
		Judge: &core.Judge{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},
	}

	jd := &core.Judge{
		Name:    "joe",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
	}

	db.Update(func(tx *bolt.Tx) error {
		err := SetAccount(tx, acct)
		if err != nil {
			t.Fatal(err)
		}

		err = SetJudge(tx, jd)
		if err != nil {
			t.Fatal(err)
		}

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		err := PopulateAccount(tx, acct)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(acct.Judge, jd) {
			t.Fatal("Judge incorrect", acct.Judge, jd)
		}
		return nil
	})
}

func TestSetChannel(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ch := &core.Channel{
		ChannelId: "xyz23",
		Phase:     2,

		OpeningTx:         &wire.OpeningTx{},
		OpeningTxEnvelope: &wire.Envelope{},

		LastFullUpdateTx:         &wire.UpdateTx{},
		LastFullUpdateTxEnvelope: &wire.Envelope{},

		Fulfillments: [][]byte{[]byte{80, 80}},

		Judge: &core.Judge{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},

		Accounts: []*core.Account{
			&core.Account{
				Name:   "wrong",
				Pubkey: []byte{40, 40, 40},
				Judge: &core.Judge{
					Name:    "wrong",
					Pubkey:  []byte{40, 40, 40},
					Address: "stoops.com:3004",
				},
			},
			&core.Account{
				Name:    "wrong",
				Pubkey:  []byte{50, 50, 50},
				Address: "stoops.com:3004",
				Judge: &core.Judge{
					Name:    "wrong",
					Pubkey:  []byte{40, 40, 40},
					Address: "stoops.com:3004",
				},
			},
		},
	}

	db.Update(func(tx *bolt.Tx) error {
		err := SetChannel(tx, ch)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		ch2 := &core.Channel{}
		json.Unmarshal(tx.Bucket([]byte("Channels")).Get([]byte(ch.ChannelId)), ch2)
		if !reflect.DeepEqual(ch, ch2) {
			t.Fatal("Channel incorrect")
		}

		jdJson := tx.Bucket([]byte("Judges")).Get(ch.Judge.Pubkey)
		jd := &core.Judge{}
		json.Unmarshal(jdJson, jd)

		if !reflect.DeepEqual(ch.Judge, jd) {
			t.Fatal("Judge incorrect")
		}

		for i, acct := range ch.Accounts {
			acctJson := tx.Bucket([]byte("Accounts")).Get(acct.Pubkey)
			json.Unmarshal(acctJson, acct)

			if !reflect.DeepEqual(ch.Accounts[i], acct) {
				t.Fatal("Account incorrect")
			}
		}

		return nil
	})
}

func TestPopulateChannel(t *testing.T) {
	db, err := bolt.Open("/tmp/test.db", 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test.db")

	err = MakeBuckets(db)
	if err != nil {
		t.Fatal(err)
	}

	ch := &core.Channel{
		ChannelId: "xyz23",
		Phase:     2,

		OpeningTx:         &wire.OpeningTx{},
		OpeningTxEnvelope: &wire.Envelope{},

		LastFullUpdateTx:         &wire.UpdateTx{},
		LastFullUpdateTxEnvelope: &wire.Envelope{},

		Fulfillments: [][]byte{[]byte{80, 80}},

		Judge: &core.Judge{
			Name:    "wrong",
			Pubkey:  []byte{40, 40, 40},
			Address: "stoops.com:3004",
		},

		Accounts: []*core.Account{
			&core.Account{
				Name:   "wrong",
				Pubkey: []byte{40, 40, 40},
				Judge: &core.Judge{
					Name:    "wrong",
					Pubkey:  []byte{40, 40, 40},
					Address: "stoops.com:3004",
				},
			},
			&core.Account{
				Name:    "wrong",
				Pubkey:  []byte{50, 50, 50},
				Address: "stoops.com:3004",
				Judge: &core.Judge{
					Name:    "wrong",
					Pubkey:  []byte{40, 40, 40},
					Address: "stoops.com:3004",
				},
			},
		},
	}

	jd := &core.Judge{
		Name:    "joe",
		Pubkey:  []byte{40, 40, 40},
		Address: "stoops.com:3004",
	}

	accts := []*core.Account{
		&core.Account{
			Name:   "crow",
			Pubkey: []byte{40, 40, 40},
			Judge: &core.Judge{
				Name:    "joe",
				Pubkey:  []byte{40, 40, 40},
				Address: "stoops.com:3004",
			},
		},

		&core.Account{
			Name:    "flerb",
			Pubkey:  []byte{50, 50, 50},
			Address: "stoops.com:3004",
			Judge: &core.Judge{
				Name:    "joe",
				Pubkey:  []byte{40, 40, 40},
				Address: "stoops.com:3004",
			},
		},
	}

	db.Update(func(tx *bolt.Tx) error {
		err := SetChannel(tx, ch)
		if err != nil {
			t.Fatal(err)
		}

		for _, acct := range accts {
			err = SetAccount(tx, acct)
			if err != nil {
				t.Fatal(err)
			}
		}

		err = SetJudge(tx, jd)
		if err != nil {
			t.Fatal(err)
		}

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		err = PopulateChannel(tx, ch)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(ch.Judge, jd) {
			t.Fatal("Judge incorrect", ch.Judge, jd)
		}

		for i, acct := range accts {
			if !reflect.DeepEqual(ch.Accounts[i], acct) {
				t.Fatal("Account incorrect", ch.Accounts[i], acct)
			}
		}

		return nil
	})
}
