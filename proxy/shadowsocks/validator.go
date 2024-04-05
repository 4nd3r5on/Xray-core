package shadowsocks

import (
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"hash/crc64"
	"strings"
	"sync"

	"github.com/4nd3r5on/Xray-core/common/dice"
	"github.com/4nd3r5on/Xray-core/common/protocol"
)

// Validator stores valid Shadowsocks users.
type Validator struct {
	sync.RWMutex
	Users []*protocol.MemoryUser

	BehaviorSeed  uint64
	BehaviorFused bool
}

var ErrNotFound = newError("Not Found")

// Add a Shadowsocks user.
func (v *Validator) Add(u *protocol.MemoryUser) error {
	v.Lock()
	defer v.Unlock()

	account := u.Account.(*MemoryAccount)
	if !account.Cipher.IsAEAD() && len(v.Users) > 0 {
		return newError("The cipher is not support Single-port Multi-user")
	}
	v.Users = append(v.Users, u)

	if !v.BehaviorFused {
		hashkdf := hmac.New(sha256.New, []byte("SSBSKDF"))
		hashkdf.Write(account.Key)
		v.BehaviorSeed = crc64.Update(v.BehaviorSeed, crc64.MakeTable(crc64.ECMA), hashkdf.Sum(nil))
	}

	return nil
}

// Del a Shadowsocks user with a non-empty Email.
func (v *Validator) Del(email string) error {
	if email == "" {
		return newError("Email must not be empty.")
	}

	v.Lock()
	defer v.Unlock()

	email = strings.ToLower(email)
	idx := -1
	for i, u := range v.Users {
		if strings.EqualFold(u.Email, email) {
			idx = i
			break
		}
	}

	if idx == -1 {
		return newError("User ", email, " not found.")
	}
	ulen := len(v.Users)

	v.Users[idx] = v.Users[ulen-1]
	v.Users[ulen-1] = nil
	v.Users = v.Users[:ulen-1]

	return nil
}

// Get a Shadowsocks user.
func (v *Validator) Get(bs []byte, command protocol.RequestCommand) (u *protocol.MemoryUser, aead cipher.AEAD, ret []byte, ivLen int32, err error) {
	v.RLock()
	defer v.RUnlock()

	for _, user := range v.Users {
		if account := user.Account.(*MemoryAccount); account.Cipher.IsAEAD() {
			// AEAD payload decoding requires the payload to be over 32 bytes
			if len(bs) < 32 {
				continue
			}

			aeadCipher := account.Cipher.(*AEADCipher)
			ivLen = aeadCipher.IVSize()
			iv := bs[:ivLen]
			subkey := make([]byte, 32)
			subkey = subkey[:aeadCipher.KeyBytes]
			hkdfSHA1(account.Key, iv, subkey)
			aead = aeadCipher.AEADAuthCreator(subkey)

			var matchErr error
			switch command {
			case protocol.RequestCommandTCP:
				data := make([]byte, 4+aead.NonceSize())
				ret, matchErr = aead.Open(data[:0], data[4:], bs[ivLen:ivLen+18], nil)
			case protocol.RequestCommandUDP:
				data := make([]byte, 8192)
				ret, matchErr = aead.Open(data[:0], data[8192-aead.NonceSize():8192], bs[ivLen:], nil)
			}

			if matchErr == nil {
				u = user
				err = account.CheckIV(iv)
				return
			}
		} else {
			u = user
			ivLen = user.Account.(*MemoryAccount).Cipher.IVSize()
			// err = user.Account.(*MemoryAccount).CheckIV(bs[:ivLen]) // The IV size of None Cipher is 0.
			return
		}
	}

	return nil, nil, nil, 0, ErrNotFound
}

func (v *Validator) GetBehaviorSeed() uint64 {
	v.Lock()
	defer v.Unlock()

	v.BehaviorFused = true
	if v.BehaviorSeed == 0 {
		v.BehaviorSeed = dice.RollUint64()
	}
	return v.BehaviorSeed
}
