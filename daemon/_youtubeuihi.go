package daemon

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/google/google-api-go-client/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
	"io.6f.sana/db"
	uppdb "upper.io/db"
)

var (
	youtubeVideoMutexes = map[string]*sync.Mutex{}
	youtubeService      *youtube.Service
	youtubeChannels     = map[string]string{
		"Self":           "LLahIftu-BClclAu0uPLewGQ",
		"FunkyPanda":     "UUUHhoftNnYfmFp1jvSavB-Q",
		"MrSuicideSheep": "UU5nc_ZtjKW1htCVZVRxlQAQ",
		"GalaxyMusic":    "UUIKF1msqN7lW9gplsifOPkQ",
		"Liquicity":      "UUSXm6c-n6lsjtyjvdD0bFVw",
		"TriangleMusic":  "UUDBdeEaSnlu-AU-ITBTRkeQ",
		"xKito Music":    "UUMOgdURr7d8pOVlc-alkfRg",
		"Monstercat":     "UUJ6td3C9QlPO9O_J5dF4ZzA",
		"Strobe Music":   "UUcoYD5HDg8P-gvJU-oDYq5Q",
		"TriDanceMusic":  "UU1qAm032wFDxwR2a3WSeV0w",
		"AvienCloud":     "UUKioNqOX_kOCLcSIWPL_lxQ",
		"ENM":            "PLcghjaDFUHgkFttlesFQ06nhgoXQkS1SI",
		"ENM RE":         "UUkwpDBMwfr6ETIm66sPxZ-w",
		"HyperboltEDM":   "UUF6-fdubXkcpELOB5sUbrWw",
		"EDMSoundEater":  "UUCfBdlrFOv3_PBc5SFUE4-Q",
		"MA Dance":       "UUF8QEPMImbbJD_s0IHCBNRw",
		// "Saika Music Network": "UUO94BOMyhSccGPDm3u3GJdQ",
		"Pixl Networks": "UU1iqebKNH36JIdBIjEy8-iQ",
	}
)

func init() {
	client := &http.Client{
		Transport: &transport.APIKey{Key: "AIzaSyBAdDIgUc_loht-bJyBtaRcD8aDeupAaeE"},
	}

	service, err := youtube.New(client)
	if err != nil {
		panic(err)
	}
	youtubeService = service
}

func launchYoutube() {
	for channelName, playlistId := range youtubeChannels {
		go func(channelName, playlistId string) {
			for {
				nextPageToken := ""
				for {
					playlistCall := youtubeService.PlaylistItems.List("snippet,contentDetails").
						PlaylistId(playlistId).
						MaxResults(50).
						PageToken(nextPageToken)

					playlistResponse, err := playlistCall.Do()
					if err != nil {
						fmt.Printf("[FATAL] Could not fetch playlist %s for channel %s\n",
							playlistId, channelName)
						continue
					}

					for _, playlistItem := range playlistResponse.Items {
						id := playlistItem.Snippet.ResourceId.VideoId

						lockVideo(id)
						_, err := db.GetMedia(uppdb.Cond{"originalId": id})
						if err != uppdb.ErrNoMoreRows {
							unlockVideo(id)
							continue
						}

						if err := downloadYoutubeVideo(playlistItem.Snippet.ResourceId.VideoId); err != nil {
							fmt.Printf("[%s] [FATAL] %s : %s\n",
								time.Now().Format(time.Kitchen), id, err.Error())
						}
						unlockVideo(id)
					}

					nextPageToken = playlistResponse.NextPageToken
					if nextPageToken == "" {
						break
					}
				}
				time.Sleep(5 * time.Minute)
			}
		}(channelName, playlistId)

	}
}

func lockVideo(id string) {
	if _, ok := youtubeVideoMutexes[id]; !ok {
		youtubeVideoMutexes[id] = &sync.Mutex{}
	}
	youtubeVideoMutexes[id].Lock()
}

func unlockVideo(id string) {
	youtubeVideoMutexes[id].Unlock()
}

func downloadYoutubeVideo(id string) error {
	fmt.Printf("[%s] Downloading video %s\n", time.Now().Format(time.Kitchen), id)
	defer func() {
		fmt.Printf("[%s] Finished downloading %s\n", time.Now().Format(time.Kitchen), id)
	}()

	videoCall := youtubeService.Videos.List("snippet,contentDetails").
		Id(id)

	videoResponse, err := videoCall.Do()
	if err != nil {
		return err
	}

	if len(videoResponse.Items) < 1 {
		return errors.New("youtube returned 0 videos")
	}
	video := videoResponse.Items[0]

	uploaded, err := time.Parse(time.RFC3339Nano, video.Snippet.PublishedAt)
	if err != nil {
		return err
	}
	duration, err := time.ParseDuration(strings.ToLower(video.ContentDetails.Duration[2:]))
	if err != nil {
		return err
	}
	m := db.NewMedia()
	m.OriginalId = video.Id
	m.Title = video.Snippet.Title
	m.Description = video.Snippet.Description
	m.Duration = duration
	m.Uploaded = uploaded
	m.Provider = video.Snippet.ChannelTitle
	m.ProviderId = video.Snippet.ChannelId
	m.Etag = video.Etag
	m.Source = "youtube"

	// Thumbnail
	thumbnails := video.Snippet.Thumbnails
	var thumbnail *youtube.Thumbnail
	switch {
	case thumbnails.Maxres != nil:
		thumbnail = thumbnails.Maxres
	case thumbnails.High != nil:
		thumbnail = thumbnails.High
	case thumbnails.Medium != nil:
		thumbnail = thumbnails.Medium
	case thumbnails.Standard != nil:
		thumbnail = thumbnails.Standard
	case thumbnails.Default != nil:
		thumbnail = thumbnails.Default
	default:
		return errors.New("Video does not contain a thumbnail")
	}

	res, err := http.Get(thumbnail.Url)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}
	defer res.Body.Close()

	thumbnailData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Media
	// FFMPEG
	output, err := exec.Command("youtube-dl", "https://youtu.be/"+id, "-f",
		"bestvideo[ext=mp4]+bestaudio[ext=m4a]", "--output",
		"_/media/"+fmt.Sprintf("%x", string(m.Id))+".%(ext)s").CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	if err := ioutil.WriteFile(fmt.Sprintf("_/thumbnail/%s.jpg", fmt.Sprintf("%x", string(m.Id))), thumbnailData, 0644); err != nil {
		return err
	}

	return m.Save()
}
