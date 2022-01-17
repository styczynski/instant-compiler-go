.text
.global main
.LC0:
  .string "a"
.LC1:
  .string "v"
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
# Function qpa
# Source: custom.lat:61:1
_qpa:
      push %rbp
      mov %rsp,%rbp
    qpa_block2: # Const int 8
      mov $0x8,%ecx
      mov %ecx,%eax
      pop %rbp
      ret
# End of function qpa
# Function main (Entrypoint)
# Source: custom.lat:65:1
main:
        push %rbp
        mov %rsp,%rbp
      main_block6: # Const string "a"
        mov $.LC0,%edx # Const string "v"
        mov $.LC1,%ecx
        mov %edx,%edi
        mov %ecx,%esi
        call AddStrings
        mov %edx,%eax # Preserve all registries for temp_18 call
        push %rax
        push %rbx
        push %rcx
        push %rdx # Load function argument temp_10
        push %rdx
        call _qpa # Store function result into temp_18
        mov %eax,%eax
        pop %rdx
        pop %rcx
        pop %rbx
        pop %rax
        mov $0x1,%ebx
        xchg %eax,%ebx
        int $0x80
        ret
# End of function main