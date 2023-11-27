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
# Source: ./tests/good/core014.lat:6:1
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
# Source: ./tests/good/core014.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block5: # Const string "%s\n"
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
# Source: ./tests/good/core014.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block3:
        mov %rsi,-0x8(%rbp)
        mov %rdi,-0x10(%rbp)
        call strlen
        mov %eax,%ecx
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%rdi
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        mov %rcx,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r11d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%rcx # Const int 1
        mov $0x1,%r12d
        add %r12d,%r11d
        add %r11d,%ecx
        mov %rsi,-0x8(%rbp)
        mov %rdi,-0x10(%rbp)
        xchg %ecx,%edi
        call malloc
        mov %eax,%r11d
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%rdi
        mov %r11,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r11d,%edi
        xchg %r11d,%esi
        call strcpy
        mov -0x8(%rbp),%r11
        mov -0x10(%rbp),%rsi
        mov %r11,-0x8(%rbp)
        xchg %r11d,%edi
        call strcat
        mov -0x8(%rbp),%r11
        mov %r11d,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core014.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block3: # Const string "Error: %s\n"
          mov $.LC2,%r10d
          mov %rax,-0x8(%rbp)
          xchg %r10d,%edi
          xchg %r10d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core014.lat:53:1
_assert:
          assert_block6: # If condition
            cmp $0x0,%edi
            je assert_block2
          assert_block5: # Const int 2
            mov $0x2,%r9d # Assign variable x
            mov %r9d,%r12d
          assert_block7:
            mov $0x0,%eax
            ret
          assert_block2: # Const string "FAILED ASSERTION"
            mov $.LC3,%edx
            xchg %edx,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core014.lat:62:1
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
# Source: ./tests/good/core014.lat:70:1
main:
                push %rbp
                mov %rsp,%rbp
                sub $0x18,%rsp
              main_block3: # Const int 0
                mov $0x0,%edx # Const int 0
                mov $0x0,%r15d # Const int 0
                mov $0x0,%r15d # Const int 1
                mov $0x1,%r10d # Assign variable hi
                mov %r10d,%r9d # Const int 5000000
                mov $0x4c4b40,%ecx
                mov %r9,-0x8(%rbp)
                mov %rcx,-0x10(%rbp)
                mov %r10,-0x18(%rbp)
                xchg %r10d,%edi
                call _printInt
                mov -0x8(%rbp),%r9
                mov -0x10(%rbp),%rcx
                mov -0x18(%rbp),%r10
                mov %r9d,%r8d
              main_block12:
                cmp %ecx,%r8d
                setl %r15b
                movzbl %r15b,%r15d # While condition
                cmp $0x0,%r15d
                je main_block13
              main_block8:
                mov %rcx,-0x8(%rbp)
                mov %r8,-0x10(%rbp)
                mov %r10,-0x18(%rbp)
                xchg %r8d,%edi
                call _printInt
                mov -0x8(%rbp),%rcx
                mov -0x10(%rbp),%r8
                mov -0x18(%rbp),%r10
                mov %r10d,%r8d
                add %r8d,%r8d
                sub %r10d,%r8d
                mov %r8d,%r10d # While loop return to block_12
                jmp main_block12
                mov $0x0,%eax
                leave
                ret
              main_block13: # Const int 0
                mov $0x0,%r12d
                mov %r12d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main