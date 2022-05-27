# Deploying your first application

As you're using the OSS version of Vite, we'll assume that you have a running container registry with an image of your
app already. If not, we're working on some guides to help you get started. In the meantime, you can
check [vite.cloud](https://vite.cloud) that manages all of that for you (including creating dockerfiles).

---

A service represents your application, at the core, it's a name and an image on a registry with some configuration
options. Let's create a service called `my_nginx` with the image `nginx:1.15.8` for the domain `example.com`. 
That would result in the following:
```yaml
services:
  my_nginx:
    image: nginx:1.15.8
    hosts:
      - example.com
```

We've specified a specific version of the image instead of simply using `latest`. Using `latest`
as a version, unless its meaning is explicitly defined, may lead to unexpected behavior or be an open window for
supply chain attacks. Therefore, Vite prevents you from using `latest` as a version, there's no way to disable that behavior.

> Applications are deployed to be available to the world, once you specify a host name, Vite's proxy will redirect the
incoming traffic to the service automatically.

You may commit your changes and push them to the remote repository. Once you're done, tell Vite to use the new commit.

```bash
$ vite use --pull
? Select a commit  [Use arrows to move, type to filter]
> 8897a7d08a1e791418904afdce369818c19d2c3e Add example.com
  4e1aeb171526b75e0e891c924d4d2448f563cb7d Update the control plane host
  97052197e893bbc5feed19c44445cfebfdf20dae initial commit
```

Deploy time!

```bash
$ vite deploy
global(StartEvent): 1653662697016213030
global(StartLayerDeployment): {1 1}
my_nginx(PullImage): nginx:1.21.5
my_nginx(CreateContainer): <nil>
my_nginx(RunHook): []
my_nginx(StartContainer): <nil>
my_nginx(RunHook): []
my_nginx(FinishDeployment): <nil>
```

> If you're wondering why this looks so bad, it's a marketing technique to make you switch to vite.cloud! Jokes aside,
> PRs are welcome.

### Understanding the logs

There are two important events:

* `StartEvent`: The deployment's ID. You may also use `vite deployments latest` to get the latest deployment's ID.
* `StartLayerDeployment`: The number on the left is the current layer being deployed, the number on the right is
  the total number of layers. Usually, the last layers take the longest to deploy.

> **KEY TAKEAWAY**: As long as you're seeing stuff flowing up the screen, and it's not red, you're probably fine.

Interested in knowing how we layer your services to make the deployment faster? Check out this [guide](internals/layering.md)