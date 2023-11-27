.text
.global main
.LC2:
  .string "Error: %s\n"
.LC3:
  .string "FAILED ASSERTION"
.LC0:
  .string "%d\n"
.LC1:
  .string "%s\n"
# Function printInt
# Source: ./tests/good/core031.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block2: # Const string "%d\n"
    mov $.LC0,%r9d
    mov %rax,-0x8(%rbp)
    xchg %r9d,%edi
    xchg %r9d,%esi
    mov $0x0,%eax
    call printf
    mov -0x8(%rbp),%rax
    mov $0x0,%eax
    leave
    ret
# End of function printInt
# Function printString
# Source: ./tests/good/core031.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block5: # Const string "%s\n"
      mov $.LC1,%r10d
      mov %rax,-0x8(%rbp)
      xchg %r10d,%edi
      xchg %r10d,%esi
      mov $0x0,%eax
      call printf
      mov -0x8(%rbp),%rax
      mov $0x0,%eax
      leave
      ret
# End of function printString
# Function AddStrings
# Source: ./tests/good/core031.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block3:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%ecx
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rcx,-0x8(%rbp)
        mov %rdi,-0x10(%rbp)
        mov %rsi,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r11d
        mov -0x8(%rbp),%rcx
        mov -0x10(%rbp),%rdi
        mov -0x18(%rbp),%rsi # Const int 1
        mov $0x1,%r8d
        add %r8d,%r11d
        add %r11d,%ecx
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %ecx,%edi
        call malloc
        mov %eax,%r11d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rsi,-0x8(%rbp)
        mov %r11,-0x10(%rbp)
        xchg %r11d,%edi
        xchg %r11d,%esi
        call strcpy
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%r11
        mov %r11,-0x8(%rbp)
        xchg %r11d,%edi
        call strcat
        mov -0x8(%rbp),%r11
        mov %r11d,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core031.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block2: # Const string "Error: %s\n"
          mov $.LC2,%r8d
          mov %rax,-0x8(%rbp)
          xchg %r8d,%edi
          xchg %r8d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%r15d
          xchg %r15d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core031.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block4
          assert_block3: # Const int 2
            mov $0x2,%r8d # Assign variable x
            mov %r8d,%r13d
          assert_block5:
            mov $0x0,%eax
            ret
          assert_block4: # Const string "FAILED ASSERTION"
            mov $.LC3,%r13d
            xchg %r13d,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core031.lat:62:1
_assertEq:
            assertEq_block2:
              cmp %esi,%edi
              sete %r11b
              movzbl %r11b,%r11d
              xchg %r11d,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core031.lat:68:1
main:
              main_block5: # Const int 1
                mov $0x1,%r11d # Const int 1
                mov $0x1,%edi
                mov %edi,%esi
                neg %esi
                xchg %r11d,%edi
                call _f
                mov %eax,%ecx
                xchg %ecx,%edi
                call _printInt # Const int 0
                mov $0x0,%r9d
                mov %r9d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main
# Function f
# Source: ./tests/good/core031.lat:74:1
_f:
                f_block4: # Const int 0
                  mov $0x0,%edx
                  cmp %edx,%edi
                  setg %r9b
                  movzbl %r9b,%r9d # Const int 0
                  mov $0x0,%r12d
                  cmp %r12d,%esi
                  setg %cl
                  movzbl %cl,%ecx
                  mov %r9d,%r12d
                  and %ecx,%r12d # Const int 0
                  mov $0x0,%ecx
                  cmp %ecx,%edi
                  setl %r9b
                  movzbl %r9b,%r9d # Const int 0
                  mov $0x0,%r11d
                  cmp %r11d,%esi
                  setl %r14b
                  movzbl %r14b,%r14d
                  mov %r9d,%r11d
                  and %r14d,%r11d
                  mov %r12d,%edi
                  or %r11d,%edi # If condition
                  cmp $0x0,%edi
                  je f_block5
                f_block3: # Const int 7
                  mov $0x7,%edi
                  mov %edi,%eax
                  ret
                f_block5: # Const int 42
                  mov $0x2a,%edx
                  mov %edx,%eax
                  ret
# End of function f