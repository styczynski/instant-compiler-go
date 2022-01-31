.text
.global main
.LC2:
  .string "Error: %s\n"
.LC3:
  .string "FAILED ASSERTION"
.LC4:
  .string "apa"
.LC5:
  .string "true"
.LC6:
  .string "false"
.LC0:
  .string "%d\n"
.LC1:
  .string "%s\n"
# Function printInt
# Source: ./tests/good/core017.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block5: # Const string "%d\n"
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
# Source: ./tests/good/core017.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block2: # Const string "%s\n"
      mov $.LC1,%esi
      mov %rax,-0x8(%rbp)
      xchg %esi,%edi
      mov $0x0,%eax
      call printf
      mov -0x8(%rbp),%rax
      mov $0x0,%eax
      leave
      ret
# End of function printString
# Function AddStrings
# Source: ./tests/good/core017.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block6:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%r15d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        mov %r15,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r12d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%r15 # Const int 1
        mov $0x1,%edx
        add %edx,%r12d
        add %r12d,%r15d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r15d,%edi
        call malloc
        mov %eax,%ecx
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rsi,-0x8(%rbp)
        mov %rcx,-0x10(%rbp)
        xchg %ecx,%edi
        xchg %ecx,%esi
        call strcpy
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%rcx
        mov %rcx,-0x8(%rbp)
        xchg %ecx,%edi
        call strcat
        mov -0x8(%rbp),%rcx
        mov %ecx,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core017.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block2: # Const string "Error: %s\n"
          mov $.LC2,%r13d
          mov %rax,-0x8(%rbp)
          xchg %r13d,%edi
          xchg %r13d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%r15d
          xchg %r15d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core017.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block4
          assert_block3: # Const int 2
            mov $0x2,%r9d # Assign variable x
            mov %r9d,%r13d
          assert_block5:
            mov $0x0,%eax
            ret
          assert_block4: # Const string "FAILED ASSERTION"
            mov $.LC3,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core017.lat:62:1
_assertEq:
            assertEq_block5:
              cmp %esi,%edi
              sete %r14b
              movzbl %r14b,%r14d
              xchg %r14d,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core017.lat:70:1
main:
                push %rbp
                mov %rsp,%rbp
                sub $0x10,%rsp
              main_block3: # Const int 4
                mov $0x4,%r15d # Assign variable x
                mov %r15d,%r14d
              main_block4: # Const int 3
                mov $0x3,%r8d
                cmp %r14d,%r8d
                setle %r10b
                movzbl %r10b,%r10d # Const boolean true
                mov $0x1,%esi
                mov %r10d,%r13d
                and %esi,%r13d # If condition
                cmp $0x0,%r13d
                je main_block6
              main_block5: # Const boolean true
                mov $0x1,%r12d
                mov %r12,-0x8(%rbp)
                mov %r14,-0x10(%rbp)
                xchg %r12d,%edi
                call _printBool
                mov -0x8(%rbp),%r12
                mov -0x10(%rbp),%r14
              main_block7: # Const boolean true
                mov $0x1,%ecx
                mov %rcx,-0x8(%rbp)
                mov %r14,-0x10(%rbp)
                xchg %ecx,%edi
                call _printBool
                mov -0x8(%rbp),%rcx
                mov -0x10(%rbp),%r14
              main_block8: # Const int 4
                mov $0x4,%esi # Const int 5
                mov $0x5,%r12d
                mov %r12d,%r9d
                neg %r9d
                cmp %r9d,%esi
                setl %r13b
                movzbl %r13b,%r13d # Const int 2
                mov $0x2,%r12d
                mov %r13,-0x8(%rbp)
                mov %r14,-0x10(%rbp)
                xchg %r12d,%edi
                call _dontCallMe
                mov %eax,%r9d
                mov -0x8(%rbp),%r13
                mov -0x10(%rbp),%r14
                mov %r13d,%edi
                and %r9d,%edi
                mov %r14,-0x8(%rbp)
                call _printBool
                mov -0x8(%rbp),%r14 # Const int 4
                mov $0x4,%r12d
                cmp %r14d,%r12d
                sete %r8b
                movzbl %r8b,%r8d # Const boolean true
                mov $0x1,%r14d # Const boolean false
                mov $0x0,%r11d
                cmp $0x0,%r11d
                sete %r13b
                cmp %r13d,%r14d
                sete %r12b
                movzbl %r12b,%r12d # Const boolean true
                mov $0x1,%r14d
                mov %r12d,%edi
                and %r14d,%edi
                mov %r8d,%r14d
                and %edi,%r14d
                xchg %r14d,%edi
                call _printBool # Const boolean false
                mov $0x0,%r9d # Const boolean false
                mov $0x0,%r14d
                xchg %r9d,%edi
                xchg %r14d,%esi
                call _implies
                mov %eax,%edx
                xchg %edx,%edi
                call _printBool # Const boolean false
                mov $0x0,%edx # Const boolean true
                mov $0x1,%edi
                xchg %edx,%edi
                xchg %edx,%esi
                call _implies
                mov %eax,%r10d
                xchg %r10d,%edi
                call _printBool # Const boolean true
                mov $0x1,%r11d # Const boolean false
                mov $0x0,%r14d
                xchg %r11d,%edi
                xchg %r14d,%esi
                call _implies
                mov %eax,%r15d
                xchg %r15d,%edi
                call _printBool # Const boolean true
                mov $0x1,%r14d # Const boolean true
                mov $0x1,%edx
                xchg %r14d,%edi
                xchg %edx,%esi
                call _implies
                mov %eax,%r9d
                xchg %r9d,%edi
                call _printBool # Const int 0
                mov $0x0,%ecx
                mov %ecx,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
              main_block6: # Const string "apa"
                mov $.LC4,%edi
                call _printString
                mov $0x0,%eax
                leave
                ret
# End of function main
# Function dontCallMe
# Source: ./tests/good/core017.lat:91:1
_dontCallMe:
                dontCallMe_block2:
                  call _printInt # Const boolean true
                  mov $0x1,%r10d
                  mov %r10d,%eax
                  ret
# End of function dontCallMe
# Function printBool
# Source: ./tests/good/core017.lat:96:1
_printBool:
                    push %rbp
                    mov %rsp,%rbp
                    sub $0x8,%rsp
                  printBool_block6: # If condition
                    cmp $0x0,%edi
                    je printBool_block5
                  printBool_block4: # Const string "true"
                    mov $.LC5,%r9d
                    mov %r9,-0x8(%rbp)
                    xchg %r9d,%edi
                    call _printString
                    mov -0x8(%rbp),%r9
                  printBool_block7:
                    mov $0x0,%eax
                    leave
                    ret
                  printBool_block5: # Const string "false"
                    mov $.LC6,%edx
                    xchg %edx,%edi
                    call _printString
                    mov $0x0,%eax
                    leave
                    ret
# End of function printBool
# Function implies
# Source: ./tests/good/core017.lat:105:1
_implies:
                    implies_block3:
                      cmp $0x0,%edi
                      sete %r15b
                      cmp %esi,%edi
                      sete %r10b
                      movzbl %r10b,%r10d
                      mov %r15d,%r12d
                      or %r10d,%r12d
                      mov %r12d,%eax
                      ret
# End of function implies