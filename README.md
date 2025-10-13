# external-assets

External assets for Meergo

## Committing to this repository

### Go through a pull request

Direct commits to `main` are **not allowed in this repository**. You must first open a Pull Request, which will trigger the required GitHub Actions. Once all checks pass, the PR branch can be merged into `main`.

### Run local checks

Before committing locally, run:

```bash
(cd potential-connectors/fmt-and-validate && go run .)
```
