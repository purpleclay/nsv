package nsv

import git "github.com/purpleclay/gitz"

func checkAndHealRepository(gitc *git.Client, opts Options) error {
	repo, err := gitc.Repository()
	if err != nil {
		return err
	}

	if repo.DetachedHead {
		opts.Logger.Warn("repository has a detached head - check your CI documentation")
	}

	if repo.ShallowClone {
		opts.Logger.Warn("repository is a shallow clone - check your CI documentation")

		if opts.FixShallow {
			opts.Logger.Info("fixing shallow clone by restoring history and tags")
			if _, err := gitc.Fetch(git.WithUnshallow(), git.WithTags()); err != nil {
				return err
			}
		}
	}

	return nil
}
