# Local Simulation
This module is used to perform a simulation of the Alea-BFT consensus protocol. To launch the 4 nodes, open 4 consoles and run
```console
go run main.go launch -n N
```
with N = 1,2,3,4.  

The four nodes, after initalization, promptly begins the agreement round. Also, each create proposals for each self periodically (3 sec).

The protocol stops after the first proposal is delivered. The batch size is set to 1 and, thus, one proposal delivered is enoughp to stop the protocol.

To manually send a proposal to node N, run
```console
go run maing.go propose -n N
```

Otherwise, you can run
```console
./script.sh
```
to run a script that launches 4 consoles with the nodes.