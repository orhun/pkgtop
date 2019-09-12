![Logo](https://user-images.githubusercontent.com/24392180/63693894-dd110e00-c81d-11e9-8f51-e00d5bd7d6a6.png)

**pkgtop** is an **interactive package manager** & **resource monitor** tool designed for the GNU/Linux.

[![Release](https://img.shields.io/github/release/orhun/pkgtop.svg?style=flat-square)](https://github.com/orhun/pkgtop/releases)
[![AUR](https://img.shields.io/aur/version/pkgtop-git.svg?style=flat-square)](https://aur.archlinux.org/packages/pkgtop-git/)
[![Travis Build](https://img.shields.io/travis/orhun/pkgtop.svg?style=flat-square)](https://travis-ci.org/orhun/pkgtop) [![Docker Build](https://img.shields.io/docker/cloud/build/orhunp/pkgtop.svg?style=flat-square)](https://hub.docker.com/r/orhunp/pkgtop/builds) [![Codacy Badge](https://img.shields.io/codacy/grade/f83f3a6b0bb042f39f799cb372405094.svg?style=flat-square)](https://www.codacy.com/app/orhun/pkgtop?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=orhun/pkgtop&amp;utm_campaign=Badge_Grade) [![Go Report Card](https://goreportcard.com/badge/github.com/orhun/pkgtop?style=flat-square)](https://goreportcard.com/report/github.com/orhun/pkgtop) [![Stars](https://img.shields.io/github/stars/orhun/pkgtop.svg?style=flat-square)](https://github.com/orhun/pkgtop/stargazers) [![License](https://img.shields.io/github/license/orhun/pkgtop.svg?color=blue&style=flat-square)](./LICENSE)

![pkgtop](https://user-images.githubusercontent.com/24392180/63897168-edcbaa80-c9fc-11e9-9092-32a55323fcf1.gif)

Package management (install/upgrade/remove etc.) can be a problem if the user is not familiar with the operating system or the required command for that operation. So pkgtop tries to solve this problem with an easy-to-use terminal interface and shortcut keys. Briefly, **pkgtop aims to provide a terminal dashboard for managing packages on GNU/Linux systems.** Using the terminal dashboard, it's possible to list installed packages by size (or alphabetically with `-a` argument), show information about the package, install/upgrade/remove packages and search package. Also, there are other handy shortcuts for easing the package management process which mentioned in the [usage information](https://github.com/orhun/pkgtop#usage).

In addition to the package management features, there's a section at the top of the dashboard that shows disk usages and general system information. For example, this section can be used as a resource monitor and help decide whether the system should be cleaned or not.  
Another useful section is the '`executed`' or '`confirm to execute`' command list which is placed below the installed packages. Thus, the user can see which command executed recently or confirm & execute the selected command. (The commands that need confirmation to execute exist in the list with a prefix like "`[y]`".) 
After scrolling the commands list with "`c`" key for selecting the command to execute, press "`y`" for executing it. pkgtop will execute the command and restart the terminal dashboard afterwards.

pkgtop uses the advantage of mainly used package managers being installed on most of the preferred GNU/Linux distributions. As an example, it works on a [Manjaro](https://manjaro.org/) based system as it works on [Arch Linux](https://www.archlinux.org/) systems since both distributions use the [Pacman](https://wiki.archlinux.org/index.php/pacman) package manager. You can use pkgtop if you have one of the package managers listed below.

* pacman
* apt
* zypp
* dnf
* xbps
* portage
* nix
* guix

If the distribution is not defined in the source but has the required package manager for running the pkgtop, `-d` argument can be used for specifying a distribution that has the same package manager. Current defined and supported distributions are `arch, manjaro, debian, ubuntu, mint, suse, fedora, centos, redhat, void, gentoo, nixos, guix`.

  * [Installation](#installation)
     * [• Dependencies](#-dependencies)
     * [• AUR (<a href="https://aur.archlinux.org/packages/pkgtop-git" rel="nofollow">pkgtop-git</a>)](#-aur-pkgtop-git)
     * [• Manual Insallation](#-manual-insallation)
  * [Command-Line Arguments](#command-line-arguments)
  * [Usage](#usage)
     * [List Installed Packages &amp; Show Package Information](#list-installed-packages--show-package-information)
     * [Search, Go-to Package](#search-go-to-package)
     * [Install, Upgrade, Remove Package](#install-upgrade-remove-package)
     * [Show Disk Usage Information](#show-disk-usage-information)
     * [Confirm Command to Execute](#confirm-command-to-execute)
     * [Show Help](#show-help)
  * [Docker](#docker)
     * [Build Docker Image](#build-docker-image)
     * [Run the Container](#run-the-container)
     * [Start a shell in the Container](#start-a-shell-in-the-container)
  * [Screenshots](#screenshots)
  * [Todo(s)](#todos)
  * [Sponsor](#sponsor)
  * [License](#license)
  * [Credit](#credit)

## Installation

### • Dependencies
* [gizak/termui](https://github.com/gizak/termui/)
* [atotto/clipboard](https://github.com/atotto/clipboard)
* [dustin/go-humanize](https://github.com/dustin/go-humanize)
* [mattn/go-runewidth](https://github.com/mattn/go-runewidth)
* [mitchellh/go-wordwrap](https://github.com/mitchellh/go-wordwrap)
* [nsf/termbox-go](https://github.com/nsf/termbox-go)

### • AUR ([pkgtop-git](https://aur.archlinux.org/packages/pkgtop-git))

```
git clone https://aur.archlinux.org/pkgtop-git.git && cd pkgtop-git/
makepkg --install
```

### • Manual Installation

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
-c, main color of the dashboard (default: blue)
   [red, green, yellow, blue, magenta, cyan, white]
-a, sort packages alphabetically
-r, reverse the package list
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
| `s,/`                      	| search package                           	|
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

## Docker

### Build Docker Image

```
docker build -f docker/Dockerfile -t pkgtop-docker .
```

### Run the Container

```
docker run pkgtop-docker
```

### Start a shell in the Container

```
docker run -it pkgtop-docker /bin/ash
```

## Screenshots

![Fedora Screenshot](https://user-images.githubusercontent.com/24392180/63807819-2ef18b00-c927-11e9-85b6-59917283a4f8.png)

![Manjaro-Mint Screenshot](https://user-images.githubusercontent.com/24392180/63795183-158f1580-c90c-11e9-8343-2dc24798c086.jpg)

![Debian-Ubuntu Screenshot](https://user-images.githubusercontent.com/24392180/63795189-17f16f80-c90c-11e9-96cc-dcd9bb660efe.jpg)

## Todo(s)
* Add 'paste' feature

## Sponsor

If you would like to support the development of pkgtop or other projects on my [GitHub](https://github.com/orhun), consider becoming a sponsor. Very appreciated!

[![Become a Patron!](https://c5.patreon.com/external/logo/become_a_patron_button.png)](https://www.patreon.com/bePatron?u=23697306)

## License

GNU General Public License v3. (see [gpl](https://www.gnu.org/licenses/gpl.txt))

## Credit

Copyright (C) 2019 by [orhun](https://www.github.com/orhun)
