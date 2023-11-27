.text
.global main
.LC0:
  .string "%d\n"
.LC1:
  .string "%s\n"
.LC2:
  .string "Error: %s\n"
.LC3:
  .string "FAILED ASSERTION"
# Function printInt
# Source: ./tests/good/core010.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block2: # Const string "%d\n"
    mov $.LC0,%r15d
    mov %rax,-0x8(%rbp)
    xchg %r15d,%edi
    xchg %r15d,%esi
    mov $0x0,%eax
    call printf
    mov -0x8(%rbp),%rax
    mov $0x0,%eax
    leave
    ret
# End of function printInt
# Function printString
# Source: ./tests/good/core010.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block4: # Const string "%s\n"
      mov $.LC1,%r11d
      mov %rax,-0x8(%rbp)
      xchg %r11d,%edi
      xchg %r11d,%esi
      mov $0x0,%eax
      call printf
      mov -0x8(%rbp),%rax
      mov $0x0,%eax
      leave
      ret
# End of function printString
# Function AddStrings
# Source: ./tests/good/core010.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block6:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%r15d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        mov %r15,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r9d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%r15 # Const int 1
        mov $0x1,%edx
        add %edx,%r9d
        add %r9d,%r15d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r15d,%edi
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
# Source: ./tests/good/core010.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block5: # Const string "Error: %s\n"
          mov $.LC2,%edx
          mov %rax,-0x8(%rbp)
          xchg %edx,%edi
          xchg %edx,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%edx
          xchg %edx,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core010.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block6
          assert_block2: # Const int 2
            mov $0x2,%r8d # Assign variable x
            mov %r8d,%edx
          assert_block3:
            mov $0x0,%eax
            ret
          assert_block6: # Const string "FAILED ASSERTION"
            mov $.LC3,%r12d
            xchg %r12d,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core010.lat:62:1
_assertEq:
            assertEq_block2:
              cmp %esi,%edi
              sete %dl
              movzbl %dl,%edx
              xchg %edx,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core010.lat:70:1
main:
              main_block2: # Const int 5
                mov $0x5,%ecx
                xchg %ecx,%edi
                call _fac
                mov %eax,%r15d
                xchg %r15d,%edi
                call _printInt # Const int 0
                mov $0x0,%edi
                mov %edi,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main
# Function fac
# Source: ./tests/good/core010.lat:75:1
_fac:
                fac_block7: # Const int 0
                  mov $0x0,%r8d # Const int 0
                  mov $0x0,%r14d # Const int 1
                  mov $0x1,%r13d
                  mov %r13d,%esi
                fac_block10: # Const int 0
                  mov $0x0,%r11d
                  cmp %r11d,%edi
                  setg %r12b
                  movzbl %r12b,%r12d # While condition
                  cmp $0x0,%r12d
                  je fac_block11
                fac_block5:
                  mov %rdx,%rbx
                  mov %esi,%eax
                  imul %edi
                  mov %eax,%esi
                  mov %rbx,%rdx # Const int 1
                  mov $0x1,%r12d
                  sub %r12d,%edi # Assign variable n
                  nop
                  nop # While loop return to block_10
                  jmp fac_block10
                  mov $0x0,%eax
                  ret
                fac_block11:
                  mov %esi,%eax
                  ret
# End of function fac