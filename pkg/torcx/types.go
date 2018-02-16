// Copyright 2017 CoreOS Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package torcx

import (
	"encoding/json"
	"fmt"
)

const (
	// SealUpperProfile is the key label for user profile name
	SealUpperProfile = "TORCX_UPPER_PROFILE"
	// SealLowerProfiles is the key label for vendor profile path
	SealLowerProfiles = "TORCX_LOWER_PROFILES"
	// SealRunProfilePath is the key label for vendor profile path
	SealRunProfilePath = "TORCX_PROFILE_PATH"
	// SealBindir is the key label for seal bindir
	SealBindir = "TORCX_BINDIR"
	// SealUnpackdir is the key label for seal runtime unpackdir path
	SealUnpackdir = "TORCX_UNPACKDIR"
	// ProfileManifestV0K - profile manifest kind, v0
	ProfileManifestV0K = "profile-manifest-v0"
	// ImageManifestV0K - image manifest kind, v0
	ImageManifestV0K = "image-manifest-v0"
	// CommonConfigV0K - common torcx config kind, v0
	CommonConfigV0K = "torcx-config-v0"
)

// ConfigV0 holds common torcx configuration in JSON format
type ConfigV0 struct {
	Kind  string       `json:"kind"`
	Value CommonConfig `json:"value"`
}

// CommonConfig contains runtime configuration items common to all
// torcx subcommands
type CommonConfig struct {
	BaseDir    string   `json:"base_dir,omitempty"`
	RunDir     string   `json:"run_dir,omitempty"`
	ConfDir    string   `json:"conf_dir,omitempty"`
	StorePaths []string `json:"store_paths,omitempty"`
}

// ApplyConfig contains runtime configuration items specific to
// the `apply` subcommand
type ApplyConfig struct {
	CommonConfig
	LowerProfiles []string
	UpperProfile  string
}

// ProfileConfig contains runtime configuration items specific to
// the `profile` subcommand
type ProfileConfig struct {
	CommonConfig
	LowerProfileNames  []string
	UserProfileName    string
	CurrentProfilePath string
	NextProfile        string
}

// ProfileManifestV0 holds JSON profile manifest
type ProfileManifestV0 struct {
	Kind  string `json:"kind"`
	Value Images `json:"value"`
}

// Archive represents a .torcx.squashfs or .torcx.tgz on disk
type Archive struct {
	Image
	Filepath string `json:"filepath"`
}

// Image is an archive name + Reference
type Image struct {
	Name      string      `json:"name"`
	Reference string      `json:"reference"`
	Format    ImageFormat `json:"format"`
}

// ImageFormat is a torcx archive format, either 'tgz' or 'squashfs'
type ImageFormat string

const (
	// ImageFormatUnknown is the zero value of ImageFormat. It indicates the image format is unknown
	ImageFormatUnknown ImageFormat = ""
	// ImageFormatTgz indicates a tar-gzipped image
	ImageFormatTgz = "tgz"
	// ImageFormatSquashfs indicates a squashfs image archive
	ImageFormatSquashfs = "squashfs"
)

// UnmarshalJSON unmarshals the "format" field of an "Image"
func (imgf *ImageFormat) UnmarshalJSON(b []byte) error {
	s := ""
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "", ImageFormatTgz:
		// Defaulting for 'format' to 'tgz', backwards compatibility
		*imgf = ImageFormatTgz
	case ImageFormatSquashfs:
		*imgf = ImageFormatSquashfs
	default:
		return fmt.Errorf("could not unmarshal into ImageFormat: must be one of %q, %q", ImageFormatTgz, ImageFormatSquashfs)
	}
	return nil
}

// FileSuffix returns the file extension this image format must have.
func (imgf ImageFormat) FileSuffix() string {
	return fmt.Sprintf(".torcx.%s", imgf)
}

type Images struct {
	Images []Image `json:"images"`
}

// ImageManifestV0 holds JSON image manifest
type ImageManifestV0 struct {
	Kind  string `json:"kind"`
	Value Assets `json:"value"`
}

// Assets holds lists of assets propagated from an image to the system
type Assets struct {
	Binaries []string `json:"bin,omitempty"`
	Network  []string `json:"network,omitempty"`
	Units    []string `json:"units,omitempty"`
	Sysusers []string `json:"sysusers,omitempty"`
	Tmpfiles []string `json:"tmpfiles,omitempty"`
}
