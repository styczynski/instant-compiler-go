.text
.global main
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
.LC6:
  .string "yes"
.LC7:
  .string "NOOO"
.LC0:
  .string "%d\n"
# Function printInt
# Source: ./tests/good/core024.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block3: # Const string "%d\n"
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
# Source: ./tests/good/core024.lat:11:1
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
# Function rawLoadString
# Source: ./tests/good/core024.lat:16:1
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
# Source: ./tests/good/core024.lat:20:1
_rawLoadInt:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        rawLoadInt_block3: # Const string "%d"
          mov $.LC3,%r8d
          push %rax
          xchg %r8d,%edi
          xchg %r8d,%esi
          mov $0x0,%eax
          call scanf
          pop %rax
          mov $0x0,%eax
          leave
          ret
# End of function rawLoadInt
# Function readString
# Source: ./tests/good/core024.lat:24:1
_readString:
            push %rbp
            mov %rsp,%rbp
            sub $0x8,%rsp
          readString_block2: # Const int 100
            mov $0x64,%r10d
            xchg %r10d,%edi
            call malloc
            mov %eax,%r13d
            push %r13
            xchg %r13d,%edi
            call _rawLoadString
            pop %r13
            mov %r13d,%eax
            leave
            ret
# End of function readString
# Function readInt
# Source: ./tests/good/core024.lat:30:1
_readInt:
              push %rbp
              mov %rsp,%rbp
              sub $0x8,%rsp
            readInt_block5: # Const int 16
              mov $0x10,%edi
              call malloc
              mov %eax,%esi
              push %rsi
              xchg %esi,%edi
              call _rawLoadInt
              pop %rsi
              mov (%rsi),%edi
              mov %edi,%eax
              leave
              ret
# End of function readInt
# Function AddStrings
# Source: ./tests/good/core024.lat:46:1
_AddStrings:
                push %rbp
                mov %rsp,%rbp
                sub $0x18,%rsp
              AddStrings_block5:
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
                mov %eax,%r12d
                mov -0x8(%rbp),%rsi
                mov -0x10(%rbp),%r11
                mov -0x18(%rbp),%rdi # Const int 1
                mov $0x1,%r15d
                add %r15d,%r12d
                add %r12d,%r11d
                mov %rsi,-0x8(%rbp)
                mov %rdi,-0x10(%rbp)
                xchg %r11d,%edi
                call malloc
                mov %eax,%edx
                mov -0x8(%rbp),%rsi
                mov -0x10(%rbp),%rdi
                mov %rsi,-0x8(%rbp)
                mov %rdx,-0x10(%rbp)
                xchg %edx,%edi
                xchg %edx,%esi
                call strcpy
                mov -0x8(%rbp),%rsi
                mov -0x10(%rbp),%rdx
                mov %rdx,-0x8(%rbp)
                xchg %edx,%edi
                call strcat
                mov -0x8(%rbp),%rdx
                mov %edx,%eax
                leave
                ret
# End of function AddStrings
# Function error
# Source: ./tests/good/core024.lat:63:1
_error:
                  push %rbp
                  mov %rsp,%rbp
                  sub $0x8,%rsp
                error_block3: # Const string "Error: %s\n"
                  mov $.LC4,%ecx
                  mov %rax,-0x8(%rbp)
                  xchg %ecx,%edi
                  xchg %ecx,%esi
                  mov $0x0,%eax
                  call printf
                  mov -0x8(%rbp),%rax # Const int 1
                  mov $0x1,%edi
                  call exit
# End of function error
# Function assert
# Source: ./tests/good/core024.lat:75:1
_assert:
                  assert_block7: # If condition
                    cmp $0x0,%edi
                    je assert_block4
                  assert_block3: # Const int 2
                    mov $0x2,%r14d # Assign variable x
                    mov %r14d,%r12d
                  assert_block5:
                    mov $0x0,%eax
                    ret
                  assert_block4: # Const string "FAILED ASSERTION"
                    mov $.LC5,%r10d
                    xchg %r10d,%edi
                    call _error
                    mov $0x0,%eax
                    ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core024.lat:84:1
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
# Source: ./tests/good/core024.lat:90:1
main:
                      main_block2: # Const int 1
                        mov $0x1,%r11d # Const int 2
                        mov $0x2,%r14d
                        xchg %r11d,%edi
                        xchg %r14d,%esi
                        call _f # Const int 0
                        mov $0x0,%r13d
                        mov %r13d,%eax
                        mov $0x1,%ebx
                        xchg %eax,%ebx
                        int $0x80
                        ret
# End of function main
# Function f
# Source: ./tests/good/core024.lat:95:1
_f:
                          push %rbp
                          mov %rsp,%rbp
                          sub $0x8,%rsp
                        f_block4:
                          cmp %edi,%esi
                          setg %r11b
                          movzbl %r11b,%r11d
                          mov %r11,-0x8(%rbp)
                          call _e
                          mov %eax,%ecx
                          mov -0x8(%rbp),%r11
                          mov %r11d,%edx
                          or %ecx,%edx # If condition
                          cmp $0x0,%edx
                          je f_block1
                        f_block3: # Const string "yes"
                          mov $.LC6,%esi
                          mov %rsi,-0x8(%rbp)
                          xchg %esi,%edi
                          call _printString
                          mov -0x8(%rbp),%rsi
                        f_block1:
                          mov $0x0,%eax
                          leave
                          ret
# End of function f
# Function e
# Source: ./tests/good/core024.lat:100:1
_e:
                          e_block2: # Const string "NOOO"
                            mov $.LC7,%r13d
                            xchg %r13d,%edi
                            call _printString # Const boolean false
                            mov $0x0,%esi
                            mov %esi,%eax
                            ret
# End of function e