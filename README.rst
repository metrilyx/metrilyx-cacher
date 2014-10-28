===============
Metrilyx Cacher
===============
Caches OpenTSDB metadata so full regex searches could be supported.

Currently the following data is collected:

* Metrics
* Tag keys
* Tag values


Build
-----::

	$ mkdir -p MetrilyxCacher/{src,bin,pkg}
	$ cd MetrilyxCacher
	$ export GOPATH=$(pwd)

	$ go get github.com/euforia/metrilyx-cacher
	$ go install github.com/euforia/metrilyx-cacher

This will produce a binary under **bin/metrilyx-cacher**

Installation
------------
Copy the generated binary to the desired location.