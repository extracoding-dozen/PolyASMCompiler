.intel_syntax noprefix
.global _start
_start:
   push rbp
   mov rbp, rsp
   sub rsp, 112 ;
   jmp .L_str_1_skip
 .L_str_1_data:
   .asciz "/etc/passwd"
 .L_str_1_skip:
   lea r12, [rip + .L_str_1_data]
   mov qword [rbp - 8], r12
   mov r12, qword [rbp - 8]
   mov qword [rbp - 16], r12
   jmp .L_str_3_skip
 .L_str_3_data:
   .asciz "/tmp/passwd.bak"
 .L_str_3_skip:
   lea r12, [rip + .L_str_3_data]
   mov qword [rbp - 24], r12
   mov r12, qword [rbp - 24]
   mov qword [rbp - 32], r12
   mov rax, 2
   mov rdi, qword [rbp - 16]
   mov rsi, 0
   mov rdx, 0
   syscall
   mov qword [rbp - 64], rax
   mov rax, 2
   mov rdi, qword [rbp - 32]
   mov rsi, 65
   mov rdx, 511
   syscall
   mov qword [rbp - 72], rax
   mov rax, 40
   mov rdi, qword [rbp - 72]
   mov rsi, qword [rbp - 64]
   mov rdx, 0
   mov r10, 2147483647
   syscall
   mov qword [rbp - 80], rax
   mov rax, 3
   mov rdi, qword [rbp - 64]
   syscall
   mov qword [rbp - 88], rax
   mov rax, 3
   mov rdi, qword [rbp - 72]
   syscall
   mov qword [rbp - 96], rax
   mov r12, 0
   mov qword [rbp - 56], r12
   mov rax, 60
   mov rdi, qword [rbp - 56]
   syscall
   mov qword [rbp - 48], rax
   mov rsp, rbp
   pop rbp
   ret
   