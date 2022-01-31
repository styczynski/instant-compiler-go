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
# Source: ./tests/good/core005.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block2: # Const string "%d\n"
    mov $.LC0,%esi
    mov %rax,-0x8(%rbp)
    xchg %esi,%edi
    mov $0x0,%eax
    call printf
    mov -0x8(%rbp),%rax
    mov $0x0,%eax
    leave
    ret
# End of function printInt
# Function printString
# Source: ./tests/good/core005.lat:11:1
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
# Function AddStrings
# Source: ./tests/good/core005.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block2:
        mov %rsi,-0x8(%rbp)
        mov %rdi,-0x10(%rbp)
        call strlen
        mov %eax,%ecx
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%rdi
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        mov %rcx,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%edx
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%rcx # Const int 1
        mov $0x1,%r10d
        add %r10d,%edx
        add %edx,%ecx
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %ecx,%edi
        call malloc
        mov %eax,%r10d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %r10,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r10d,%edi
        xchg %r10d,%esi
        call strcpy
        mov -0x8(%rbp),%r10
        mov -0x10(%rbp),%rsi
        mov %r10,-0x8(%rbp)
        xchg %r10d,%edi
        call strcat
        mov -0x8(%rbp),%r10
        mov %r10d,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core005.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block6: # Const string "Error: %s\n"
          mov $.LC2,%r12d
          mov %rax,-0x8(%rbp)
          xchg %r12d,%edi
          xchg %r12d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%r9d
          xchg %r9d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core005.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block4
          assert_block3: # Const int 2
            mov $0x2,%r13d # Assign variable x
            mov %r13d,%r9d
          assert_block5:
            mov $0x0,%eax
            ret
          assert_block4: # Const string "FAILED ASSERTION"
            mov $.LC3,%r11d
            xchg %r11d,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core005.lat:62:1
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
# Source: ./tests/good/core005.lat:70:1
main:
                push %rbp
                mov %rsp,%rbp
                sub $0x8,%rsp
              main_block6: # Const int 0
                mov $0x0,%ecx # Const int 56
                mov $0x38,%ecx # Assign variable y
                mov %ecx,%esi
              main_block8: # Const int 45
                mov $0x2d,%r11d
                add %r11d,%esi # Const int 2
                mov $0x2,%r15d
                cmp %r15d,%esi
                setle %r13b
                movzbl %r13b,%r13d # If condition
                cmp $0x0,%r13d
                je main_block4
              main_block7: # Const int 1
                mov $0x1,%r11d
                mov %r11d,%r13d
              main_block9:
                mov %r13,-0x8(%rbp)
                xchg %r13d,%edi
                call _printInt
                mov -0x8(%rbp),%r13
              main_block10: # Const int 0
                mov $0x0,%r10d
                mov %r10d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
              main_block4: # Const int 2
                mov $0x2,%r13d
                mov $0x0,%eax
                leave
                ret
# End of function main