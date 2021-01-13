package pr

import (
	"fmt"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/apis/updatebot/v1alpha1"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/jx-helpers/v3/pkg/helmer"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"
	"github.com/jenkins-x/jx-helpers/v3/pkg/versionstream"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

// ApplyVersionStream applies the version stream change
func (o *Options) ApplyVersionStream(dir string, gitURL string, change v1alpha1.Change, vs *v1alpha1.VersionStreamChange) error {
	kind := vs.Kind
	if kind == "" {
		return options.MissingOption("kind")
	}
	if stringhelpers.StringArrayIndex(versionstream.KindStrings, kind) < 0 {
		return options.InvalidOption("kind", kind, versionstream.KindStrings)
	}

	if kind == string(versionstream.KindChart) {
		err := o.applyVersionStreamCharts(dir, gitURL, change, vs, kind)
		if err != nil {
			return errors.Wrapf(err, "failed to apply kind %s", kind)
		}
	}

	return nil
}

func (o *Options) applyVersionStreamCharts(dir string, url string, change v1alpha1.Change, vs *v1alpha1.VersionStreamChange, kindStr string) error {
	prefixes, err := versionstream.GetRepositoryPrefixes(dir)
	if err != nil {
		return errors.Wrapf(err, "failed to load chart repository prefixes")
	}

	kindDir := filepath.Join(dir, kindStr)
	glob := filepath.Join(kindDir, "*", "defaults.yaml")
	paths, err := filepath.Glob(glob)
	if err != nil {
		return errors.Wrapf(err, "bad glob pattern %s", glob)
	}
	glob = filepath.Join(kindDir, "*", "*", "defaults.yaml")
	morePaths, err := filepath.Glob(glob)
	if err != nil {
		return errors.Wrapf(err, "bad glob pattern %s", glob)
	}
	paths = append(paths, morePaths...)

	o.CommitTitle = "chore: upgrade charts"

	chartInfos := map[string]*chartInfo{}
	for _, path := range paths {
		rel, err := filepath.Rel(kindDir, path)
		if err != nil {
			return errors.Wrapf(err, "failed to get relative path of %s", path)
		}

		paths := strings.Split(rel, string(os.PathSeparator))
		if len(paths) < 3 {
			continue
		}

		repoPrefix := paths[0]
		chartName := paths[1]

		name := scm.Join(repoPrefix, chartName)
		if !stringhelpers.StringMatchesAny(name, vs.Includes, vs.Excludes) {
			continue
		}

		ci := chartInfos[repoPrefix]
		if ci == nil {
			ci = &chartInfo{}
			chartInfos[repoPrefix] = ci
		}
		ci.Names = append(ci.Names, chartName)
	}

	for repoPrefix, ci := range chartInfos {
		urls := prefixes.URLsForPrefix(repoPrefix)
		if len(urls) == 0 {
			log.Logger().Warnf("repository prefix %s has no URL in charts/repositories.yml", repoPrefix)
			continue
		}

		ci.RepoURL = urls[0]
		log.Logger().Infof("updating helm repository %s at %s", repoPrefix, ci.RepoURL)

		_, err = helmer.AddHelmRepoIfMissing(o.Helmer, ci.RepoURL, repoPrefix, "", "")
		if err != nil {
			return errors.Wrapf(err, "failed to add helm repository %s for prefix %s", ci.RepoURL, repoPrefix)
		}
	}

	err = o.Helmer.UpdateRepo()
	if err != nil {
		log.Logger().Warnf("failed to update helm repositories: %s", err.Error())
	}

	for repoPrefix, ci := range chartInfos {
		if ci.RepoURL == "" {
			continue
		}

		for _, n := range ci.Names {
			name := scm.Join(repoPrefix, n)
			info, err := o.Helmer.SearchCharts(name, true)
			if err != nil {
				return errors.Wrapf(err, "failed to search for chart %s", name)
			}
			if len(info) == 0 {
				log.Logger().Warnf("no version found for chart %s", name)
				continue
			}
			chartSummary := info[0]
			version := chartSummary.ChartVersion
			if version == "" {
				log.Logger().Warnf("no chart version found for chart %s", name)
				continue
			}

			sv, err := versionstream.LoadStableVersion(dir, versionstream.VersionKind(kindStr), name)
			if err != nil {
				return errors.Wrapf(err, "failed to load stable version for %s", name)
			}

			oldVersion := sv.Version
			if oldVersion != version {
				_, err := versionstream.UpdateStableVersion(dir, kindStr, name, version)
				if err != nil {
					return errors.Wrapf(err, "failed to upgrade version of %s to %s", name, version)
				}
				log.Logger().Infof("updated chart %s from %s to %s", name, oldVersion, version)

				if o.CommitTitle != "" {
					o.CommitTitle += "\n"
				}
				chartText := name
				chartURL := sv.GitURL
				if chartURL == "" {
					chartURL = sv.URL
				}
				if chartURL != "" {
					chartText = fmt.Sprintf("[%s](%s)", name, chartURL)
				}
				o.CommitTitle += fmt.Sprintf("* updated chart %s from `%s` to `%s`", chartText, oldVersion, version)
			}
		}
	}
	return nil
}

type chartInfo struct {
	RepoURL string
	Names   []string
}
