## Usage

This demonstrates a crash in go-tree-sitter, which I suspect is due to upstream tree-sitter.

When scanning the JSON tests from this Perl module, [Cpanel-JSON-XS-4.37](https://metacpan.org/dist/Cpanel-JSON-XS), I found it reliably crashes my Go program. I suspect this is an issue in tree-sitter-javascript.


## To recreate the bug:

```
go mod tidy
go run main.go -path ./crashers
```

It should reliably crash. At least it does on my M1 macbook. You should see:

```
[......truncated.....]
goroutine 34 [finalizer wait]:
runtime.gopark(0x10?, 0x104bed7e0?, 0x0?, 0x0?, 0x104bf5e80?)
        /opt/homebrew/Cellar/go/1.21.3/libexec/src/runtime/proc.go:398 +0xc8 fp=0x14000052580 sp=0x14000052560 pc=0x104b08f18
runtime.runfinq()
        /opt/homebrew/Cellar/go/1.21.3/libexec/src/runtime/mfinal.go:193 +0x108 fp=0x140000527d0 sp=0x14000052580 pc=0x104ae9f58
runtime.goexit()
        /opt/homebrew/Cellar/go/1.21.3/libexec/src/runtime/asm_arm64.s:1197 +0x4 fp=0x140000527d0 sp=0x140000527d0 pc=0x104b35e64
created by runtime.createfing in goroutine 1
        /opt/homebrew/Cellar/go/1.21.3/libexec/src/runtime/mfinal.go:163 +0x80

r0      0x104d90680
r1      0x104dbb200
r2      0x170268000
r3      0x0
r4      0x10
r5      0x1
r6      0x2
r7      0x0
r8      0xffffffff
r9      0x100
r10     0xffffffff
r11     0x3e
r12     0x104dbb210
r13     0x104dbb218
r14     0x0
r15     0x170000000
r16     0x170268000
r17     0x1a41ef1e0
r18     0x0
r19     0x104d90680
r20     0x104dbb200
r21     0x4
r22     0x170000000
r23     0x0
r24     0x4c
r25     0x170268000
r26     0x0
r27     0x1ff442160
r28     0x170268000
r29     0x16b32a300
```

## To demonstrate it used to work
To see it not crash, open `go.mod` and un-comment the last line with comment of "does NOT crash":
```
// require github.com/smacker/go-tree-sitter v0.0.0-20230720070738-0d0a9f78d8f8 //crashes
// require github.com/smacker/go-tree-sitter v0.0.0-20221023091341-2009a4db91e4 //crashes
require github.com/smacker/go-tree-sitter v0.0.0-20220829074436-0a7a807924f2 //does NOT crash
```

Then re-run: 

```
go mod tidy
go run main.go -path ./crashers
```
