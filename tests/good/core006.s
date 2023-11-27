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
# Source: ./tests/good/core006.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block5: # Const string "%d\n"
    mov $.LC0,%r13d
    mov %rax,-0x8(%rbp)
    xchg %r13d,%edi
    xchg %r13d,%esi
    mov $0x0,%eax
    call printf
    mov -0x8(%rbp),%rax
    mov $0x0,%eax
    leave
    ret
# End of function printInt
# Function printString
# Source: ./tests/good/core006.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block2: # Const string "%s\n"
      mov $.LC1,%r8d
      mov %rax,-0x8(%rbp)
      xchg %r8d,%edi
      xchg %r8d,%esi
      mov $0x0,%eax
      call printf
      mov -0x8(%rbp),%rax
      mov $0x0,%eax
      leave
      ret
# End of function printString
# Function AddStrings
# Source: ./tests/good/core006.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block3:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%r8d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        mov %r8,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r11d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%r8 # Const int 1
        mov $0x1,%r9d
        add %r9d,%r11d
        add %r11d,%r8d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r8d,%edi
        call malloc
        mov %eax,%edx
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rdx,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %edx,%edi
        xchg %edx,%esi
        call strcpy
        mov -0x8(%rbp),%rdx
        mov -0x10(%rbp),%rsi
        mov %rdx,-0x8(%rbp)
        xchg %edx,%edi
        call strcat
        mov -0x8(%rbp),%rdx
        mov %edx,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core006.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block3: # Const string "Error: %s\n"
          mov $.LC2,%r11d
          mov %rax,-0x8(%rbp)
          xchg %r11d,%edi
          xchg %r11d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%edx
          xchg %edx,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core006.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block3
          assert_block2: # Const int 2
            mov $0x2,%r9d # Assign variable x
            mov %r9d,%r12d
          assert_block4:
            mov $0x0,%eax
            ret
          assert_block3: # Const string "FAILED ASSERTION"
            mov $.LC3,%ecx
            xchg %ecx,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core006.lat:62:1
_assertEq:
            assertEq_block5:
              cmp %esi,%edi
              sete %dl
              movzbl %dl,%edx
              xchg %edx,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core006.lat:70:1
main:
                push %rbp
                mov %rsp,%rbp
                sub $0x8,%rsp
              main_block6: # Const int 0
                mov $0x0,%r15d # Const int 0
                mov $0x0,%edi # Const int 45
                mov $0x2d,%r10d # Const int 36
                mov $0x24,%edx
                mov %edx,%r9d
                neg %r9d
                mov %r9,-0x8(%rbp)
                xchg %r10d,%edi
                call _printInt
                mov -0x8(%rbp),%r9
                xchg %r9d,%edi
                call _printInt # Const int 0
                mov $0x0,%esi
                mov %esi,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main