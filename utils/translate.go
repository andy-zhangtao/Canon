/*
 * Andy ZhangTao (ztao8607@gmail.com)
 *
 * MIT License
 *
 * Copyright (c) 2017 The Po.et Foundation
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package utils

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"log"
)

// 根据不同视频源的编码规则,获取相对应的播放地址



// TouTiao 头条规则
// GET请求此url,会返回TouTiao对应的数据结构体。其中main_url 是经过base64编码的地址。 经过解码之后可以得到播放地址
func TouTiao(url string) ([]TouTiaoResult, error){

	var tts TouTiaoStruct
	var ttr []TouTiaoResult

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ttr, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return ttr, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ttr, err
	}
	
	err = json.Unmarshal(content, &tts)
	if err != nil {
		return ttr, err
	}

	for _, t := range tts.Data.ViedoList{

		data, err := base64.StdEncoding.DecodeString(t.Main)
		if err != nil{
			data, err = base64.StdEncoding.DecodeString(t.Backup)
			if err != nil{
				continue
			}
			ttr = append(ttr, TouTiaoResult{
				Vtype:t.Vtype,
				VURL:"backup_url",
				BackupURL:string(data),
			})
			continue
		}
		ttr = append(ttr, TouTiaoResult{
			Vtype:t.Vtype,
			VURL:"main_url",
			MainURL:string(data),
		})
	}

	return ttr, nil
}
