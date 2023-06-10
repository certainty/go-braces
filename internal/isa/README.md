## Instruction Set Architecture (ISA) 

This term is often used in the context of hardware processors, but it can also be applied to virtual machines. 
The ISA defines the set of instructions, their encoding, and the behavior of the VM when executing them.
The compiler generates bytecode according to the ISA, and the VM executes that bytecode based on the same specification.

We will eventually encode the assembly using protobuf. So our protobuf definitions will essentialy become the description of instructions and their encoding.
