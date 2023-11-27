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
# Source: ./tests/good/core003.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block2: # Const string "%d\n"
    mov $.LC0,%r9d
    mov %rax,-0x8(%rbp)
    xchg %r9d,%edi
    xchg %r9d,%esi
    mov $0x0,%eax
    call printf
    mov -0x8(%rbp),%rax
    mov $0x0,%eax
    leave
    ret
# End of function printInt
# Function printString
# Source: ./tests/good/core003.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block5: # Const string "%s\n"
      mov $.LC1,%r14d
      mov %rax,-0x8(%rbp)
      xchg %r14d,%edi
      xchg %r14d,%esi
      mov $0x0,%eax
      call printf
      mov -0x8(%rbp),%rax
      mov $0x0,%eax
      leave
      ret
# End of function printString
# Function AddStrings
# Source: ./tests/good/core003.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block6:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%ecx
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %rsi,-0x8(%rbp)
        mov %rcx,-0x10(%rbp)
        mov %rdi,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%edx
        mov -0x8(%rbp),%rsi
        mov -0x10(%rbp),%rcx
        mov -0x18(%rbp),%rdi # Const int 1
        mov $0x1,%r10d
        add %r10d,%edx
        add %edx,%ecx
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %ecx,%edi
        call malloc
        mov %eax,%r10d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %r10,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r10d,%edi
        xchg %r10d,%esi
        call strcpy
        mov -0x8(%rbp),%r10
        mov -0x10(%rbp),%rsi
        mov %r10,-0x8(%rbp)
        xchg %r10d,%edi
        call strcat
        mov -0x8(%rbp),%r10
        mov %r10d,%eax
        leave
        ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core003.lat:41:1
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
          mov $0x1,%r13d
          xchg %r13d,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core003.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block3
          assert_block2: # Const int 2
            mov $0x2,%r13d # Assign variable x
            mov %r13d,%edi
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
# Source: ./tests/good/core003.lat:62:1
_assertEq:
            assertEq_block5:
              cmp %esi,%edi
              sete %cl
              movzbl %cl,%ecx
              xchg %ecx,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function f
# Source: ./tests/good/core003.lat:70:1
_f:
              f_block4: # Const int 0
                mov $0x0,%r15d
                mov %r15d,%eax
                ret
# End of function f
# Function g
# Source: ./tests/good/core003.lat:77:1
_g:
                g_block2: # Const int 0
                  mov $0x0,%edi
                  mov %edi,%eax
                  ret
# End of function g
# Function p
# Source: ./tests/good/core003.lat:84:1
_p:
                  p_block2:
                    mov $0x0,%eax
                    ret
# End of function p
# Function main (Entrypoint)
# Source: ./tests/good/core003.lat:87:1
main:
                    main_block2:
                      call _p # Const int 0
                      mov $0x0,%edi
                      mov %edi,%eax
                      mov $0x1,%ebx
                      xchg %eax,%ebx
                      int $0x80
                      ret
# End of function main