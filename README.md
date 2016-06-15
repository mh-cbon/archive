# archive

bin to create / extract zip, tgz archives.

# install

```
mkdir -p $GOPATH/github.com/mh-cbon
cd $GOPATH/github.com/mh-cbon
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
   --output value, -o value	Output file
   --force, -f			Force overwrite
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
```
