// Package pkg provides functionality similar to the AWS CLI cloudformation package command
// but has greater flexibility, allowing content to be included anywhere in a template
//
// To include content into your templates, use any of the following either as YAML tags
// or as one-property objects, much as AWS instrinsic functions are used, e.g. "Fn::Join"
//
// `Rain::Include`: insert the content of the file into the template directly. The file must be in YAML or JSON format.
// `Rain::Embed`: insert the content of the file as a string
// `Rain::S3Http`: uploads the file or directory (zipping it first) to S3 and returns the HTTP URI (i.e. `https://bucket.s3.region.amazonaws.com/key`)
// `Rain::S3`: a string value uploads the file or directory (zipping it first) to S3 and returns the S3 URI (i.e. `s3://bucket/key`)
// `Rain::S3`: an object with the following properties
//    `Path`: path to the file or directory. If a directory is supplied, it will be zipped before uploading to S3
//    `BucketProperty`: Name of returned property that will contain the bucket name
//    `KeyProperty`: Name of returned property that will contain the object key
//    `VersionProperty`: (optional) Name of returned property that will contain the object version
package pkg

import (
	"github.com/aws-cloudformation/rain/cft"
	"github.com/aws-cloudformation/rain/cft/parse"
)

func transform(t cft.Template) bool {
	changed := false

	for path, fn := range registry {
		for node := range t.MatchPath(path) {
			changed = changed || fn(node)
		}
	}

	return changed
}

// Template returns a copy of the template with assets included as per the various `Include::` functions
func Template(t cft.Template) (cft.Template, error) {
	var err error

	// Keep transforming until we've recursed enough
	for {
		changed := transform(t)

		// Parse again, in case there are other things to deal with
		t, err = parse.Node(t.Node)
		if err != nil {
			return t, err
		}

		if !changed {
			break
		}
	}

	return t, nil
}