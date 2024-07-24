package vmess

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash/crc64"
	"strings"
	"sync"

	"github.com/xtls/xray-core/common/dice"
	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/proxy/vmess/aead"
)

// TimedUserValidator is a user Validator based on time.
type TimedUserValidator struct {
	sync.RWMutex
	Users []*protocol.MemoryUser

	BehaviorSeed  uint64
	BehaviorFused bool

	AeadDecoderHolder *aead.AuthIDDecoderHolder
}

// NewTimedUserValidator creates a new TimedUserValidator.
func NewTimedUserValidator() *TimedUserValidator {
	tuv := &TimedUserValidator{
		Users:             make([]*protocol.MemoryUser, 0, 16),
		AeadDecoderHolder: aead.NewAuthIDDecoderHolder(),
	}
	return tuv
}

func (v *TimedUserValidator) Add(u *protocol.MemoryUser) error {
	v.Lock()
	defer v.Unlock()

	v.Users = append(v.Users, u)

	account := u.Account.(*MemoryAccount)
	if !v.BehaviorFused {
		hashkdf := hmac.New(sha256.New, []byte("VMESSBSKDF"))
		hashkdf.Write(account.ID.Bytes())
		v.BehaviorSeed = crc64.Update(v.BehaviorSeed, crc64.MakeTable(crc64.ECMA), hashkdf.Sum(nil))
	}

	var cmdkeyfl [16]byte
	copy(cmdkeyfl[:], account.ID.CmdKey())
	v.AeadDecoderHolder.AddUser(cmdkeyfl, u)

	return nil
}

func (v *TimedUserValidator) GetAEAD(userHash []byte) (*protocol.MemoryUser, bool, error) {
	v.RLock()
	defer v.RUnlock()

	var userHashFL [16]byte
	copy(userHashFL[:], userHash)

	userd, err := v.AeadDecoderHolder.Match(userHashFL)
	if err != nil {
		return nil, false, err
	}
	return userd.(*protocol.MemoryUser), true, err
}

func (v *TimedUserValidator) Remove(email string) bool {
	v.Lock()
	defer v.Unlock()

	email = strings.ToLower(email)
	idx := -1
	for i, u := range v.Users {
		if strings.EqualFold(u.Email, email) {
			idx = i
			var cmdkeyfl [16]byte
			copy(cmdkeyfl[:], u.Account.(*MemoryAccount).ID.CmdKey())
			v.AeadDecoderHolder.RemoveUser(cmdkeyfl)
			break
		}
	}
	if idx == -1 {
		return false
	}
	ulen := len(v.Users)

	v.Users[idx] = v.Users[ulen-1]
	v.Users[ulen-1] = nil
	v.Users = v.Users[:ulen-1]

	return true
}

func (v *TimedUserValidator) GetBehaviorSeed() uint64 {
	v.Lock()
	defer v.Unlock()

	v.BehaviorFused = true
	if v.BehaviorSeed == 0 {
		v.BehaviorSeed = dice.RollUint64()
	}
	return v.BehaviorSeed
}

var ErrNotFound = errors.New("Not Found")

var ErrTainted = errors.New("ErrTainted")
