package wx_pay
import (
	"math/rand"
	"time"
)

func genRandomStr() string {
	str := "abcdefghijkl1234567mnopqrstuvwxzyA890BCDEFGHIJKLMNOPQSRTUVWXYZ"
	l := len(str)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	data := make([]byte, 20)
	for i := 0; i < 20; i++ {
		data[i] = str[r.Intn(l)]
	}
	return string(data)
}