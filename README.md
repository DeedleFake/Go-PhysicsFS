Go-PhysicsFS
======

This Go package provides bindings for the [PhysicsFS][physfs] library. PhysicsFS is a library to provide abstract access to various archives. It provides an easy method to access files in various locations, even inside archives, without having to worry about their actual locations, similar to Quake 3's file subsystem.

Prerequisites
-------------

 * [PhysicsFS][physfs].
 * A recent version of [Go][go]. (As of 2012-02-06)

Installation
------------

To install simply type:

> go get github.com/DeedleFake/Go-PhysicsFS/physfs

This will install into your GOPATH. For more information, type:

> go help importpath

and

> go help get

Usage
-----

To import, use the following:

    import "github.com/DeedleFake/Go-PhysicsFS/physfs"

Docs
----

To view the documentation locally, type:

> go doc github.com/DeedleFake/Go-PhysicsFS/physfs

If you would like to see the docs in a nice layout online, simply visit [GoPkgDoc][gpd].

Authors
-------

 * [DeedleFake](https://github.com/DeedleFake)

[physfs]: http://www.icculus.org/physfs
[go]: http://www.golang.org
[gpd]: http://gopkgdoc.appspot.com/pkg/github.com/DeedleFake/Go-PhysicsFS/physfs

<!--
    vim:ts=4 sw=4 et
-->
