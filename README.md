# Edge

Implementation of the Liteseed Bundler Node in Go

## System Requirements

The node was tested on a system with the following specification:

- CPU - 1GHz
- Memory - 1GiB
- Storage - 8GiB
- OS - Linux

The requirements are expected to grow in the future.

## Getting Started

First make a directory to store the binary

```sh
  mkdir $HOME/.edge
  cd $HOME/.edge
```

Fetch the latest release of edge from [github.com/liteseed/edge/releases](https://github.com/liteseed/edge/releases).

```sh
  wget https://github.com/liteseed/edge/latest/download/edge-linux-386
```

Set permission to execute the binary

```sh
  chmod +777 edge-linux-386
```

Export`EDGE_PATH` and source your shell file

```sh
 echo 'export PATH=$HOME/.edge/edge:$PATH' >> ~/.bash_profile
```
