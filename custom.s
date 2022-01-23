.text
.global main
.LC3:
  .string "rawsx"
.LC4:
  .string "lol"
.LC0:
  .string "Error: %s\n"
.LC1:
  .string "%s\n"
.LC2:
  .string "%d\n"
# Function AddStrings
# Source: custom.lat:67:1
_AddStrings:
    push %rbp
    mov %rsp,%rbp
  AddStrings_block7:
    push %rdi
    push %rsi
    call strlen
    mov %eax,%r9d
    pop %rsi
    pop %rdi
    push %rdi
    push %rsi
    push %r9
    xchg %esi,%edi
    call strlen
    mov %eax,%edx
    pop %r9
    pop %rsi
    pop %rdi # Const int 1
    mov $0x1,%r14d
    add %r14d,%edx
    add %edx,%r9d
    push %rsi
    push %r9
    push %rdx
    push %r14
    push %rdi
    xchg %r9d,%edi
    call malloc
    mov %eax,%r8d
    pop %rdi
    pop %r14
    pop %rdx
    pop %r9
    pop %rsi
    push %r10
    push %rdi
    push %rsi
    push %r9
    push %rdx
    push %r14
    push %r8
    xchg %r8d,%edi
    xchg %r8d,%esi
    call strcpy
    pop %r8
    pop %r14
    pop %rdx
    pop %r9
    pop %rsi
    pop %rdi
    pop %r10
    push %rdi
    push %rsi
    push %r9
    push %rdx
    push %r14
    push %r8
    push %r10
    xchg %r8d,%edi
    call strcat
    pop %r10
    pop %r8
    pop %r14
    pop %rdx
    pop %r9
    pop %rsi
    pop %rdi
    mov %r8d,%eax
    leave
    ret
# End of function AddStrings
# Function error
# Source: custom.lat:76:1
_error:
      push %rbp
      mov %rsp,%rbp
    error_block2: # Const string "Error: %s\n"
      mov $.LC0,%r10d
      push %rax
      push %rdi
      push %r10
      push %rsi
      xchg %r10d,%edi
      xchg %r10d,%esi
      mov $0x0,%eax
      call printf
      pop %rsi
      pop %r10
      pop %rdi
      pop %rax # Const int 1
      mov $0x1,%r10d
      push %rcx
      push %rdi
      push %rsi
      push %r10
      xchg %r10d,%edi
      call exit
      pop %r10
      pop %rsi
      pop %rdi
      pop %rcx
      mov $0x0,%eax
      leave
      ret
# End of function error
# Function printString
# Source: custom.lat:82:1
_printString:
        push %rbp
        mov %rsp,%rbp
      printString_block3: # Const string "%s\n"
        mov $.LC1,%r9d
        push %rax
        push %rdi
        push %r9
        push %r15
        xchg %r9d,%edi
        xchg %r9d,%esi
        mov $0x0,%eax
        call printf
        pop %r15
        pop %r9
        pop %rdi
        pop %rax
        mov $0x0,%eax
        leave
        ret
# End of function printString
# Function printInt
# Source: custom.lat:87:1
_printInt:
          push %rbp
          mov %rsp,%rbp
        printInt_block2: # Const string "%d\n"
          mov $.LC2,%esi
          push %rax
          push %r10
          push %rdi
          push %rsi
          xchg %esi,%edi
          mov $0x0,%eax
          call printf
          pop %rsi
          pop %rdi
          pop %r10
          pop %rax
          mov $0x0,%eax
          leave
          ret
# End of function printInt
# Function main (Entrypoint)
# Source: custom.lat:92:1
main:
            push %rbp
            mov %rsp,%rbp
          main_block3: # Const string "rawsx"
            mov $.LC3,%edx # Const string "lol"
            mov $.LC4,%r13d
            push %rdx
            push %r13
            xchg %r13d,%edi
            xchg %edx,%esi
            call _AddStrings
            mov %eax,%r15d
            pop %r13
            pop %rdx
            push %rdx
            push %r13
            push %r15
            xchg %edx,%edi
            xchg %r15d,%esi
            call _AddStrings
            mov %eax,%r12d
            pop %r15
            pop %r13
            pop %rdx
            push %rdx
            push %r13
            push %r15
            push %r12
            push %rsi
            xchg %r12d,%edi
            call _printString
            pop %rsi
            pop %r12
            pop %r15
            pop %r13
            pop %rdx # Const int 15
            mov $0xf,%r10d
            mov %r10d,%eax
            mov $0x1,%ebx
            xchg %eax,%ebx
            int $0x80
            ret
# End of function main