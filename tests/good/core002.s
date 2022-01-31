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
.LC4:
  .string "foo"
# Function printInt
# Source: ./tests/good/core002.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block2: # Const string "%d\n"
    mov $.LC0,%r14d
    mov %rax,-0x8(%rbp)
    xchg %r14d,%edi
    xchg %r14d,%esi
    mov $0x0,%eax
    call printf
    mov -0x8(%rbp),%rax
    mov $0x0,%eax
    leave
    ret
# End of function printInt
# Function printString
# Source: ./tests/good/core002.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block2: # Const string "%s\n"
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
# Source: ./tests/good/core002.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block7:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%r14d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %r14,-0x8(%rbp)
        mov %rdi,-0x10(%rbp)
        mov %rsi,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r11d
        mov -0x8(%rbp),%r14
        mov -0x10(%rbp),%rdi
        mov -0x18(%rbp),%rsi # Const int 1
        mov $0x1,%r12d
        add %r12d,%r11d
        add %r11d,%r14d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r14d,%edi
        call malloc
        mov %eax,%r8d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %r8,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r8d,%edi
        xchg %r8d,%esi
        call strcpy
        mov -0x8(%rbp),%r8
        mov -0x10(%rbp),%rsi
        mov %r8,-0x8(%rbp)
        xchg %r8d,%edi
        call strcat
        mov -0x8(%rbp),%r8
        mov %r8d,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core002.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block2: # Const string "Error: %s\n"
          mov $.LC2,%r14d
          mov %rax,-0x8(%rbp)
          xchg %r14d,%edi
          xchg %r14d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%ecx
          xchg %ecx,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core002.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block3
          assert_block2: # Const int 2
            mov $0x2,%edx # Assign variable x
            mov %edx,%r13d
          assert_block4:
            mov $0x0,%eax
            ret
          assert_block3: # Const string "FAILED ASSERTION"
            mov $.LC3,%r8d
            xchg %r8d,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core002.lat:62:1
_assertEq:
            assertEq_block2:
              cmp %esi,%edi
              sete %r8b
              movzbl %r8b,%r8d
              xchg %r8d,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core002.lat:70:1
main:
              main_block5:
                call _foo # Const int 0
                mov $0x0,%ecx
                mov %ecx,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main
# Function foo
# Source: ./tests/good/core002.lat:76:1
_foo:
                foo_block3: # Const string "foo"
                  mov $.LC4,%r8d
                  xchg %r8d,%edi
                  call _printString
                  mov $0x0,%eax
                  ret
# End of function foo