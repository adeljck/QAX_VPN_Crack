package module

import (
	"fmt"
	"math/rand"
	"strings"
)

func (v vpnConnect) generatePassword() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	sb := strings.Builder{}
	sb.Grow(16)
	for i := 0; i < 12; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
func (v vpnConnect) showUserList() {
	for k, v := range v.users {
		fmt.Printf("%d.%s  ", k, v)
		if k == 0 {
			continue
		} else if k%8 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}
