.text
.global main
.LC3:
  .string "FAILED ASSERTION"
.LC0:
  .string "%d\n"
.LC1:
  .string "%s\n"
.LC2:
  .string "Error: %s\n"
# Function printInt
# Source: ./tests/good/core020.lat:6:1
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
# Source: ./tests/good/core020.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block3: # Const string "%s\n"
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
# Source: ./tests/good/core020.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block5:
        mov %rsi,-0x8(%rbp)
        mov %rdi,-0x10(%rbp)
        call strlen
        mov %eax,%r8d
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%rdi
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        mov %r8,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r10d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%r8 # Const int 1
        mov $0x1,%r11d
        add %r11d,%r10d
        add %r10d,%r8d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r8d,%edi
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
# Source: ./tests/good/core020.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block5: # Const string "Error: %s\n"
          mov $.LC2,%r9d
          mov %rax,-0x8(%rbp)
          xchg %r9d,%edi
          xchg %r9d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%r12d
          xchg %r12d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core020.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block4
          assert_block3: # Const int 2
            mov $0x2,%r10d # Assign variable x
            mov %r10d,%r15d
          assert_block5:
            mov $0x0,%eax
            ret
          assert_block4: # Const string "FAILED ASSERTION"
            mov $.LC3,%esi
            xchg %esi,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core020.lat:62:1
_assertEq:
            assertEq_block4:
              cmp %esi,%edi
              sete %r8b
              movzbl %r8b,%r8d
              xchg %r8d,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core020.lat:68:1
main:
              main_block3:
                call _p # Const int 1
                mov $0x1,%r9d
                xchg %r9d,%edi
                call _printInt # Const int 0
                mov $0x0,%r12d
                mov %r12d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main
# Function p
# Source: ./tests/good/core020.lat:74:1
_p:
                p_block2:
                  mov $0x0,%eax
                  ret
# End of function p