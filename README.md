# nixeepass

## Install

```shell
helm repo add phntom phntom.kix.co.il/charts
wget https://raw.githubusercontent.com/phntom/nixeepass/main/charts/nixeepass/edit-me.yaml
nano edit-me.yaml
helm install nixeepass phntom/nixeepass -f edit-me.yaml
```
##


## Build

```shell
nvm install 18.12
cd static/ui
npm run build
```
