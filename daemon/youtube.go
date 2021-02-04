package daemon

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/6f7262/pipe"
	"github.com/google/google-api-go-client/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
	"gopkg.in/mgo.v2/bson"
	"io.6f.sana/db"
	uppdb "upper.io/db"
)

type Youtube struct {
	sync.Mutex
	p  pipe.Pipe
	ep pipe.Pipe
	m  map[string]*sync.Mutex
	s  *youtube.Service
	c  map[string]string
	d  map[string]struct{}
}

func NewYoutube() *Youtube {
	s, err := youtube.New(&http.Client{
		Transport: &transport.APIKey{
			Key: "AIzaSyBAdDIgUc_loht-bJyBtaRcD8aDeupAaeE",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return &Youtube{
		p:  pipe.New(5),
		ep: pipe.New(20),
		m:  make(map[string]*sync.Mutex),
		s:  s,
		c: map[string]string{
			"Affinity":            "UUOJEt1dxJYH1N6jnLfB6SeA",
			"AngelicBunny":        "UUu5NDvqE639JOxpqlK7UKUA",
			"Artzie Music":        "UU6hBefyLMtG7FXhZ55da3Vw",
			"AvienCloud":          "UUKioNqOX_kOCLcSIWPL_lxQ",
			"CloudKid":            "UUSa8IUd1uEjlREMa21I3ZPQ",
			"DiscoFactory":        "UU_xBnI4WIyYyKw6yJjHsCdA",
			"DiscoThrill Records": "UU3e96kvQtQD9IT0pAqMGj5w",
			"Diversity":           "UU7tD6Ifrwbiy-BoaAHEinmQ",
			"EDMSoundEater":       "UUCfBdlrFOv3_PBc5SFUE4-Q",
			"Electronic Gems":     "UUPzWlhG7QM56Y8MYB3qMVnQ",
			"ElFamosoDemon":       "UURUOfuNIb_sk__7snjK3aVg",
			"ENM":                 "PLcghjaDFUHgkFttlesFQ06nhgoXQkS1SI",
			"ENM RE":              "UUkwpDBMwfr6ETIm66sPxZ-w",
			"EvolveMusicNetwork":  "UU87ARwwcHUbd1WPaz11zw4g",
			"FunkyPanda":          "UUUHhoftNnYfmFp1jvSavB-Q",
			"Future Classic":      "UUy3DbVl0K1qj0e8jGAOYgPg",
			"FutureHype":          "PLDXNpKa6EPd1hElb15HOiHUBj-Pek3lg1",
			"GalaxyMusic":         "UUIKF1msqN7lW9gplsifOPkQ",
			"GalaxyMusicNet":      "UUbfMTDftQhRySQKFCV8PfHg",
			"HyperboltEDM":        "UUF6-fdubXkcpELOB5sUbrWw",
			"Kyra":                "UUqolymr8zonJzC08v2wXNrQ",
			"Liquicity":           "UUSXm6c-n6lsjtyjvdD0bFVw",
			"MA Dance":            "UUF8QEPMImbbJD_s0IHCBNRw",
			"MA Lite":             "UUlFjOLis3eUKAGkCKCwdYxQ",
			"Majestic Casual":     "UUXIyz409s7bNWVcM-vjfdVA",
			"majesticdnb":         "UUd3qPLVGUGKfWeuUjC0OJWQ",
			"MikuMusic":           "UUjw0SX_9OGFI7buLG4-M0Fw",
			"Monstercat":          "UUJ6td3C9QlPO9O_J5dF4ZzA",
			"MOR Network":         "UUkfMJApxxdy-h41xy_8AHNw",
			"MORindie":            "UUcTvjjmFeFDd5Ri5NajGImA",
			"MOΣDM":               "UU_dBZzj82blh0hp_w6nxrDA",
			"MrSuicideSheep":      "UU5nc_ZtjKW1htCVZVRxlQAQ",
			"NeurofunkGrid":       "UU1hPZ15rCXacep76EgDfwyw",
			"NoCopyrightSounds":   "UU_aEa8K-EOJ3D6gOs7HcyNg",
			"OneChilledPanda":     "UUkUTBwZKwA9ojYqzj6VRlMQ",
			"Paradoxium":          "PLeeznYF7ppjfH3JetE0AocdA0dI3r3fFg",
			"Pixl Networks":       "UU1iqebKNH36JIdBIjEy8-iQ",
			"Real ℒℴѵℯ ❤":         "UU8EvY5Cky6LM9tTE07kk2kQ",
			"Selected.":           "PLSr_oFUba1jtP9x5ZFs5Y0GJkb8fmC161",
			"Self":                "LLahIftu-BClclAu0uPLewGQ",
			"Strobe Music":        "UUcoYD5HDg8P-gvJU-oDYq5Q",
			"Synergy Music":       "UUrbRMQk4-CnFZLmLuz0KyHQ",
			"SyrebralVibes":       "UUi0LydWaEUy3Vx8flL29ebQ",
			"TriangleMusic":       "UUDBdeEaSnlu-AU-ITBTRkeQ",
			"TriDanceMusic":       "UU1qAm032wFDxwR2a3WSeV0w",
			"Waifu Wednesdays":    "PU-wNjNTqCfXSKd4S1tNgWUg",
			"Welcome Home Music":  "UUGBI6hbiJWpRs00cANGlx2w",
			"xKito Music":         "UUMOgdURr7d8pOVlc-alkfRg",
		},
		d: make(map[string]struct{}),
	}
}

func (y *Youtube) Start() {
	for c, l := range y.c {
		go y.loop(c, l)
	}
}

func (y *Youtube) lock(id string) {
	y.Lock()
	m, ok := y.m[id]
	if !ok {
		m = &sync.Mutex{}
		y.m[id] = m
	}
	y.Unlock()
	m.Lock()
}

func (y *Youtube) unlock(id string) {
	y.Lock()
	m, ok := y.m[id]
	y.Unlock()
	if ok {
		m.Unlock()
	}
}

func (y *Youtube) doesExist(id string) {
	y.Lock()
	defer y.Unlock()
	y.d[id] = struct{}{}
}

func (y *Youtube) isExist(id string) bool {
	y.Lock()
	_, ok := y.d[id]
	y.Unlock()
	if ok {
		return ok
	}
	y.ep.Increment()
	defer y.ep.Decrement()
	_, err := db.GetMedia(uppdb.Cond{"originalId": id})
	return err != uppdb.ErrNoMoreRows
}

func (y *Youtube) loop(channel, list string) {
	var token string
	for {
		y.p.Increment()
		res, err := y.s.PlaylistItems.List("snippet,contentDetails").
			PlaylistId(list).
			MaxResults(50).
			PageToken(token).
			Do()
		y.p.Decrement()
		if err != nil {
			log.Printf("Could not fetch playlist %s for channel %s\n", list, channel)
			break
		}
		for _, item := range res.Items {
			item := item
			go func() {
				id := item.Snippet.ResourceId.VideoId
				y.lock(id)
				defer y.unlock(id)
				if y.isExist(id) {
					return
				}
				if err := y.download(id); err != nil {
					log.Printf("Couldn't download video %s\n%s\n", id, err)
					return
				}
				y.doesExist(id)
			}()
		}
		token = res.NextPageToken
		if token == "" {
			break
		}
	}
	time.Sleep(15 * time.Minute)
	go y.loop(channel, list)
}

func (y *Youtube) download(id string) error {
	defer y.p.One()()
	fmt.Printf("Downloading video %s\n", id)
	defer fmt.Printf("Finished downloading video %s\n", id)
	res, err := y.s.Videos.List("snippet,contentDetails").
		Id(id).
		Do()
	if err != nil {
		return err
	}
	if len(res.Items) == 0 {
		return errors.New("youtube returned no videos")
	}
	v := res.Items[0]
	if v.LiveStreamingDetails != nil {
		return errors.New("video is currently streaming")
	}
	u, err := time.Parse(time.RFC3339Nano, v.Snippet.PublishedAt)
	if err != nil {
		return err
	}
	d, err := time.ParseDuration(strings.ToLower(v.ContentDetails.Duration[2:]))
	if err != nil {
		return err
	}
	m := &db.Media{
		Id:          bson.NewObjectId(),
		OriginalId:  id,
		Title:       v.Snippet.Title,
		Description: v.Snippet.Description,
		Duration:    d,
		Uploaded:    u,
		Provider:    v.Snippet.ChannelTitle,
		ProviderId:  v.Snippet.ChannelId,
		Etag:        v.Etag,
		Source:      "youtube",
	}
	o, err := exec.Command("youtube-dl", "https://youtu.be/"+id, "-f",
		"bestvideo[ext=mp4]+bestaudio[ext=m4a]", "--output",
		filepath.Join("_", "media", fmt.Sprintf("%x.%%(ext)s", string(m.Id)))).CombinedOutput()
	if err != nil {
		return errors.New(strings.TrimSpace(string(o)))
	}
	ts := v.Snippet.Thumbnails
	var t *youtube.Thumbnail
	switch {
	case ts.Maxres != nil:
		t = ts.Maxres
	case ts.High != nil:
		t = ts.High
	case ts.Medium != nil:
		t = ts.Medium
	case ts.Standard != nil:
		t = ts.Standard
	case ts.Default != nil:
		t = ts.Default
	default:
		return errors.New("video does not contain a thumbnail")
	}
	resp, err := http.Get(t.Url)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	f, err := os.Create(filepath.Join("_", "thumbnail", fmt.Sprintf("%x.jpg", string(m.Id))))
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}
	return m.Save()
}
