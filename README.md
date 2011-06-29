Go-PhysicsFS
======

This Go package provides bindings for the [PhysicsFS][physfs] library. PhysicsFS is a library to provide abstract access to various archives. It provides an easy method to access files in various locations, even inside archives, without having to worry about their actual locations, similar to Quake 3's file subsystem.

Prerequisites
-------------

 * [PhysicsFS][physfs].
 * A recent version of [Go][go]. (As of 2011-04-27)

Installation
------------

To install simply type:

> make install

If you don't have write permission for GOROOT, you may need to run the previous command as root. If you get errors while trying to run it using sudo, it's possible that the GOROOT/GOOS/GOARCH/GOBIN variables are not available to the make command. You can try using '-E' to preserve the environment:

> sudo -E make install

Usage:
------

To import, simply use the following:

    import "physfs"

Authors
-------

 * [DeedleFake](https://github.com/DeedleFake)

[physfs]: http://www.icculus.org/physfs/
[go]: http://www.golang.org

<!--
    vim:ts=4 sw=4 et
-->
