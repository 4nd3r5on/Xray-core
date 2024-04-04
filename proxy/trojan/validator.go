package trojan

import (
	"strings"
	"sync"

	"github.com/4nd3r5on/Xray-core/common/protocol"
)

// Validator stores valid trojan users.
type Validator struct {
	// Considering email's usage here, map + sync.Mutex/RWMutex may have better performance.
	Email sync.Map
	Users sync.Map
}

// Add a trojan user, Email must be empty or unique.
func (v *Validator) Add(u *protocol.MemoryUser) error {
	if u.Email != "" {
		_, loaded := v.Email.LoadOrStore(strings.ToLower(u.Email), u)
		if loaded {
			return newError("User ", u.Email, " already exists.")
		}
	}
	v.Users.Store(hexString(u.Account.(*MemoryAccount).Key), u)
	return nil
}

// Del a trojan user with a non-empty Email.
func (v *Validator) Del(e string) error {
	if e == "" {
		return newError("Email must not be empty.")
	}
	le := strings.ToLower(e)
	u, _ := v.Email.Load(le)
	if u == nil {
		return newError("User ", e, " not found.")
	}
	v.Email.Delete(le)
	v.Users.Delete(hexString(u.(*protocol.MemoryUser).Account.(*MemoryAccount).Key))
	return nil
}

// Get a trojan user with hashed key, nil if user doesn't exist.
func (v *Validator) Get(hash string) *protocol.MemoryUser {
	u, _ := v.Users.Load(hash)
	if u != nil {
		return u.(*protocol.MemoryUser)
	}
	return nil
}
