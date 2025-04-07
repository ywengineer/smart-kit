package oauths

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/ywengineer/smart-kit/pkg/nets"
	"github.com/ywengineer/smart-kit/pkg/utilk"
	"time"
)

func NewGameCenterAuth(bundleId string) AuthFacade {
	return &gameCenterAuth{appBundleId: bundleId, rootCertPath: "https://www.apple.com/certificateauthority/AppleRootCA-G3.cer"}
}

type gameCenterAuth struct {
	appBundleId  string
	rootCertPath string
	rootCertPool *x509.CertPool
}

func (g *gameCenterAuth) Validate(metadata string) (AuthFacade, error) {
	if len(g.appBundleId) == 0 {
		return nil, errors.New("missing prop [app-bundle-id] for game center")
	}
	sc, body, err := cli.Get(context.Background(), g.rootCertPath)
	if err != nil || !nets.Is2xx(sc) {
		return nil, fmt.Errorf("failed to download Apple root certificate: %w", err)
	}
	// 解析根证书
	rootCert, err := x509.ParseCertificate(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Apple root certificate: %w", err)
	}
	// 步骤 2: 构建证书池
	g.rootCertPool = x509.NewCertPool()
	g.rootCertPool.AddCert(rootCert)
	return g, nil
}

func (g *gameCenterAuth) GetToken(code string) (*AccessToken, error) {
	return &AccessToken{
		AccessToken: code,
	}, nil
}

func (g *gameCenterAuth) RefreshToken(refreshToken string) (*RefreshToken, error) {
	return nil, errors.New("unsupported")
}

func (g *gameCenterAuth) validateCert(publicKeyData []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKeyData)
	if block == nil {
		return nil, errors.New("failed to decode public key PEM")
	}
	publicKeyCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key certificate: %w", err)
	}
	// 步骤 4: 验证证书链
	_, err = publicKeyCert.Verify(x509.VerifyOptions{Roots: g.rootCertPool})
	if err != nil {
		return nil, fmt.Errorf("public key certificate verification failed: %w", err)
	}
	return publicKeyCert.PublicKey.(*rsa.PublicKey), nil
}

func (g *gameCenterAuth) GetUserInfo(_ string, accessToken string) (*UserInfo, error) {
	// 解析票据
	var ticket GameCenterTicket
	err := sonic.UnmarshalString(accessToken, &ticket)
	if err != nil {
		return nil, err
	}
	if ticket.AppBundleID != g.appBundleId {
		return nil, errors.New("invalid app bundle id")
	}
	// 验证时间戳
	now := time.Now().Unix()
	// 步骤 1: 检查时间戳以减轻重放攻击
	// 这里可以根据实际情况调整时间范围，例如允许 5 分钟的误差
	if now-ticket.Timestamp > 300 || ticket.Timestamp > now {
		return nil, errors.New("timestamp is not recent")
	}
	// 步骤 2: 下载公钥
	sc, body, err := cli.Get(context.Background(), ticket.PublicKey)
	if err != nil || !nets.Is2xx(sc) {
		return nil, fmt.Errorf("failed to download public key: %w, sc = %d", err, sc)
	}
	// 步骤 3: 验证公钥是否由 Apple 签名
	pubKey, err := g.validateCert(body)
	if err != nil {
		return nil, err
	}
	// 步骤 4: 拼接数据
	buf := utilk.NewLinkBuffer([]byte{})
	_, _ = buf.WriteBinary([]byte(ticket.PlayerID))
	_, _ = buf.WriteBinary([]byte(ticket.AppBundleID))
	// 将时间戳转换为大端序的 UInt64 格式
	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(ticket.Timestamp))
	_, _ = buf.WriteBinary(timestampBytes)
	_, _ = buf.WriteBinary([]byte(ticket.Salt))
	// 步骤 5: 计算数据的哈希
	hash := sha256.Sum256(buf.Bytes())
	_ = buf.Release()
	// 步骤 6: 解码签名
	sigBytes, err := base64.StdEncoding.DecodeString(ticket.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}
	// 步骤 7: 使用公钥验证签名
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], sigBytes)
	if err != nil {
		return nil, fmt.Errorf("signature verification failed: %w", err)
	}
	//-------
	return &UserInfo{
		Openid:     ticket.PlayerID,
		Nickname:   ticket.PlayerID,
		Sex:        0,
		Province:   "",
		City:       "",
		Country:    "",
		HeadImgUrl: "",
		Privilege:  nil,
		UnionId:    ticket.PlayerID,
	}, nil
}

// GameCenterTicket 定义 Game Center 票据结构体
type GameCenterTicket struct {
	AppBundleID string `json:"app_bundle_id"`
	Timestamp   int64  `json:"timestamp"`
	PlayerID    string `json:"sub"`
	PublicKey   string `json:"public_key"`
	Salt        string `json:"salt"`
	Signature   string `json:"signature"`
}
