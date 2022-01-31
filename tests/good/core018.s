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
# Source: ./tests/good/core018.lat:6:1
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
# Source: ./tests/good/core018.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block3: # Const string "%s\n"
      mov $.LC1,%r13d
      mov %rax,-0x8(%rbp)
      xchg %r13d,%edi
      xchg %r13d,%esi
      mov $0x0,%eax
      call printf
      mov -0x8(%rbp),%rax
      mov $0x0,%eax
      leave
      ret
# End of function printString
# Function AddStrings
# Source: ./tests/good/core018.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block3:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%r9d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rsi,-0x8(%rbp)
        mov %r9,-0x10(%rbp)
        mov %rdi,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r15d
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%r9
        mov -0x18(%rbp),%rdi # Const int 1
        mov $0x1,%ecx
        add %ecx,%r15d
        add %r15d,%r9d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r9d,%edi
        call malloc
        mov %eax,%r15d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %r15,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r15d,%edi
        xchg %r15d,%esi
        call strcpy
        mov -0x8(%rbp),%r15
        mov -0x10(%rbp),%rsi
        mov %r15,-0x8(%rbp)
        xchg %r15d,%edi
        call strcat
        mov -0x8(%rbp),%r15
        mov %r15d,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core018.lat:41:1
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
          mov $0x1,%r12d
          xchg %r12d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core018.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block6
          assert_block2: # Const int 2
            mov $0x2,%esi # Assign variable x
            mov %esi,%r15d
          assert_block3:
            mov $0x0,%eax
            ret
          assert_block6: # Const string "FAILED ASSERTION"
            mov $.LC3,%esi
            xchg %esi,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core018.lat:62:1
_assertEq:
            assertEq_block2:
              cmp %esi,%edi
              sete %r10b
              movzbl %r10b,%r10d
              xchg %r10d,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core018.lat:70:1
main:
                push %rbp
                mov %rsp,%rbp
                sub $0x10,%rsp
              main_block3:
                call _readInt
                mov %eax,%r13d
                mov %r13,-0x8(%rbp)
                call _readString
                mov %eax,%esi
                mov -0x8(%rbp),%r13
                mov %rsi,-0x8(%rbp)
                mov %r13,-0x10(%rbp)
                call _readString
                mov %eax,%edx
                mov -0x8(%rbp),%rsi
                mov -0x10(%rbp),%r13 # Const int 5
                mov $0x5,%r8d
                sub %r8d,%r13d
                mov %rdx,-0x8(%rbp)
                mov %rsi,-0x10(%rbp)
                xchg %r13d,%edi
                call _printInt
                mov -0x8(%rbp),%rdx
                mov -0x10(%rbp),%rsi
                xchg %esi,%edi
                xchg %edx,%esi
                call _AddStrings
                mov %eax,%r10d
                xchg %r10d,%edi
                call _printString # Const int 0
                mov $0x0,%r13d
                mov %r13d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main