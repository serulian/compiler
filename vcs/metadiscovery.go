package vcs

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// discoveryUrlParamter holds the name of the parameter added to the URL to request it
// to add the golang <meta> discovery tag.
const discoveryUrlParamter = "go-get"

// discoveryUrlValue holds the value of the paramter added to the URL to request it
// to add the golang <meta> discovery tag.
const discoveryUrlValue = "1"

// discoveryMetaTagName holds the name of the <meta> tag holding the VCS discovery
// information.
const discoveryMetaTagName = "go-import"

// discoverInformationForVCSUrl attempts to download the given URL, find the discovery <meta> tag,
// and return the VCSUrlInformation found.
func discoverInformationForVCSUrl(vcsUrl string) (VCSUrlInformation, error) {
	log.Printf("Parsing VCS URL %s", vcsUrl)

	// Parse the VCS url.
	parsedUrl, err := url.Parse(vcsUrl)
	if err != nil {
		return VCSUrlInformation{}, fmt.Errorf("Could not parse VCS url '%s': %v", vcsUrl, err)
	}

	// Add the discovery parameter.
	query := parsedUrl.Query()
	query.Set(discoveryUrlParamter, discoveryUrlValue)
	parsedUrl.RawQuery = query.Encode()

	// Make sure we have a known scheme.
	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "http"
	}

	// Attempt to download the URL.
	log.Printf("Retrieving VCS URL %s", parsedUrl.String())
	response, err := http.Get(parsedUrl.String())
	if err != nil {
		return VCSUrlInformation{}, fmt.Errorf("Could not download VCS url '%s': %v", vcsUrl, err)
	}

	// Parse the contents into HTML and search for the <meta> tag.
	doc, err := html.Parse(response.Body)
	if err != nil {
		return VCSUrlInformation{}, fmt.Errorf("VCS url '%s' returned an invalid body", vcsUrl)
	}

	notFoundErr := fmt.Errorf("Could not discover VCS url '%s'", vcsUrl)
	discoveryTag, ok := findVCSDiscoveryTag(doc)
	if !ok {
		log.Printf("Discovery <meta> tag not found for VCS URL %v", vcsUrl)
		return VCSUrlInformation{}, notFoundErr
	}

	// Find the <meta> tag's value.
	var metaTagValue string
	for _, attr := range discoveryTag.Attr {
		if attr.Key == "content" {
			metaTagValue = attr.Val
			break
		}
	}

	if metaTagValue == "" {
		log.Printf("<meta> tag value not found for VCS URL %v", vcsUrl)
		return VCSUrlInformation{}, notFoundErr
	}

	// <meta name="go-import" content="github.com/repo/path git https://github.com/repo/path">
	pieces := strings.SplitN(metaTagValue, " ", 3)
	if pieces == nil {
		log.Printf("<meta> tag value could not be parsed for VCS URL %v", vcsUrl)
		return VCSUrlInformation{}, notFoundErr
	}

	// Ensure that the path component (index 0) matches or is a prefix of the import URL.
	if strings.Index(vcsUrl, pieces[0]) != 0 {
		log.Printf("<meta> tag prefix '%v' does not match VCS URL %v", pieces[0], vcsUrl)
		return VCSUrlInformation{}, notFoundErr
	}

	// Find the associated VCS.
	handler, ok := vcsById[pieces[1]]
	if !ok {
		err := fmt.Errorf("VCS url '%s' requires engine '%s', which is not currently supported", vcsUrl, pieces[1])
		return VCSUrlInformation{}, err
	}

	return VCSUrlInformation{pieces[0], handler.kind, pieces[2]}, nil
}

// findVCSDiscoveryTag searches the parsed HTML DOM from the given start node, downward,
// looking for the VCS <meta> discovery tag. If found, the tag is returned.
func findVCSDiscoveryTag(n *html.Node) (*html.Node, bool) {
	if n.Type == html.ElementNode && n.Data == "meta" {
		for _, attr := range n.Attr {
			if attr.Key == "name" && attr.Val == discoveryMetaTagName {
				return n, true
			}
		}
	}

	// Search recursively.
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		found, ok := findVCSDiscoveryTag(c)
		if ok {
			return found, true
		}
	}

	return nil, false
}
