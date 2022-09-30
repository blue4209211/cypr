package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
)

func NewAesCommand() *AesCommand {
	gc := &AesCommand{
		fs: flag.NewFlagSet("aes", flag.ContinueOnError),
	}

	gc.fs.StringVar(&gc.op, "op", "encrypt", "encrypt/decrypt values")
	gc.fs.StringVar(&gc.key, "key", "", "Aes Key (Hex), required")
	gc.fs.StringVar(&gc.nonce, "nonce", "", "12 bytes nonce value (optional) if not provided then random nonce will be used and appended to cypher text")

	return gc
}

type AesCommand struct {
	fs *flag.FlagSet

	op    string
	key   string
	nonce string
}

func (g *AesCommand) Flag() flag.FlagSet {
	return *g.fs
}

func (g *AesCommand) Name() string {
	return g.fs.Name()
}

func (g *AesCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *AesCommand) Run() (err error) {
	if g.op == "encrypt" {
		if g.fs.NArg() == 0 {
			return errors.New("data not provided")
		}

		s, err := g.encryptAES(g.fs.Arg(0))
		if err != nil {
			return err
		}
		fmt.Println(s)
	} else if g.op == "decrypt" {
		s, err := g.decryptAES(g.fs.Arg(0))
		if err != nil {
			return err
		}
		fmt.Println(s)
	} else {
		err = errors.New("Unknown Op - " + g.op)
	}
	return err
}

func (g *AesCommand) getKey() ([]byte, error) {
	if g.key == "" {
		return nil, errors.New("key cannot be empty")
	}
	return hex.DecodeString(g.key)
}

func (g *AesCommand) encryptAES(stringToEncrypt string) (s string, err error) {
	key, err := g.getKey()
	if err != nil {
		return s, err
	}

	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return s, err
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return s, err
	}

	//Create a nonce. Nonce should be from GCM
	var nonce []byte
	var prefix []byte

	if g.nonce != "" {
		nonce = []byte(nonce)
	} else {
		nonce = make([]byte, aesGCM.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return s, err
		}
		prefix = nonce
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(prefix, nonce, plaintext, nil)
	s = hex.EncodeToString(ciphertext)

	return s, err
}

func (g *AesCommand) decryptAES(encryptedString string) (s string, err error) {
	key, err := g.getKey()
	if err != nil {
		return s, err
	}

	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return s, err
	}

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return s, err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return s, err
	}

	var nonce []byte
	var ciphertext []byte

	if g.nonce != "" {
		nonce = []byte(g.nonce)
		ciphertext = enc
	} else {
		//Get the nonce size
		nonceSize := aesGCM.NonceSize()
		//Extract the nonce from the encrypted data
		nonce, ciphertext = enc[:nonceSize], enc[nonceSize:]

	}

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return s, err
	}

	s = string(plaintext)

	return s, err
}
