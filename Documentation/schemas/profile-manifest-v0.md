# Profile Manifest - v0

A "profile manifest" is a JSON data structure consumed by torcx and usually provided by an external party (e.g. an user) as a configuration file with `.json` extension.
It contains an ordered list of images (name + reference) to a be applied on a system.

## Schema

- kind (string, required)
- value (object, required)
  - images (array, required, fixed-type, not-nil, min-lenght=0)
    -(object)
      - image (string, required)
      - reference (string, required)
      - format (string, optional, "tgz" or "squashfs", default "tgz")

## Entries

- kind: hardcoded to `profile-manifest-v0` for this schema revision.
  The type+version of this JSON manifest.
- value: object containing a single typed key-value.
  Manifest content.
- value/images: array of single-type objects, arbitrary length.
  List of packages to be unpacked and set up.
- value/images/#: anonymous array entry, object
- value/images/#/image: string, compatible with OCI image name specs.
  Name of the image to unpack.
- value/images/#/reference: string, compatible with OCI image reference specs.
  Referenced image must be available in the storepath, as a file name `${image}:${reference}.torcx.${format}`
- value/images/#/format : string, either "tgz" or "squashfs" to indicate the image is stored as either a gzipped tarball or squashfs archive.
  For backwards compatibility reasons, format will default to "tgz", but "squashfs" should be preferred for new images.

## JSON schema

```json

{
  "$schema": "http://json-schema.org/draft-05/schema#",
  "type": "object",
  "properties": {
    "kind": {
      "type": "string",
      "enum": ["profile-manifest-v0"]
    },
    "value": {
      "type": "object",
      "properties": {
        "images": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "name": {
                "type": "string"
              },
              "reference": {
                "type": "string"
              },
              "format": {
                "type": "string",
                "default": "tgz",
                "enum": ["tgz", "squashfs"]
              }
            },
            "required": [
              "name",
              "reference"
            ]
          }
        }
      },
      "required": [
        "images"
      ]
    }
  },
  "required": [
    "kind",
    "value"
  ]
}

```
