# Sisyphus: Intelligent Junk Mail Handler
As we all know too well, many mails we receive are undesired for various
reasons. Sometimes, we just do not want to be part of a scam, sometimes we
really prefer no to have this latest joke mail sent by our beloved friends --
even though we are happy to exchange serious messages with them.

**Sisyphus** is a junk mail handler of the latest generation. It has the
following features:

* requires zero configuration, neither on the server nor on the client
* works with any MTA and any client
* learns about your preferences based on all messages in your inbox and your
  junk folder
* can handle multiple mail accounts with independant junk mail preferences
* requires minimal resources, e.g. learning over 50000 mails and keeping track of roughly 90000 words requires only 10MB of storage

[![Build Status](https://travis-ci.org/carlostrub/sisyphus.svg?branch=master)](https://travis-ci.org/carlostrub/sisyphus)
[![Go Report Card](https://goreportcard.com/badge/github.com/carlostrub/sisyphus)](https://goreportcard.com/report/github.com/carlostrub/sisyphus)
[![GoDoc](https://godoc.org/github.com/carlostrub/sisyphus?status.svg)](https://godoc.org/github.com/carlostrub/sisyphus)
[![Documentation](https://readthedocs.org/projects/sisyphus/badge/?version=latest)](http://sisyphus.readthedocs.org/en/latest/?badge=latest)
[![Codebeat](https://codebeat.co/badges/64615809-e3c4-4267-a049-eaec20ad63b5)](https://codebeat.co/projects/github-com-carlostrub-sisyphus-master)
[![Coverage](https://gocover.io/_badge/github.com/carlostrub/sisyphus?0 "Coverage")](http://gocover.io/github.com/carlostrub/sisyphus)

## How it works
Sisyphus analyzes each mail in the inbox and the junk folder and uses its
subject and text to improve the learning algorithm. Whenever a new mail arrives
in the `Maildir/new` directory, Sisyphus classifies this mail based on its
content. Junk mails are then moved automatically to the `Maildir/.Junk`
directory, while good mails are left untouched. See the following [blog
post](http://carlostrub.ch/code/security/sisyphus/) on a rather non-technical
explanation.

Technically, Sisyphus applies a classic [Bayesian Update
algorithm](https://en.wikipedia.org/wiki/Bayesian_inference) to classify mails.
However, in contrast to many traditional junk mail filters, classification is
based on all mails ever received. This includes mails that are classified by
the user as junk by moving them manually into the junk folder, or mails that
have been correctly classified by Sisyphus previously. This is only possible
with limited resources by applying the [HyperLogLog
algorithm](https://en.wikipedia.org/wiki/HyperLogLog) to store the learned
mails.

The learned information is stored in a local database called `sisyphus.db`
which is located in each `Maildir` directory.

## Install
Sisyphus can be installed by downloading the released [binary
package.](https://github.com/carlostrub/sisyphus/releases)

To build from source, you can
1. Clone this repository into `$GOPATH/src/github.com/carlostrub/sisyphus` and
   change directory into it
2. Run `make build`

This will leave you with `./sisyphus` in the `sisyphus` directory, which you
can put in your `$PATH`. (You can also take a look at `make install` to install
for you.)

## Usage
First, set the environment variable necessary for operation:
```
$ setenv SISYPHUS_DIRS PATHTOMAILDIR
```
or
```
$ export SISYPHUS_DIRS=PATHTOMAILDIR
```
or for Windows
```
$ set SISYPHUS_DIRS=PATHTOMAILDIR
```

For all other configuration options, please consult the help. It can
be started by running
```
$ sisyphus help
```

To start sisyphus, do
```
$ sisyphus run
```

To display various statistics, do
```
$ sisyphus stats
```
(caveat: run at least one learning cycle)

See the help for more details.

## License
Sisyphus is licensed under the 3-Clause BSD license. See the LICENSE file for
detailed information.
