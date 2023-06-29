package arkose

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

var (
	jar     = tls_client.NewCookieJar()
	options = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(360),
		tls_client.WithClientProfile(tls_client.Safari_IOS_16_0),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}
	client *tls_client.HttpClient
	bv     = "Mozilla/5.0 (X11; Linux x86_64; rv:114.0) Gecko/20100101 Firefox/114.0"
)

func init() {
	cli, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	client = &cli
	proxy := os.Getenv("http_proxy")
	if proxy != "" {
		(*client).SetProxy(proxy)
	}
}

func SetTLSClient(cli *tls_client.HttpClient) {
	client = cli
}

func GetOpenAIToken() (string, error) {

	bda := getBDA()
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
	resp, err := (*client).Do(req)
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
	return arkose.Token, nil
}

func getBDA() string {
	// Generate initial fingerprint
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1000000000)
	bx := fmt.Sprintf(`[{"key":"api_type","value":"js"},{"key":"p","value":1},{"key":"f","value":"03877290fb0ed162ef8f3b3919883429"},{"key":"n","value":"%s"},{"key":"wh","value":"085833e2c4895769fc0b03fd32281b80|72627afbfd19a741c7da1732218301ac"},{"key":"enhanced_fp","value":[{"key":"webgl_extensions","value":"ANGLE_instanced_arrays;EXT_blend_minmax;EXT_color_buffer_half_float;EXT_disjoint_timer_query;EXT_float_blend;EXT_frag_depth;EXT_shader_texture_lod;EXT_texture_compression_bptc;EXT_texture_compression_rgtc;EXT_texture_filter_anisotropic;WEBKIT_EXT_texture_filter_anisotropic;EXT_sRGB;KHR_parallel_shader_compile;OES_element_index_uint;OES_fbo_render_mipmap;OES_standard_derivatives;OES_texture_float;OES_texture_float_linear;OES_texture_half_float;OES_texture_half_float_linear;OES_vertex_array_object;WEBGL_color_buffer_float;WEBGL_compressed_texture_s3tc;WEBKIT_WEBGL_compressed_texture_s3tc;WEBGL_compressed_texture_s3tc_srgb;WEBGL_debug_renderer_info;WEBGL_debug_shaders;WEBGL_depth_texture;WEBKIT_WEBGL_depth_texture;WEBGL_draw_buffers;WEBGL_lose_context;WEBKIT_WEBGL_lose_context;WEBGL_multi_draw"},{"key":"webgl_extensions_hash","value":"5e287cfa20ad57f8652f0b2320f0de21"},{"key":"webgl_renderer","value":"WebKit WebGL"},{"key":"webgl_vendor","value":"WebKit"},{"key":"webgl_version","value":"WebGL 1.0 (OpenGL ES 2.0 Chromium)"},{"key":"webgl_shading_language_version","value":"WebGL GLSL ES 1.0 (OpenGL ES GLSL ES 1.0 Chromium)"},{"key":"webgl_aliased_line_width_range","value":"[1, 1]"},{"key":"webgl_aliased_point_size_range","value":"[1, 1024]"},{"key":"webgl_antialiasing","value":"yes"},{"key":"webgl_bits","value":"8,8,24,8,8,0"},{"key":"webgl_max_params","value":"16,32,16384,1024,16384,16,16384,30,16,16,4095"},{"key":"webgl_max_viewport_dims","value":"[32767, 32767]"},{"key":"webgl_unmasked_vendor","value":"Google Inc. (NVIDIA)"},{"key":"webgl_unmasked_renderer","value":"ANGLE (NVIDIA, NVIDIA GeForce RTX 3060 Ti Direct3D11 vs_5_0 ps_5_0, D3D11-31.0.15.2698)"},{"key":"webgl_vsf_params","value":"23,127,127,23,127,127,23,127,127"},{"key":"webgl_vsi_params","value":"0,31,30,0,31,30,0,31,30"},{"key":"webgl_fsf_params","value":"23,127,127,23,127,127,23,127,127"},{"key":"webgl_fsi_params","value":"0,31,30,0,31,30,0,31,30"},{"key":"webgl_hash_webgl","value":"90a4bbe26aaf9f91bcb94e74592eb705"},{"key":"user_agent_data_brands","value":"Chromium, Not;A Brand"},{"key":"user_agent_data_mobile","value":false},{"key":"navigator_connection_downlink","value":1.0},{"key":"navigator_connection_downlink_max","value":null},{"key":"network_info_rtt","value":400},{"key":"network_info_save_data","value":false},{"key":"network_info_rtt_type","value":null},{"key":"screen_pixel_depth","value":24},{"key":"navigator_device_memory","value":8},{"key":"navigator_languages","value":"zh-CN,zh"},{"key":"window_inner_width","value":0},{"key":"window_inner_height","value":0},{"key":"window_outer_width","value":3338},{"key":"window_outer_height","value":1440},{"key":"browser_detection_firefox","value":false},{"key":"browser_detection_brave","value":false},{"key":"audio_codecs","value":"{\"ogg\":\"probably\",\"mp3\":\"probably\",\"wav\":\"probably\",\"m4a\":\"maybe\",\"aac\":\"probably\"}"},{"key":"video_codecs","value":"{\"ogg\":\"probably\",\"h264\":\"probably\",\"webm\":\"probably\",\"mpeg4v\":\"\",\"mpeg4a\":\"\",\"theora\":\"\"}"},{"key":"media_query_dark_mode","value":false},{"key":"headless_browser_phantom","value":false},{"key":"headless_browser_selenium","value":false},{"key":"headless_browser_nightmare_js","value":false},{"key":"document__referrer","value":""},{"key":"window__ancestor_origins","value":["https://chat.openai.com"]},{"key":"window__tree_index","value":[1]},{"key":"window__tree_structure","value":"[[],[]]"},{"key":"window__location_href","value":"https://tcr9i.chat.openai.com/v2/1.5.2/enforcement.64b3a4e29686f93d52816249ecbf9857.html#35536E1E-65B4-4D96-9D97-6ADB7EFF8147"},{"key":"client_config__sitedata_location_href","value":"https://chat.openai.com/"},{"key":"client_config__surl","value":"https://tcr9i.chat.openai.com"},{"key":"mobile_sdk__is_sdk"},{"key":"client_config__language","value":null},{"key":"navigator_battery_charging","value":true},{"key":"audio_fingerprint","value":"124.04347527516074"}]},{"key":"fe","value":["DNT:unknown","L:zh-CN","D:24","PR:1","S:3440,1440","AS:3338,1440","TO:-480","SS:true","LS:true","IDB:true","B:false","ODB:true","CPUC:unknown","PK:Win32","CFP:-1034319970","FR:false","FOS:false","FB:false","JSF:Arial,Arial Black,Arial Narrow,Calibri,Cambria,Cambria Math,Comic Sans MS,Consolas,Courier,Courier New,Georgia,Helvetica,Impact,Lucida Console,Lucida Sans Unicode,Microsoft Sans Serif,MS Gothic,MS PGothic,MS Sans Serif,MS Serif,Palatino Linotype,Segoe Print,Segoe Script,Segoe UI,Segoe UI Light,Segoe UI Semibold,Segoe UI Symbol,Tahoma,Times,Times New Roman,Trebuchet MS,Verdana,Wingdings","P:360SoftMgrPlugin,AliSSOLogin plugin,Alipay Security Control  3,Alipay security control,BaiduYunGuanjia Application,Chromium PDF Plugin,Chromium PDF Viewer,Java Deployment Toolkit 8.0.3410.10,Java(TM) Platform SE 8 U341,Microsoft® Windows Media Player Firefox Plugin,Tencent QQ,iTrusChina iTrusPTA,XEnroll,iEnroll,hwPTA,UKeyInstalls Firefox Plugin","T:0,false,false","H:16","SWF:false"]},{"key":"ife_hash","value":"50f49e991788c8789817eede6e3e9ffa"},{"key":"cs","value":1},{"key":"jsbd","value":"{\"HL\":9,\"NCE\":true,\"DT\":\"\",\"NWD\":\"false\",\"DOTO\":1,\"DMTO\":1}"}]`,
		base64.StdEncoding.EncodeToString([]byte(timestamp)))
	bt := time.Now().UnixMicro() / 1000000
	bw := strconv.FormatInt(bt-(bt%21600), 10)
	return encrypt(bx, bv+bw)
}
