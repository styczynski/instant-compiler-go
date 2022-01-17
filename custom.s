.text
.global main
# Function r
# Source: custom.lat:1:1
_r:
    push %rbp
    mov %rsp,%rbp
  r_block2: # Const int 1
    mov $0x1,%ecx
    mov %ecx,%eax
    pop %rbp
    ret
# End of function r
# Function main (Entrypoint)
# Source: custom.lat:49:1
main:
      push %rbp
      mov %rsp,%rbp
    main_block3: # Const boolean true
      mov $0x1,%ecx # Assign variable x
      mov %ecx,%edx
    main_block5: # If condition
      cmp $0x0,%edx
      jne main_block4
      jmp main_block6
    main_block4: # Const int 4
      mov $0x4,%edx
      mov %edx,%eax
      mov $0x1,%ebx
      xchg %eax,%ebx
      int $0x80
      ret
    main_block6: # Const int 9
      mov $0x9,%ecx
      mov %ecx,%eax
      mov $0x1,%ebx
      xchg %eax,%ebx
      int $0x80
      ret
# End of function main