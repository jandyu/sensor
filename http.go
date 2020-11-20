package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func SetResponseHead(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Expires", "-1")
}

func writeConfigFile(s string) error {
	fileObj, err := os.OpenFile("./data/config.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Failed to open the file", err.Error())
		return err
	}
	if _, err := io.WriteString(fileObj, s); err != nil {
		fmt.Println("error write log", err.Error())
		return err
	}

	fileObj.Close()

	fileObj, err = os.OpenFile("./data/config"+fmt.Sprint(time.Now().Unix())+".json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Failed to open the file", err.Error())
		return err
	}
	if _, err := io.WriteString(fileObj, s); err != nil {
		fmt.Println("error write log", err.Error())
		return err
	}

	fileObj.Close()
	return nil
}

func readConfigFile() (string, error) {

	data, err := ioutil.ReadFile("./data/config.json")
	if err != nil {
		fmt.Println("readConfigFile", err)
		return "", err
	}

	//读取的数据为json格式，需要进行解码
	var f map[string]string
	err = json.Unmarshal(data, &f)
	if err != nil {
		fmt.Println("readConfigFile", err)
		return "", err
	}
	//fmt.Println("readConfigFile",string(data))
	return string(data), nil
}

func FormatPostJsonData(r *http.Request) (map[string]interface{}, error) {
	if r.Method == "POST" {
		//取数据
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		var f interface{}
		err := json.Unmarshal(result, &f)
		if err != nil {
			fmt.Errorf("参数格式错误", err)
			return nil, err
		}

		params := f.(map[string]interface{})
		// 输出日志
		for key, val := range params {
			fmt.Println("params", key, val)
		}
		//写配置
		err = writeConfigFile(string(result))
		//返回
		return params, err
	} else {
		return nil, errors.New("无效GET请求")
	}
}
func SetConfigHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHead(w)
	fmt.Println(r.URL)
	_, err := FormatPostJsonData(r)
	if err != nil {
		return
	}
	//

	if err != nil {
		fmt.Fprint(w, string("{\"message\":\"保存配置失败,"+err.Error()+"\"}"))
	} else {
		fmt.Fprint(w, string("{\"message\":\"保存成功，下一个显示周期（默认1分钟）对应配置生效\"}"))
	}
}

func getConfigHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHead(w)
	fmt.Println(r.URL)
	dat, err := readConfigFile()
	if err != nil {
		fmt.Fprint(w, string("{\"message\":\"读取配置失败,"+err.Error()+"\"}"))
	} else {
		//
		//修改监控数据
		for i := 1; i < 10; i++ {
			prefix := "LED" + fmt.Sprint(i) + "_tmp_m"
			dat = strings.Replace(dat, "#"+prefix+"#", getMapData(prefix), 1)
			prefix = "LED" + fmt.Sprint(i) + "_sal_m"
			dat = strings.Replace(dat, "#"+prefix+"#", getMapData(prefix), 1)
		}
		fmt.Fprint(w, dat)
	}
}

func getMapData(key string) string {
	v, ok := MDATA.Load(key)
	if ok {
		return v.(string)
	} else {
		return ""
	}
}

func dataConfigHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHead(w)
	fmt.Println(r.URL)
	dat, err := readConfigFile()

	var f map[string]string
	err = json.Unmarshal([]byte(dat), &f)
	if err != nil {
		fmt.Println("readConfigFile", err)
		fmt.Fprint(w, string("[]"))
	}

	prefix := "LED"
	dest := make([]string, 0)
	for i := 1; i < 10; i++ {
		prefix = "LED" + fmt.Sprint(i)
		fmt.Println(prefix)
		if _, ok := f[prefix+"_title"]; ok {
			//修改监控数据
			f[prefix+"_tmp_m"] = getMapData(prefix + "_tmp_m")
			f[prefix+"_sal_m"] = getMapData(prefix + "_sal_m")
			strArr := []string{
				`{"LEDID": "` + prefix + `","DATA": "` + f[prefix+"_title"] + `"}`,
				`{"LEDID": "` + prefix + `","DATA": "` + f[prefix+"_area"] + `"}`,
				`{"LEDID": "` + prefix + `","DATA": "` + f[prefix+"_temp"] + `"}`,
				`{"LEDID": "` + prefix + `","DATA": "` + f[prefix+"_tmp_stand"] + `"}`,
				`{"LEDID": "` + prefix + `","DATA": "` + f[prefix+"_tmp_m"] + `"}`,
				`{"LEDID": "` + prefix + `","DATA": "` + f[prefix+"_salinity"] + `"}`,
				`{"LEDID": "` + prefix + `","DATA": "` + f[prefix+"_sal_stand"] + `"}`,
				`{"LEDID": "` + prefix + `","DATA": "` + f[prefix+"_sal_m"] + `"}`,
				`{"LEDID": "` + prefix + `","DATA": "` + f[prefix+"_tm1"] + `"}`,
				`{"LEDID": "` + prefix + `","DATA": "` + f[prefix+"_store1"] + `"}`,
				//`{"LEDID": "`+prefix+`","DATA": "`+f[prefix+"_tm2"]+`"}`,
				//`{"LEDID": "`+prefix+`","DATA": "`+f[prefix+"_store2"]+`"}`,
				//`{"LEDID": "`+prefix+`","DATA": "`+f[prefix+"_tm3"]+`"}`,
				//`{"LEDID": "`+prefix+`","DATA": "`+f[prefix+"_store3"]+`"}`,
			}
			if f[prefix+"_tm2"] != "" {
				strArr = append(strArr, `{"LEDID": "`+prefix+`","DATA": "`+f[prefix+"_tm2"]+`"}`, `{"LEDID": "`+prefix+`","DATA": "`+f[prefix+"_store2"]+`"}`)
			}
			if f[prefix+"_tm3"] != "" {
				strArr = append(strArr, `{"LEDID": "`+prefix+`","DATA": "`+f[prefix+"_tm3"]+`"}`, `{"LEDID": "`+prefix+`","DATA": "`+f[prefix+"_store3"]+`"}`)
			}

			dest = append(dest, strArr...)
		} else {
			break
		}
	}

	fmt.Fprint(w, "["+strings.Join(dest, ",")+"]")

}

func startHttp(port int) {
	//log.Info("http listen ", *addr)
	sp := ":" + fmt.Sprint(port)
	fmt.Println("http listen ", sp)

	http.HandleFunc("/monitor/set", SetConfigHandler)
	http.HandleFunc("/monitor/get", getConfigHandler)
	http.HandleFunc("/monitor/data", dataConfigHandler)

	//static dir
	http.Handle("/monitor/config/", http.StripPrefix("/monitor/config/", http.FileServer(http.Dir("./data"))))

	err := http.ListenAndServe(sp, nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
