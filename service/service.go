package service

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/seb7887/heimdallr/storage"
)

type Service interface {
	Create(ctx context.Context, clientId string) (*string, error)
	Authenticate(ctx context.Context, clientId string, token string) bool
	UpdateBlacklist(ctx context.Context, clientId string) error
	ReadBlacklist(ctx context.Context) ([]string, error)
	Delete(ctx context.Context, clientId string) error
}

type service struct {
	repository storage.Repository
}

type keyPair struct {
	privateKey string
	publicKey  string
}

func NewService(repo storage.Repository) Service {
	return &service{repository: repo}
}

func generateKeys() (*keyPair, error) {
	pubkeyCurve := elliptic.P256()
	// generate public and private keys
	privateKey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.PublicKey

	// x509 serialization
	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, err
	}

	// Parse to string
	privBlock := pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	}
	var privKeyRow bytes.Buffer
	err = pem.Encode(&privKeyRow, &privBlock)
	if err != nil {
		return nil, err
	}

	pubBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	var pubKeyRow bytes.Buffer
	err = pem.Encode(&pubKeyRow, &pubBlock)
	if err != nil {
		return nil, err
	}

	return &keyPair{
		privateKey: privKeyRow.String(),
		publicKey:  pubKeyRow.String(),
	}, nil
}

func (s *service) Create(ctx context.Context, clientId string) (*string, error) {
	keyPair, err := generateKeys()
	if err != nil {
		return nil, err
	}

	client := &storage.Client{
		ClientId:  clientId,
		PublicKey: keyPair.publicKey,
	}

	err = s.repository.CreateClient(ctx, client)
	if err != nil {
		return nil, err
	}

	return &keyPair.privateKey, nil
}

func validateToken(token string, publicKey []byte) bool {
	key, err := jwt.ParseECPublicKeyFromPEM(publicKey)
	if err != nil {
		return false
	}

	parts := strings.Split(token, ".")
	err = jwt.SigningMethodES256.Verify(strings.Join(parts[0:2], "."), parts[2], key)
	if err != nil {
		return false
	}

	return true
}

func (s *service) Authenticate(ctx context.Context, clientId string, token string) bool {
	// get client public key
	publicKey, err := s.repository.GetClientKey(ctx, clientId)
	if err != nil {
		return false
	}

	// check if client exists in blacklist
	blacklist, err := s.repository.GetBlacklist(ctx)
	if exists(blacklist, clientId) || err != nil {
		return false
	}

	// validate token
	return validateToken(token, publicKey)
}

func (s *service) UpdateBlacklist(ctx context.Context, clientId string) error {
	blacklist, err := s.repository.GetBlacklist(ctx)
	if err != nil {
		return err
	}

	if exists(blacklist, clientId) {
		return fmt.Errorf("Client is already a member")
	}
	blacklist = append(blacklist, clientId)

	return s.repository.UpsertBlacklist(ctx, blacklist)
}

func (s *service) ReadBlacklist(ctx context.Context) ([]string, error) {
	return s.repository.GetBlacklist(ctx)
}

func (s *service) Delete(ctx context.Context, clientId string) error {
	return s.repository.DeleteClient(ctx, clientId)
}

func exists(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}
	return false
}
