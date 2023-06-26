package main

import "encoding/json"

// import (
// 	arkose "github.com/flyingpot/funcaptcha"
// )

// func get_arkose_token() (string, error) {
// 	options := arkose.GetTokenOptions{
// 		PKey: "35536E1E-65B4-4D96-9D97-6ADB7EFF8147",
// 		SURL: "https://tcr9i.chat.openai.com",
// 		Headers: map[string]string{
// 			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36",
// 		},
// 		Site: "https://chat.openai.com",
// 	}
// 	result, err := arkose.GetToken(&options)
// 	println(result.Token)
// 	return result.Token, err
// }

func get_arkose_token() (string, error) {
	type arkose struct {
		Token string `json:"token"`
	}
	var result arkose
	resp, err := client.Get("https://arkose-token.linweiyuan.com/backup")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Token, nil
}
