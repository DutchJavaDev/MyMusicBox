import requests
import json
import time

# === CONFIGURATION ===
SPOTIFY_ACCESS_TOKEN = ""
YOUTUBE_ACCESS_TOKEN = ""
SPOTIFY_PLAYLIST_ID = ""  # Only the ID, not full URL
MAX_SONGS = 100  # Adjust as needed

SPOTIFY_API_BASE = "https://api.spotify.com/v1"
YOUTUBE_API_BASE = "https://www.googleapis.com/youtube/v3"

# === STEP 1: Get Songs from a Spotify Playlist ===
def get_spotify_playlist_songs(playlist_id, max_songs=None):
    headers = {"Authorization": f"Bearer {SPOTIFY_ACCESS_TOKEN}"}
    url = f"{SPOTIFY_API_BASE}/playlists/{playlist_id}/tracks?limit=100"
    songs = []

    while url:
        response = requests.get(url, headers=headers)
        if response.status_code != 200:
            print(f"Failed to fetch playlist songs: {response.status_code} — {response.text}")
            break

        data = response.json()
        for item in data.get('items', []):
            track = item.get('track')
            if track:
                name = track['name']
                artists = ', '.join(artist['name'] for artist in track['artists'])
                query = f"{artists} {name}"
                songs.append(query)

                # Optional cap
                if max_songs and len(songs) >= max_songs:
                    return songs

        url = data.get('next')  # Spotify handles pagination via `next`
        time.sleep(0.1)  # Throttle to be API-friendly

    return songs

# === STEP 2: Search YouTube Video ===
def search_youtube_video(query):
    headers = {"Authorization": f"Bearer {YOUTUBE_ACCESS_TOKEN}"}
    params = {
        'part': 'snippet',
        'q': query,
        'maxResults': 1,
        'type': 'video'
    }

    response = requests.get(f"{YOUTUBE_API_BASE}/search", headers=headers, params=params)
    if response.status_code != 200:
        print(f"Failed YouTube search: {response.status_code} — {query}")
        return None

    items = response.json().get('items', [])
    if items:
        return items[0]['id']['videoId']
    return None

# === STEP 3: Create YouTube Playlist ===
def create_youtube_playlist(title, description="Imported from Spotify Playlist"):
    headers = {
        "Authorization": f"Bearer {YOUTUBE_ACCESS_TOKEN}",
        "Content-Type": "application/json"
    }
    payload = {
        "snippet": {
            "title": title,
            "description": description
        },
        "status": {
            "privacyStatus": "private"
        }
    }

    response = requests.post(
        f"{YOUTUBE_API_BASE}/playlists?part=snippet,status",
        headers=headers,
        data=json.dumps(payload)
    )

    if response.status_code != 200:
        print(f"Failed to create playlist: {response.status_code} — {response.text}")
        return None

    return response.json().get("id")

# === STEP 4: Add Video to YouTube Playlist ===
def add_video_to_playlist(video_id, playlist_id):
    headers = {
        "Authorization": f"Bearer {YOUTUBE_ACCESS_TOKEN}",
        "Content-Type": "application/json"
    }
    payload = {
        "snippet": {
            "playlistId": playlist_id,
            "resourceId": {
                "kind": "youtube#video",
                "videoId": video_id
            }
        }
    }

    response = requests.post(
        f"{YOUTUBE_API_BASE}/playlistItems?part=snippet",
        headers=headers,
        data=json.dumps(payload)
    )

    if response.status_code not in (200, 201):
        print(f"Failed to add video {video_id} to playlist: {response.status_code}")
        return False

    return True

# === MAIN EXECUTION ===
def main():
    print("▶ Fetching songs from Spotify playlist...")
    songs = get_spotify_playlist_songs(SPOTIFY_PLAYLIST_ID)
    print(f"🎵 Found {len(songs)} songs in playlist.")
    with open("playlist.txt", "a") as f:
        for line in songs:
            f.write(f"{line}\n")
    # print("📺 Creating YouTube playlist...")
    # playlist_title = "Spotify Playlist to YouTube"
    # playlist_id = create_youtube_playlist(playlist_title)
    # if not playlist_id:
    #     print("❌ Could not create YouTube playlist.")
    #     return

    print("🔁 Searching on YouTube and adding to playlist...")
    for idx, song in enumerate(songs, 1):
        print(f"[{idx}/{len(songs)}] Searching: {song}")
        # video_id = search_youtube_video(song)
        # if video_id:
        #     if add_video_to_playlist(video_id, playlist_id):
        #         print("  ✅ Added successfully")
        #     else:
        #         print("  ❌ Failed to add")
        # else:
        #     print("  ❌ No video found")
        # time.sleep(1.1)  # Be gentle with YouTube's API

    print("✅ Done!")

if __name__ == "__main__":
    main()

