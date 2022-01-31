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
# Source: ./tests/good/core021.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block2: # Const string "%d\n"
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
# Source: ./tests/good/core021.lat:11:1
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
# Source: ./tests/good/core021.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block3:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%r11d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        mov %r11,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r9d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%r11 # Const int 1
        mov $0x1,%r12d
        add %r12d,%r9d
        add %r9d,%r11d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r11d,%edi
        call malloc
        mov %eax,%r12d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %r12,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r12d,%edi
        xchg %r12d,%esi
        call strcpy
        mov -0x8(%rbp),%r12
        mov -0x10(%rbp),%rsi
        mov %r12,-0x8(%rbp)
        xchg %r12d,%edi
        call strcat
        mov -0x8(%rbp),%r12
        mov %r12d,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core021.lat:41:1
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
          mov $0x1,%esi
          xchg %esi,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core021.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block3
          assert_block2: # Const int 2
            mov $0x2,%r10d # Assign variable x
            mov %r10d,%r11d
          assert_block4:
            mov $0x0,%eax
            ret
          assert_block3: # Const string "FAILED ASSERTION"
            mov $.LC3,%r13d
            xchg %r13d,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core021.lat:62:1
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
# Source: ./tests/good/core021.lat:68:1
main:
              main_block5: # Const int 1
                mov $0x1,%r10d
                xchg %r10d,%edi
                call _printInt # Const int 0
                mov $0x0,%r11d
                mov %r11d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main