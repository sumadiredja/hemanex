package registry

import (
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"

	helper "hemanex/helper"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

const ACCEPT_HEADER = "application/vnd.docker.distribution.manifest.v2+json"

type Registry struct {
	Host           string `toml:"nexus_host"`
	Username       string `toml:"nexus_username"`
	Password       string `toml:"nexus_password"`
	RepositoryPort string `toml:"nexus_repository_port"`
	Repository     string `toml:"nexus_repository"`
	NexusPort      string `toml:"nexus_host_port"`
	Namespace      string `toml:"nexus_namespace"`
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
	var CREDENTIALS_FILE string
	var err error

	r := Registry{}

	if CREDENTIALS_FILE, err = helper.GetCredentialPath(); err != nil {
		return r, err
	}

	if _, err := os.Stat(CREDENTIALS_FILE); os.IsNotExist(err) {
		return r, helper.CliErrorGen(fmt.Errorf("user not logged in"), 1)
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
	var repositories Repositories

	client := &http.Client{}
	var host_port = ""

	if r.NexusPort != "443" && r.NexusPort != "80" {
		host_port = ":" + r.NexusPort
	}

	url := fmt.Sprintf("%s%s/repository/%s/v2/_catalog", r.Host, host_port, r.Repository)
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
		return nil, helper.CliErrorGen(fmt.Errorf("HTTP Code: %d", resp.StatusCode), 1)
	}

	err = json.NewDecoder(resp.Body).Decode(&repositories)
	if err != nil {
		return nil, err
	}

	return repositories.Images, nil
}

func (r Registry) ListTagsByImage(image string) ([]string, error) {
	var imageTags ImageTags

	client := &http.Client{}
	var host_port = ""

	if r.NexusPort != "443" && r.NexusPort != "80" {
		host_port = ":" + r.NexusPort
	}
	url := fmt.Sprintf("%s%s/repository/%s/v2/%s/tags/list", r.Host, host_port, r.Repository, image)
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
		return nil, helper.CliErrorGen(fmt.Errorf("HTTP Code: %d", resp.StatusCode), 1)
	}

	err = json.NewDecoder(resp.Body).Decode(&imageTags)
	if err != nil {
		return nil, err
	}

	return imageTags.Tags, nil
}

func (r Registry) ImageManifest(image string, tag string) (ImageManifest, error) {
	var imageManifest ImageManifest
	client := &http.Client{}
	var host_port = ""

	if r.NexusPort != "443" && r.NexusPort != "80" {
		host_port = ":" + r.NexusPort
	}

	url := fmt.Sprintf("%s%s/repository/%s/v2/%s/manifests/%s", r.Host, host_port, r.Repository, image, tag)
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
		return imageManifest, helper.CliErrorGen(fmt.Errorf("HTTP Code: %d", resp.StatusCode), 1)
	}

	err = json.NewDecoder(resp.Body).Decode(&imageManifest)
	if err != nil {
		return imageManifest, err
	}

	return imageManifest, nil

}

func (r Registry) DeleteImageByTag(image string, tag string) error {
	sha, err := r.getImageSHA(image, tag)
	if err != nil {
		return err
	}
	client := &http.Client{}

	var host_port = ""

	if r.NexusPort != "443" && r.NexusPort != "80" {
		host_port = ":" + r.NexusPort
	}

	url := fmt.Sprintf("%s%s/repository/%s/v2/%s/manifests/%s", r.Host, host_port, r.Repository, image, sha)
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
		return helper.CliErrorGen(fmt.Errorf("HTTP Code: %d", resp.StatusCode), 1)
	}

	helper.CliSuccessVerbose(fmt.Sprintf("%s:%s has been successful deleted\n", image, tag))

	return nil
}

func (r Registry) getImageSHA(image string, tag string) (string, error) {
	client := &http.Client{}
	var host_port = ""

	if r.NexusPort != "443" && r.NexusPort != "80" {
		host_port = ":" + r.NexusPort
	}

	url := fmt.Sprintf("%s%s/repository/%s/v2/%s/manifests/%s", r.Host, host_port, r.Repository, image, tag)
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
		return "", helper.CliErrorGen(fmt.Errorf("HTTP Code: %d", resp.StatusCode), 1)
	}

	return resp.Header.Get("docker-content-digest"), nil
}
