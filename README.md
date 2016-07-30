# archive

cross-platform bin to create / extract zip, tgz archives.

## Install

Pick an msi package [here](https://github.com/mh-cbon/archive/releases)!

__deb/ubuntu/rpm repositories__

```sh
wget -O - https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/archive sh -xe
# or
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/archive sh -xe
```

__deb/ubuntu/rpm packages__

```sh
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/archive sh -xe
# or
wget -q -O - --no-check-certificate \
https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/archive sh -xe
```

__chocolatey__

```sh
choco install archive -y
```

__go__

```sh
mkdir -p $GOPATH/src/github.com/mh-cbon
cd $GOPATH/src/github.com/mh-cbon
git clone https://github.com/mh-cbon/archive.git
cd archive
glide install
go install
```

# Usage

```
NAME:
   archive - Command line to create and extract archive files

USAGE:
   archive <cmd> <options>

VERSION:
   0.0.0

COMMANDS:
     create     Create a new archive
     extract    Extract an archive file

GLOBAL OPTIONS:
   --help, -h      show help
   --version, -v   print the version
```

#### create

```
NAME:
   archive create - Create a new archive

USAGE:
   archive create [command options] [arguments...]

OPTIONS:
   --output value, -o value	      Output file
   --change-dir value, -C value	  Change directory before archiving files
   --force, -f			              Force overwrite

EXAMPLES:
  archive create -o test.zip README.md
  archive create -o test.zip uncompress
  archive create -o test.zip README.md uncompress
  archive create -o build/test.zip README.md uncompress
  archive create -o build/test.tgz README.md uncompress
  archive create -o build/test.tar.gz README.md uncompress
  archive create -o test.zip -C fixtures/create/test1/ .
```

#### extract

```
NAME:
   archive extract - Extract an archive file

USAGE:
   archive extract [command options] [arguments...]

OPTIONS:
   --dest value, -d value	Destination path
   --force, -f			Force overwrite

EXAMPLES:
  archive extract -d res build/test.zip
  archive extract -d res build/test.tar.gz

```
