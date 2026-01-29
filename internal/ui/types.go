package ui

const (
	focusSearch = iota
	focusSidebar
	focusMain
	focusSong
	focusPlaylist = 90
)

const (
	viewList = iota
	viewQueue
	viewLogin = 99
)

const (
	filterSongs = iota
	filterAlbums
	filterArtist
)

const (
	displaySongs = iota
	displayAlbums
	displayArtist
)

const (
	LoopNone = 0
	LoopAll  = 1
	LoopOne  = 2
)
