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
# Source: ./tests/good/core032.lat:6:1
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
# Source: ./tests/good/core032.lat:11:1
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
# Source: ./tests/good/core032.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block3:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%r12d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        mov %r12,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r15d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%r12 # Const int 1
        mov $0x1,%r14d
        add %r14d,%r15d
        add %r15d,%r12d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r12d,%edi
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
# Source: ./tests/good/core032.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block2: # Const string "Error: %s\n"
          mov $.LC2,%edx
          mov %rax,-0x8(%rbp)
          xchg %edx,%edi
          xchg %edx,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%r10d
          xchg %r10d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core032.lat:53:1
_assert:
          assert_block6: # If condition
            cmp $0x0,%edi
            je assert_block5
          assert_block4: # Const int 2
            mov $0x2,%r9d # Assign variable x
            mov %r9d,%r8d
          assert_block7:
            mov $0x0,%eax
            ret
          assert_block5: # Const string "FAILED ASSERTION"
            mov $.LC3,%esi
            xchg %esi,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core032.lat:62:1
_assertEq:
            assertEq_block5:
              cmp %esi,%edi
              sete %r9b
              movzbl %r9b,%r9d
              xchg %r9d,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core032.lat:70:1
main:
              main_block4: # Const int 42
                mov $0x2a,%r14d
                mov %r14d,%r10d
                neg %r10d # Const int 1
                mov $0x1,%r8d
                mov %r8d,%r13d
                neg %r13d
                mov %rdx,%rbx
                mov $0x0,%rdx
                mov %r10d,%eax
                idiv %r13d
                mov %eax,%r10d
                mov %rbx,%rdx
                xchg %r10d,%edi
                call _printInt # Const int 0
                mov $0x0,%r10d
                mov %r10d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main