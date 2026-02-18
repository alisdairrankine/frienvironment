# Capercaillie Stack Machine Spec

The Capercaillie is a stack machine operating on 8 bit memory cells.
Many operations have both 16 bit and 8 bit modes.

## Memory Mapping

Capercaillie machines have 64kb of addressable memory space.
Some areas are reserved for device mapping and stacks as described below:

[0x0000 .. 0x0001] - Address of entrypoint
[0x0002]           - Stack Pointer
[0x0003]           - Return Stack Pointer
[0x0004]           - System Status Flags
[0x0005 .. 0x00FF] - Reserved
[0x0100 .. 0x01FF] - Stack
[0x0200 .. 0x02FF] - Return Stack
[0x0300 .. 0x03FF] - Device Mapping
[0x0400 .. 0xFFFF] - Working Memory

### System status flags
The system status byte represents many statuses of the cpu 
[0x0004] -> [7 6 5 4 3 2 1 0]
0 (LSB) - system fault
1 - waiting
2 - stack overflow
3 - stack underflow
4 - return stack overflow
5 - return stack underflow
6 - divide by 0
7 (MSB) - unused

When a fault is experienced by the machine, execution halts, and the system fault flag is set. other flags will be set to determine the cause of the fault.

When the machine is in a waiting state (triggered bu the yield instruction), the waiting flag is set and execution is paused until an interrupt from a device.

### Stacks

Both stacks are empty stacks. Stack underflows will result in a fault system state, setting the relevant overflow/underflow flags
Both stack pointers point to the next possible value, initialised to 0.


## Devices

Capercaillie supports up to 16 connected devices. Each device has 16 addressable bytes.
The first byte of each device block is a device type identifier. Pointers are 1 byte offsets as thestacks are only 256 cells in length

Some devices support interupts which can trigger a callback specified in the device mapped memory.
If the system is in a waiting state when an interrupt is triggered, execution will resume.
Interrupts are pre-emptive and will push the current PC into the return stack before jumping to the callback.

Device Types:

0x00 - No device connected
0x01 - Terminal
0x02 - System
0x03 - Clock

Specs TBD.
Other device possibilities: graphics, cryptography, networking, file system, keyboard, mouse etc.

## Instructions

If stack operations attempt to add more or less bytes than available in either stack, a relevant overflow/underflow flag will be set, resulting in a system fault.

0x00 YIELD  - ( -- )                      "Wait for next callback, sets waiting flag to true"
0x01 HALT   - ( -- )                      "Halt machine, fault flag set to 0"
0x02 DUP    - ( a -- a a )
0x03 DROP   - ( a -- )
0x04 SWAP   - ( a b -- b a )
0x05 ROT    - ( a b c -- b c a )
0x06 OVER   - ( a b -- a b a )
0x07 NIP    - ( a b -- b )
0x08 TUCK   - ( a b -- b a b )
0x09 TOR    - ( a b -- ) r( -- a b )      "push two bytes to return stack"
0x0A FROMR  - ( -- a b ) r( a b -- )      "pop two bytes from return stack"
0x0B FETCHR - ( -- a b ) r( a b -- a b )  "fetch two bytes to return stack"


0x10 ADD   - ( a b -- c )        c = a + b
0x11 ADD16 - ( a b c d -- e f )  (e f) = (a b) + (c d)
0x12 SUB   - ( a b -- c )        c = b - a
0x13 SUB16 - ( a b c d -- e f )  (e f) = (c d) - (a b)
0x14 MUL   - ( a b -- c )        c = a * b
0x15 MUL16 - ( a b c d -- e f )  (e f) = (a b) * (c d)
0x16 DIV   - ( a b -- c )        c = b / a              "dividing by 0 result in a fault system state, setting the divide by zero flag. remainders are dropped"
0x17 DIV16 - ( a b c d -- e f )  (e f) = (c d) / (a b)  "dividing by 0 result in a fault system state, setting the divide by zero flag. remainders are dropped"
0x18 MOD   - ( a b -- c )        c = b % a              "dividing by 0 result in a fault system state, setting the divide by zero flag. remainders are dropped"
0x19 MOD16 - ( a b c d -- e f )  (e f) = (c d) % (a b)  "dividing by 0 result in a fault system state, setting the divide by zero flag. remainders are dropped"



0x20 AND   - ( a b -- c )        c = a & b
0x21 AND16 - ( a b c d -- e f )  (e f) = (a b) & (c d)
0x22 OR  -   ( a b -- c )        c = a | b
0x23 OR16  - ( a b c d -- e f )  (e f) = (a b) | (c d)
0x24 XOR   - ( a b -- c )        c = a ^ b
0x25 XOR16 - ( a b c d -- e f )  (e f) = (a b) ^ (c d)
0x26 NOT   - ( a -- b )          b  = ~a
0x27 NOT16 - ( a b -- c d )      (c d) = ~(a b)
0x28 INC   - ( a -- b )          b = a + 1
0x29 INC16 - ( a b -- c d )      (c d) = (a b) + 1
0x2A DEC   - ( a -- b )          b = a - 1
0x2B DEC16 - ( a b -- c d )      (c d) = (a b) - 1
0x2C SHL   - ( a b -- c )        c = a << b
0x2D SHL16 - ( a b c -- d e )    (d e) = (a b) << c
0x2E SHR   - ( a b -- c )        c = a >> b
0x2F SHR16  - ( a b c -- d e )    (d e) = (a b) >> c


0x40 JZ   - ( a b c -- )     "jump to location in ( a b ) if c is 0"
0x41 JNZ  - ( a b c -- )     "jump to location in ( a b ) if c is not 0"
0x42 CALL - ( a b -- )       "jump to location in ( a b ), push pc+1 to rp stack"
0x43 RET  - ( -- )           "pop 2 bytes from rp stack, jump to that location"
0x44 EQ   - ( a b -- c)      "push 1 if a == b, otherwise 0"
0x45 EQ16 - ( a b c d -- e)  "push 1 if (a b) == (c d), otherwise 0"
0x46 NQ   - ( a b -- c)      "push 1 if a != b, otherwise 0"
0x47 NQ16 - ( a b c d -- e)  "push 1 if (a b) != (c d), otherwise 0"
0x48 GT   - ( a b -- c)      "push 1 if b > a, otherwise 0"
0x49 GT16 - ( a b c d -- e)  "push 1 if (c d) > (a b), otherwise 0"
0x4A LT   - ( a b -- c)      "push 1 if b < a, otherwise 0"
0x4B LT16 - ( a b c d -- e)  "push 1 if (c d) < (a b), otherwise 0"

0x50 PUSH    - ( -- a )        "Push next byte instructions to stack, pc++"
0x51 PUSH16    - ( -- b a )    "Push next short (two bytes) in instructions to stack pc+=2"
0x52 STORE  - ( a b c -- )     addr = (b << 8 | c) MEM[addr] = a
0x53 STORE16 - ( a b c d -- )  addr = (c << 8 | d) MEM[addr] = a; MEM[addr+1] = b
0x54 LOAD   - ( a b -- c )     addr = (a << 8 | b) c = MEM[addr]
0x55 LOAD16  - ( a b -- c d )  addr = (a << 8 | b) c = MEM[addr]; d = MEM[addr+1]

