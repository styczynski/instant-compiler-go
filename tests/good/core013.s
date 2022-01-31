.text
.global main
.LC9:
  .string "false"
.LC0:
  .string "%d\n"
.LC3:
  .string "%d"
.LC5:
  .string "FAILED ASSERTION"
.LC6:
  .string "&&"
.LC8:
  .string "!"
.LC10:
  .string "true"
.LC1:
  .string "%s\n"
.LC2:
  .string "%s"
.LC4:
  .string "Error: %s\n"
.LC7:
  .string "||"
# Function printInt
# Source: ./tests/good/core013.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block2: # Const string "%d\n"
    mov $.LC0,%r11d
    mov %rax,-0x8(%rbp)
    xchg %r11d,%edi
    xchg %r11d,%esi
    mov $0x0,%eax
    call printf
    mov -0x8(%rbp),%rax
    mov $0x0,%eax
    leave
    ret
# End of function printInt
# Function printString
# Source: ./tests/good/core013.lat:11:1
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
# Function rawLoadString
# Source: ./tests/good/core013.lat:16:1
_rawLoadString:
        push %rbp
        mov %rsp,%rbp
        sub $0x8,%rsp
      rawLoadString_block3: # Const string "%s"
        mov $.LC2,%esi
        push %rax
        xchg %esi,%edi
        mov $0x0,%eax
        call scanf
        pop %rax
        mov $0x0,%eax
        leave
        ret
# End of function rawLoadString
# Function rawLoadInt
# Source: ./tests/good/core013.lat:20:1
_rawLoadInt:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        rawLoadInt_block4: # Const string "%d"
          mov $.LC3,%esi
          push %rax
          xchg %esi,%edi
          mov $0x0,%eax
          call scanf
          pop %rax
          mov $0x0,%eax
          leave
          ret
# End of function rawLoadInt
# Function readString
# Source: ./tests/good/core013.lat:24:1
_readString:
            push %rbp
            mov %rsp,%rbp
            sub $0x8,%rsp
          readString_block3: # Const int 100
            mov $0x64,%r8d
            xchg %r8d,%edi
            call malloc
            mov %eax,%r11d
            push %r11
            xchg %r11d,%edi
            call _rawLoadString
            pop %r11
            mov %r11d,%eax
            leave
            ret
# End of function readString
# Function readInt
# Source: ./tests/good/core013.lat:30:1
_readInt:
              push %rbp
              mov %rsp,%rbp
              sub $0x8,%rsp
            readInt_block6: # Const int 16
              mov $0x10,%r10d
              xchg %r10d,%edi
              call malloc
              mov %eax,%r15d
              push %r15
              xchg %r15d,%edi
              call _rawLoadInt
              pop %r15
              mov (%r15),%esi
              mov %esi,%eax
              leave
              ret
# End of function readInt
# Function AddStrings
# Source: ./tests/good/core013.lat:46:1
_AddStrings:
                push %rbp
                mov %rsp,%rbp
                sub $0x18,%rsp
              AddStrings_block7:
                mov %rsi,-0x8(%rbp)
                mov %rdi,-0x10(%rbp)
                call strlen
                mov %eax,%r10d
                mov -0x8(%rbp),%rsi
                mov -0x10(%rbp),%rdi
                mov %rdi,-0x8(%rbp)
                mov %rsi,-0x10(%rbp)
                mov %r10,-0x18(%rbp)
                xchg %esi,%edi
                call strlen
                mov %eax,%ecx
                mov -0x8(%rbp),%rdi
                mov -0x10(%rbp),%rsi
                mov -0x18(%rbp),%r10 # Const int 1
                mov $0x1,%r14d
                add %r14d,%ecx
                add %ecx,%r10d
                mov %rdi,-0x8(%rbp)
                mov %rsi,-0x10(%rbp)
                xchg %r10d,%edi
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
# Source: ./tests/good/core013.lat:63:1
_error:
                  push %rbp
                  mov %rsp,%rbp
                  sub $0x8,%rsp
                error_block2: # Const string "Error: %s\n"
                  mov $.LC4,%r9d
                  mov %rax,-0x8(%rbp)
                  xchg %r9d,%edi
                  xchg %r9d,%esi
                  mov $0x0,%eax
                  call printf
                  mov -0x8(%rbp),%rax # Const int 1
                  mov $0x1,%r8d
                  xchg %r8d,%edi
                  call exit
# End of function error
# Function assert
# Source: ./tests/good/core013.lat:75:1
_assert:
                  assert_block7: # If condition
                    cmp $0x0,%edi
                    je assert_block6
                  assert_block2: # Const int 2
                    mov $0x2,%r11d # Assign variable x
                    mov %r11d,%r8d
                  assert_block3:
                    mov $0x0,%eax
                    ret
                  assert_block6: # Const string "FAILED ASSERTION"
                    mov $.LC5,%r8d
                    xchg %r8d,%edi
                    call _error
                    mov $0x0,%eax
                    ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core013.lat:84:1
_assertEq:
                    assertEq_block4:
                      cmp %esi,%edi
                      sete %r15b
                      movzbl %r15b,%r15d
                      xchg %r15d,%edi
                      call _assert
                      mov $0x0,%eax
                      ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core013.lat:92:1
main:
                        push %rbp
                        mov %rsp,%rbp
                        sub $0x8,%rsp
                      main_block3: # Const string "&&"
                        mov $.LC6,%r13d
                        xchg %r13d,%edi
                        call _printString # Const int 1
                        mov $0x1,%edi
                        mov %edi,%r13d
                        neg %r13d
                        xchg %r13d,%edi
                        call _test
                        mov %eax,%ecx # Const int 0
                        mov $0x0,%r13d
                        mov %rcx,-0x8(%rbp)
                        xchg %r13d,%edi
                        call _test
                        mov %eax,%edi
                        mov -0x8(%rbp),%rcx
                        mov %ecx,%r13d
                        and %edi,%r13d
                        xchg %r13d,%edi
                        call _printBool # Const int 2
                        mov $0x2,%r13d
                        mov %r13d,%ecx
                        neg %ecx
                        xchg %ecx,%edi
                        call _test
                        mov %eax,%esi # Const int 1
                        mov $0x1,%r8d
                        mov %rsi,-0x8(%rbp)
                        xchg %r8d,%edi
                        call _test
                        mov %eax,%r11d
                        mov -0x8(%rbp),%rsi
                        mov %esi,%r15d
                        and %r11d,%r15d
                        xchg %r15d,%edi
                        call _printBool # Const int 3
                        mov $0x3,%r15d
                        xchg %r15d,%edi
                        call _test
                        mov %eax,%edx # Const int 5
                        mov $0x5,%r13d
                        mov %r13d,%r9d
                        neg %r9d
                        mov %rdx,-0x8(%rbp)
                        xchg %r9d,%edi
                        call _test
                        mov %eax,%r13d
                        mov -0x8(%rbp),%rdx
                        mov %edx,%r15d
                        and %r13d,%r15d
                        xchg %r15d,%edi
                        call _printBool # Const int 234234
                        mov $0x392fa,%r11d
                        xchg %r11d,%edi
                        call _test
                        mov %eax,%r13d # Const int 21321
                        mov $0x5349,%r9d
                        mov %r13,-0x8(%rbp)
                        xchg %r9d,%edi
                        call _test
                        mov %eax,%r14d
                        mov -0x8(%rbp),%r13
                        mov %r13d,%edi
                        and %r14d,%edi
                        call _printBool # Const string "||"
                        mov $.LC7,%r13d
                        xchg %r13d,%edi
                        call _printString # Const int 1
                        mov $0x1,%ecx
                        mov %ecx,%r13d
                        neg %r13d
                        xchg %r13d,%edi
                        call _test
                        mov %eax,%r8d
                        cmp $0x0,%r8d
                        jne main_3_local_lazy_0 # Const int 0
                        mov $0x0,%r13d
                        mov %r8,-0x8(%rbp)
                        xchg %r13d,%edi
                        call _test
                        mov %eax,%ecx
                        mov -0x8(%rbp),%r8
                        mov %r8d,%r13d
                        or %ecx,%r13d
                      main_3_local_lazy_0: # True value for lazy expression
                        mov $0x1,%r13d
                        xchg %r13d,%edi
                        call _printBool # Const int 2
                        mov $0x2,%r11d
                        mov %r11d,%r13d
                        neg %r13d
                        xchg %r13d,%edi
                        call _test
                        mov %eax,%esi
                        cmp $0x0,%esi
                        jne main_3_local_lazy_1 # Const int 1
                        mov $0x1,%r15d
                        mov %rsi,-0x8(%rbp)
                        xchg %r15d,%edi
                        call _test
                        mov %eax,%r13d
                        mov -0x8(%rbp),%rsi
                        mov %esi,%ecx
                        or %r13d,%ecx
                      main_3_local_lazy_1: # True value for lazy expression
                        mov $0x1,%ecx
                        xchg %ecx,%edi
                        call _printBool # Const int 3
                        mov $0x3,%ecx
                        xchg %ecx,%edi
                        call _test
                        mov %eax,%edx
                        cmp $0x0,%edx
                        jne main_3_local_lazy_2 # Const int 5
                        mov $0x5,%r8d
                        mov %r8d,%r11d
                        neg %r11d
                        mov %rdx,-0x8(%rbp)
                        xchg %r11d,%edi
                        call _test
                        mov %eax,%esi
                        mov -0x8(%rbp),%rdx
                        mov %edx,%r9d
                        or %esi,%r9d
                      main_3_local_lazy_2: # True value for lazy expression
                        mov $0x1,%r9d
                        xchg %r9d,%edi
                        call _printBool # Const int 234234
                        mov $0x392fa,%r9d
                        xchg %r9d,%edi
                        call _test
                        mov %eax,%r14d
                        cmp $0x0,%r14d
                        jne main_3_local_lazy_3 # Const int 21321
                        mov $0x5349,%esi
                        mov %r14,-0x8(%rbp)
                        xchg %esi,%edi
                        call _test
                        mov %eax,%r13d
                        mov -0x8(%rbp),%r14
                        mov %r14d,%edx
                        or %r13d,%edx
                      main_3_local_lazy_3: # True value for lazy expression
                        mov $0x1,%edx
                        xchg %edx,%edi
                        call _printBool # Const string "!"
                        mov $.LC8,%ecx
                        xchg %ecx,%edi
                        call _printString # Const boolean true
                        mov $0x1,%esi
                        xchg %esi,%edi
                        call _printBool # Const boolean false
                        mov $0x0,%esi
                        xchg %esi,%edi
                        call _printBool # Const int 0
                        mov $0x0,%r9d
                        mov %r9d,%eax
                        mov $0x1,%ebx
                        xchg %eax,%ebx
                        int $0x80
                        ret
# End of function main
# Function printBool
# Source: ./tests/good/core013.lat:110:1
_printBool:
                          push %rbp
                          mov %rsp,%rbp
                          sub $0x8,%rsp
                        printBool_block7:
                          cmp $0x0,%edi
                          sete %r9b # If condition
                          cmp $0x0,%r9d
                          je printBool_block6
                        printBool_block2: # Const string "false"
                          mov $.LC9,%r12d
                          mov %r12,-0x8(%rbp)
                          xchg %r12d,%edi
                          call _printString
                          mov -0x8(%rbp),%r12
                        printBool_block3:
                          mov $0x0,%eax
                          leave
                          ret
                        printBool_block6: # Const string "true"
                          mov $.LC10,%r10d
                          xchg %r10d,%edi
                          call _printString
                          mov $0x0,%eax
                          leave
                          ret
# End of function printBool
# Function test
# Source: ./tests/good/core013.lat:119:1
_test:
                            push %rbp
                            mov %rsp,%rbp
                            sub $0x8,%rsp
                          test_block5:
                            mov %rdi,-0x8(%rbp)
                            call _printInt
                            mov -0x8(%rbp),%rdi # Const int 0
                            mov $0x0,%ecx
                            cmp %ecx,%edi
                            setg %r14b
                            movzbl %r14b,%r14d
                            mov %r14d,%eax
                            leave
                            ret
# End of function test