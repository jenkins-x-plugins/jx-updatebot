package pr

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/jenkins-x-plugins/jx-updatebot/pkg/apis/updatebot/v1alpha1"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"github.com/yargevad/filepathx"
)

// ApplyRegex applies the regex change
func (o *Options) ApplyRegex(dir, gitURL string, change v1alpha1.Change, regex *v1alpha1.Regex) error {
	pattern := regex.Pattern
	if pattern == "" {
		return errors.Errorf("no pattern for regex change %#v", change)
	}
	r, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Wrapf(err, "failed to parse change regex: %s", pattern)
	}

	namedCaptures := make([]bool, 0)
	namedCapture := false
	for i, n := range r.SubexpNames() {
		if i == 0 {
			continue
		} else if n == "version" {
			namedCaptures = append(namedCaptures, true)
			namedCapture = true
		} else {
			namedCaptures = append(namedCaptures, false)
		}
	}

	for _, g := range regex.Globs {
		path := filepath.Join(dir, g)
		matches, err := filepathx.Glob(path)
		if err != nil {
			return errors.Wrapf(err, "failed to evaluate glob %s", path)
		}
		for _, f := range matches {
			log.Logger().Infof("found file %s", f)

			data, err := os.ReadFile(f)
			if err != nil {
				return errors.Wrapf(err, "failed to load file %s", f)
			}

			text := string(data)
			version := o.Version
			if change.VersionTemplate != "" {
				version, err = o.EvaluateVersionTemplate(change.VersionTemplate, gitURL)
				if err != nil {
					return errors.Wrapf(err, "failed to valuate version template %s", change.VersionTemplate)
				}
			}

			oldVersions := make([]string, 0)

			text2 := stringhelpers.ReplaceAllStringSubmatchFunc(r, text, func(groups []stringhelpers.Group) []string {
				answer := make([]string, 0)
				for i, group := range groups {
					if namedCapture {
						// If we are using named capture, then replace only the named captures that have the right name
						if namedCaptures[i] {
							oldVersions = append(oldVersions, group.Value)
							answer = append(answer, version)
						} else {
							answer = append(answer, group.Value)
						}
					} else {
						oldVersions = append(oldVersions, group.Value)
						answer = append(answer, version)
					}
				}
				return answer
			})

			if text2 != text {
				err = os.WriteFile(f, []byte(text2), files.DefaultFileWritePermissions)
				if err != nil {
					return errors.Wrapf(err, "failed to save file %s", f)
				}
				log.Logger().Infof("modified file %s", info(f))
			}
		}
	}
	return nil
}
