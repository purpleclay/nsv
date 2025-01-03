package nsv

import git "github.com/purpleclay/gitz"

func checkAndHealRepository(gitc *git.Client, opts Options) error {
	repo, err := gitc.Repository()
	if err != nil {
		return err
	}
	opts.Logger.Debug("repository summary", "detached", repo.DetachedHead, "shallow",
		repo.ShallowClone, "ref", repo.Ref)

	if repo.ShallowClone {
		opts.Logger.Warn("repository is a shallow clone and history may be missing", "depth", repo.CloneDepth)

		if opts.FixShallow {
			opts.Logger.Info("fixing shallow clone by restoring history and tags")
			if _, err := gitc.Fetch(git.WithUnshallow(), git.WithTags()); err != nil {
				return err
			}

			opts.Logger.Info("history and tags restored")
		} else {
			opts.Logger.Info("please check your ci documentation on how to resolve it")
		}
	}

	return nil
}
