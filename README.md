# Proxy Bastard [![Build Status](https://travis-ci.org/andystanton/proxybastard.svg?branch=master)](https://travis-ci.org/andystanton/proxybastard)

> **bastard /ˈbɑːstəd; ˈbæs-/ (noun)**

> (informal) something extremely difficult or unpleasant: 'that job is a real bastard'

> **proxy /ˈprɒksi/ (noun)**

> see bastard

A simple command line interface for enabling and disabling proxy settings in shell environment and other application settings.

**Warning:** this tool is for people who know what they are doing but are too lazy to do it themselves.

## Usage

```sh
$ proxybastard on|off
```

Proxy settings can then be applied to your current shell session either by sourcing your shell profile/rc or running

```sh
$ $(proxybastard env)
```

## Supported configurations

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

See [CONTRIBUTING.md](CONTRIBUTING.md) for how you can add your favourite configuration.
