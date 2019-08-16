package in_toto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

/*
ParseRSAPublicKeyFromPEM parses the passed pemBytes as e.g. read from a PEM
formatted file, and instantiates and returns the corresponding RSA public key.
If the no RSA public key can be parsed, the first return value is nil and the
second return value is the error.
*/
func ParseRSAPublicKeyFromPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	// TODO: There could be more key data in _, which we silently ignore here.
	// Should we handle it / fail / say something about it?
	data, _ := pem.Decode([]byte(pemBytes))
	if data == nil {
		return nil, fmt.Errorf("Could not find a public key PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	//ParsePKIXPublicKey might return an rsa, dsa, or ecdsa public key
	rsaPub, isRsa := pub.(*rsa.PublicKey)
	if !isRsa {
		return nil, fmt.Errorf("We currently only support rsa keys: got '%s'",
			reflect.TypeOf(pub))
	}

	return rsaPub, nil
}

/*
LoadPublicKey parses an RSA public key from a PEM formatted file at the passed
path into the Key object on which it was called.  It returns an error if the
file at path does not exist or is not a PEM formatted RSA public key.
*/
func (k *Key) LoadPublicKey(path string) error {
	keyFile, err := os.Open(path)
	defer keyFile.Close()
	if err != nil {
		return err
	}

	// Read key bytes and decode PEM
	keyBytes, err := ioutil.ReadAll(keyFile)
	if err != nil {
		return err
	}

	// We only parse to see if this is indeed a pem formatted rsa public key,
	// but don't use the returned *rsa.PublicKey. Instead, we continue with
	// the original keyBytes from above.
	_, err = ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return err
	}

	// Strip leading and trailing data from PEM file like securesystemslib does
	// TODO: Should we instead use the parsed public key to reconstruct the PEM?
	keyHeader := "-----BEGIN PUBLIC KEY-----"
	keyFooter := "-----END PUBLIC KEY-----"
	keyStart := strings.Index(string(keyBytes), keyHeader)
	keyEnd := strings.Index(string(keyBytes), keyFooter) + len(keyFooter)
	// Successful call to ParseRSAPublicKeyFromPEM already guarantees that
	// header and footer are present, i.e. `!(keyStart == -1 || keyEnd == -1)`
	keyBytesStripped := keyBytes[keyStart:keyEnd]

	// Declare values for key
	// TODO: Do not hardcode here, but define defaults elsewhere and add support
	// for parametrization
	keyType := "rsa"
	scheme := "rsassa-pss-sha256"
	keyIdHashAlgorithms := []string{"sha256", "sha512"}

	// Create partial key map used to create the keyid
	// Unfortunately, we can't use the Key object because this also carries
	// yet unwanted fields, such as KeyId and KeyVal.Private and therefore
	// produces a different hash
	var keyToBeHashed = map[string]interface{}{
		"keytype":               keyType,
		"scheme":                scheme,
		"keyid_hash_algorithms": keyIdHashAlgorithms,
		"keyval": map[string]string{
			"public": string(keyBytesStripped),
		},
	}

	// Canonicalize key and get hex representation of hash
	keyCanonical, err := encodeCanonical(keyToBeHashed)
	if err != nil {
		return err
	}
	keyHashed := sha256.Sum256(keyCanonical)

	// Unmarshalling the canonicalized key into the Key object would seem natural
	// Unfortunately, our mandated canonicalization function produces a byte
	// slice that cannot be unmarshalled by Golang's json decoder, hence we have
	// to manually assign the values
	k.KeyType = keyType
	k.KeyVal = KeyVal{
		Public: string(keyBytesStripped),
	}
	k.Scheme = scheme
	k.KeyIdHashAlgorithms = keyIdHashAlgorithms
	k.KeyId = fmt.Sprintf("%x", keyHashed)

	return nil
}

/*
VerifySignature uses the passed Key to verify the passed Signature over the
passed data.  It returns an error if the key is not a valid RSA public key or
if the signature is not valid for the data.
*/
func VerifySignature(key Key, sig Signature, data []byte) error {
	// Create rsa.PublicKey object from DER encoded public key string as
	// found in the public part of the keyval part of a securesystemslib key dict
	keyReader := strings.NewReader(key.KeyVal.Public)
	pemBytes, err := ioutil.ReadAll(keyReader)
	if err != nil {
		return err
	}
	rsaPub, err := ParseRSAPublicKeyFromPEM(pemBytes)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256(data)

	// Create hex bytes from the signature hex string
	sigHex, _ := hex.DecodeString(sig.Sig)

	// SecSysLib uses a SaltLength of `hashes.SHA256().digest_size`, i.e. 32
	if err := rsa.VerifyPSS(rsaPub, crypto.SHA256, hashed[:], sigHex,
		&rsa.PSSOptions{SaltLength: sha256.Size, Hash: crypto.SHA256}); err != nil {
		return err
	}

	return nil
}
