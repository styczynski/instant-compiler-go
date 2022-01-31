.text
.global main
.LC2:
  .string "Error: %s\n"
.LC3:
  .string "FAILED ASSERTION"
.LC0:
  .string "%d\n"
.LC1:
  .string "%s\n"
# Function printInt
# Source: ./tests/good/core016.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block2: # Const string "%d\n"
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
# Source: ./tests/good/core016.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block4: # Const string "%s\n"
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
# Source: ./tests/good/core016.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block7:
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
        mov $0x1,%r10d
        add %r10d,%r15d
        add %r15d,%r12d
        mov %rsi,-0x8(%rbp)
        mov %rdi,-0x10(%rbp)
        xchg %r12d,%edi
        call malloc
        mov %eax,%r11d
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%rdi
        mov %rsi,-0x8(%rbp)
        mov %r11,-0x10(%rbp)
        xchg %r11d,%edi
        xchg %r11d,%esi
        call strcpy
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%r11
        mov %r11,-0x8(%rbp)
        xchg %r11d,%edi
        call strcat
        mov -0x8(%rbp),%r11
        mov %r11d,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core016.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block2: # Const string "Error: %s\n"
          mov $.LC2,%r11d
          mov %rax,-0x8(%rbp)
          xchg %r11d,%edi
          xchg %r11d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%r8d
          xchg %r8d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core016.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block4
          assert_block3: # Const int 2
            mov $0x2,%r10d # Assign variable x
            mov %r10d,%r9d
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
# Source: ./tests/good/core016.lat:62:1
_assertEq:
            assertEq_block4:
              cmp %esi,%edi
              sete %r10b
              movzbl %r10b,%r10d
              xchg %r10d,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core016.lat:70:1
main:
              main_block3: # Const int 17
                mov $0x11,%r15d
                mov %r15d,%edx
              main_block4: # Const int 0
                mov $0x0,%esi
                cmp %esi,%edx
                setg %r13b
                movzbl %r13b,%r13d # While condition
                cmp $0x0,%r13d
                je main_block5
              main_block9: # Const int 2
                mov $0x2,%r11d
                sub %r11d,%edx # Assign variable y
                nop # While loop return to block_4
                jmp main_block4
                mov $0x0,%eax
                ret
              main_block5: # Const int 0
                mov $0x0,%r14d
                cmp %r14d,%edx
                setl %r12b
                movzbl %r12b,%r12d # If condition
                cmp $0x0,%r12d
                je main_block11
              main_block10: # Const int 0
                mov $0x0,%edi
                call _printInt # Const int 0
                mov $0x0,%r10d
                mov %r10d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
              main_block11: # Const int 1
                mov $0x1,%r15d
                xchg %r15d,%edi
                call _printInt # Const int 0
                mov $0x0,%r14d
                mov %r14d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main