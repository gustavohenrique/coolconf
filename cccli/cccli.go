package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gustavohenrique/coolconf/aes"
)

func main() {
	var output, secret string
	encrypt := flag.String("encrypt", "", "plain file to be encrypted")
	decrypt := flag.String("decrypt", "", "encrypted file to be decrypted")
	flag.StringVar(&output, "output", "", "path to file output")
	flag.StringVar(&secret, "secret", "", "secret key")
	flag.Usage = func() {
		program := os.Args[0]
		fmt.Fprintf(os.Stderr, `To encrypt:
%s -encrypt plain.yaml -output encrypted.yaml -secret strongpass

To decrypt:
%s -decrypt encrypted.yaml -output plain.yaml -secret strongpass
`,
			program,
			program)
	}
	flag.Parse()

	var err error

	input := getFileToEncryptOrDecrypt(encrypt, decrypt)
	if input == "" || output == "" || secret == "" {
		flag.Usage()
		os.Exit(1)
	}
	content, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalln("[ERROR] Cannot read file", input, ":", err)
	}

	if *encrypt != "" {
		encrypted, err := aes.Encrypt(secret, content)
		quitIfError(err)
		log.Println("[INFO] Writing file", output)
		err = ioutil.WriteFile(output, []byte(encrypted), 0644)
	}
	if *decrypt != "" {
		message, _ := hex.DecodeString(string(content))
		decrypted, err := aes.Decrypt(secret, message)
		quitIfError(err)
		log.Println("[INFO] Writing file", output)
		err = ioutil.WriteFile(output, []byte(decrypted), 0644)
	}
	quitIfError(err)
	log.Println("[INFO] done!")
}

func quitIfError(err error) {
	if err != nil {
		log.Fatalln("[ERROR]", err)
	}
}

func getFileToEncryptOrDecrypt(s1, s2 *string) string {
	if *s1 != "" {
		return *s1
	}
	return *s2
}
