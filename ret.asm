.data
.text
.global main
main:
mov $5, %eax
mov %eax, %ebx
mov $1, %eax
int $0x80
ret
