# How do we bring Filesystem (and Pipe) analogy to Collector?

## Unified Resource Model
The idea is that, similar to how Unix treats (everything as a file)[https://en.wikipedia.org/wiki/Everything_is_a_file], we should try to have everything as a resource. This will make so many things possible for us - piping operations (composability), standard input output, uniform interface, etc.

In a Unix system, everything is a file of a different type. Regular files, directories, devices, sockets, pipes, etc. are all files. Different file types have different behaviors (directory contins fixed-size records of (i-node, filename) but cannot be written to directly, etc), but share a common interface. And operations like read(), write(), open() work uniformly across all the types.