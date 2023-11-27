.text
.global main
.LC1:
  .string "%s\n"
.LC2:
  .string "Error: %s\n"
.LC3:
  .string "FAILED ASSERTION"
.LC0:
  .string "%d\n"
# Function printInt
# Source: ./tests/good/core009.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block3: # Const string "%d\n"
    mov $.LC0,%r10d
    mov %rax,-0x8(%rbp)
    xchg %r10d,%edi
    xchg %r10d,%esi
    mov $0x0,%eax
    call printf
    mov -0x8(%rbp),%rax
    mov $0x0,%eax
    leave
    ret
# End of function printInt
# Function printString
# Source: ./tests/good/core009.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block4: # Const string "%s\n"
      mov $.LC1,%edx
      mov %rax,-0x8(%rbp)
      xchg %edx,%edi
      xchg %edx,%esi
      mov $0x0,%eax
      call printf
      mov -0x8(%rbp),%rax
      mov $0x0,%eax
      leave
      ret
# End of function printString
# Function AddStrings
# Source: ./tests/good/core009.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block5:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%r15d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rsi,-0x8(%rbp)
        mov %r15,-0x10(%rbp)
        mov %rdi,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%ecx
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%r15
        mov -0x18(%rbp),%rdi # Const int 1
        mov $0x1,%r10d
        add %r10d,%ecx
        add %ecx,%r15d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r15d,%edi
        call malloc
        mov %eax,%r10d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rsi,-0x8(%rbp)
        mov %r10,-0x10(%rbp)
        xchg %r10d,%edi
        xchg %r10d,%esi
        call strcpy
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%r10
        mov %r10,-0x8(%rbp)
        xchg %r10d,%edi
        call strcat
        mov -0x8(%rbp),%r10
        mov %r10d,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core009.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block5: # Const string "Error: %s\n"
          mov $.LC2,%r8d
          mov %rax,-0x8(%rbp)
          xchg %r8d,%edi
          xchg %r8d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%r12d
          xchg %r12d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core009.lat:53:1
_assert:
          assert_block6: # If condition
            cmp $0x0,%edi
            je assert_block5
          assert_block4: # Const int 2
            mov $0x2,%r12d # Assign variable x
            mov %r12d,%r10d
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
# Source: ./tests/good/core009.lat:62:1
_assertEq:
            assertEq_block2:
              cmp %esi,%edi
              sete %r13b
              movzbl %r13b,%r13d
              xchg %r13d,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core009.lat:70:1
main:
              main_block6:
                call _foo
                mov %eax,%r13d
                xchg %r13d,%edi
                call _printInt # Const int 0
                mov $0x0,%edx
                mov %edx,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main
# Function foo
# Source: ./tests/good/core009.lat:77:1
_foo:
                foo_block2: # Const int 10
                  mov $0xa,%ecx
                  mov %ecx,%eax
                  ret
# End of function foo