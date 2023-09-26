### Run

```sh
make dev
```

- `local.db` file will be created and it will sync
- check `siteMaps` array in `const.go` file for setting the list of site-map to be crawled
- `SyncContents` is commented which is the "whole" data crawl
- `SyncContentsV2` is active, which is very limited data crawl (only title and description)
- Error logs will be printed to `log.log` file

