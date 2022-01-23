.text
.global main
.LC0:
  .string "%s\n"
.LC1:
  .string "a"
# Function printStringX
# Source: custom.lat:66:1
_printStringX:
    push %rbp
    mov %rsp,%rbp
  printStringX_block5: # Const string "%s\n"
    mov $.LC0,%ecx
    # push %rcx
    # push %rdi
    xchg %ecx,%edi
    xchg %ecx,%esi
    movl $0, %eax
    call printf
    # pop %rcx
    # pop %rdi
    leave
    ret
# End of function printStringX
# Function main (Entrypoint)
# Source: custom.lat:71:1
main:
      push %rbp
      mov %rsp,%rbp
      subq    $48, %rsp
    main_block5: # Const string "a"
      mov $.LC1,%r9d
      push %r9
      xchg %r9d,%edi
      call _printStringX
      pop %r9 # Const int 15
      mov $0xf,%r14d
      mov %r14d,%eax
      mov $0x1,%ebx
      xchg %eax,%ebx
      int $0x80
      ret
# End of function main