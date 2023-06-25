package main

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"

	http "github.com/bogdanfinn/fhttp"
)

type arkose_response struct {
	Token string `json:"token"`
}

func random_num(length int) string {
	chars := "0123456789"
	var result string
	for i := 0; i < length; i++ {
		result += string(chars[rand.Intn(len(chars))])
	}
	return result
}

// func random_int(min, max int) int {
// 	return min + rand.Intn(max-min)
// }

// func get_random_hex(length int) string {
// 	chars := "0123456789abcdef"
// 	var result string
// 	for i := 0; i < length; i++ {
// 		result += string(chars[random_int(0, len(chars))])
// 	}
// 	return result
// }
// func base64_encode(data string) string {
// 	return base64.StdEncoding.EncodeToString([]byte(data))
// }

// func random_float() string {
// 	return fmt.Sprintf("%f", rand.Float64())
// }

func get_arkose_token() (string, error) {
	resp, err := client.Get("https://ai.fakeopen.com/api/arkose/params")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	payload, _ := io.ReadAll(resp.Body)
	url := "https://tcr9i.chat.openai.com/fc/gt2/public_key/35536E1E-65B4-4D96-9D97-6ADB7EFF8147"
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	req.Header.Set("Host", "tcr9i.chat.openai.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:114.0) Gecko/20100101 Firefox/114.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://tcr9i.chat.openai.com")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://tcr9i.chat.openai.com/v2/1.5.2/enforcement.64b3a4e29686f93d52816249ecbf9857.html")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("TE", "trailers")
	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var arkose arkose_response
	err = json.NewDecoder(resp.Body).Decode(&arkose)
	if err != nil {
		return "", err
	}
	println(arkose.Token)
	return arkose.Token, nil
}
