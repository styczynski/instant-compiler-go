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
# Source: ./tests/good/core022.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block3: # Const string "%d\n"
    mov $.LC0,%edx
    mov %rax,-0x8(%rbp)
    xchg %edx,%edi
    xchg %edx,%esi
    mov $0x0,%eax
    call printf
    mov -0x8(%rbp),%rax
    mov $0x0,%eax
    leave
    ret
# End of function printInt
# Function printString
# Source: ./tests/good/core022.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block2: # Const string "%s\n"
      mov $.LC1,%esi
      mov %rax,-0x8(%rbp)
      xchg %esi,%edi
      mov $0x0,%eax
      call printf
      mov -0x8(%rbp),%rax
      mov $0x0,%eax
      leave
      ret
# End of function printString
# Function AddStrings
# Source: ./tests/good/core022.lat:24:1
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
        mov %rsi,-0x8(%rbp)
        mov %rcx,-0x10(%rbp)
        mov %rdi,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r13d
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%rcx
        mov -0x18(%rbp),%rdi # Const int 1
        mov $0x1,%r10d
        add %r10d,%r13d
        add %r13d,%ecx
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %ecx,%edi
        call malloc
        mov %eax,%r13d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %r13,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r13d,%edi
        xchg %r13d,%esi
        call strcpy
        mov -0x8(%rbp),%r13
        mov -0x10(%rbp),%rsi
        mov %r13,-0x8(%rbp)
        xchg %r13d,%edi
        call strcat
        mov -0x8(%rbp),%r13
        mov %r13d,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core022.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block2: # Const string "Error: %s\n"
          mov $.LC2,%r15d
          mov %rax,-0x8(%rbp)
          xchg %r15d,%edi
          xchg %r15d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core022.lat:53:1
_assert:
          assert_block6: # If condition
            cmp $0x0,%edi
            je assert_block5
          assert_block4: # Const int 2
            mov $0x2,%ecx # Assign variable x
            mov %ecx,%r10d
          assert_block7:
            mov $0x0,%eax
            ret
          assert_block5: # Const string "FAILED ASSERTION"
            mov $.LC3,%r10d
            xchg %r10d,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core022.lat:62:1
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
# Source: ./tests/good/core022.lat:68:1
main:
              main_block3: # Const int 0
                mov $0x0,%esi
                xchg %esi,%edi
                call _printInt # Const int 0
                mov $0x0,%r9d
                mov %r9d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main