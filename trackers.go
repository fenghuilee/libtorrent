package libtorrent

import (
	"time"
)

type Tracker struct {
	// Tracker URI or DHT, LSD, PE
	Addr         string
	Error        string
	LastAnnounce int64
	NextAnnounce int64
	Peers        int

	// scrape info
	LastScrape int64
	Seeders    int
	Leechers   int
	Downloaded int
}

func TorrentTrackersCount(i int) int {
	mu.Lock()
	defer mu.Unlock()

	t := torrents[i]
	fs := filestorage[t.InfoHash()]
	fs.Trackers = nil
	for _, v := range t.Trackers() {
		e := ""
		if v.Err != nil {
			e = v.Err.Error()
		}
		fs.Trackers = append(fs.Trackers, Tracker{v.Url,
			e,
			(time.Duration(v.LastAnnounce) * time.Second).Nanoseconds(),
			(time.Duration(v.NextAnnounce) * time.Second).Nanoseconds(),
			v.Peers,
			0, 0, 0, 0})
	}
	fs.Trackers = append(fs.Trackers, Tracker{"LPD",
		"",
		0,
		0,
		lpdCount(t.InfoHash()),
		0, 0, 0, 0})
	return len(fs.Trackers)
}

func TorrentTrackers(i int, p int) *Tracker {
	mu.Lock()
	defer mu.Unlock()

	t := torrents[i]
	f := filestorage[t.InfoHash()]
	return &f.Trackers[p]
}

func TorrentTrackerRemove(i int, url string) {
	mu.Lock()
	defer mu.Unlock()

	t := torrents[i]
	t.RemoveTracker(url)
}

func TorrentTrackerAdd(i int, addr string) {
	mu.Lock()
	defer mu.Unlock()

	t := torrents[i]
	t.AddTrackers([][]string{[]string{addr}})
}
