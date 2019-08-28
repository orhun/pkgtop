![Logo](https://user-images.githubusercontent.com/24392180/63693894-dd110e00-c81d-11e9-8f51-e00d5bd7d6a6.png)

# pkgtop [![Release](https://img.shields.io/github/release/keylo99/pkgtop.svg?style=flat-square)](https://github.com/keylo99/pkgtop/releases)
[![AUR](https://img.shields.io/aur/version/pkgtop-git.svg?style=flat-square)](https://aur.archlinux.org/packages/pkgtop-git/)
[![Travis Build](https://img.shields.io/travis/keylo99/pkgtop.svg?style=flat-square)](https://travis-ci.org/keylo99/pkgtop) [![Docker Build](https://img.shields.io/docker/cloud/build/keylo99/pkgtop.svg?style=flat-square)](https://hub.docker.com/r/keylo99/pkgtop/builds) [![Codacy Badge](https://img.shields.io/codacy/grade/f83f3a6b0bb042f39f799cb372405094.svg?style=flat-square)](https://www.codacy.com/app/keylo99/pkgtop?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=keylo99/pkgtop&amp;utm_campaign=Badge_Grade) [![Go Report Card](https://goreportcard.com/badge/github.com/keylo99/pkgtop?style=flat-square)](https://goreportcard.com/report/github.com/keylo99/pkgtop) [![Stars](https://img.shields.io/github/stars/keylo99/pkgtop.svg?style=flat-square)](https://github.com/keylo99/pkgtop/stargazers) [![License](https://img.shields.io/github/license/keylo99/pkgtop.svg?color=blue&style=flat-square)](./LICENSE)

pkgtop is an interactive package manager and resource monitor tool designed for the GNU/Linux.

## Installation

### • Dependencies
*  [gizak/termui](https://github.com/gizak/termui/)
*  [atotto/clipboard](https://github.com/atotto/clipboard)
*  [dustin/go-humanize](https://github.com/dustin/go-humanize)

### • AUR ([pkgtop-git](https://aur.archlinux.org/packages/pkgtop-git))

### • Manual Insallation

```
go get ./...
go build src/pkgtop.go
sudo mv pkgtop /usr/local/bin/
```
Preferably, [go install](https://golang.org/cmd/go/#hdr-Compile_and_install_packages_and_dependencies) command can be used.

## Command-Line Arguments
```
-h, show help message
-d, select linux distribution
-s, sort packages alphabetically
-v, print version
```

## Usage

| Key                      	| Action                                   	|
|--------------------------	|------------------------------------------	|
| `?`                      	| help                                     	|
| `enter, space, tab`      	| show package information                 	|
| `i`                      	| install package                          	|
| `u/ctrl-u`               	| upgrade package/with input               	|
| `r/ctrl-r`               	| remove package/with input                	|
| `s`                      	| search package                           	|
| `g`                      	| go to package (index)                    	|
| `y`                      	| confirm and execute the selected command 	|
| `p`                      	| copy selected package                    	|
| `e`                      	| copy selected command                    	|
| `c`                      	| scroll executed commands list            	|
| `j/k, down/up`           	| scroll down/up (packages)                	|
| `ctrl-j/ctrl-k`          	| scroll to bottom/top (packages)          	|
| `l/h, right/left`        	| scroll down/up (disk usage)              	|
| `backspace`              	| go back                                  	|
| `q, esc, ctrl-c, ctrl-d` 	| exit                                     	|

### List Installed Packages & Show Package Information

![List Packages & Show Information](https://user-images.githubusercontent.com/24392180/63809280-98bf6400-c92a-11e9-960f-8c50257babdd.gif)

```
pressed keys: down, enter, backspace
```

### Search, Go-to Package

![Search, Go-to Package](https://user-images.githubusercontent.com/24392180/63809733-c35dec80-c92b-11e9-9a99-09317741a86c.gif)

```
pressed keys: s, (type), enter, g, (type), enter
```

### Install, Upgrade, Remove Package

![Install, Upgrade, Remove Package](https://user-images.githubusercontent.com/24392180/63811379-f3a78a00-c92f-11e9-9551-430d2437b69c.gif)

```
pressed keys:
i, (type), enter, y -> install
ctrl-u, (type), enter, y -> upgrade
ctrl-r, (type), enter, y -> remove
```

### Show Disk Usage Information

![Show Disk Usage Information](https://user-images.githubusercontent.com/24392180/63811686-d9ba7700-c930-11e9-9067-b0e412b5797f.gif)

```
pressed keys: right, left
```

### Confirm Command to Execute

![Confirm Command to Execute](https://user-images.githubusercontent.com/24392180/63812019-be03a080-c931-11e9-9732-de8bdcf75204.gif)

```
pressed keys: c, y
```

### Show Help

![Show Help](https://user-images.githubusercontent.com/24392180/63812128-15a20c00-c932-11e9-8ffd-7e222c78b588.gif)

```
pressed key: ?
```

## Screenshots

![Fedora Screenshot](https://user-images.githubusercontent.com/24392180/63807819-2ef18b00-c927-11e9-85b6-59917283a4f8.png)

![Manjaro-Mint Screenshot](https://user-images.githubusercontent.com/24392180/63795183-158f1580-c90c-11e9-8343-2dc24798c086.jpg)

![Debian-Ubuntu Screenshot](https://user-images.githubusercontent.com/24392180/63795189-17f16f80-c90c-11e9-96cc-dcd9bb660efe.jpg)

## Todo(s)
*  Add 'paste' feature

## License

GNU General Public License v3. (see [gpl](https://www.gnu.org/licenses/gpl.txt))

## Credit

Copyright (C) 2019 by [keylo99](https://www.github.com/keylo99)