package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func main() {
	passwordlenght := 10

	databagList := map[string]string{"my": "test", "test": "test", "file": "test"}

	password := Generate(passwordlenght)

	for databag, value := range databagList {
		file := generateFile(value, password)
		runKnife(databag, file)
	}

}

var allowedChars = []rune("QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm0987654321")

// Generate random password
// strong - include to password special characters
func Generate(length int) string {
	newRandom(time.Now().UTC().UnixNano())
	charsArray := make([]rune, length)
	for i := range charsArray {
		charsArray[i] = allowedChars[next(0, len(allowedChars))]
	}
	return string(charsArray)
}

// newRandom return random int64
func newRandom(seed int64) {
	rand.Seed(seed)
}

// next return random between min and max value
func next(min, max int) int {
	return min + rand.Intn(max-min)
}

func generateFile(id string, value string) string {

	filename := id + ".json"

	var databag = map[string]string{"id": id, "value": value}

	databagjson, _ := json.Marshal(databag)
	err := ioutil.WriteFile(filename, databagjson, 0644)
	if err != nil {
		fmt.Print(err)
	}
	return filename

}

func runKnife(my string, filename string) {

	encrypteFile := "/etc/chef/encrypted_data_bag_secret"
	if _, err := os.Stat(encrypteFile); os.IsNotExist(err) {
		fmt.Print("Can't find the key")
	}

	cmd := exec.Command("knife", "data", "bag", "from", "file", my, filename, "--secret-file", encrypteFile)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())
}
