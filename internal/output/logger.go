package output

import (
	"fmt"
)

func Log(jobs []Job) {
	for _, j := range jobs {
		if j.IsTarget() {
			fmt.Println(j)
		}
	}
}
