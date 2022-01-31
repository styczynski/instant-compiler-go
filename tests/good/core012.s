.text
.global main
.LC1:
  .string "%s\n"
.LC2:
  .string "Error: %s\n"
.LC3:
  .string "FAILED ASSERTION"
.LC4:
  .string "string concatenation"
.LC5:
  .string "true"
.LC6:
  .string "false"
.LC0:
  .string "%d\n"
# Function printInt
# Source: ./tests/good/core012.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block5: # Const string "%d\n"
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
# Source: ./tests/good/core012.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block3: # Const string "%s\n"
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
# Source: ./tests/good/core012.lat:24:1
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
        mov %eax,%r12d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%r11 # Const int 1
        mov $0x1,%r15d
        add %r15d,%r12d
        add %r12d,%r11d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r11d,%edi
        call malloc
        mov %eax,%r8d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rsi,-0x8(%rbp)
        mov %r8,-0x10(%rbp)
        xchg %r8d,%edi
        xchg %r8d,%esi
        call strcpy
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%r8
        mov %r8,-0x8(%rbp)
        xchg %r8d,%edi
        call strcat
        mov -0x8(%rbp),%r8
        mov %r8d,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core012.lat:41:1
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
          mov $0x1,%r14d
          xchg %r14d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core012.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block6
          assert_block2: # Const int 2
            mov $0x2,%r11d # Assign variable x
            mov %r11d,%r13d
          assert_block3:
            mov $0x0,%eax
            ret
          assert_block6: # Const string "FAILED ASSERTION"
            mov $.LC3,%r10d
            xchg %r10d,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core012.lat:62:1
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
# Source: ./tests/good/core012.lat:70:1
main:
                push %rbp
                mov %rsp,%rbp
                sub $0x10,%rsp
              main_block7: # Const int 56
                mov $0x38,%edi # Const int 23
                mov $0x17,%r15d
                mov %r15d,%r11d
                neg %r11d
                mov %edi,%r9d
                add %r11d,%r9d
                mov %rdi,-0x8(%rbp)
                mov %r11,-0x10(%rbp)
                xchg %r9d,%edi
                call _printInt
                mov -0x8(%rbp),%rdi
                mov -0x10(%rbp),%r11
                mov %edi,%r15d
                sub %r11d,%r15d
                mov %rdi,-0x8(%rbp)
                mov %r11,-0x10(%rbp)
                xchg %r15d,%edi
                call _printInt
                mov -0x8(%rbp),%rdi
                mov -0x10(%rbp),%r11
                mov %edi,%r14d
                mov %rdx,%rbx
                mov %r14d,%eax
                imul %r11d
                mov %eax,%r14d
                mov %rbx,%rdx
                mov %rdi,-0x8(%rbp)
                mov %r11,-0x10(%rbp)
                xchg %r14d,%edi
                call _printInt
                mov -0x8(%rbp),%rdi
                mov -0x10(%rbp),%r11 # Const int 22
                mov $0x16,%r9d
                mov %rdi,-0x8(%rbp)
                mov %r11,-0x10(%rbp)
                xchg %r9d,%edi
                call _printInt
                mov -0x8(%rbp),%rdi
                mov -0x10(%rbp),%r11 # Const int 0
                mov $0x0,%esi
                mov %r11,-0x8(%rbp)
                mov %rdi,-0x10(%rbp)
                xchg %esi,%edi
                call _printInt
                mov -0x8(%rbp),%r11
                mov -0x10(%rbp),%rdi
                mov %edi,%r8d
                sub %r11d,%r8d
                mov %edi,%r14d
                add %r11d,%r14d
                cmp %r14d,%r8d
                setg %r13b
                movzbl %r13b,%r13d
                mov %rdi,-0x8(%rbp)
                mov %r11,-0x10(%rbp)
                xchg %r13d,%edi
                call _printBool
                mov -0x8(%rbp),%rdi
                mov -0x10(%rbp),%r11
                mov %edi,%esi
                mov %rdx,%rbx
                mov $0x0,%rdx
                mov %esi,%eax
                idiv %r11d
                mov %eax,%esi
                mov %rbx,%rdx
                mov %rdx,%rbx
                mov %edi,%eax
                imul %r11d
                mov %eax,%edi
                mov %rbx,%rdx
                cmp %edi,%esi
                setle %r9b
                movzbl %r9b,%r9d
                xchg %r9d,%edi
                call _printBool # Const string "string concatenation"
                mov $.LC4,%r13d
                xchg %r13d,%edi
                call _printString # Const int 0
                mov $0x0,%r9d
                mov %r9d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main
# Function printBool
# Source: ./tests/good/core012.lat:84:1
_printBool:
                printBool_block4: # If condition
                  cmp $0x0,%edi
                  je printBool_block6
                printBool_block3: # Const string "true"
                  mov $.LC5,%r11d
                  xchg %r11d,%edi
                  call _printString
                  mov $0x0,%eax
                  ret
                printBool_block6: # Const string "false"
                  mov $.LC6,%r8d
                  xchg %r8d,%edi
                  call _printString
                  mov $0x0,%eax
                  ret
# End of function printBool