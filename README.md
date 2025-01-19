# Functions functions

- Internal packages have no dependencies, they just define their ports as functions which will be filled at main.
- Ports are function signatures, adapters are functions with those signatures.
- No explicit types on callee side, only callee and caller are dependent on function signature level.

- Package naming or folder structure is not important, idea is code structure and how layers use each others.
- You can achieve feature based folder structure without circular dependency since every layer only can have concrete dependency to business modal.

## Idea - Composition

```
fn1 = fn(a, b)
fn2 = fn(a, b, f(a, b))

f1 = (a, b) => a + b + 2
f2 = (a, b) => a * b - 2
f3 = (a, b, fn1) => a + b + fn1(a, b)
f4 = (a, b, fn1) => a * b + fn1(a, b)
```

If we call with `a = 1, b = 2, fn1 = f1`

```
f3(1, 2, f1) => 1 + 2 + f1(1, 2)
             => 1 + 2 + (1 + 2 + 2)
             => 8
```

If we call with `a = 1, b = 2, fn1 = f2`

```
f3(1, 2, f2) => 1 + 2 + f2(1, 2)
             => 1 + 2 + (1 * 2 - 2)
             => 3
```

We can use whatever implements fn1 as parameter
We can use it on any parameter accepts fn1 signature

```
f5 = (a, b) => 0 && do db call

f3(1, 2, f5) => 1 + 2 + 0 && do db call
f4(1, 2, f5) => 1 * 2 + 0 2 && do db call
```

## Start

```sh
make prepare-db
go run .
```

Swagger in
http://127.0.0.1:8000/docs
