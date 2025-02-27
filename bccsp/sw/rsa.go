/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sw

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/hyperledger/fabric/bccsp"
)

type rsaSigner struct{}

func (s *rsaSigner) Sign(k bccsp.Key, digest []byte, opts bccsp.SignerOpts) ([]byte, error) {
	if opts == nil {
		opts = crypto.SHA256
	}

	return k.(*rsaPrivateKey).privKey.Sign(rand.Reader, digest, opts)
}

type rsaPrivateKeyVerifier struct{}

func (v *rsaPrivateKeyVerifier) Verify(k bccsp.Key, signature, digest []byte, opts bccsp.SignerOpts) (bool, error) {
	if opts == nil {
		opts = &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: crypto.SHA256}
	}
	switch opts.(type) {
	case *rsa.PSSOptions:
		err := rsa.VerifyPSS(&(k.(*rsaPrivateKey).privKey.PublicKey),
			(opts.(*rsa.PSSOptions)).Hash,
			digest, signature, opts.(*rsa.PSSOptions))

		return err == nil, err
	default:
		return false, fmt.Errorf("Opts type not recognized [%s]", opts)
	}
}

type rsaPublicKeyKeyVerifier struct{}

func (v *rsaPublicKeyKeyVerifier) Verify(k bccsp.Key, signature, digest []byte, opts bccsp.SignerOpts) (bool, error) {
	switch k.(type) {
	case *rsaPrivateKey:
		if opts != nil {
			//return false, errors.New("Invalid options. It must not be nil.")
			switch opts.(type) {
			case *rsa.PSSOptions:
				err := rsa.VerifyPSS(&(k.(*rsaPrivateKey).privKey.PublicKey),
					(opts.(*rsa.PSSOptions)).Hash,
					digest, signature, opts.(*rsa.PSSOptions))
				return err == nil, err
			default:
				return false, fmt.Errorf("Opts type not recognized [%s]", opts)
			}
		}
		err := rsa.VerifyPKCS1v15(&(k.(*rsaPrivateKey).privKey.PublicKey), crypto.SHA256, digest, signature)
		return err == nil, err
	case *rsaPublicKey:
		if opts != nil {
			switch opts.(type) {
			case *rsa.PSSOptions:
				err := rsa.VerifyPSS((k.(*rsaPublicKey)).pubKey,
					(opts.(*rsa.PSSOptions)).Hash,
					digest, signature, opts.(*rsa.PSSOptions))

				return err == nil, err
			default:
				return false, fmt.Errorf("Opts type not recognized [%s]", opts)
			}
		}

		err := rsa.VerifyPKCS1v15((k.(*rsaPublicKey)).pubKey, crypto.SHA256, digest, signature)
		return err == nil, err
	}
	return true, nil
}
