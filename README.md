# Go-Braces
A comprehensive compiler and VM implementation for programming language enthusiasts.

## Overview

Go-Braces is a dialect of the Scheme programming language, offering a compiler and virtual machine (VM) designed as a learning platform for individuals interested in building programming languages and compilers.

The primary objectives of the tools provided here are:

* Introspection - Enable visibility into intricate details for better understanding.
* Real-world implementation - Avoid using toy versions, and handle real-world scenarios.

## Objectives

* Develop a VM that is reasonably performant and provides runtime introspection for a clear view of program execution.
  * Implement through a Text-based User Interface (TUI) for access to VM internals and interactive runtime system engagement.
* Create a modular compiler that:
  * Features a modular design and supports language extensions.
  * Allows for an optional type system.
  * Supports both ahead-of-time and just-in-time compilation.
  * Incorporates introspection capabilities to examine every phase of the compilation process in detail.
  * Supports macros with an explicit renaming macro expander.
* Offer a gradual type system that can be enabled/disabled via a compiler flag (possibly on a per-module basis).

## Non-Objectives

* Develop a fully compliant r7rs Scheme - Although it will be used as a foundation.
* Create the most efficient/fast/resource-optimized VM or compiler.
  * Strive for a balance between comprehensibility and performance.

## Getting Started

### Building the Project

The following will build all executables in each of the packages. 
The binaries will be placed in `$package/target/bin`. 

```
make
```

You can also build individual packages separately, by changing into the corresponding directory and then invoking `make`.

## References

This implementation is heavily based on an earlier version written in rust, which can be found here: https://github.com/certainty/braces
