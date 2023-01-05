package registry

import (
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hemanex/helper"
	"html"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

const ACCEPT_HEADER = "application/vnd.docker.distribution.manifest.v2+json"

type Registry struct {
	Host       string `toml:"nexus_host"`
	Username   string `toml:"nexus_username"`
	Password   string `toml:"nexus_password"`
	Repository string `toml:"nexus_repository"`
}

type Repositories struct {
	Images []string `json:"repositories"`
}

type ImageTags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ImageManifest struct {
	SchemaVersion int64       `json:"schemaVersion"`
	MediaType     string      `json:"mediaType"`
	Config        LayerInfo   `json:"config"`
	Layers        []LayerInfo `json:"layers"`
}
type LayerInfo struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

func NewRegistry(cli_context *cli.Context) (Registry, error) {
	r := Registry{}
	CREDENTIALS_FILE, err := helper.GetCredentialPath()

	if err != nil {
		return r, err
	}

	if _, err := os.Stat(CREDENTIALS_FILE); os.IsNotExist(err) {
		return r, fmt.Errorf("User not logged in")
	} else if err != nil {
		return r, err
	}

	if _, err := toml.DecodeFile(CREDENTIALS_FILE, &r); err != nil {
		return r, err
	}

	unescapePassword := html.UnescapeString(r.Password)
	decodePassword, _ := b64.StdEncoding.DecodeString(unescapePassword)
	r.Password = string(decodePassword)

	if cli_context.Bool("insecure-registry") {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return r, nil
}

func (r Registry) ListImages() ([]string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/_catalog", r.Host, r.Repository)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", ACCEPT_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Code: %d", resp.StatusCode))
	}

	var repositories Repositories
	json.NewDecoder(resp.Body).Decode(&repositories)

	return repositories.Images, nil
}

func (r Registry) ListTagsByImage(image string) ([]string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/tags/list", r.Host, r.Repository, image)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", ACCEPT_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP Code: %d", resp.StatusCode))
	}

	var imageTags ImageTags
	json.NewDecoder(resp.Body).Decode(&imageTags)

	return imageTags.Tags, nil
}

func (r Registry) ImageManifest(image string, tag string) (ImageManifest, error) {
	var imageManifest ImageManifest
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return imageManifest, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", ACCEPT_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return imageManifest, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return imageManifest, errors.New(fmt.Sprintf("HTTP Code: %d", resp.StatusCode))
	}

	json.NewDecoder(resp.Body).Decode(&imageManifest)

	return imageManifest, nil

}

func (r Registry) DeleteImageByTag(image string, tag string) error {
	sha, err := r.getImageSHA(image, tag)
	if err != nil {
		return err
	}
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, sha)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", ACCEPT_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return errors.New(fmt.Sprintf("HTTP Code: %d", resp.StatusCode))
	}

	fmt.Printf("%s:%s has been successful deleted\n", image, tag)

	return nil
}

func (r Registry) getImageSHA(image string, tag string) (string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", ACCEPT_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("HTTP Code: %d", resp.StatusCode))
	}

	return resp.Header.Get("docker-content-digest"), nil
}
