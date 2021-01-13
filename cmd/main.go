package main

import (
	"encoding/hex"
	"flag"
	"io/ioutil"
	"log"

	"coolconf/aes"
)

func main() {
	var action, input, output, secret string
	flag.StringVar(&action, "action", "", "encrypt or decrypt")
	flag.StringVar(&input, "input", "", "path to file input")
	flag.StringVar(&output, "output", "", "path to file output")
	flag.StringVar(&secret, "secret", "", "secret key")
	flag.Parse()

	var err error

	if secret == "" {
		log.Fatalln("[error] secret cannot be empty")
	}

	content, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalln("[error] Invalid input file:", err)
	}

	if action == "encrypt" {
		encrypted, err := aes.Encrypt(secret, content)
		if err != nil {
			log.Fatalln("[error]", err)
		}
		log.Println("[info] writing file", output)
		err = ioutil.WriteFile(output, []byte(encrypted), 0644)
	}
	if action == "decrypt" {
		message, _ := hex.DecodeString(string(content))
		decrypted, err := aes.Decrypt(secret, message)
		if err != nil {
			log.Fatalln("[error]", err)
		}
		log.Println("[info] writing file", output)
		err = ioutil.WriteFile(output, []byte(decrypted), 0644)
	}
	if err != nil {
		log.Fatalln("[error]", err)
	}
	log.Println("[info] done!")
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
