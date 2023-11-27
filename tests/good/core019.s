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
# Source: ./tests/good/core019.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block3: # Const string "%d\n"
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
# Source: ./tests/good/core019.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block4: # Const string "%s\n"
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
# Source: ./tests/good/core019.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block3:
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
        mov %eax,%r14d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%r8 # Const int 1
        mov $0x1,%ecx
        add %ecx,%r14d
        add %r14d,%r8d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r8d,%edi
        call malloc
        mov %eax,%ecx
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rcx,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %ecx,%edi
        xchg %ecx,%esi
        call strcpy
        mov -0x8(%rbp),%rcx
        mov -0x10(%rbp),%rsi
        mov %rcx,-0x8(%rbp)
        xchg %ecx,%edi
        call strcat
        mov -0x8(%rbp),%rcx
        mov %ecx,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core019.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block6: # Const string "Error: %s\n"
          mov $.LC2,%r13d
          mov %rax,-0x8(%rbp)
          xchg %r13d,%edi
          xchg %r13d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%r9d
          xchg %r9d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core019.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block2
          assert_block6: # Const int 2
            mov $0x2,%r10d # Assign variable x
            mov %r10d,%r15d
          assert_block3:
            mov $0x0,%eax
            ret
          assert_block2: # Const string "FAILED ASSERTION"
            mov $.LC3,%r13d
            xchg %r13d,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core019.lat:62:1
_assertEq:
            assertEq_block2:
              cmp %esi,%edi
              sete %cl
              movzbl %cl,%ecx
              xchg %ecx,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core019.lat:68:1
main:
                push %rbp
                mov %rsp,%rbp
                sub $0x8,%rsp
              main_block3: # Const int 78
                mov $0x4e,%esi # Const int 1
                mov $0x1,%r9d
                mov %r9,-0x8(%rbp)
                xchg %r9d,%edi
                call _printInt
                mov -0x8(%rbp),%r9
                mov %r9,-0x8(%rbp)
                xchg %r9d,%edi
                call _printInt
                mov -0x8(%rbp),%r9
                mov %r9d,%edx
              main_block8: # Const int 76
                mov $0x4c,%r8d
                cmp %r8d,%edx
                setg %r9b
                movzbl %r9b,%r9d # While condition
                cmp $0x0,%r9d
                je main_block13
              main_block9: # Const int 1
                mov $0x1,%esi
                nop
                sub %esi,%edx
                mov %rdx,-0x8(%rbp)
                xchg %edx,%edi
                call _printInt
                mov -0x8(%rbp),%rdx # Const int 7
                mov $0x7,%r8d
                add %r8d,%edx
                mov %rdx,-0x8(%rbp)
                xchg %edx,%edi
                call _printInt
                mov -0x8(%rbp),%rdx # While loop return to block_8
                jmp main_block8
                mov $0x0,%eax
                leave
                ret
              main_block13:
                mov %rdx,-0x8(%rbp)
                xchg %edx,%edi
                call _printInt
                mov -0x8(%rbp),%rdx
              main_block14: # Const int 4
                mov $0x4,%r8d
                cmp %r8d,%edx
                setg %r11b
                movzbl %r11b,%r11d # If condition
                cmp $0x0,%r11d
                je main_block17
              main_block15: # Const int 4
                mov $0x4,%edx
                mov %rdx,-0x8(%rbp)
                xchg %edx,%edi
                call _printInt
                mov -0x8(%rbp),%rdx
              main_block18:
                mov %rdx,-0x8(%rbp)
                xchg %edx,%edi
                call _printInt
                mov -0x8(%rbp),%rdx
              main_block19: # Const int 0
                mov $0x0,%r13d
                mov %r13d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
              main_block17: # Const string "foo"
                mov $.LC4,%edx
                xchg %edx,%edi
                call _printString
                mov $0x0,%eax
                leave
                ret
# End of function main