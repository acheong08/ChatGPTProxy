package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"

	arkose "github.com/acheong08/funcaptcha"
)

var (
	bx string = `[{"key":"api_type","value":"js"},{"key":"p","value":1},{"key":"f","value":"9711bd3695defe0844fb8fd8a722f38b"},{"key":"n","value":"MTY4Nzc4ODQyMw=="},{"key":"wh","value":"80b13fd48b8da8e4157eeb6f9e9fbedb|5ab5738955e0611421b686bc95655ad0"},{"key":"enhanced_fp","value":[{"key":"webgl_extensions","value":null},{"key":"webgl_extensions_hash","value":null},{"key":"webgl_renderer","value":null},{"key":"webgl_vendor","value":null},{"key":"webgl_version","value":null},{"key":"webgl_shading_language_version","value":null},{"key":"webgl_aliased_line_width_range","value":null},{"key":"webgl_aliased_point_size_range","value":null},{"key":"webgl_antialiasing","value":null},{"key":"webgl_bits","value":null},{"key":"webgl_max_params","value":null},{"key":"webgl_max_viewport_dims","value":null},{"key":"webgl_unmasked_vendor","value":null},{"key":"webgl_unmasked_renderer","value":null},{"key":"webgl_vsf_params","value":null},{"key":"webgl_vsi_params","value":null},{"key":"webgl_fsf_params","value":null},{"key":"webgl_fsi_params","value":null},{"key":"webgl_hash_webgl","value":null},{"key":"user_agent_data_brands","value":null},{"key":"user_agent_data_mobile","value":null},{"key":"navigator_connection_downlink","value":null},{"key":"navigator_connection_downlink_max","value":null},{"key":"network_info_rtt","value":null},{"key":"network_info_save_data","value":null},{"key":"network_info_rtt_type","value":null},{"key":"screen_pixel_depth","value":24},{"key":"navigator_device_memory","value":null},{"key":"navigator_languages","value":"en-US,en"},{"key":"window_inner_width","value":0},{"key":"window_inner_height","value":0},{"key":"window_outer_width","value":0},{"key":"window_outer_height","value":0},{"key":"browser_detection_firefox","value":true},{"key":"browser_detection_brave","value":false},{"key":"audio_codecs","value":"{\\"ogg\\":\\"probably\\",\\"mp3\\":\\"maybe\\",\\"wav\\":\\"probably\\",\\"m4a\\":\\"maybe\\",\\"aac\\":\\"maybe\\"}"},{"key":"video_codecs","value":"{\\"ogg\\":\\"probably\\",\\"h264\\":\\"probably\\",\\"webm\\":\\"probably\\",\\"mpeg4v\\":\\"\\",\\"mpeg4a\\":\\"\\",\\"theora\\":\\"\\"}"},{"key":"media_query_dark_mode","value":false},{"key":"headless_browser_phantom","value":false},{"key":"headless_browser_selenium","value":false},{"key":"headless_browser_nightmare_js","value":false},{"key":"document__referrer","value":""},{"key":"window__ancestor_origins","value":null},{"key":"window__tree_index","value":[1]},{"key":"window__tree_structure","value":"[[],[]]"},{"key":"window__location_href","value":"https://tcr9i.chat.openai.com/v2/1.5.2/enforcement.64b3a4e29686f93d52816249ecbf9857.html#35536E1E-65B4-4D96-9D97-6ADB7EFF8147"},{"key":"client_config__sitedata_location_href","value":"https://chat.openai.com/"},{"key":"client_config__surl","value":"https://tcr9i.chat.openai.com"},{"key":"mobile_sdk__is_sdk"},{"key":"client_config__language","value":null},{"key":"audio_fingerprint","value":"35.73833402246237"}]},{"key":"fe","value":["DNT:1","L:en-US","D:24","PR:1","S:0,0","AS:false","TO:0","SS:true","LS:true","IDB:true","B:false","ODB:false","CPUC:unknown","PK:Linux x86_64","CFP:330110783","FR:false","FOS:false","FB:false","JSF:Arial,Arial Narrow,Bitstream Vera Sans Mono,Bookman Old Style,Century Schoolbook,Courier,Courier New,Helvetica,MS Gothic,MS PGothic,Palatino,Palatino Linotype,Times,Times New Roman","P:Chrome PDF Viewer,Chromium PDF Viewer,Microsoft Edge PDF Viewer,PDF Viewer,WebKit built-in PDF","T:0,false,false","H:2","SWF:false"]},{"key":"ife_hash","value":"2a007a5daef41ee943d5fc73a0a8c312"},{"key":"cs","value":1},{"key":"jsbd","value":"{\\"HL\\":2,\\"NCE\\":true,\\"DT\\":\\"\\",\\"NWD\\":\\"false\\",\\"DOTO\\":1,\\"DMTO\\":1}"}]`
	bv string = "Mozilla/5.0 (X11; Linux x86_64; rv:114.0) Gecko/20100101 Firefox/114.0"
	bw string
)

func init() {
	// var bt = new Date() ['getTime']() / 1000
	bt := time.Now().UnixMicro() / 1000000
	// bw = Math.round(bt - (bt % 21600)
	bw = strconv.FormatInt(bt-(bt%21600), 10)
	token, err := get_arkose_token()
	if err != nil {
		log.Fatal(err)
	}
	println(token)
}

func get_arkose_token() (string, error) {
	bda := arkose.Encrypt(bx, bv+bw)
	bda = base64.StdEncoding.EncodeToString([]byte(bda))
	form := url.Values{
		"bda":          {bda},
		"public_key":   {"35536E1E-65B4-4D96-9D97-6ADB7EFF8147"},
		"site":         {"https://chat.openai.com"},
		"userbrowser":  {bv},
		"capi_version": {"1.5.2"},
		"capi_mode":    {"lightbox"},
		"style_theme":  {"default"},
		"rnd":          {strconv.FormatFloat(rand.Float64(), 'f', -1, 64)},
	}
	req, _ := http.NewRequest(http.MethodPost, "https://tcr9i.chat.openai.com/fc/gt2/public_key/35536E1E-65B4-4D96-9D97-6ADB7EFF8147", strings.NewReader(form.Encode()))
	req.Header.Set("Host", "tcr9i.chat.openai.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:114.0) Gecko/20100101 Firefox/114.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://tcr9i.chat.openai.com")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://tcr9i.chat.openai.com/v2/1.5.2/enforcement.64b3a4e29686f93d52816249ecbf9857.html")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("TE", "trailers")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("status code " + resp.Status)
	}
	type arkose_response struct {
		Token string `json:"token"`
	}
	var arkose arkose_response
	err = json.NewDecoder(resp.Body).Decode(&arkose)
	if err != nil {
		return "", err
	}
	println(arkose.Token)
	return arkose.Token, nil

}

// func get_arkose_token() (string, error) {
// 	type arkose struct {
// 		Token string `json:"token"`
// 	}
// 	var result arkose
// 	resp, err := client.Get("https://arkose-token.linweiyuan.com/backup")
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()
// 	err = json.NewDecoder(resp.Body).Decode(&result)
// 	if err != nil {
// 		return "", err
// 	}
// 	return result.Token, nil
// }
