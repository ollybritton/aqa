package repl

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	au "github.com/logrusorgru/aurora"
)

// Buffer opens up a buffer in vim which a multi-line input can be recieved in.
func Buffer(clear bool) string {
	if clear {
		err := ioutil.WriteFile("/tmp/.aqa-buf.aqa", []byte{}, 0777)
		if err != nil {
			fmt.Println(au.Red("Error clearing buffer:").Bold())
			fmt.Println(au.Red(err))

			return ""
		}
	}

	cmd := exec.Command("vim", "/tmp/.aqa-buf.aqa")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	if err != nil {
		fmt.Println(au.Red("Error opening vim buffer:").Bold())
		fmt.Println(au.Red(err))

		return ""
	}

	bytes, err := ioutil.ReadFile("/tmp/.aqa-buf.aqa")
	if err != nil {
		fmt.Println(au.Red("Error opening vim buffer:").Bold())
		fmt.Println(au.Red(err))

		return ""
	}

	return string(bytes)
}
