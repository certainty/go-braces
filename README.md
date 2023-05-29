[![Build Status](https://github.com/certainty/go-braces/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/certainty/go-braces/actions/workflows/ci.yml)

# Go-Braces
A comprehensive compiler and VM implementation for programming language enthusiasts.

## Project state

While I'm making good progress most of this is in very very early development and many of the features only exist in my head.
I totally expect this to take maybe a year till this is somewhere in a state I can show it, at the current rate of development.
Well I guess that's just it. If you still feel interested or even want to contribute, please get in contact. I'm happy to
nerd out about this :) 

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


## Language deviations

* No first class continuations 
* Gradual type system

## Getting Started

The first thing you'll have to do currently is building everything.

### Building the Project

Build the entire project with a simple invocation of make:

```
make
```

#### Testing

Again testing is simple and you can run the harness with:

```
make test
```

### Introspecting a compiler run

The introspector and the compiler not packed into one executable. This is so that you can attach to arbitrary 
compilations including the ones taking place during interaction with the repl. And this is in fact also
the easiest (at the time of this writing only way to do it).

You will need two terminals, either one separate instancs or if you're using tmux two panes would also do.
In one pane you're going to start the repl like so:

```
./target/braces-vm repl -c
```

This will open up the repl and it will also start it in introspection mode.
You'll be greeted by a little banner saying that it's waiting for a client to connect.
So let's do that now. In a separate terminal or pane, execute:

```
./target/braces-introspect compiler 
```

This starts the introspector for the compiler. If everything went smoothly, the repl should have gotton out of 
the waiting state an dropped you into the prompt. Here you can interact with the repl now as usual, but what you
will observe is that the introspector picks up the compilation events and gives you insights into what's happening.

By default the compiler is in single-stepping mode, meaning that it will stop at strategic points during compilation
and allow you to interact with it and see its state. You will see that also in the introspector.

## References

This implementation is heavily based on an earlier version written in rust, which can be found here: https://github.com/certainty/braces
