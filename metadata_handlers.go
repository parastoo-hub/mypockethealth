package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/suyashkumar/dicom"
)

// We accept two formats. Either a simple integer, which presents group and element tag,
// a more standard format of (group, element), where the numbers are in hex.
var hexCoordPattern = regexp.MustCompile(`^\([0-9A-Fa-f]{4},[0-9A-Fa-f]{4}\)$`)
var intPattern = regexp.MustCompile(`^\d+$`)

// GetDicomMetadata handles the request to retrieve DICOM metadata.
func GetDicomMetadata(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	tag := r.URL.Query().Get("tag")

	dataset, err := dicom.ParseFile(serverStorageDir+fileName, nil)
	if err != nil {
		http.Error(w, "Unable to parse DICOM file.", http.StatusBadRequest)
		return
	}

	// If the query doesn't specify any tag, return all tag values.
	// This is for debug purposes.
	if tag == "" {
		allTags := getAllTagsFromDICOM(dataset)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(allTags))
		return
	}

	value, err := getTagFromDICOM(dataset, tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"tag value": "%v"}`, value)))
}

func getTagFromDICOM(dataset dicom.Dataset, tag string) (string, error) {
	switch {
	case hexCoordPattern.MatchString(tag):
		group, element, err := parseHexTags(tag)
		if err != nil {
			return "", fmt.Errorf("tag is of a wrong hex format: %s", tag)
		}
		return findTagInDataset(dataset, group, element)
	case intPattern.MatchString(tag):
		group, element, err := parseIntTags(tag)
		if err != nil {
			return "", fmt.Errorf("tag is of a wrong integer format: %s", tag)
		}
		return findTagInDataset(dataset, group, element)
	default:
		return "", fmt.Errorf("tag is of an unknown format: %s", tag)
	}
}

func findTagInDataset(dataset dicom.Dataset, tagGroup uint16, tagElement uint16) (string, error) {
	for _, elem := range dataset.Elements {
		if elem.Tag.Group == tagGroup && elem.Tag.Element == tagElement {
			return fmt.Sprintf("%v", elem.Value), nil
		}
	}
	return "", fmt.Errorf("tag not found in DICOM file")
}

func parseIntTags(tag string) (uint16, uint16, error) {
	tagNumber, err := strconv.Atoi(tag)
	// This should never happen as we already check for number format before calling this.
	if err != nil {
		return 0, 0, fmt.Errorf("tag %s can't be parsed to number", tag)
	}

	// The four left hex digits are for group, and the right ones for element.
	tagGroup := uint16(tagNumber >> 16)
	tagElement := uint16(tagNumber & 0xFFFF)

	return uint16(tagGroup), uint16(tagElement), nil
}

func parseHexTags(tag string) (uint16, uint16, error) {
	tag = strings.Trim(tag, "()")
	parts := strings.Split(tag, ",")

	// The four left hex digits are for group, and the right ones for element.
	tagGroup, err := strconv.ParseUint(parts[0], 16, 16)
	if err != nil {
		return 0, 0, fmt.Errorf("group tag is in wrong format: %s", tag)
	}

	tagElement, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return 0, 0, fmt.Errorf("element tag is in wrong format: %s", tag)
	}

	return uint16(tagGroup), uint16(tagElement), nil
}

// Debug: This function is for debugging.
func getAllTagsFromDICOM(dataset dicom.Dataset) string {
	result := "{\n"
	for tag, elem := range dataset.Elements {
		result += fmt.Sprintf(`  "0x%04X": "%v",\n`, tag, elem.Value)
	}
	result += "}\n"
	return result
}
