.intel_syntax noprefix
.global _start
_start:
  push rbp
  mov rbp, rsp
  sub rsp, 192
  jmp .L_str_1_skip
.L_str_1_data:
  .asciz "/tmp/hack.log"
.L_str_1_skip:
  lea r12, [rip + .L_str_1_data]
  mov qword [rbp - 8], r12
  mov r12, qword [rbp - 8]
  mov qword [rbp - 16], r12
  jmp .L_str_3_skip
.L_str_3_data:
  .asciz "/tmp/hack.bak"
.L_str_3_skip:
  lea r12, [rip + .L_str_3_data]
  mov qword [rbp - 24], r12
  mov r12, qword [rbp - 24]
  mov qword [rbp - 32], r12
  jmp .L_str_5_skip
.L_str_5_data:
  .asciz "phantom"
.L_str_5_skip:
  lea r12, [rip + .L_str_5_data]
  mov qword [rbp - 40], r12
  mov r12, qword [rbp - 40]
  mov qword [rbp - 48], r12
  mov r12, -1
  mov qword [rbp - 56], r12
  mov rax, 2
  mov rdi, qword [rbp - 16]
  mov rsi, 0
  mov rdx, 0
  syscall
  mov qword [rbp - 128], rax
  mov r12, qword [rbp - 128]
  cmp r12, r13
  jl .L3
  mov rax, 8
  mov rdi, qword [rbp - 128]
  mov rsi, 0
  mov rdx, 2
  syscall
  mov qword [rbp - 56], rax
  mov rax, 3
  mov rdi, qword [rbp - 128]
  syscall
  mov qword [rbp - 136], rax
.L3:
  mov r12, qword [rbp - 56]
  mov qword [rbp - 64], r12
  mov r12, 0
  mov qword [rbp - 72], r12
  mov r12, qword [rbp - 64]
  mov qword [rbp - 80], r12
  mov r12, qword [rbp - 80]
  mov r13, 0
  cmp r12, r13
  je .L1
  jmp .L_str_12_skip
.L_str_12_data:
  .asciz "System compromised\n"
.L_str_12_skip:
  lea r12, [rip + .L_str_12_data]
  mov qword [rbp - 96], r12
  mov rax, 2
  mov rdi, qword [rbp - 16]
  mov rsi, 65
  mov rdx, 420
  syscall
  mov qword [rbp - 144], rax
  mov rax, 8
  mov rdi, qword [rbp - 144]
  mov rsi, qword [rbp - 64]
  mov rdx, 0
  syscall
  mov qword [rbp - 152], rax
  mov rdi, qword [rbp - 96]
  sub rcx, rcx
  not rcx
  sub rax, rax
  cld
  repne scasb
  not rcx
  dec rcx
  mov qword [rbp - 160], rcx
  mov rax, 1
  mov rdi, qword [rbp - 144]
  mov rsi, qword [rbp - 96]
  mov rdx, qword [rbp - 160]
  syscall
  mov qword [rbp - 168], rax
  mov rax, 3
  mov rdi, qword [rbp - 144]
  syscall
  mov qword [rbp - 176], rax
  mov rax, 82
  mov rdi, qword [rbp - 16]
  mov rsi, qword [rbp - 32]
  syscall
  mov qword [rbp - 104], rax
  jmp .L2
.L1:
.L2:
  mov r12, 0
  mov qword [rbp - 120], r12
  mov rax, 60
  mov rdi, qword [rbp - 120]
  syscall
  mov qword [rbp - 112], rax
  mov rsp, rbp
  pop rbp
  ret
