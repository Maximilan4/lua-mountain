package nexus

import "time"

const (
	RawRepositoryFormat = "raw"
)

type (
	Repository struct {
		Name       string         `json:"name"`
		Format     string         `json:"format"`
		Type       string         `json:"type"`
		Url        string         `json:"url"`
		Attributes map[string]any `json:"attributes"`
	}
	AssetChecksum struct {
		Sha1 string `json:"sha1"`
		Md5  string `json:"md5"`
	}

	Asset struct {
		DownloadUrl    string        `json:"downloadUrl"`
		Path           string        `json:"path"`
		Id             string        `json:"id"`
		Repository     string        `json:"repository"`
		Format         string        `json:"format"`
		Checksum       AssetChecksum `json:"checksum"`
		ContentType    string        `json:"contentType"`
		LastModified   time.Time     `json:"lastModified"`
		LastDownloaded *time.Time    `json:"lastDownloaded"`
		Uploader       string        `json:"uploader"`
		UploaderIp     string        `json:"uploaderIp"`
		FileSize       int           `json:"fileSize"`
		BlobCreated    time.Time     `json:"blobCreated"`
	}

	AssetList struct {
		Items []Asset `json:"items"`
		Next  string  `json:"continuationToken"`
	}
)
