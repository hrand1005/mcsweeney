package main

import ()

type Clip struct {
	Author      string
	Broadcaster string
	MP4URL      string
	Title       string
	URL         string
}

type ClipService interface {
	Get(int) []*Clip
}
