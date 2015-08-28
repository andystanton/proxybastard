# Proxy Bastard [![Build Status](https://travis-ci.org/andystanton/proxybastard.svg?branch=master)](https://travis-ci.org/andystanton/proxybastard)[ ![Download](https://api.bintray.com/packages/andystanton/generic/proxybastard/images/download.svg) ](https://bintray.com/andystanton/generic/proxybastard/_latestVersion)

> **bastard /ˈbɑːstəd; ˈbæs-/ (noun)**

> (informal) something extremely difficult or unpleasant: 'that job is a real bastard'

> **proxy /ˈprɒksi/ (noun)**

> see bastard

A command line interface for enabling and disabling proxy settings in the shell and other applications.

**Warning:** this tool is for people who know what they are doing but are too lazy to do it themselves.

## Usage

```sh
$ proxybastard on|off
```

The first time you run Proxy Bastard you'll be presented with some configuration options. Proxy Bastard scans predefined locations for existing proxy settings and applications that utilise proxy settings, and suggests a configuration based on this. Once you've reviewed the suggestions, it will write a file called ```config.json``` in the folder ```~/.proxybastard``` and you're ready to go.

Proxy settings can then be applied to your current shell session either by sourcing your shell profile/rc or running ```$(proxybastard env)```.

You can re-run the initial configuration at any time with ```proxybastard setup``` or by manually editing ```~/.proxybastard/config.json```.

## Installation

### via Brew

```sh
$ brew install andystanton/formula/proxybastard
```

### via Go

You can clone the repository and build from source locally. This approach assumes a working installation of Go including a valid ```GOPATH``` environment variable and ```$GOPATH/bin``` added to your path.

```sh
$ git clone https://github.com/andystanton/proxybastard.git
$ cd proxybastard
$ go get
$ go install
```

### Binary download

* Download the binary for your OS and architecture: https://bintray.com/andystanton/generic/proxybastard
* Unzip the file
* Copy ```proxybastard``` to ```/usr/local/bin```

## Supported applications

* Atom Package Manager
* Boot2Docker
* Docker Machine
* Git
* Maven
* NPM
* Shell profile/rc
* SSH
* Stunnel
* Subversion

See [CONTRIBUTING.md](CONTRIBUTING.md) for how you can contribute other configurations.
