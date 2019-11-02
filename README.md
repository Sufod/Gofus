# Gofus

Gofus is a Dofus 1.30.0e full socket Bot written in Golang.

The repository also contains a MITM proxy used to simplify the bot development process.

Gofus is still in early development phase.

## Prerequisites

- Go >= 1.12

## Installation
### Compilation
```bash
git clone git@github.com:Sufod/Gofus.git
cd Gofus/cmd/client
go build
```
### Configuration

Copy the [example config file](https://github.com/Sufod/Gofus/blob/dev/configs/config.yml) and put it in the same folder as the previously generated binary (eg. Gofus/cmd/client/config.yml). Edit the values as desired.

## Usage

```bash
./client
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
### Community tools

Join us on [Slack](https://join.slack.com/t/gofus/shared_invite/enQtNzc1ODk2NDYzMjgzLTU0NzM2Mzk1YTBlZjhiNGE3YTQyMjc1NmYyYzY4ODVmMzVjZTI0MTk5Mzk5ZGIxNDQwNGE3ZDM2ZWFiM2I0NmY) or [Discord](https://discord.gg/xgNfzVJ) !

Check our [Trello](https://trello.com/b/VZBgmYfO/gofus) ! [Team Invitation Link](https://trello.com/invite/b/VZBgmYfO/9de3394da289a6e6d72063e44314e4dd/gofus)

Check the [Gofus Proposal](https://docs.google.com/document/d/1KUc18DVd6pT7niy608H0nvGO8Rt7bsbxklc0YHvisH0/edit?usp=sharing) specification document (wip) to learn more about the project.

### Using the proxy
Modify the file config.xml in your Dofus installation folder as following:
```xml
<conf name="En ligne">
	<connserver name="GofusProxy" ip="127.0.0.1" port="8081" />
	<dataserver url="data/" type="local" priority="3" />
	<dataserver url="http://staticns.ankama.com/dofus/gamedata/dofus/" priority="1" />
	<dataserver url="http://gamedata.ankama-games.com/dofus/" priority="0" />
</conf>
```

Follow the same procedure as for the client, don't forget to add a config.yml in the binary folder.

```bash
./proxy
```

After launching the proxy, you can use the official client and play normally, the proxy is logging every packets between client and server.

## License
[MIT](https://choosealicense.com/licenses/mit/)
