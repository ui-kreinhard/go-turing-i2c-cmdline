# go-turing-i2c-cmdline
## What is it?
Controlling the i2c management bus of the turing pi with i2c works fine. But it's not easy to remember the correct bit mask for the node. So this commandline easens a little the pain :-)


It's just a wrapper around i2c. 
## Enable i2c management bus access

see: https://docs.turingpi.com/cluster-management-bus-i2c

## Build

```
go build
```

## How to use

When you run this cmdline from a turing pi node you need to 'lock' this node. Without it you can easily turn off the node where you're running the cmdline(and lock you out)

Node numbering starts with 1

```
echo "NODENUMBER" > /etc/lockedNodes
```
e.g. ignore 1 and 4 for commands

```
echo "1,4" > /etc/lockedNodes
```

```
# turn off node 2
go-turing-i2c-cmdline turnOff 2
# turn on node 3
go-turing-i2c-cmdline turnOn 3

# turn off all nodes(except locked nodes)
go-turing-i2c-cmdline turnOffAll
# turn on all nodes(except locked nodes)
go-turing-i2c-cmdline turnOnAll

```
## State
I've just coded it in 2h and just tested it. 

! Use with caution :) !

# Resources

https://docs.turingpi.com/cluster-management-bus-i2c