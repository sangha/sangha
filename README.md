### Build instructions

```
export GOPATH=~/go
mkdir -p $GOPATH/src/gitlab.techcultivation.org/techcultivation
cd $GOPATH/src/gitlab.techcultivation.org/techcultivation

git clone git@gitlab.techcultivation.org:techcultivation/sangha.git
cd sangha
go get -u
go build
```

### Reference

Visit Swagger: http://localhost:9991/apidocs/
