# Local Process Simulation

The main.go file is a script to launch 4 operators that communicate through the process itself (i.e. there is not communication using any network, even localhost:port).
The purpose of this script is to test the behavior of the consensus protocol. and to get some metrics that exclude the communication time but increase processing time (since they all run in the same computer).

To run the script, execute the following:
```
go run main.go
```

The proposals to each client is created periodically in the own alea/instance.go file. The execution ends when the node reaches 40 delivered proposals.