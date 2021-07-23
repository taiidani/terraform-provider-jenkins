package jenkins

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type folder struct {
	XMLName       xml.Name         `xml:"com.cloudbees.hudson.plugins.folder.Folder"`
	Description   string           `xml:"description"`
	DisplayName   string           `xml:"displayName,omitempty"`
	Properties    folderProperties `xml:"properties"`
	FolderViews   xmlRawProperty   `xml:"folderViews"`
	HealthMetrics xmlRawProperty   `xml:"healthMetrics"`
}

type folderProperties struct {
	Security *folderSecurity  `xml:"com.cloudbees.hudson.plugins.folder.properties.AuthorizationMatrixProperty,omitempty"`
	Other    []xmlRawProperty `xml:",any"`
}

type folderSecurity struct {
	InheritanceStrategy folderPermissionInheritanceStrategy `xml:"inheritanceStrategy"`
	Permission          []string                            `xml:"permission"`
}

type folderPermissionInheritanceStrategy struct {
	Class string `xml:"class,attr"`
}

type xmlRawProperty struct {
	XMLName xml.Name
	Plugin  string `xml:"plugin,attr,omitempty"`
	Raw     string `xml:",innerxml"`
}

func parseFolder(config string) (*folder, error) {
	ret := &folder{}

	doc := handleXml(config)
	if err := xml.Unmarshal(doc, &ret); err != nil {
		return ret, fmt.Errorf("could not parse job XML: %w", err)
	}

	return ret, nil
}

func (j *folder) Render() ([]byte, error) {
	return xml.MarshalIndent(j, "", "\t")
}

func handleXml(def string) []byte {
	// This is a horrible practice...but Go doesn't seem to have any mature
	// support for the XML 1.1 specification. As long as Jenkins doesn't make
	// use of any 1.1 additions then this should still parse.
	def = strings.ReplaceAll(def, `<?xml version='1.1' encoding='UTF-8'?>`, `<?xml version='1.0' encoding='UTF-8'?>`)
	return []byte(def)
}
