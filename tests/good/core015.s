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
# Source: ./tests/good/core015.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block5: # Const string "%d\n"
    mov $.LC0,%ecx
    mov %rax,-0x8(%rbp)
    xchg %ecx,%edi
    xchg %ecx,%esi
    mov $0x0,%eax
    call printf
    mov -0x8(%rbp),%rax
    mov $0x0,%eax
    leave
    ret
# End of function printInt
# Function printString
# Source: ./tests/good/core015.lat:11:1
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
# Source: ./tests/good/core015.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block3:
        mov %rsi,-0x8(%rbp)
        mov %rdi,-0x10(%rbp)
        call strlen
        mov %eax,%r11d
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%rdi
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        mov %r11,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r8d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%r11 # Const int 1
        mov $0x1,%r12d
        add %r12d,%r8d
        add %r8d,%r11d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r11d,%edi
        call malloc
        mov %eax,%edx
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rdx,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %edx,%edi
        xchg %edx,%esi
        call strcpy
        mov -0x8(%rbp),%rdx
        mov -0x10(%rbp),%rsi
        mov %rdx,-0x8(%rbp)
        xchg %edx,%edi
        call strcat
        mov -0x8(%rbp),%rdx
        mov %edx,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core015.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block4: # Const string "Error: %s\n"
          mov $.LC2,%r10d
          mov %rax,-0x8(%rbp)
          xchg %r10d,%edi
          xchg %r10d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%r10d
          xchg %r10d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core015.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block4
          assert_block3: # Const int 2
            mov $0x2,%edx # Assign variable x
            mov %edx,%r14d
          assert_block5:
            mov $0x0,%eax
            ret
          assert_block4: # Const string "FAILED ASSERTION"
            mov $.LC3,%r9d
            xchg %r9d,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core015.lat:62:1
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
# Source: ./tests/good/core015.lat:70:1
main:
              main_block2: # Const int 17
                mov $0x11,%r12d
                xchg %r12d,%edi
                call _ev
                mov %eax,%r8d
                xchg %r8d,%edi
                call _printInt # Const int 0
                mov $0x0,%r15d
                mov %r15d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main
# Function ev
# Source: ./tests/good/core015.lat:75:1
_ev:
                ev_block4: # Const int 0
                  mov $0x0,%esi
                  cmp %esi,%edi
                  setg %r9b
                  movzbl %r9b,%r9d # If condition
                  cmp $0x0,%r9d
                  je ev_block5
                ev_block3: # Const int 2
                  mov $0x2,%r11d
                  sub %r11d,%edi
                  call _ev
                  mov %eax,%edx
                  mov %edx,%eax
                  ret
                ev_block5: # Const int 0
                  mov $0x0,%r14d
                  cmp %r14d,%edi
                  setl %cl
                  movzbl %cl,%ecx # If condition
                  cmp $0x0,%ecx
                  je ev_block8
                ev_block6: # Const int 0
                  mov $0x0,%r9d
                  mov %r9d,%eax
                  ret
                ev_block8: # Const int 1
                  mov $0x1,%edx
                  mov %edx,%eax
                  ret
# End of function ev