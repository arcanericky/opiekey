
# OPIEKey

[![Build Status](https://travis-ci.com/arcanericky/opiekey.svg?branch=master)](https://travis-ci.com/arcanericky/opiekey)
[![codecov](https://codecov.io/gh/arcanericky/opiekey/branch/master/graph/badge.svg)](https://codecov.io/gh/arcanericky/opiekey)
[![GoDoc](https://img.shields.io/badge/docs-GoDoc-brightgreen.svg)](https://godoc.org/github.com/arcanericky/opiekey)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

A Go package and command-line interface to generate [OPIE (One-time Passwords In Everything)](https://en.wikipedia.org/wiki/OPIE_Authentication_System) challenge responses.

`opiekey`, but written in Go! _S/KEY authentication like it's 1996_

### Command Line Usage
```
$ opiekey 499 testseed testpassphrase
LAIR FUME GIBE FROM JIG COP

$ opiekey 499 testseed
Using the MD5 algorithm to compute response.
Reminder: Don't use opiekey from telnet or dial-in sessions.
Enter secret pass phrase: 
LAIR FUME GIBE FROM JIG COP

$ OPIE_PASSPHRASE="testpassphrase" opiekey 499 testseed
LAIR FUME GIBE FROM JIG COP

$ opiekey --version
opiekey version 1.0.0 linux/amd64
```

### Package Usage
```
package main

import "fmt"
import "github.com/arcanericky/opiekey"

func main() {
	fmt.Println(opiekey.ComputeWordResponse(499, "testseed" "testpassphrase", opiekey.MD5))
}
```

### Help
```
$ opiekey --help
opiekey - Program for computing responses to OTP challenges.

opiekey takes the optional count of the number of responses to print
along with a (maximum) sequence number, seed and optional secret pass
phrase as command line args then produces an OPIE response as six
words or hexadecimal numbers. If the OPIE_PASSPHRASE environment
variable is set it will be used for the secret pass phrase. If no
secret pass phrase was specified as a command line argument or
environment variable, the program will prompt for it.

Usage:
  opiekey sequence_number seed passphrase [flags]

Flags:
  -h, --help         help for opiekey
  -x, --hex          output the OTPs as hexadecimal numbers instead of six words
  -4, --md4          selects MD4 as the response generation algorithm
  -5, --md5          selects MD5 as the response generation algorithm (default true)
  -n, --number int   the number of one time access passwords to print (default 1)
  -s, --sha1         selects SHA1 as the response generation algorithm
      --version      version for opiekey
```

### Warnings
A particular OPIE implementation might limit the value of the sequence number and the characters and lengths of the seed and passphrase. Neither this utility or package validate this data.

The utility and package support the [MD5](https://en.wikipedia.org/wiki/MD5) and [MD4](https://en.wikipedia.org/wiki/MD4) algorithms.

The [SHA1](https://en.wikipedia.org/wiki/SHA-1) implementation is questionable and probably doesn't produce correct output.

### Credits
The output of this `opiekey` utility was tested against the output of the now obsolete [opie-client 2.40 Ubuntu package](https://launchpad.net/ubuntu/jaunty/amd64/opie-client/2.40~dfsg-0ubuntu1). The word list was lifted from the integer-word translation dictionary in the `btoe.c` module of the same package and documented as part of [RFC-2289](https://tools.ietf.org/html/rfc2289#page-4-19). Most of my comprehension of the OPIE algorithm came from [ruby-otp](https://github.com/moumar/ruby-otp) so this translation may not be ideal, but it produces the results I require.

### Inspiration
My day job has deployed a few machines that require responses to OPIE challenges for logins. I very rarely need to login to these machines, but when I do, it's a pain to find an opiekey utility, [mobile device app](https://m.apkpure.com/opiekey/nl.f00d.android.opiekey), or [web page](https://math.berkeley.edu/~vojta/opiekey.html) to generate these responses. I finally decided to code up a small library and utility that can run on most any OS I'm using. This very basic package and utility is the result.

### Other Resources
The OPIE Authentication System is a dinosaur and its usage is rapidly dwindling. Listed below are a few useful resources, noting that I don't endorse, recommmend, or support any executables you dare to run:
* A copy of the [OPIE archive](http://www.tifosi.com/OTP/) including programs for various operating systems and a [README](http://www.tifosi.com/OTP/README)
* [FreeBSD's OPIE implementation](https://github.com/freebsd/freebsd/tree/master/contrib/opie)
* [OPIE host setup](http://seann.herdejurgen.com/resume/samag.com/html/v08/i13/a4.htm)
* [Source for OTPDroid](https://github.com/felixb/otpdroid)
* Old [OPIEKey Android app](https://m.apkpure.com/opiekey/nl.f00d.android.opiekey)
* Source for the defunct [OPIE 2.40 for Ubuntu](https://launchpad.net/ubuntu/+source/opie/2.40~dfsg-0ubuntu1)