package hot

import (
	"log"
)

func (h *Hot) log(v ...interface{}) {
	if h.option.Log {
		log.Println(v...)
	}
}
