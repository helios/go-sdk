package datautils

import (
	"mime"
	"strings"

	"golang.org/x/exp/slices"
)

var excludedTypes = []string{	"audio", "image", "multipart", "video" }
var excludedTextSubTypes = []string{	"css", "html", "javascript" }
var excludedApplicationSubTypes = []string{	"javascript" }

func ShouldSkipContentCollectionByContentType(contentType string) (bool, error) {
	if contentType == "" {
		return false, nil
	}

	mediaType, _, err := mime.ParseMediaType(contentType) 
	if err != nil {
		return true, err
	}
	
	mainType, subType, _ := strings.Cut(mediaType, "/")

	if slices.Contains(excludedTypes, mainType) {
		return true, nil;
	}

	if (mainType == "text" && (slices.Contains(excludedTextSubTypes, subType) || strings.HasPrefix(subType, "vnd"))) ||
		(mainType == "application" && slices.Contains(excludedApplicationSubTypes, subType)) {
		return true, nil;
	}

	return false, nil;
}
