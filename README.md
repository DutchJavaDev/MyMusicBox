# MyMusicBox

This project aims to create not just a music player but a universal music search that enables you to **import** music from different sources and play them using different platforms like mobile, web, desktop etc.

[![Auto Bump Version](https://github.com/DutchJavaDev/MyMusicBox/actions/workflows/frontend.yml/badge.svg)](https://github.com/DutchJavaDev/MyMusicBox/actions/workflows/frontend.yml)

[![Api Tests](https://github.com/DutchJavaDev/MyMusicBox/actions/workflows/backend.yml/badge.svg)](https://github.com/DutchJavaDev/MyMusicBox/actions/workflows/backend.yml)

### Database schema
<img width="546" height="467" alt="image" src="https://github.com/user-attachments/assets/b1f3813b-997b-4423-84d8-886970106500" />

# MusicBoxApi API
## Version: 1.0


### /api/v1/download:sourceId

#### GET
##### Description:

Enables playback for song/file using http 206 partial content

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| download | body | message/rfc8259 see models.DownloadRequestModel | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 0 |  | [models.DownloadRequestModel](#models.DownloadRequestModel) |
| 200 | serve song/file with range request (http 206) |  |
| 500 | Internal Server Error | [models.ApiResponseModel](#models.ApiResponseModel) |

### /api/v1/play/:sourceId

#### GET
##### Description:

Enables playback for song/file using http 206 partial content

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| sourceId | path | Id of song/file to serve using http 206 partial content | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | serve song/file with range request (http 206) |  |
| 500 | Internal Server Error | [models.ApiResponseModel](#models.ApiResponseModel) |

### /api/v1/playlist

#### GET
##### Description:

Returns data for all playlist, if lastKnowPlaylistId then only the playlist after lastKnowPlaylistId

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| lastKnowPlaylistId | path | Last know playlist id by the client, default is 0 | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Playlist](#models.Playlist) |
| 500 | Internal Server Error | [models.ApiResponseModel](#models.ApiResponseModel) |

### /api/v1/playlist/:playlistId

#### GET
##### Description:

Returns data for a playlist, if lastKnowSongPosition then only songs added after lastKnowSongPosition

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| playlistId | path | Id of playlist | Yes | integer |
| lastKnowSongPosition | path | Last song that is know by the client, pass this in to only get the latest songs | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Song](#models.Song) |
| 500 | Internal Server Error | [models.ApiResponseModel](#models.ApiResponseModel) |

### /api/v1/songs

#### GET
##### Description:

Returns data for all songs

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Song](#models.Song) |
| 500 | Internal Server Error | [models.ApiResponseModel](#models.ApiResponseModel) |

### Models


#### models.ApiResponseModel

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| data |  |  | No |
| message | string |  | No |

#### models.DownloadRequestModel

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| url | string |  | No |

#### models.Playlist

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| creationDate | string |  | No |
| description | string |  | No |
| id | integer |  | No |
| isPublic | boolean |  | No |
| name | string |  | No |
| thumbnailPath | string |  | No |
| updatedAt | string |  | No |

#### models.Song

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created_at | string |  | No |
| duration | integer |  | No |
| id | integer |  | No |
| name | string |  | No |
| path | string |  | No |
| source_id | string |  | No |
| thumbnail_path | string |  | No |
| updated_at | string |  | No |
