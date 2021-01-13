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

func getFileToEncryptOrDecrypt(s1, s2 *string) string {
	if *s1 != "" {
		return *s1
	}
	return *s2
}

func main() {
	var output, secret string
	encrypt := flag.String("encrypt", "", "plain file to be encrypted")
	decrypt := flag.String("decrypt", "", "encrypted file to be decrypted")
	flag.StringVar(&output, "output", "", "path to file output")
	flag.StringVar(&secret, "secret", "", "secret key")
	flag.Usage = func() {
		program := os.Args[0]
		fmt.Fprintf(os.Stderr, `To encrypt:
%s --encrypt plain.yaml --output encrypted.yaml --secret strongpass

To decrypt:
%s --decrypt encrypted.yaml --output plain.yaml --secret strongpass
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

	if encrypt != nil {
		if secret == "" {
			log.Fatalln("[ERROR] The secret cannot be empty")
		}
		encrypted, err := aes.Encrypt(secret, content)
		if err != nil {
			log.Fatalln("[ERROR]", err)
		}
		log.Println("[INFO] Writing file", output)
		err = ioutil.WriteFile(output, []byte(encrypted), 0644)
	}
	if decrypt != nil {
		if secret == "" {
			log.Fatalln("[ERROR] The secret cannot be empty")
		}
		message, _ := hex.DecodeString(string(content))
		decrypted, err := aes.Decrypt(secret, message)
		if err != nil {
			log.Fatalln("[ERROR]", err)
		}
		log.Println("[INFO] Writing file", output)
		err = ioutil.WriteFile(output, []byte(decrypted), 0644)
	}
	if err != nil {
		log.Fatalln("[ERROR]", err)
	}
	log.Println("[INFO] done!")
	/*
		router := web.NewRouter()
		http.HandleFunc("/", router.Index())
		port := os.Getenv("COOLCONF_PORT")
		if port == "" {
			port = "10987"
		}
		log.Println("Listening on " + port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	*/
}
