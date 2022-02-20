# Vite

[![Tests](https://github.com/vite-cloud/vite/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/vite-cloud/vite/actions/workflows/tests.yml)
[![Static](https://github.com/vite-cloud/vite/actions/workflows/static.yml/badge.svg)](https://github.com/vite-cloud/vite/actions/workflows/static.yml)
[![CodeBeat badge](https://codebeat.co/badges/7171e9ea-53d7-4c81-82bf-a9a2f222b027)](https://codebeat.co/projects/github-com-vite-cloud-vite-main)
[![Go Report Card](https://goreportcard.com/badge/github.com/vite-cloud/vite)](https://goreportcard.com/report/github.com/vite-cloud/vite)
[![codecov](https://codecov.io/gh/vite-cloud/vite/branch/main/graph/badge.svg?token=DWSP4O0YO8)](https://codecov.io/gh/vite-cloud/vite)
![PRs not welcome](https://img.shields.io/badge/PRs-not%20welcome-red)

#### Documentation Status

The goal is to write a lot and then eventually make it more concise and improve upon it.

**VERY MUCH WIP, JUST RANDOM THINGS**

## Requirements

* docker
* git

## What is Vite?

Vite is a tool to help you manage many applications (called "services" from now on) on a single server. You can think of
it as a supercharged docker-compose.

Features:

* zero downtime deployments
* built-in reverse proxy
* versioned configuration
* powerful configuration diagnosis (if anything looks wrong in your configuration, Vite will SCREAM LOUDLY)
* an api to trigger deployments automatically (CD [What's Continous Deployment (link needed)]() with a single api call)

## Why use Vite?

Vite is the perfect middle ground between messy configuration files all over your server and kubernetes.

## When not to use Vite?

* You have more than two servers

  If you have exactly two servers, you can still use Vite very effectively and make your architecture redundant by
  running them in a Active-Active configuration (or Active-Passive if one is less powerful)
  . [(link needed)]()
