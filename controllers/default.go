package controllers

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	//"io/ioutil"
	//"reflect"

	//"net/http"
	"os"
	//"runtime"
	"strings"
	"sync"

	"math/rand"
	"time"

	"encoding/json"
	//"flag"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var adinfo_map map[int]adinfo
var adinfo_map_lock sync.RWMutex

// Adinfo 信息
type adinfo struct {
	Adid     int
	Price    int
	Level    int
	Weight   int
	Is_https int
	Banner   banner
	Video    video
	Native   native
	Ext      adinfo_ext
	Ad_type  int
}

type adinfo_ext struct {
	Action              int
	Imptrackers         []string
	Clktrackers         []string
	Clkurl              string
	Html_snippet        string
	Inventory_type      int
	Fallback_url        string
	Fallback_action     int
	Download_file_name  string
	Title               string
	Desc                string
	Play_start_trackers []string
	Play_end_trackers   []string
	Nurl                string
	Adm                 string
}

type banner struct {
	Weight       int
	Height       int
	Src          string
	BannerAdType int
	Mime         string
}

type video struct {
	Weight    int
	Height    int
	Src       string
	Mime      string
	Linearity bool
	Duration  int
	Protocol  int
}

type native struct {
	ver    string
	Assets []asset
}

type asset struct {
	Id          int
	Asset_oneof int
	Title       native_title
	Image       native_image
	Video       native_video
	Data        native_data
}

type native_title struct {
	Len   int
	Title string
}
type native_image struct {
	Width          int
	Height         int
	Src            string
	Mime           string
	ImageAssetType int
}
type native_video struct {
	Width     int
	Height    int
	Src       string
	Mime      string
	Linearity bool
	Duration  int
}
type native_data struct {
	DataAssetType int
	Len           int
	Data          string
}

// Adinfo 信息

// request 信息
type request struct {
	Id     string         `json:"id"`
	Imp    []request_imp  `json:"imp"`
	App    request_app    `json:"app"`
	Device request_device `json:"device"`
	Ext    request_ext    `json:"ext"`
}

type request_imp struct {
	Id          string             `json:"id"`
	Banner      request_imp_banner `json:"banner"`
	Video       request_imp_video  `json:"video"`
	Native      request_imp_native `json:"native"`
	Instl       bool               `json:"instl"`
	Tagid       string             `json:"tagid"`
	Bidfloor    int                `json:"bidfloor"`
	Bidfloorcur string             `json:"bidfloorcur"`
	Ext         request_imp_ext    `json:"ext"`
}
type request_imp_banner struct {
	W   int `json:"w"`
	H   int `json:"h"`
	Pos int `json:"pos"`
}

type request_imp_ext struct {
	Is_splash_screen bool   `json:"is_splash_screen"`
	Inventory_types  []int  `json:"inventory_types"`
	Ad_type          int    `json:"ad_type"`
	Tag_name         string `json:"tag_name"`
}
type request_imp_video struct {
	Mimes          []string `json:"mimes"`
	Linearity      int      `json:"linearity"`
	Minduration    int      `json:"minduration"`
	Maxduration    int      `json:"maxduration"`
	Protocol       int      `json:"protocol"`
	Protocols      []int    `json:"protocols"`
	W              int      `json:"w"`
	H              int      `json:"h"`
	Startdelay     int      `json:"startdelay"`
	Sequence       int      `json:"sequence"`
	Battr          []int    `json:"battr"`
	Minextended    int      `json:"minextended"`
	Maxextended    int      `json:"maxextended"`
	Minbitrate     int      `json:"minbitrate"`
	Maxbitrate     int      `json:"maxbitrate"`
	Boxingallowed  bool     `json:"boxingallowed"`
	Playbackmethod []int    `json:"playbackmethod"`
	Delivery       []int    `json:"delivery"`
	Pos            int      `json:"pos"`
	Companionad    []int    `json:"companionad"`
	Companionad_21 int      `json:"companionad_21"`
	Api            []int    `json:"api"`
	Companiontype  []int    `json:"companiontype"`
}

type request_imp_native struct {
	RequestOneof request_imp_native_RequestOneof `json:"RequestOneof"`
}

type request_imp_native_RequestOneof struct {
	RequestNative request_imp_native_requestnative `json:"RequestNative"`
	Request       string                           `json:"request"`
	Ver           string                           `json:"ver"`
	Api           []int                            `json:"api"`
	Battr         []int                            `json:"battr"`
}
type request_imp_native_requestnative struct {
	Ver      string                                   `json:"ver"`
	Layout   int                                      `json:"layout"`
	Adunit   int                                      `json:"adunit"`
	Plcmtcnt int                                      `json:"plcmtcnt"`
	Seq      int                                      `json:"seq"`
	Assets   []request_imp_native_requestnative_asset `json:"assets"`
}
type request_imp_native_requestnative_asset struct {
	Id         int                                               `json:"id"`
	Required   bool                                              `json:"required"`
	AssetOneof request_imp_native_requestnative_asset_AssetOneof `json:"AssetOneof"`
}

type request_imp_native_requestnative_asset_AssetOneof struct {
	Title request_imp_native_requestnative_asset_title `json:"title"`
	Img   request_imp_native_requestnative_asset_image `json:"img"`
	Video request_imp_video                            `json:"video"`
	Data  request_imp_native_requestnative_asset_data  `json:"data"`
}

type request_imp_native_requestnative_asset_title struct {
	Len int `json:"len"`
}
type request_imp_native_requestnative_asset_image struct {
	Type  int      `json:"type"`
	W     int      `json:"w"`
	H     int      `json:"h"`
	Wmin  int      `json:"wmin"`
	Hmin  int      `json:"hmin"`
	Wmax  int      `json:"wmax"`
	Hmax  int      `json:"hmax"`
	Mimes []string `json:"mimes"`
}
type request_imp_native_requestnative_asset_data struct {
	Type int `json:"type"`
	Len  int `json:"len"`
}

type request_app struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Ver    string `json:"ver"`
	Bundle string `json:"bundle"`
}
type request_device struct {
	Dnt             bool               `json:"dnt"`
	Ua              string             `json:"ua"`
	Ip              string             `json:"ip"`
	Geo             request_device_geo `json:"geo"`
	Didsha1         string             `json:"didsha1"`
	Dpidsha1        string             `json:"dpidsha1"`
	Make            string             `json:"make"`
	Model           string             `json:"model"`
	Os              string             `json:"os"`
	Osv             string             `json:"osv"`
	W               int                `json:"w"`
	H               int                `json:"h"`
	Ppi             int                `json:"ppi"`
	Connectyiontype int                `json:"connectiontype"`
	Devicetype      int                `json:"devicetype"`
	Macsha1         string             `json:"macsha1"`
	Ext             request_device_ext `json:"ext"`
}

type request_device_ext struct {
	Plmm        string `json:"plmm"`
	Imei        string `json:"imei"`
	Imsi        string `json:"imsi"`
	Mac         string `json:"mac"`
	Android_id  string `json:"android_id"`
	Adid        string `json:"adid"`
	Orientation int    `json:"orientation"`
}
type request_device_geo struct {
	Lat     float64                `json:"lat"`
	Lon     float64                `json:"lon"`
	Country string                 `json:"country"`
	Region  string                 `json:"region"`
	City    string                 `json:"city"`
	Type    int                    `json:"type"`
	Ext     request_device_geo_ext `json:"ext"`
}

type request_device_geo_ext struct {
	Accu int `json:"accu"`
}
type request_ext struct {
	Version    int  `json:"version"`
	Need_https bool `json:"need_https"`
}

// request 信息

// response 信息
type response struct {
	Id      string             `json:"id"`
	Seatbid []response_seatbid `json:"seatbid"`
	Ts      int64              `json:"ts"`
}

type response_seatbid struct {
	Bid []response_seatbid_bid `json:"bid"`
}
type response_seatbid_bid struct {
	Id                  string                        `json:"id"`
	Impid               string                        `json:"impid"`
	Price               float64                       `json:"price"`
	Adid                string                        `json:"adid"`
	Admoneof            response_seatbid_bid_admoneof `json:"admoneof,omitempty"`
	W                   int                           `json:"w,omitempty"`
	H                   int                           `json:"h,omitempty"`
	Iurl                string                        `json:"iurl"`
	Nurl                string                        `json:"nurl,omitempty"`
	Adm                 string                        `json:"adm,omitempty"`
	Ext                 response_seatbid_bid_ext      `json:"ext,omitempty"`
	Fallback_url        string                        `json:"fallback_url,omitempty"`
	Fallback_action     int                           `json:"fallback_action,omitempty"`
	Play_start_trackers []string                      `json:"play_start_trackers,omitempty"`
	Play_end_trackers   []string                      `json:"play_end_trackers,omitempty"`
}

type response_seatbid_bid_ext struct {
	Clkurl         string   `json:"clkurl,omitempty"`
	Imptrackers    []string `json:"imptrackers,omitempty"`
	Clktrackers    []string `json:"clktrackers,omitempty"`
	Title          string   `json:"title,omitempty"`
	Desc           string   `json:"desc,omitempty"`
	Action         int      `json:"action,omitempty"`
	Html_snippet   string   `json:"html_snappet,omitempty"`
	Inventory_type int      `json:"inventory_type,omitempty"`
}

type response_seatbid_bid_admoneof struct {
	Adm       string                                  `json:"adm,omitempty"`
	Admnative response_seatbid_bid_admoneof_admnative `json:"admnative,omitempty"`
}

type response_seatbid_bid_admoneof_admnative struct {
	Assets []response_seatbid_bid_admoneof_admnative_asset `json:"assets,omitempty"`
	Link   response_seatbid_bid_admoneof_admnative_link    `json:"link,omitempty"`
}

type response_seatbid_bid_admoneof_admnative_asset struct {
	Id    int                                                 `json:"id"`
	Title response_seatbid_bid_admoneof_admnative_asset_title `json:"title,omitempty"`
	Img   response_seatbid_bid_admoneof_admnative_asset_image `json:"image,omitempty"`
	Video response_seatbid_bid_admoneof_admnative_asset_video `json:"video,omitempty"`
	Data  response_seatbid_bid_admoneof_admnative_asset_data  `json:"data,omitempty"`
}

type response_seatbid_bid_admoneof_admnative_asset_title struct {
	Text string `json:"text,omitempty"`
}
type response_seatbid_bid_admoneof_admnative_asset_image struct {
	Url string `json:"url,omitempty"`
	W   int    `json:"w,omitempty"`
	H   int    `json:"h,omitempty"`
}

type response_seatbid_bid_admoneof_admnative_asset_video struct {
	Url string `json:"url,omitempty"`
	W   int    `json:"w,omitempty"`
	H   int    `json:"h,omitempty"`
}

type response_seatbid_bid_admoneof_admnative_asset_data struct {
	Value string `json:"value,omitempty"`
}

type response_seatbid_bid_admoneof_admnative_link struct {
	Url string `json:"url"`
}

func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Ctx.WriteString("this is hello")
}

func initLogger() (err error) {

	config := make(map[string]interface{})
	config["filename"] = beego.AppConfig.String("log_path")

	// map 转 json
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("initLogger failed, marshal err:", err)
		return
	}
	// log 的配置
	beego.SetLogger(logs.AdapterFile, string(configStr))
	// log打印文件名和行数
	beego.SetLogFuncCall(true)
	fmt.Println(string(configStr))
	return
}

func (c *MainController) Load() {
	c.Ctx.WriteString("this is Load \n")
	c.Ctx.WriteString(beego.AppConfig.String("adinfo_file") + "\n")

	inputFile, inputError := os.Open(beego.AppConfig.String("adinfo_file"))
	if inputError != nil {
		c.Ctx.WriteString("An error occurred on opening: adinfo_file")
		return
	}
	defer inputFile.Close()
	inputReader := bufio.NewReader(inputFile)
	adinfo_map_lock.Lock()
	adinfo_map = make(map[int]adinfo)
	i := 0
	for {
		inputString, readerError := inputReader.ReadString('\n')
		inputString = strings.Replace(inputString, "\n", "", 1)
		inputString = strings.Replace(inputString, "\r", "", 1)
		inputString = strings.Replace(inputString, "\t", "", 1)
		if readerError == io.EOF {
			break
		}
		if len(inputString) > 0 {
			//adinfo_map[i] = inputString
		}
		var someOne adinfo

		if err := json.Unmarshal([]byte(inputString), &someOne); err == nil {
			//fmt.Println("someOne:", someOne)
			//fmt.Println("someOne.adid:", someOne.Adid)
		} else {
			//fmt.Println("someOne.adid err :", err)
		}
		c.Ctx.WriteString(" i:" + strconv.Itoa(i) + " inputString:" + string(inputString) + "\n")

		if len(someOne.Native.Assets) > 0 {
			someOne.Ad_type = 5
		} else if someOne.Video.Weight > 0 {
			someOne.Ad_type = 3
		} else if someOne.Banner.Weight > 0 {
			someOne.Ad_type = 1
		}

		adinfo_map[someOne.Adid] = someOne
		i++
	}
	adinfo_map_lock.Unlock()

	return
}

func (c *MainController) List() {
	c.Ctx.WriteString("this is List \n")
	for key, value := range adinfo_map {
		c.Ctx.WriteString(" key:" + strconv.Itoa(key) + " value:")
		jsonStr, err := json.Marshal(value)
		if err == nil {
			c.Ctx.WriteString(string(jsonStr))
		} else {
			fmt.Println("json err ", err)
			c.Ctx.WriteString("json error")
		}
		c.Ctx.WriteString("\n")
	}
	return
}
func (c *MainController) Clear() {
	c.Ctx.WriteString("this is Clear \n")
	adinfo_map_lock.Lock()
	adinfo_map = make(map[int]adinfo)
	adinfo_map_lock.Unlock()

	return
}

func (c *MainController) GetAdJson() {

	//c.Ctx.WriteString("this is GetAdJson \n")
	//c.Ctx.WriteString(string(c.Ctx.Input.RequestBody) + "\n")
	//c.Ctx.WriteString("c.Ctx.Input.RequestBody end \n")
	inputString := c.Ctx.Input.RequestBody
	var requestJson request
	var responseJson response
	need_ad_type := 0
	res_adid := 0
	start := time.Now().Nanosecond()
	fmt.Println("requestJson.Id:", requestJson.Id)
	fmt.Println("responseJson.Id:", responseJson.Id)

	responseJson.Ts = 0
	responseJson.Seatbid = []response_seatbid{}
	if err := json.Unmarshal([]byte(inputString), &requestJson); err == nil {

		responseJson.Id = requestJson.Id
		if requestJson.Imp[0].Native.RequestOneof.RequestNative.Layout > 0 {
			need_ad_type = 5
		} else if requestJson.Imp[0].Video.W > 0 {
			need_ad_type = 3
		} else if requestJson.Imp[0].Banner.W > 0 {
			need_ad_type = 1
		}
		//responseJson = searchAd(requestJson)
		if Adinfo, err := searchAd(requestJson); err == nil {

			if Adinfo.Adid > 0 {
				res_adid = Adinfo.Adid
				// 这里继续填写广告信息

				var responseJson_seatbid response_seatbid
				var responseJson_seatbid_bid response_seatbid_bid
				var responseJson_seatbid_bid_ext response_seatbid_bid_ext
				responseJson_seatbid_bid.Id = "id-" + strconv.Itoa(Adinfo.Adid)
				responseJson_seatbid_bid.Impid = "impid-" + strconv.Itoa(Adinfo.Adid)
				responseJson_seatbid_bid.Price = float64(Adinfo.Price) / 100
				responseJson_seatbid_bid.Adid = "Adid-" + strconv.Itoa(Adinfo.Adid)
				responseJson_seatbid_bid.W = Adinfo.Banner.Weight
				responseJson_seatbid_bid.H = Adinfo.Banner.Height
				responseJson_seatbid_bid.Iurl = Adinfo.Banner.Src
				responseJson_seatbid_bid.Adm = Adinfo.Ext.Adm
				responseJson_seatbid_bid.Fallback_url = Adinfo.Ext.Fallback_url
				responseJson_seatbid_bid.Fallback_action = Adinfo.Ext.Fallback_action
				responseJson_seatbid_bid.Nurl = Adinfo.Ext.Nurl
				responseJson_seatbid_bid.Play_start_trackers = Adinfo.Ext.Play_start_trackers
				responseJson_seatbid_bid.Play_end_trackers = Adinfo.Ext.Play_end_trackers

				responseJson_seatbid_bid_ext.Clkurl = Adinfo.Ext.Clkurl
				/*
					for _, tmp := range Adinfo.Ext.Imptrackers {
						responseJson_seatbid_bid_ext.Imptrackers = append(responseJson_seatbid_bid_ext.Imptrackers, tmp)
					}
					for _, tmp := range Adinfo.Ext.Clktrackers {
						responseJson_seatbid_bid_ext.Clktrackers = append(responseJson_seatbid_bid_ext.Clktrackers, tmp)
					}
				*/
				responseJson_seatbid_bid_ext.Imptrackers = responseJson_seatbid_bid_ext.Imptrackers
				responseJson_seatbid_bid_ext.Clktrackers = responseJson_seatbid_bid_ext.Clktrackers
				responseJson_seatbid_bid_ext.Action = Adinfo.Ext.Action
				responseJson_seatbid_bid_ext.Inventory_type = Adinfo.Ext.Inventory_type

				responseJson_seatbid_bid.Ext = responseJson_seatbid_bid_ext

				// 原生广告，填充assets信息
				if requestJson.Imp[0].Ext.Ad_type == 3 {
					var admoneof response_seatbid_bid_admoneof
					var admnative response_seatbid_bid_admoneof_admnative
					//
					for _, nv := range Adinfo.Native.Assets {
						//jsonStrtemp, _ := json.MarshalIndent(nv, "", "\t") //格式化编码
						//fmt.Println("jsonStrtemp:", string(jsonStrtemp))
						var admnative_asset response_seatbid_bid_admoneof_admnative_asset

						if nv.Asset_oneof == 3 {
							var admnative_asset_title response_seatbid_bid_admoneof_admnative_asset_title
							admnative_asset_title.Text = nv.Title.Title
							//fmt.Println("title:", nv.Title.Title)
							admnative_asset.Id = nv.Id
							admnative_asset.Title = admnative_asset_title
						} else if nv.Asset_oneof == 4 {
							var admnative_asset_image response_seatbid_bid_admoneof_admnative_asset_image
							admnative_asset_image.Url = nv.Image.Src
							admnative_asset.Id = nv.Id
							admnative_asset.Img = admnative_asset_image
						} else if nv.Asset_oneof == 5 {
							var admnative_asset_video response_seatbid_bid_admoneof_admnative_asset_video
							admnative_asset_video.Url = nv.Video.Src
							//fmt.Println("title:", nv.Title.Title)
							admnative_asset.Id = nv.Id
							admnative_asset.Video = admnative_asset_video
						} else if nv.Asset_oneof == 6 {
							var admnative_asset_data response_seatbid_bid_admoneof_admnative_asset_data
							admnative_asset_data.Value = nv.Data.Data
							//fmt.Println("title:", nv.Title.Title)
							admnative_asset.Id = nv.Id
							admnative_asset.Data = admnative_asset_data
						}
						admnative.Assets = append(admnative.Assets, admnative_asset)
					}
					admnative.Link.Url = Adinfo.Ext.Clkurl
					admoneof.Admnative = admnative
					responseJson_seatbid_bid.Admoneof = admoneof
					fmt.Println(" Admoneof admoneof ")
				}

				responseJson_seatbid.Bid = append(responseJson_seatbid.Bid, responseJson_seatbid_bid)
				responseJson.Seatbid = append(responseJson.Seatbid, responseJson_seatbid)
				fmt.Printf("add adid \n")
			} else {
				//fmt.Printf(responseJson.Seatbid)
				fmt.Printf("add not adid \n")
			}
		}

	} else {
		//c.Ctx.WriteString("requestJson.id err : \n")
	}

	//time.Sleep(100 * time.Nanosecond)
	end := time.Now().Nanosecond()
	//fmt.Fprintln(writer," responseJson is ",responseJson)
	responseJson.Ts = int64(end - start)
	fmt.Println("start:", start)
	fmt.Println("end:", end)

	c.Data["json"] = responseJson
	c.ServeJSON()
	/*
		//jsonStr, err := json.MarshalIndent(responseJson, "", "\t") //格式化编码
		jsonStr, err := json.Marshal(responseJson) //格式化编码
		if err != nil {
			c.Ctx.WriteString(" responseJson is failed \n")
		}

		//fmt.Println("response json : ", jsonStr)

		c.Ctx.WriteString(string(jsonStr))
	*/
	beego.Info(need_ad_type, "\x02", res_adid)

	return
}

func searchAd(requestJson request) (adinfo adinfo, err error) {

	// 筛选广告
	need_ad_type := 0
	var adinfo_map_res = []int{}
	//map_adid_res := make(map[int])
	if requestJson.Imp[0].Native.RequestOneof.RequestNative.Layout > 0 {
		fmt.Println(" is native ")
		need_ad_type = 5
	} else if requestJson.Imp[0].Video.W > 0 {
		fmt.Println(" is video ")
		need_ad_type = 3
	} else if requestJson.Imp[0].Banner.W > 0 {
		fmt.Println(" is banner ")
		need_ad_type = 1
	}

	for _, value := range adinfo_map {

		//jsonStr, err := json.Marshal(value)
		if err == nil {

			if value.Ad_type == need_ad_type {
				//fmt.Println(" this " + strconv.Itoa(need_ad_type))
				adinfo_map_res = append(adinfo_map_res, value.Adid)
			}
		} else {
		}
	}

	//fmt.Println(adinfo_map_res)
	//fmt.Println(len(adinfo_map_res))
	//fmt.Println(GenerateRangeNum(0, len(adinfo_map_res)))
	//fmt.Println("assets:", requestJson.Imp[0].Native.RequestOneof.RequestNative.Assets[0].Id)
	//

	// 级别筛选

	// 比重随机

	// 获取最终结果
	if len(adinfo_map_res) > 0 {
		res_adid := adinfo_map_res[GenerateRangeNum(0, len(adinfo_map_res))]

		if adinfo, err := getAdInfo(res_adid); err == nil {

			fmt.Println("searchAd success ", adinfo)
			return adinfo, nil
		} else {
			fmt.Println("searchAd failed ")
			var err error
			err = errors.New("searchAd failed ")
			return adinfo, err
		}
	} else {
		return adinfo, err

	}

}

func getAdInfo(adid int) (adinfo adinfo, err error) {
	fmt.Println("get adinfo :", adid)
	if _, ok := adinfo_map[adid]; ok {
		return adinfo_map[adid], nil
	} else {
		var err error
		err = errors.New("adinfo null")
		return adinfo, err
	}
}
