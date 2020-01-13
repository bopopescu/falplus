package message

import "fmt"

var StatusMsg = make(map[int32][2]string)

func registerMessage(mMap map[int32][2]string) {
	for code, message := range mMap {
		if _, exist := StatusMsg[code]; exist {
			panic(fmt.Sprintf("error code %d is reused", code))
		}
		StatusMsg[code] = message
	}
}
