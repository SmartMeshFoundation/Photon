package contracttest

import "math/rand"

func (env *Env) getTwoRandomAccount() (*Account, *Account) {
	var index1, index2 int
	n := len(env.Accounts)
	index1 = rand.Intn(n)
	index2 = rand.Intn(n)
	for index1 == index2 {
		index2 = rand.Intn(n)
	}
	return env.Accounts[index1], env.Accounts[index2]
}
