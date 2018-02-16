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
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var testTgzManifest = ProfileManifestV0{
	Kind: "profile-manifest-v0",
	Value: Images{
		[]Image{
			Image{
				Name:      "test-name",
				Reference: "test-reference",
				Format:    "tgz",
			},
		},
	},
}

var testSquashfsManifest = ProfileManifestV0{
	Kind: "profile-manifest-v0",
	Value: Images{
		[]Image{
			Image{
				Name:      "test-name",
				Reference: "test-reference",
				Format:    "squashfs",
			},
		},
	},
}

func TestGetProfile(t *testing.T) {
	// Schema of profile v0 is described in
	// https://github.com/coreos/torcx/blob/master/Documentation/schemas/profile-manifest-v0.md
	// with each possible value for 'format'
	testItems := []struct {
		name     string
		path     string
		manifest ProfileManifestV0
	}{
		{
			name:     "tgz",
			path:     "../../fixtures/test-no-format-profile-v0.json",
			manifest: testTgzManifest,
		},
		{
			name:     "default",
			path:     "../../fixtures/test-tgz-profile-v0.json",
			manifest: testTgzManifest,
		},
		{
			name:     "squashfs",
			path:     "../../fixtures/test-squashfs-profile-v0.json",
			manifest: testSquashfsManifest,
		},
	}

	for _, item := range testItems {
		t.Run(item.name, func(t *testing.T) {
			if _, err := os.Stat(item.path); err != nil {
				t.Fatal(err)
			}

			outManifest, err := getProfile(item.path)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(item.manifest, outManifest) {
				t.Fatalf("manifests do not match with each other.\nin:%v\nout:%v\n", item.manifest, outManifest)
			}
		})
	}
}

func TestPutGetProfile(t *testing.T) {
	tmp, err := ioutil.TempFile("", "test-torcx-putget-profile")
	if err != nil {
		t.Fatal(err)
	}
	profilePath := tmp.Name()
	defer os.Remove(profilePath)

	if err := tmp.Close(); err != nil {
		t.Fatal(err)
	}

	if err := putProfile(profilePath, os.FileMode(0755), testTgzManifest); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(profilePath); err != nil {
		t.Fatal(err)
	}

	outManifest, err := getProfile(profilePath)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testTgzManifest, outManifest) {
		t.Fatalf("manifests do not match with each other.\nin:%v\nout:%v\n", testTgzManifest, outManifest)
	}
}

func TestMergeImages(t *testing.T) {
	testCases := []struct {
		desc  string
		lower []Image
		upper []Image

		merged []Image
	}{
		{
			"empty total",
			[]Image{},
			[]Image{},

			[]Image{},
		}, {
			"empty lower",
			[]Image{},
			[]Image{
				Image{
					Name:      "foo0",
					Reference: "0",
				},
			},

			[]Image{
				Image{
					Name:      "foo0",
					Reference: "0",
				},
			},
		}, {
			"empty upper",
			[]Image{
				Image{
					Name:      "foo1",
					Reference: "1",
				},
			},
			[]Image{},

			[]Image{
				Image{
					Name:      "foo1",
					Reference: "1",
				},
			},
		}, {
			"ordered merge",
			[]Image{
				Image{
					Name:      "foo1",
					Reference: "1",
				},
			},
			[]Image{
				Image{
					Name:      "foo2",
					Reference: "2",
				},
			},

			[]Image{
				Image{
					Name:      "foo1",
					Reference: "1",
				},
				Image{
					Name:      "foo2",
					Reference: "2",
				},
			},
		}, {
			"re-order and override reference in lower",
			[]Image{
				Image{
					Name:      "foo2",
					Reference: "3",
				},
				Image{
					Name:      "foo1",
					Reference: "1",
				},
			},
			[]Image{
				Image{
					Name:      "foo2",
					Reference: "2",
				},
			},

			[]Image{
				Image{
					Name:      "foo1",
					Reference: "1",
				},
				Image{
					Name:      "foo2",
					Reference: "2",
				},
			},
		}, {
			"remove image in lower",
			[]Image{
				Image{
					Name:      "foo3",
					Reference: "3",
				},
			},
			[]Image{
				Image{
					Name:      "foo3",
					Reference: "",
				},
			},

			[]Image{},
		},
	}

	for _, tt := range testCases {
		expected := Images{tt.merged}
		lower := Images{Images: tt.lower}
		upper := Images{Images: tt.upper}
		result := mergeImages(lower, upper)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("testcase %q failed: got %#v, expected %#v", tt.desc, result.Images, expected.Images)
		}

	}
}
