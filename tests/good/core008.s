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
# Source: ./tests/good/core008.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block2: # Const string "%d\n"
    mov $.LC0,%r12d
    mov %rax,-0x8(%rbp)
    xchg %r12d,%edi
    xchg %r12d,%esi
    mov $0x0,%eax
    call printf
    mov -0x8(%rbp),%rax
    mov $0x0,%eax
    leave
    ret
# End of function printInt
# Function printString
# Source: ./tests/good/core008.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block2: # Const string "%s\n"
      mov $.LC1,%r15d
      mov %rax,-0x8(%rbp)
      xchg %r15d,%edi
      xchg %r15d,%esi
      mov $0x0,%eax
      call printf
      mov -0x8(%rbp),%rax
      mov $0x0,%eax
      leave
      ret
# End of function printString
# Function AddStrings
# Source: ./tests/good/core008.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block3:
        mov %rsi,-0x8(%rbp)
        mov %rdi,-0x10(%rbp)
        call strlen
        mov %eax,%r14d
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%rdi
        mov %rsi,-0x8(%rbp)
        mov %r14,-0x10(%rbp)
        mov %rdi,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r11d
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%r14
        mov -0x18(%rbp),%rdi # Const int 1
        mov $0x1,%r15d
        add %r15d,%r11d
        add %r11d,%r14d
        mov %rsi,-0x8(%rbp)
        mov %rdi,-0x10(%rbp)
        xchg %r14d,%edi
        call malloc
        mov %eax,%r8d
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%rdi
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
# Source: ./tests/good/core008.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block6: # Const string "Error: %s\n"
          mov $.LC2,%edx
          mov %rax,-0x8(%rbp)
          xchg %edx,%edi
          xchg %edx,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%r12d
          xchg %r12d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core008.lat:53:1
_assert:
          assert_block5: # If condition
            cmp $0x0,%edi
            je assert_block4
          assert_block3: # Const int 2
            mov $0x2,%r13d # Assign variable x
            mov %r13d,%r14d
          assert_block6:
            mov $0x0,%eax
            ret
          assert_block4: # Const string "FAILED ASSERTION"
            mov $.LC3,%r13d
            xchg %r13d,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core008.lat:62:1
_assertEq:
            assertEq_block5:
              cmp %esi,%edi
              sete %r13b
              movzbl %r13b,%r13d
              xchg %r13d,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core008.lat:71:1
main:
                push %rbp
                mov %rsp,%rbp
                sub $0x8,%rsp
              main_block3: # Const int 0
                mov $0x0,%r12d # Const int 7
                mov $0x7,%r12d # Const int 1234234
                mov $0x12d53a,%r11d
                mov %r11d,%r8d
                neg %r8d
                mov %r12,-0x8(%rbp)
                xchg %r8d,%edi
                call _printInt
                mov -0x8(%rbp),%r12
                xchg %r12d,%edi
                call _printInt # Const int 0
                mov $0x0,%ecx
                mov %ecx,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main