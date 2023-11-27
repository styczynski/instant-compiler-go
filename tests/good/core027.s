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
  .string "bad"
.LC5:
  .string "good"
# Function printInt
# Source: ./tests/good/core027.lat:6:1
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
# Source: ./tests/good/core027.lat:11:1
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
# Source: ./tests/good/core027.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block3:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%r13d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        mov %r13,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%r12d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov -0x18(%rbp),%r13 # Const int 1
        mov $0x1,%ecx
        add %ecx,%r12d
        add %r12d,%r13d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r13d,%edi
        call malloc
        mov %eax,%r8d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
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
# Source: ./tests/good/core027.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block6: # Const string "Error: %s\n"
          mov $.LC2,%r11d
          mov %rax,-0x8(%rbp)
          xchg %r11d,%edi
          xchg %r11d,%esi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%r11d
          xchg %r11d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core027.lat:53:1
_assert:
          assert_block6: # If condition
            cmp $0x0,%edi
            je assert_block5
          assert_block4: # Const int 2
            mov $0x2,%r8d # Assign variable x
            mov %r8d,%edi
          assert_block7:
            mov $0x0,%eax
            ret
          assert_block5: # Const string "FAILED ASSERTION"
            mov $.LC3,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core027.lat:62:1
_assertEq:
            assertEq_block3:
              cmp %esi,%edi
              sete %r15b
              movzbl %r15b,%r15d
              xchg %r15d,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core027.lat:68:1
main:
              main_block5: # Const string "bad"
                mov $.LC4,%esi
                xchg %esi,%edi
                call _f # Const int 0
                mov $0x0,%r8d
                mov %r8d,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main
# Function f
# Source: ./tests/good/core027.lat:73:1
_f:
                f_block4: # Const string "good"
                  mov $.LC5,%ecx
                  xchg %ecx,%edi
                  call _printString
                  mov $0x0,%eax
                  ret
# End of function f