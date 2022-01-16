# Remote file

This is example code for referencing remote scaffold file.

```shell
# go get style (currently only public repository is supported)
manaita -c github.com/igtm/manaita/example/go-ddd-api/SCAFFOLD.md -p name=company # default branch
manaita -c github.com/igtm/manaita/example/go-ddd-api/SCAFFOLD.md@master -p name=company # branch
manaita -c github.com/igtm/manaita/example/go-ddd-api/SCAFFOLD.md@v1.0.6 -p name=company # version
manaita -c github.com/igtm/manaita/example/go-ddd-api/SCAFFOLD.md@a06f1da -p name=company # commit hash
# http url
manaita -c https://raw.githubusercontent.com/igtm/manaita/master/example/go-ddd-api/SCAFFOLD.md -p name=company
```
