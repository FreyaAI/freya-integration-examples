package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Credentials struct {
	UserId      string `json:"user_id"`
	CompanyCode string `json:"company_code"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func readPEMFile(filePath string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func signMessage(message string, privateKey *rsa.PrivateKey) (string, error) {
	hashed := sha256.Sum256([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func postAuthenticationDetails(companyCode, userId, signature, timestamp string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	url := "https://api.freyafashion.ai/api/v1/authenticate"
	body, _ := json.Marshal(map[string]string{
		"company_code": companyCode,
		"user_id":      userId,
		"signature":    signature,
		"timestamp":    timestamp,
	})
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var res AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.Token, nil
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if creds.UserId == "" || creds.CompanyCode == "" {
		http.Error(w, "Missing user_id or company_code", http.StatusBadRequest)
		return
	}
	keyPath := creds.CompanyCode + ".cer"
	privateKey, err := readPEMFile(keyPath)
	if err != nil {
		http.Error(w, "Company code not found", http.StatusNotFound)
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	signature, err := signMessage(timestamp, privateKey)
	if err != nil {
		http.Error(w, "Failed to sign message", http.StatusInternalServerError)
		return
	}
	token, err := postAuthenticationDetails(creds.CompanyCode, creds.UserId, signature, timestamp)
	if err != nil {
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}

func main() {
	http.HandleFunc("/demo/v1/authenticate", authenticate)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
