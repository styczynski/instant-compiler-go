.text
.global main
.LC0:
  .string "%d\n"
.LC1:
  .string "%s\n"
.LC2:
  .string "%s"
.LC3:
  .string "%d"
.LC4:
  .string "Error: %s\n"
.LC5:
  .string "FAILED ASSERTION"
# Function printInt
# Source: ./tests/good/core020.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block2: # Const string "%d\n"
    mov $.LC0,%esi
    mov %rax,-0x8(%rbp)
    xchg %esi,%edi
    mov $0x0,%eax
    call printf
    mov -0x8(%rbp),%rax
    mov $0x0,%eax
    leave
    ret
# End of function printInt
# Function printString
# Source: ./tests/good/core020.lat:11:1
_printString:
      push %rbp
      mov %rsp,%rbp
      sub $0x8,%rsp
    printString_block2: # Const string "%s\n"
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
# Function rawLoadString
# Source: ./tests/good/core020.lat:16:1
_rawLoadString:
        push %rbp
        mov %rsp,%rbp
        sub $0x8,%rsp
      rawLoadString_block4: # Const string "%s"
        mov $.LC2,%r11d
        push %rax
        xchg %r11d,%edi
        xchg %r11d,%esi
        mov $0x0,%eax
        call scanf
        pop %rax
        mov $0x0,%eax
        leave
        ret
# End of function rawLoadString
# Function rawLoadInt
# Source: ./tests/good/core020.lat:20:1
_rawLoadInt:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        rawLoadInt_block3: # Const string "%d"
          mov $.LC3,%edx
          push %rax
          xchg %edx,%edi
          xchg %edx,%esi
          mov $0x0,%eax
          call scanf
          pop %rax
          mov $0x0,%eax
          leave
          ret
# End of function rawLoadInt
# Function readString
# Source: ./tests/good/core020.lat:24:1
_readString:
            push %rbp
            mov %rsp,%rbp
            sub $0x8,%rsp
          readString_block6: # Const int 100
            mov $0x64,%r12d
            xchg %r12d,%edi
            call malloc
            mov %eax,%edx
            push %rdx
            xchg %edx,%edi
            call _rawLoadString
            pop %rdx
            mov %edx,%eax
            leave
            ret
# End of function readString
# Function readInt
# Source: ./tests/good/core020.lat:30:1
_readInt:
              push %rbp
              mov %rsp,%rbp
              sub $0x8,%rsp
            readInt_block2: # Const int 16
              mov $0x10,%r14d
              xchg %r14d,%edi
              call malloc
              mov %eax,%esi
              push %rsi
              xchg %esi,%edi
              call _rawLoadInt
              pop %rsi
              mov (%rsi),%ecx
              mov %ecx,%eax
              leave
              ret
# End of function readInt
# Function AddStrings
# Source: ./tests/good/core020.lat:46:1
_AddStrings:
                push %rbp
                mov %rsp,%rbp
                sub $0x18,%rsp
              AddStrings_block9:
                mov %rdi,-0x8(%rbp)
                mov %rsi,-0x10(%rbp)
                call strlen
                mov %eax,%r11d
                mov -0x8(%rbp),%rdi
                mov -0x10(%rbp),%rsi
                mov %rsi,-0x8(%rbp)
                mov %r11,-0x10(%rbp)
                mov %rdi,-0x18(%rbp)
                xchg %esi,%edi
                call strlen
                mov %eax,%r8d
                mov -0x8(%rbp),%rsi
                mov -0x10(%rbp),%r11
                mov -0x18(%rbp),%rdi # Const int 1
                mov $0x1,%edx
                add %edx,%r8d
                add %r8d,%r11d
                mov %rdi,-0x8(%rbp)
                mov %rsi,-0x10(%rbp)
                xchg %r11d,%edi
                call malloc
                mov %eax,%r13d
                mov -0x8(%rbp),%rdi
                mov -0x10(%rbp),%rsi
                mov %r13,-0x8(%rbp)
                mov %rsi,-0x10(%rbp)
                xchg %r13d,%edi
                xchg %r13d,%esi
                call strcpy
                mov -0x8(%rbp),%r13
                mov -0x10(%rbp),%rsi
                mov %r13,-0x8(%rbp)
                xchg %r13d,%edi
                call strcat
                mov -0x8(%rbp),%r13
                mov %r13d,%eax
                leave
                ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core020.lat:63:1
_error:
                  push %rbp
                  mov %rsp,%rbp
                  sub $0x8,%rsp
                error_block4: # Const string "Error: %s\n"
                  mov $.LC4,%r14d
                  mov %rax,-0x8(%rbp)
                  xchg %r14d,%edi
                  xchg %r14d,%esi
                  mov $0x0,%eax
                  call printf
                  mov -0x8(%rbp),%rax # Const int 1
                  mov $0x1,%r8d
                  xchg %r8d,%edi
                  call exit
# End of function error
# Function assert
# Source: ./tests/good/core020.lat:75:1
_assert:
                  assert_block4: # If condition
                    cmp $0x0,%edi
                    je assert_block3
                  assert_block2: # Const int 2
                    mov $0x2,%r11d # Assign variable x
                    mov %r11d,%r10d
                  assert_block5:
                    mov $0x0,%eax
                    ret
                  assert_block3: # Const string "FAILED ASSERTION"
                    mov $.LC5,%r13d
                    xchg %r13d,%edi
                    call _error
                    mov $0x0,%eax
                    ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core020.lat:84:1
_assertEq:
                    assertEq_block4:
                      cmp %esi,%edi
                      sete %r13b
                      movzbl %r13b,%r13d
                      xchg %r13d,%edi
                      call _assert
                      mov $0x0,%eax
                      ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core020.lat:90:1
main:
                      main_block6:
                        call _p # Const int 1
                        mov $0x1,%r10d
                        xchg %r10d,%edi
                        call _printInt # Const int 0
                        mov $0x0,%edx
                        mov %edx,%eax
                        mov $0x1,%ebx
                        xchg %eax,%ebx
                        int $0x80
                        ret
# End of function main
# Function p
# Source: ./tests/good/core020.lat:96:1
_p:
                        p_block2:
                          mov $0x0,%eax
                          ret
# End of function p