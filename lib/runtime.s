.text
.global AddStrings
.global PrintString
.LC_Print_Format_String:
  .string "\n"
AddStrings:
        pushq   %rbp
        movq    %rsp, %rbp
        subq    $32, %rsp
        movl    %edi, -20(%rbp)
        movl    %esi, -24(%rbp)
        movl    -20(%rbp), %eax
        cltq
        movq    %rax, %rdi
        call    strlen
        movl    %eax, -4(%rbp)
        movl    -24(%rbp), %eax
        cltq
        movq    %rax, %rdi
        call    strlen
        movl    %eax, -8(%rbp)
        movl    -4(%rbp), %edx
        movl    -8(%rbp), %eax
        addl    %edx, %eax
        cltq
        movq    %rax, %rdi
        call    malloc
        movq    %rax, -16(%rbp)
        movl    -20(%rbp), %eax
        cltq
        movq    %rax, %rdx
        movq    -16(%rbp), %rax
        movq    %rdx, %rsi
        movq    %rax, %rdi
        call    strcpy
        movl    -24(%rbp), %eax
        cltq
        movq    %rax, %rdx
        movq    -16(%rbp), %rax
        movq    %rdx, %rsi
        movq    %rax, %rdi
        call    strcat
        movq    -16(%rbp), %rax
        leave
        ret
PrintString:
        pushq   %rbp
        movq    %rsp, %rbp
        subq    $16, %rsp
        movl    %edi, -4(%rbp)
        movl    -4(%rbp), %eax
        movl    %eax, %esi
        movl    $.LC_Print_Format_String, %edi
        movl    $0, %eax
        call    printf
        nop
        leave
        ret