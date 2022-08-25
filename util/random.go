package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet="abcdefghijklmnopqrstuvwxyz"

func init(){
	rand.Seed(time.Now().UnixNano()) //make sure every random number generated is not the same
}

func RandomInt(min,max int64) int64 {
	return min+rand.Int63n(max-min+1) //the generated number is not included max-min+1, so +1 is needed
}

func RandomString(n int) string{
	var sb strings.Builder // can concat rune or byte
	ra :=  []rune(alphabet)
	l := len(ra)
	for i:=0 ; i<n;i++{
		c := ra[rand.Intn(l)]
		sb.WriteRune(c)
	}
	return sb.String()
}

func RandomOwner() string{
	return RandomString(6)
}

func RandomMoney() int64{
	return RandomInt(0,1000)
}

func RandomCurrency() string{
	currencies := []string{"EUR","USD","CAD"}
	l := len(currencies)
	return currencies[rand.Intn(l)]
}

