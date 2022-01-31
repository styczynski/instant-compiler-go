.text
.global main
.LC5:
  .string "hello */"
.LC6:
  .string "/* world"
.LC7:
  .string ""
.LC0:
  .string "%d\n"
.LC1:
  .string "%s\n"
.LC2:
  .string "Error: %s\n"
.LC3:
  .string "FAILED ASSERTION"
.LC4:
  .string "="
# Function printInt
# Source: ./tests/good/core001.lat:6:1
_printInt:
    push %rbp
    mov %rsp,%rbp
    sub $0x8,%rsp
  printInt_block2: # Const string "%d\n"
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
# Source: ./tests/good/core001.lat:11:1
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
# Function AddStrings
# Source: ./tests/good/core001.lat:24:1
_AddStrings:
        push %rbp
        mov %rsp,%rbp
        sub $0x18,%rsp
      AddStrings_block7:
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        call strlen
        mov %eax,%r8d
        mov -0x8(%rbp),%rdi
        mov -0x10(%rbp),%rsi
        mov %r8,-0x8(%rbp)
        mov %rdi,-0x10(%rbp)
        mov %rsi,-0x18(%rbp)
        xchg %esi,%edi
        call strlen
        mov %eax,%edx
        mov -0x8(%rbp),%r8
        mov -0x10(%rbp),%rdi
        mov -0x18(%rbp),%rsi # Const int 1
        mov $0x1,%ecx
        add %ecx,%edx
        add %edx,%r8d
        mov %rdi,-0x8(%rbp)
        mov %rsi,-0x10(%rbp)
        xchg %r8d,%edi
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
# Source: ./tests/good/core001.lat:41:1
_error:
          push %rbp
          mov %rsp,%rbp
          sub $0x8,%rsp
        error_block6: # Const string "Error: %s\n"
          mov $.LC2,%esi
          mov %rax,-0x8(%rbp)
          xchg %esi,%edi
          mov $0x0,%eax
          call printf
          mov -0x8(%rbp),%rax # Const int 1
          mov $0x1,%esi
          xchg %esi,%edi
          call exit
# End of function error
# Function assert
# Source: ./tests/good/core001.lat:53:1
_assert:
          assert_block7: # If condition
            cmp $0x0,%edi
            je assert_block3
          assert_block2: # Const int 2
            mov $0x2,%r13d # Assign variable x
            mov %r13d,%ecx
          assert_block4:
            mov $0x0,%eax
            ret
          assert_block3: # Const string "FAILED ASSERTION"
            mov $.LC3,%r15d
            xchg %r15d,%edi
            call _error
            mov $0x0,%eax
            ret
# End of function assert
# Function assertEq
# Source: ./tests/good/core001.lat:62:1
_assertEq:
            assertEq_block2:
              cmp %esi,%edi
              sete %r12b
              movzbl %r12b,%r12d
              xchg %r12d,%edi
              call _assert
              mov $0x0,%eax
              ret
# End of function assertEq
# Function main (Entrypoint)
# Source: ./tests/good/core001.lat:68:1
main:
              main_block4: # Const int 10
                mov $0xa,%ecx
                xchg %ecx,%edi
                call _fac
                mov %eax,%edi
                call _printInt # Const int 10
                mov $0xa,%r11d
                xchg %r11d,%edi
                call _rfac
                mov %eax,%r13d
                xchg %r13d,%edi
                call _printInt # Const int 10
                mov $0xa,%r11d
                xchg %r11d,%edi
                call _mfac
                mov %eax,%r13d
                xchg %r13d,%edi
                call _printInt # Const int 10
                mov $0xa,%edx
                xchg %edx,%edi
                call _ifac
                mov %eax,%r15d
                xchg %r15d,%edi
                call _printInt # Const int 0
                mov $0x0,%ecx # Const int 10
                mov $0xa,%edi # Const int 1
                mov $0x1,%r8d
                mov %edi,%edx
              main_block11: # Const int 0
                mov $0x0,%r11d
                cmp %r11d,%edx
                setg %r10b
                movzbl %r10b,%r10d # While condition
                cmp $0x0,%r10d
                je main_block14
              main_block12:
                mov %rdx,%rbx
                mov %r8d,%eax
                imul %edx
                mov %eax,%r8d
                mov %rbx,%rdx # Const int 1
                mov $0x1,%r12d
                sub %r12d,%edx
                nop
                nop # While loop return to block_11
                jmp main_block11
                mov $0x0,%eax
                ret
              main_block14:
                xchg %r8d,%edi
                call _printInt # Const string "="
                mov $.LC4,%r12d # Const int 60
                mov $0x3c,%edi
                xchg %r12d,%edi
                xchg %r12d,%esi
                call _repStr
                mov %eax,%r13d
                xchg %r13d,%edi
                call _printString # Const string "hello */"
                mov $.LC5,%r11d
                xchg %r11d,%edi
                call _printString # Const string "/* world"
                mov $.LC6,%ecx
                xchg %ecx,%edi
                call _printString # Const int 0
                mov $0x0,%ecx
                mov %ecx,%eax
                mov $0x1,%ebx
                xchg %eax,%ebx
                int $0x80
                ret
# End of function main
# Function fac
# Source: ./tests/good/core001.lat:89:1
_fac:
                fac_block3: # Const int 0
                  mov $0x0,%ecx # Const int 0
                  mov $0x0,%edx # Const int 1
                  mov $0x1,%esi
                  mov %edi,%r9d
                fac_block9: # Const int 0
                  mov $0x0,%r8d
                  cmp %r8d,%r9d
                  setg %dil
                  movzbl %dil,%edi # While condition
                  cmp $0x0,%edi
                  je fac_block11
                fac_block10:
                  mov %rdx,%rbx
                  mov %esi,%eax
                  imul %r9d
                  mov %eax,%esi
                  mov %rbx,%rdx # Const int 1
                  mov $0x1,%r10d
                  sub %r10d,%r9d
                  nop
                  nop # While loop return to block_9
                  jmp fac_block9
                  mov $0x0,%eax
                  ret
                fac_block11:
                  mov %esi,%eax
                  ret
# End of function fac
# Function rfac
# Source: ./tests/good/core001.lat:102:1
_rfac:
                    push %rbp
                    mov %rsp,%rbp
                    sub $0x8,%rsp
                  rfac_block5: # Const int 0
                    mov $0x0,%r14d
                    cmp %r14d,%edi
                    sete %sil
                    movzbl %sil,%esi # If condition
                    cmp $0x0,%esi
                    je rfac_block6
                  rfac_block2: # Const int 1
                    mov $0x1,%r11d
                    mov %r11d,%eax
                    leave
                    ret
                  rfac_block6: # Const int 1
                    mov $0x1,%r13d
                    mov %edi,%r11d
                    sub %r13d,%r11d
                    mov %rdi,-0x8(%rbp)
                    xchg %r11d,%edi
                    call _rfac
                    mov %eax,%r10d
                    mov -0x8(%rbp),%rdi
                    mov %rdx,%rbx
                    mov %edi,%eax
                    imul %r10d
                    mov %eax,%edi
                    mov %rbx,%rdx
                    mov %edi,%eax
                    leave
                    ret
# End of function rfac
# Function mfac
# Source: ./tests/good/core001.lat:109:1
_mfac:
                      push %rbp
                      mov %rsp,%rbp
                      sub $0x8,%rsp
                    mfac_block4: # Const int 0
                      mov $0x0,%r13d
                      cmp %r13d,%edi
                      sete %sil
                      movzbl %sil,%esi # If condition
                      cmp $0x0,%esi
                      je mfac_block5
                    mfac_block3: # Const int 1
                      mov $0x1,%r9d
                      mov %r9d,%eax
                      leave
                      ret
                    mfac_block5: # Const int 1
                      mov $0x1,%r13d
                      mov %edi,%esi
                      sub %r13d,%esi
                      mov %rdi,-0x8(%rbp)
                      xchg %esi,%edi
                      call _nfac
                      mov %eax,%r11d
                      mov -0x8(%rbp),%rdi
                      mov %rdx,%rbx
                      mov %edi,%eax
                      imul %r11d
                      mov %eax,%edi
                      mov %rbx,%rdx
                      mov %edi,%eax
                      leave
                      ret
# End of function mfac
# Function nfac
# Source: ./tests/good/core001.lat:116:1
_nfac:
                        push %rbp
                        mov %rsp,%rbp
                        sub $0x8,%rsp
                      nfac_block5: # Const int 0
                        mov $0x0,%edx
                        cmp %edx,%edi
                        setne %cl
                        movzbl %cl,%ecx # If condition
                        cmp $0x0,%ecx
                        je nfac_block6
                      nfac_block4: # Const int 1
                        mov $0x1,%edx
                        mov %edi,%r9d
                        sub %edx,%r9d
                        mov %rdi,-0x8(%rbp)
                        xchg %r9d,%edi
                        call _mfac
                        mov %eax,%edx
                        mov -0x8(%rbp),%rdi
                        nop
                        mov %edx,%eax
                        imul %edi
                        mov %eax,%edx
                        nop
                        mov %edx,%eax
                        leave
                        ret
                      nfac_block6: # Const int 1
                        mov $0x1,%r13d
                        mov %r13d,%eax
                        leave
                        ret
# End of function nfac
# Function ifac
# Source: ./tests/good/core001.lat:123:1
_ifac:
                        ifac_block3: # Const int 1
                          mov $0x1,%r11d
                          xchg %r11d,%edi
                          xchg %r11d,%esi
                          call _ifac2f
                          mov %eax,%ecx
                          mov %ecx,%eax
                          ret
# End of function ifac
# Function ifac2f
# Source: ./tests/good/core001.lat:125:1
_ifac2f:
                            push %rbp
                            mov %rsp,%rbp
                            sub $0x10,%rsp
                          ifac2f_block4:
                            cmp %esi,%edi
                            sete %r8b
                            movzbl %r8b,%r8d # If condition
                            cmp $0x0,%r8d
                            je ifac2f_block6
                          ifac2f_block3:
                            mov %edi,%eax
                            leave
                            ret
                          ifac2f_block6:
                            cmp %esi,%edi
                            setg %r12b
                            movzbl %r12b,%r12d # If condition
                            cmp $0x0,%r12d
                            je ifac2f_block7
                          ifac2f_block5: # Const int 1
                            mov $0x1,%r13d
                            mov %r13d,%eax
                            leave
                            ret
                          ifac2f_block7: # Const int 0
                            mov $0x0,%r9d
                            mov %edi,%r13d
                            add %esi,%r13d # Const int 2
                            mov $0x2,%r15d
                            mov %rdx,%rbx
                            mov $0x0,%rdx
                            mov %r13d,%eax
                            idiv %r15d
                            mov %eax,%r13d
                            mov %rbx,%rdx
                            mov %rsi,-0x8(%rbp)
                            mov %r13,-0x10(%rbp)
                            xchg %r13d,%esi
                            call _ifac2f
                            mov %eax,%r15d
                            mov -0x8(%rbp),%rsi
                            mov -0x10(%rbp),%r13 # Const int 1
                            mov $0x1,%r14d
                            add %r14d,%r13d
                            mov %r15,-0x8(%rbp)
                            xchg %r13d,%edi
                            call _ifac2f
                            mov %eax,%ecx
                            mov -0x8(%rbp),%r15
                            mov %rdx,%rbx
                            mov %r15d,%eax
                            imul %ecx
                            mov %eax,%r15d
                            mov %rbx,%rdx
                            mov %r15d,%eax
                            leave
                            ret
# End of function ifac2f
# Function repStr
# Source: ./tests/good/core001.lat:135:1
_repStr:
                              push %rbp
                              mov %rsp,%rbp
                              sub $0x18,%rsp
                            repStr_block4: # Const string ""
                              mov $.LC7,%r15d # Const int 0
                              mov $0x0,%ecx
                              mov %ecx,%r8d
                            repStr_block8:
                              cmp %esi,%r8d
                              setl %r11b
                              movzbl %r11b,%r11d # While condition
                              cmp $0x0,%r11d
                              je repStr_block9
                            repStr_block6:
                              mov %rdi,-0x8(%rbp)
                              mov %rsi,-0x10(%rbp)
                              mov %r8,-0x18(%rbp)
                              xchg %r15d,%edi
                              xchg %r15d,%esi
                              call _AddStrings
                              mov %eax,%r15d
                              mov -0x8(%rbp),%rdi
                              mov -0x10(%rbp),%rsi
                              mov -0x18(%rbp),%r8 # Const int 1
                              mov $0x1,%ecx
                              add %ecx,%r8d
                              nop # While loop return to block_8
                              jmp repStr_block8
                              mov $0x0,%eax
                              leave
                              ret
                            repStr_block9:
                              mov %r15d,%eax
                              leave
                              ret
# End of function repStr