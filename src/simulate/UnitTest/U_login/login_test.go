package U_login

import (
	"testing"
)

func TestLogin(t *testing.T) {
	LoginRun("127.0.0.1:51001", "login")
}

func LoginRun(host string, module string) {
	UserRegister(host, module)
}
