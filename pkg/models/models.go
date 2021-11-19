package models

type Threads []struct {
	Page    int `json:"page"`
	Threads []struct {
		No           int `json:"no"`
		LastModified int `json:"last_modified"`
		Replies      int `json:"replies"`
	} `json:"threads"`
}

type Thread struct {
	Posts []struct {
		No          int    `json:"no"`
		Sticky      int    `json:"sticky,omitempty"`
		Closed      int    `json:"closed,omitempty"`
		Now         string `json:"now"`
		Name        string `json:"name"`
		Sub         string `json:"sub,omitempty"`
		Com         string `json:"com"`
		Filename    string `json:"filename"`
		Ext         string `json:"ext"`
		W           int    `json:"w"`
		H           int    `json:"h"`
		TnW         int    `json:"tn_w"`
		TnH         int    `json:"tn_h"`
		Tim         int64  `json:"tim"`
		Time        int    `json:"time"`
		Md5         string `json:"md5"`
		Fsize       int    `json:"fsize"`
		Resto       int    `json:"resto"`
		Capcode     string `json:"capcode"`
		SemanticURL string `json:"semantic_url,omitempty"`
		Replies     int    `json:"replies,omitempty"`
		Images      int    `json:"images,omitempty"`
		UniqueIps   int    `json:"unique_ips,omitempty"`
	} `json:"posts"`
}
