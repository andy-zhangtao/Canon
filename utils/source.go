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


type TouTiaoResult struct {
	Vtype string `json:"vtype"`
	VURL string `json:"vurl"`
	MainURL string `json:"main_url"`
	BackupURL string `json:"backup_url"`
}

type TouTiaoStruct struct{
	// Total 用来标记返回的URL总数
	Total int `json:"total"`
	Data TouTiaoS_Data `json:"data"`
}

// TouTiaoS_Data 保存缩略图和播放地址
type TouTiaoS_Data struct{
	// PosterURL 缩略图地址
	PosterURL string `json:"poster_url"`
	ViedoList map[string]TouTiaoS_List `json:"video_list"`
}

// TouTiaoS_List 保存视频格式和地址
type TouTiaoS_List struct{
	// Vtype 视频格式
	Vtype string `json:"vtype"`
	// Main  主播放地址
	Main string `json:"main_url"`
	// Backup 备用播放地址
	Backup string `json:"backup_url_1"`
}


type UCResult struct{
	Vtype string `json:"vtype"`
	VURL string `json:"vurl"`
	MainURL string `json:"main_url"`
	BackupURL string `json:"backup_url"`
}

type UCStruct struct{
	Message string `json:"message"`
	Data UC_Data `json:"data"`
}

type UC_Data struct{
	VideoList []UC_List `json:"videoList"`
}

type UC_List struct{
	Format string `json:"format"`
	Fragment []UC_Fragment `json:"fragment"`
}

type UC_Fragment struct{
	Url string `json:"url"`
}