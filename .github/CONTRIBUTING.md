# Contributing

## How the mirror works

[`updater.yml`](./workflows/updater.yml) publishes each new upstream version as a release:

```mermaid
flowchart TD
    A["updater.yml, every hour"] --> B["fetch.sh version"]
    B --> C{"Release with<br/>that tag exists?"}
    C -- yes --> F["Nothing to do"]
    C -- no --> D["fetch.sh download<br/>every file, verified"]
    D --> E["fetch.sh notes"]
    E --> G["Publish as a release<br/>with those notes"]
```

Everything specific to the upstream lives in [`scripts/fetch.sh`](../scripts/fetch.sh) and the READMEs. To mirror another project with this repository, rewrite the script's three commands, described at its top, and the three READMEs.
