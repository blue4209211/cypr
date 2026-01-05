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

	gc.fs.Usage = func() {
		fmt.Printf("Usage of %s : <encrypt/decrypt> [args] <value> \n Args:\n", gc.fs.Name())
		gc.fs.PrintDefaults()
	}

	gc.fs.StringVar(&gc.key, "key", "", "Aes Key (Hex), required")
	gc.fs.StringVar(&gc.nonce, "nonce", "", "12 bytes nonce value (optional) if not provided then random nonce will be used and appended to cypher text")

	return gc
}

type AesCommand struct {
	fs *flag.FlagSet

	op     string
	opArgs []string
	key    string
	nonce  string
}

func (g *AesCommand) Flag() flag.FlagSet {
	return *g.fs
}

func (g *AesCommand) Name() string {
	return g.fs.Name()
}

func (g *AesCommand) Init(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("invalid args - %v", args)
	}
	g.op = args[0]

	err := g.fs.Parse(args[1:])
	if err != nil {
		return err
	}

	if g.fs.NArg() != 1 {
		return fmt.Errorf("invalid args - %v", g.fs.Args())
	}
	g.opArgs = g.fs.Args()

	return nil
}

func (g *AesCommand) Run() (err error) {
	switch g.op {
	case "encrypt":
		s, err := g.encrypt(g.opArgs[0])
		if err != nil {
			return err
		}
		fmt.Println(s)
	case "decrypt":
		s, err := g.decrypt(g.opArgs[0])
		if err != nil {
			return err
		}
		fmt.Println(s)
	default:
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

func (g *AesCommand) getNonce(size int) (nonce []byte, err error) {
	if g.nonce == "" {
		return nonce, err
	}
	nonce, err = hex.DecodeString(g.nonce)
	if err != nil {
		return nonce, errors.New("invalid nonce - " + err.Error())
	} else if len(nonce) != size {
		return nonce, fmt.Errorf("nonce size (%d) is not equal to agm size (%d)", len(nonce), size)
	}
	return nonce, err
}

func (g *AesCommand) encrypt(stringToEncrypt string) (s string, err error) {
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
	var prefix []byte
	nonce, err := g.getNonce(aesGCM.NonceSize())
	if err != nil {
		return s, err
	}

	if len(nonce) == 0 {
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

func (g *AesCommand) decrypt(encryptedString string) (s string, err error) {
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

	var ciphertext []byte
	nonce, err := g.getNonce(aesGCM.NonceSize())
	if err != nil {
		return s, err
	}

	if len(nonce) == 0 {
		//Get the nonce size
		nonceSize := aesGCM.NonceSize()
		//Extract the nonce from the encrypted data
		nonce, ciphertext = enc[:nonceSize], enc[nonceSize:]
	} else {
		ciphertext = enc
	}

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return s, err
	}

	s = string(plaintext)

	return s, err
}
