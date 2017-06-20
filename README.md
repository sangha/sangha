### Build instructions

```
export GOPATH=~/go
mkdir -p $GOPATH/src/gitlab.techcultivation.org/sangha
cd $GOPATH/src/gitlab.techcultivation.org/sangha

git clone git@gitlab.techcultivation.org:sangha/sangha.git
cd sangha
go get -u
go build
```

### Run sangha

```
./sangha
```

### Reference

Visit Swagger: http://localhost:9991/apidocs/
