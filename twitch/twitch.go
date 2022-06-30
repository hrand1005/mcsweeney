package twitch

import (
	"github.com/nicklaw5/helix"
)

func GetClipChannel(c helix.Clip) string {
	return "https://twitch.tv/" + c.BroadcasterName
}
