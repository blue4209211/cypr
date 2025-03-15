package internal

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	HS256SecretEnvVar = "JWT_HS256_SECRET"
)

func NewJwtCommand() *JwtCommand {
	gc := &JwtCommand{
		fs: flag.NewFlagSet("jwt", flag.ContinueOnError),
	}

	gc.fs.Usage = func() {
		fmt.Printf("Usage of %s : <encode/decode> [args] \n Args:\n", gc.fs.Name())
		fmt.Printf(" encode: <header json string> <payload json string> <signing-algorithm(HS256/RS256)> \n decode: <jwt token>\n")
		fmt.Printf(" For RS256 - privatekey and public key path needs to provided via flag -private-key-path=<path-to-private-key> -public-key-path=<path-to-public-key>\n")
		gc.fs.PrintDefaults()
	}
	gc.fs.StringVar(&gc.privateKeyPath, "private-key-path", "", "Path to the RSA private key file")
	gc.fs.StringVar(&gc.publicKeyPath, "public-key-path", "", "Path to the RSA public key file")

	return gc
}

type JwtCommand struct {
	fs             *flag.FlagSet
	op             string
	opArgs         []string
	privateKeyPath string
	publicKeyPath  string
}

func (g *JwtCommand) Flag() flag.FlagSet {
	return *g.fs
}

func (g *JwtCommand) Name() string {
	return g.fs.Name()
}

func (g *JwtCommand) Init(args []string) error {
	if len(args) < 2 {
		return errors.New("invalid args, op and data required")
	}
	g.op = args[0]

	err := g.fs.Parse(args[1:])
	if err != nil {
		return err
	}

	if g.op == "encode" {
		if g.fs.NArg() != 3 {
			return errors.New("invalid args, header, payload and algo required")
		}
		if strings.ToUpper(g.fs.Arg(2)) == "RS256" {
			if g.privateKeyPath == "" || g.publicKeyPath == "" {
				return errors.New("private key and public key paths are required for RS256")
			}
		}
	} else if g.op == "decode" {
		if g.fs.NArg() != 1 {
			return errors.New("invalid args, token required")
		}
	} else {
		return fmt.Errorf("unknown operation: %s", g.op)
	}

	g.opArgs = g.fs.Args()
	return nil
}

func (g *JwtCommand) Run() (err error) {
	if g.op == "encode" {
		s, err := g.encode(g.opArgs[0], g.opArgs[1], g.opArgs[2])
		if err != nil {
			return err
		}
		fmt.Println(s)
	} else if g.op == "decode" {
		s, err := g.decode(g.opArgs[0])
		if err != nil {
			return err
		}
		fmt.Println(s)
	} else {
		err = errors.New("Unknown Op - " + g.op)
	}
	return err
}

func (g *JwtCommand) encode(headerJSON, payloadJSON, signingAlgo string) (s string, err error) {
	// Decode Header
	var header map[string]any
	err = json.Unmarshal([]byte(headerJSON), &header)
	if err != nil {
		return "", fmt.Errorf("invalid header JSON: %w", err)
	}

	// Decode Payload
	var claims jwt.MapClaims
	err = json.Unmarshal([]byte(payloadJSON), &claims)
	if err != nil {
		return "", fmt.Errorf("invalid payload JSON: %w", err)
	}

	// create token
	token := jwt.NewWithClaims(getHeaderAlgo(header), claims)

	var signingKey any
	switch strings.ToUpper(signingAlgo) {
	case "HS256":
		signingKey, err = getHS256SigningKey()
		if err != nil {
			return "", err
		}
	case "RS256":
		signingKey, err = g.getRS256SigningKey()
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("invalid signing algorithm : %v", signingAlgo)
	}

	// Sign and get the complete encoded token as a string
	s, err = token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}
	return s, nil
}

func (g *JwtCommand) decode(tokenString string) (s string, err error) {
	var token *jwt.Token
	if g.publicKeyPath != "" {
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
				return getHS256VerificationKey()
			} else if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
				return g.getRS256VerificationKey()
			}
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		})

	} else {
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return jwt.UnsafeAllowNoneSignatureType, nil
		})

		if err != nil && !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return "", fmt.Errorf("invalid token : %w", err)
		}

	}

	if err != nil && !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return "", fmt.Errorf("invalid token : %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		header, err := json.MarshalIndent(token.Header, "", "  ")
		if err != nil {
			return "", fmt.Errorf("error in marshal header: %v", err)
		}
		payload, err := json.MarshalIndent(claims, "", "  ")
		if err != nil {
			return "", fmt.Errorf("error in marshal claims: %v", err)
		}
		s = fmt.Sprintf("Header: %s\nPayload: %s", header, payload)
		if g.publicKeyPath != "" && !token.Valid {
			s = fmt.Sprintf("%s\nWarning: token is invalid/expired", s)
		} else if g.publicKeyPath == "" {
			s = fmt.Sprintf("%s\nWarning: provide publicKey to validate token", s)
		}

		return s, nil
	} else {
		return "", fmt.Errorf("invalid token claims")
	}
}

// getHeaderAlgo get header alg field and returns jwt.SigningMethod.
func getHeaderAlgo(header map[string]any) jwt.SigningMethod {
	switch header["alg"] {
	case "HS256":
		return jwt.SigningMethodHS256
	case "RS256":
		return jwt.SigningMethodRS256
	default:
		return nil
	}
}

// getHS256SigningKey get secret from env and returns as []byte
func getHS256SigningKey() ([]byte, error) {
	key := os.Getenv(HS256SecretEnvVar)
	if key == "" {
		return nil, fmt.Errorf("environment variable %s not set", HS256SecretEnvVar)
	}
	return []byte(key), nil
}

// getHS256VerificationKey get secret from env and returns as []byte
func getHS256VerificationKey() ([]byte, error) {
	key, err := getHS256SigningKey()
	if err != nil {
		return nil, err
	}
	return key, nil
}

// getRS256SigningKey get private key from file and returns as *rsa.PrivateKey
func (g *JwtCommand) getRS256SigningKey() (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(g.privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}
	return privateKey, nil
}

// getRS256VerificationKey get public key from file and returns as *rsa.PublicKey
func (g *JwtCommand) getRS256VerificationKey() (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(g.publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading public key file: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}
	return publicKey, nil
}
