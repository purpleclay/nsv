package ci

import (
	"os"
	"sync"
)

// Environment captures details of a continous integration (CI) platform
type Environment struct {
	// SkipPipelineTag defines a tag that can be injected into the first
	// line of a commit message to prevent the CI platform from running
	// an unnecessary build.
	//
	// Supported platforms and their corresponding tag
	// 	- [GitHub] [skip ci]
	// 	- [GitLab] [skip ci]
	// 	- [CircleCI] [skip ci]
	// 	- [Travis CI] [skip ci]
	// 	- [Drone] [CI SKIP]
	// 	- [Semaphore] [skip ci]
	// 	- [Codefresh] [skip ci]
	// 	- [Cirrus CI] [skip ci]
	// 	- [Buildkite] [skip ci]
	// 	- [Jenkins] [ci skip]
	// 	- [Bitbucket] [skip ci]
	//
	// [GitHub]: https://github.blog/changelog/2021-02-08-github-actions-skip-pull-request-and-push-workflows-with-skip-ci/
	// [GitLab]: https://docs.gitlab.com/ee/ci/pipelines/#skip-a-pipeline
	// [CircleCI]: https://circleci.com/docs/skip-build/
	// [Travis CI]: https://docs.travis-ci.com/user/customizing-the-build#skipping-a-build
	// [Drone]: https://docs.drone.io/pipeline/skipping/
	// [Semaphore]: https://docs.semaphoreci.com/essentials/skip-building-some-commits-with-ci-skip/
	// [Codefresh]: https://codefresh.io/docs/docs/pipelines/triggers/git-triggers/#skip-triggering-pipeline-on-commit
	// [Cirrus CI]: https://cirrus-ci.org/guide/writing-tasks/#conditional-task-execution
	// [Buildkite]: https://buildkite.com/docs/pipelines/skipping#ignore-a-commit
	// [Jenkins]: https://plugins.jenkins.io/scmskip/
	// [Bitbucket]: https://confluence.atlassian.com/bbkb/how-to-skip-triggering-an-automatic-pipeline-build-using-skip-ci-label-1207188270.html
	SkipPipelineTag string
}

func droneCI(res chan<- Environment) {
	// https://docs.drone.io/pipeline/environment/reference/
	if os.Getenv("DRONE") == "true" {
		res <- Environment{SkipPipelineTag: "[CI SKIP]"}
	}
}

func jenkinsCI(res chan<- Environment) {
	// https://www.jenkins.io/doc/book/pipeline/jenkinsfile/#using-environment-variables
	if os.Getenv("JENKINS_URL") != "" {
		res <- Environment{SkipPipelineTag: "[ci skip]"}
	}
}

// Detect will attempt to identify the current continous integration (CI) platform
// by checking for predefined environment variables. Once detected, details about
// the CI platform will be collated
func Detect() Environment {
	chn := make(chan Environment, 1)

	detectors := sync.WaitGroup{}
	detectors.Add(2)

	goDetect := func(detectCI func(res chan<- Environment)) {
		defer detectors.Done()
		detectCI(chn)
	}

	go goDetect(droneCI)
	go goDetect(jenkinsCI)

	detectors.Wait()
	close(chn)

	env := Environment{SkipPipelineTag: "[skip ci]"}
	for denv := range chn {
		env = denv
	}

	return env
}
