.text
.global main
.LC0:
  .string "%d\n"
.LC1:
  .string "%s\n"
# Function AddStrings
# Source: customtest.lat:1:1
_AddStrings:
    push %rbp
    mov %rsp,%rbp
    sub $0x18,%rsp
  AddStrings_block2:
    mov %rdi,-0x8(%rbp)
    mov %rsi,-0x10(%rbp)
    call strlen
    mov %eax,%r15d
    mov -0x8(%rbp),%rdi
    mov -0x10(%rbp),%rsi
    mov %r15,-0x8(%rbp)
    mov %rdi,-0x10(%rbp)
    mov %rsi,-0x18(%rbp)
    xchg %esi,%edi
    call strlen
    mov %eax,%ecx
    mov -0x8(%rbp),%r15
    mov -0x10(%rbp),%rdi
    mov -0x18(%rbp),%rsi # Const int 1
    mov $0x1,%r9d
    add %r9d,%ecx
    add %ecx,%r15d
    mov %rsi,-0x8(%rbp)
    mov %rdi,-0x10(%rbp)
    xchg %r15d,%edi
    call malloc
    mov %eax,%r12d
    mov -0x8(%rbp),%rsi
    mov -0x10(%rbp),%rdi
    mov %r12,-0x8(%rbp)
    mov %rsi,-0x10(%rbp)
    xchg %r12d,%edi
    xchg %r12d,%esi
    call strcpy
    mov -0x8(%rbp),%r12
    mov -0x10(%rbp),%rsi
    mov %r12,-0x8(%rbp)
    xchg %r12d,%edi
    call strcat
    mov -0x8(%rbp),%r12
    mov %r12d,%eax
    leave
    ret
# End of function AddStrings
# Function printInt
# Source: customtest.lat:10:1
_printInt:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printInt_block2: # Const string "%d\n"
      mov $.LC0,%ecx
      mov %rax,-0x8(%rbp)
      xchg %ecx,%edi
      xchg %ecx,%esi
      mov $0x0,%eax
      call printf
      mov -0x8(%rbp),%rax
      mov $0x0,%eax
      leave
      ret
# End of function printInt
# Function printString
# Source: customtest.lat:15:1
_printString:
        push %rbp
        mov %rsp,%rbp
        sub $0x8,%rsp
      printString_block2: # Const string "%s\n"
        mov $.LC1,%ecx
        mov %rax,-0x8(%rbp)
        xchg %ecx,%edi
        xchg %ecx,%esi
        mov $0x0,%eax
        call printf
        mov -0x8(%rbp),%rax
        mov $0x0,%eax
        leave
        ret
# End of function printString
# Function main (Entrypoint)
# Source: customtest.lat:20:1
main:
        main_block5: # Const int 10
          mov $0xa,%esi
          xchg %esi,%edi
          call _rfac
          mov %eax,%r10d
          xchg %r10d,%edi
          call _printInt # Const int 0
          mov $0x0,%edx
          mov %edx,%eax
          mov $0x1,%ebx
          xchg %eax,%ebx
          int $0x80
          ret
# End of function main
# Function rfac
# Source: customtest.lat:25:1
_rfac:
          rfac_block3: # Const int 0
            mov $0x0,%r9d
            cmp %r9d,%edi
            sete %r14b
            movzbl %r14b,%r14d # If condition
            cmp $0x0,%r14d
            jne rfac_block6
          rfac_block2: # Const int 1
            mov $0x1,%esi
            mov %esi,%eax
            ret
          rfac_block6: # Const int 1
            mov $0x1,%r10d
            mov %edi,%esi
            sub %r10d,%esi
            xchg %esi,%edi
            call _rfac
            mov %eax,%r12d
            imul %r12d,%edi
            mov %edi,%eax
            ret
# End of function rfac