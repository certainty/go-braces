package isa

// Values are not actullay part of the ISA normally
// and this file will go away eventually, but in the
// beginning it's very helpful to take some shortcuts
// during development. The compiler doesn't need to care too
// much about layout and encoding of the values in the constant pool
// for example.

type Value interface{}
