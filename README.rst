===============
Metrilyx Cacher
===============

Caches OpenTSDB metadata so full regex searches could be supported.

Currently the following data is collected:

    * Metrics
    * Tag keys
    * Tag values


Usage
-----
::
    
    metrilyx-cacher -url http://my.opentsdb.inst/api/suggest


By default the startup script looks for the metrilyx.conf

TODO
----

Cache dashboard listing.