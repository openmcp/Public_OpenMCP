package main

type jsonErr struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
	Text   string `json:"text"`
}
