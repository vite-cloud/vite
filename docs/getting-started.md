# Getting started
First, you need to set up a remote repository that will host your configuration. You can use any git provider
you'd like but keep in mind that [vite.cloud](https://vite.cloud) only supports GitHub at the time of writing.

Once you created a repository, you can initialize a new Vite project by running the following command:

```bash
$ vite setup
```
```bash
? Select your provider:  [Use arrows to move, type to filter]
> github
  gitlab
  bitbucket

? Select your protocol:  [Use arrows to move, type to filter]
> ssh
  https
  auto
  
? Enter your repository: 
? Enter your branch: (main) 
? Enter a sub-path (optional): 
```

You may always re-run `vite setup` to update the configuration but editing manually is not recommended.


Second, you need to choose which commit you want Vite to use in order to read your configuration:

```bash
$ vite use
? Select a commit  [Use arrows to move, type to filter]
  4e1aeb171526b75e0e891c924d4d2448f563cb7d Update the control plane host
  97052197e893bbc5feed19c44445cfebfdf20dae initial commit
```

> Can't see your commit? Run `vite use --pull` to pull the latest commits from remote.

And you're done!

### What's next?

* [Deploying your first service](deploying-your-first-service.md)
* [Vite in Production](vite-in-production.md)
