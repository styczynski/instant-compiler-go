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
# Source: ./tests/good/core007.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block3: # Const string "%d\n"
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
# Source: ./tests/good/core007.lat:11:1
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
# Source: ./tests/good/core007.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block3:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%r14d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        mov %r14,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r12d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%r14 # Const int 1
        mov $0x1,%ecx
        add %ecx,%r12d
        add %r12d,%r14d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r14d,%edi
        call malloc
        mov %eax,%r15d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rsi,-0x8(%rbp)
        mov %r15,-0x10(%rbp)
        xchg %r15d,%edi
        xchg %r15d,%esi
        call strcpy
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%r15
        mov %r15,-0x8(%rbp)
        xchg %r15d,%edi
        call strcat
        mov -0x8(%rbp),%r15
        mov %r15d,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core007.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block2: # Const string "Error: %s\n"
          mov $.LC2,%r12d
          mov %rax,-0x8(%rbp)
          xchg %r12d,%edi
          xchg %r12d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%r15d
          xchg %r15d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core007.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block6
          assert_block2: # Const int 2
            mov $0x2,%ecx # Assign variable x
            mov %ecx,%edx
          assert_block3:
            mov $0x0,%eax
            ret
          assert_block6: # Const string "FAILED ASSERTION"
            mov $.LC3,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core007.lat:62:1
_assertEq:
            assertEq_block3:
              cmp %esi,%edi
              sete %r9b
              movzbl %r9b,%r9d
              xchg %r9d,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core007.lat:70:1
main:
              main_block6: # Const int 7
                mov $0x7,%ecx
                xchg %ecx,%edi
                call _printInt # Const int 0
                mov $0x0,%r13d
                mov %r13d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main