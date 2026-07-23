package pr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"github.com/jenkins-x-plugins/jx-updatebot/pkg/apis/updatebot/v1alpha1"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/jx-helpers/v3/pkg/helmer"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"
	"github.com/jenkins-x/jx-helpers/v3/pkg/versionstream"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/credentials"
	"oras.land/oras-go/v2/registry/remote/retry"
)

// ApplyVersionStream applies the version stream change
func (o *Options) ApplyVersionStream(dir string, vs *v1alpha1.VersionStreamChange) error {
	kind := vs.Kind
	if kind == "" {
		return options.MissingOption("kind")
	}
	if stringhelpers.StringArrayIndex(versionstream.KindStrings, kind) < 0 {
		return options.InvalidOption("kind", kind, versionstream.KindStrings)
	}

	if kind == string(versionstream.KindChart) {
		err := o.applyVersionStreamCharts(dir, vs, kind)
		if err != nil {
			return fmt.Errorf("failed to apply kind %s: %w", kind, err)
		}
	}

	return nil
}

func (o *Options) applyVersionStreamCharts(dir string, vs *v1alpha1.VersionStreamChange, kindStr string) error {
	prefixes, err := versionstream.GetRepositoryPrefixes(dir)
	if err != nil {
		return fmt.Errorf("failed to load chart repository prefixes: %w", err)
	}

	kindDir := filepath.Join(dir, kindStr)
	glob := filepath.Join(kindDir, "*", "defaults.yaml")
	paths, err := filepath.Glob(glob)
	if err != nil {
		return fmt.Errorf("bad glob pattern %s: %w", glob, err)
	}
	glob = filepath.Join(kindDir, "*", "*", "defaults.yaml")
	morePaths, err := filepath.Glob(glob)
	if err != nil {
		return fmt.Errorf("bad glob pattern %s: %w", glob, err)
	}
	paths = append(paths, morePaths...)

	o.CommitTitle = "chore: upgrade charts"
	o.CommitMessage = ""

	chartInfos := map[string]*chartInfo{}
	for _, path := range paths {
		rel, err := filepath.Rel(kindDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path of %s: %w", path, err)
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
			return fmt.Errorf("failed to add helm repository %s for prefix %s: %w", ci.RepoURL, repoPrefix, err)
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
			sv, err := versionstream.LoadStableVersion(dir, versionstream.VersionKind(kindStr), name)
			if err != nil {
				return fmt.Errorf("failed to load stable version for %s: %w", name, err)
			}

			oldVersion := sv.Version
			if oldVersion == "" {
				log.Logger().Debugf("no upgrade is done of chart %s since no version is set", name)
				continue
			}
			var upperLimit *semver.Version
			if sv.UpperLimit != "" {
				upperLimit = &semver.Version{}
				*upperLimit, err = semver.ParseTolerant(sv.UpperLimit)
				if err != nil {
					log.Logger().WithError(err).Errorf("upperLimit '%s' cannot be parsed. Skipping", sv.UpperLimit)
					continue
				}
			}
			version := ""
			if strings.HasPrefix(ci.RepoURL, "oci://") {
				// shim for lack of support for searching OCI charts in helm cli
				ociRepo := scm.Join(ci.RepoURL, n)
				version, err = ociFindLatestVersion(ociRepo, upperLimit)
				if err != nil {
					return fmt.Errorf("failed to search for chart %s: %w", ociRepo, err)
				}
			} else {
				info, err := o.Helmer.SearchCharts(name, true)
				if err != nil {
					return fmt.Errorf("failed to search for chart %s: %w", name, err)
				}
				if len(info) == 0 {
					log.Logger().Warnf("no version found for chart %s", name)
					continue
				}
				for i := range info {
					chartSummary := info[i]
					if upperLimit != nil {
						parsedVersion, err := semver.ParseTolerant(chartSummary.ChartVersion)
						if err != nil {
							log.Logger().WithError(err).
								Debugf("ignore version %s since it is malformed and upperLimit is set for chart %s",
									chartSummary.ChartVersion, name)
							continue
						}
						if parsedVersion.GE(*upperLimit) {
							log.Logger().Debugf("ignore version %s since it conflicts with upperLimit for chart %s",
								chartSummary.ChartVersion, name)
							continue
						}
					}
					version = chartSummary.ChartVersion
					break
				}
			}
			if version == "" {
				log.Logger().Warnf("no chart version found for chart %s", name)
				continue
			}

			if oldVersion != version {
				_, err := versionstream.UpdateStableVersion(dir, kindStr, name, version)
				if err != nil {
					return fmt.Errorf("failed to upgrade version of %s to %s: %w", name, version, err)
				}
				log.Logger().Infof("updated chart %s from %s to %s", name, oldVersion, version)

				if o.CommitMessage != "" {
					o.CommitMessage += "\n"
				}
				chartText := name
				chartURL := sv.GitURL
				if chartURL == "" {
					chartURL = sv.URL
				}
				if chartURL != "" {
					chartText = fmt.Sprintf("[%s](%s)", name, chartURL)
				}
				o.CommitMessage += fmt.Sprintf("* updated chart %s from `%s` to `%s`", chartText, oldVersion, version)
			}
		}
	}
	return nil
}

// This method only returns the minimal answer needed
func ociFindLatestVersion(ociRepo string, upperLimit *semver.Version) (string, error) {
	repo, err := remote.NewRepository(strings.TrimPrefix(ociRepo, "oci://"))
	if err != nil {
		return "", err
	}

	docker, err := credentials.NewStoreFromDocker(credentials.StoreOptions{
		AllowPlaintextPut:        false,
		DetectDefaultNativeStore: false,
	})
	if err != nil {
		return "", err
	}
	// Note: The below code can be omitted if authentication is not required.
	repo.Client = &auth.Client{
		Client:     retry.DefaultClient,
		Cache:      auth.NewCache(),
		Credential: docker.Get,
	}
	latestVersion := ""
	var latestFound semver.Version
	err = repo.Tags(context.Background(), "", func(tags []string) error {
		for _, tag := range tags {
			// Change underscore (_) back to plus (+) for Helm
			version, err := semver.ParseTolerant(strings.ReplaceAll(tag, "_", "+"))
			if err != nil {
				log.Logger().WithError(err).Debugf("ignore tag that doesn't look like version: %s", tag)
				continue
			}
			log.Logger().Debugf("considering tag that does look like version: %s", tag)
			if version.GT(latestFound) && (upperLimit == nil || version.LT(*upperLimit)) {
				latestFound = version
				latestVersion = strings.ReplaceAll(tag, "_", "+")
			}
		}
		return nil
	})
	return latestVersion, err
}

type chartInfo struct {
	RepoURL string
	Names   []string
}
