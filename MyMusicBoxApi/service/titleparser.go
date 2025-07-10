package service

import (
	"fmt"
	"regexp"
	"strings"
)

// V2 downloader using json

// testRun := ytdlp.New().
// 	DownloadArchive(archiveFileName).
// 	ForceIPv4().
// 	NoKeepVideo().
// 	SkipDownload().
// 	// ExtractAudio().
// 	// AudioQuality("0").
// 	// AudioFormat(fileExtension).
// 	// PostProcessorArgs("FFmpegExtractAudio:-b:a 160k").
// 	DownloadArchive(archiveFileName).
// 	// WriteThumbnail().
// 	// ConcurrentFragments(10).
// 	// ConvertThumbnails("jpg").
// 	// ForceIPv4().
// 	// Downloader("aria2c").
// 	// DownloaderArgs("aria2c:-x 16 -s 16 -j 16").
// 	// Output(storageFolderName+"/%(id)s.%(ext)s").
// 	ReplaceInMetadata("title", `["()]`, "").
// 	PrintToFile(`{"id": "%(id)s","title": "%(title)s", "playlist":"%(playlist_title)s", "duration": %(duration)s, "index":%(playlist_autonumber)s}`, idsFileName).
// 	Cookies("selenium/cookies_netscape")

// testRun.Run(context.Background(), downloadRequest.Url)

// rawJsons, _ := readLines(idsFileName)

// for index := range rawJsons {
// 	var jsonResult models.YtdlpJsonResult

// 	parseError := json.Unmarshal([]byte(FixJSONStringValues(rawJsons[index])), &jsonResult)

// 	if parseError != nil {
// 		logging.Error(parseError.Error())
// 		continue
// 	}

// 	logging.Info(fmt.Sprintf("%d %s %d %s", index, jsonResult.Title, jsonResult.PlaylistIndex, jsonResult.Playlist))
// }

// _dlp := ytdlp.New().
// 	DownloadArchive(archiveFileName).
// 	ForceIPv4().
// 	NoKeepVideo().
// 	SkipDownload().
// 	FlatPlaylist().
// 	WriteThumbnail().
// 	PrintToFile("%(id)s", idsFileName).
// 	PrintToFile("%(title)s", namesFileName).
// 	PrintToFile("%(duration)s", durationFileName).
// 	PrintToFile("%(playlist_title)s", playlistTitleFileName).
// 	PrintToFile("%(playlist_id)s", playlistIdFileName).
// 	Cookies("selenium/cookies_netscape")

// _dlp.Run(context.Background(), downloadRequest.Url)

// return

// FixJSONStringValues scans all JSON string values and escapes invalid characters.
func FixJSONStringValues(input string) string {
	// Regex pattern to match string values like: "key": "value"
	re := regexp.MustCompile(`"(?:[^"\\]|\\.)*?"\s*:\s*"(.*?)"`)

	return strings.TrimSpace(re.ReplaceAllStringFunc(input, func(match string) string {
		parts := strings.SplitN(match, ":", 2)
		if len(parts) != 2 {
			return match
		}
		key := strings.TrimSpace(parts[0])
		rawValue := strings.TrimSpace(parts[1])

		// Remove surrounding quotes
		unquoted := strings.Trim(rawValue, `"`)

		// Escape invalid characters inside the value
		escaped := escapeUnsafeJSONCharacters(unquoted)

		return fmt.Sprintf(`%s: "%s"`, key, escaped)
	}))
}

// escapeUnsafeJSONCharacters escapes characters that would break JSON string parsing.
func escapeUnsafeJSONCharacters(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 0x20 {
				// Escape control characters as unicode
				fmt.Fprintf(&b, `\u%04x`, r)
			} else {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}
