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
		No          int    `json:"no" bson:"no"`
		Sticky      int    `json:"sticky,omitempty" bson:"sticky,omitempty"`
		Closed      int    `json:"closed,omitempty" bson:"closed,omitempty"`
		Now         string `json:"now" bson:"now"`
		Name        string `json:"name" bson:"name"`
		Sub         string `json:"sub,omitempty" bson:"sub,omitempty"`
		Com         string `json:"com" bson:"com"`
		Filename    string `json:"filename" bson:"filename"`
		Ext         string `json:"ext" bson:"ext"`
		W           int    `json:"w" bson:"w"`
		H           int    `json:"h" bson:"h"`
		TnW         int    `json:"tn_w" bson:"tn_w"`
		TnH         int    `json:"tn_h" bson:"tn_h"`
		Tim         int64  `json:"tim" bson:"tim"`
		Time        int    `json:"time" bson:"time"`
		Md5         string `json:"md5" bson:"md5"`
		Fsize       int    `json:"fsize" bson:"fsize"`
		Resto       int    `json:"resto" bson:"resto"`
		Capcode     string `json:"capcode" bson:"capcode"`
		SemanticURL string `json:"semantic_url,omitempty" bson:"semantic_url,omitempty"`
		Replies     int    `json:"replies,omitempty" bson:"replies,omitempty"`
		Images      int    `json:"images,omitempty" bson:"images,omitempty"`
		UniqueIps   int    `json:"unique_ips,omitempty" bson:"unique_ips,omitempty"`
	} `json:"posts" bson:"posts"`
}
