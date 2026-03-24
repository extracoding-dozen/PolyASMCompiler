.intel_syntax noprefix
.global _start
_start:
  push rbp
  mov rbp, rsp
  sub rsp, 1248
  jmp .L_str_1_skip
.L_str_1_data:
  .byte 0x6d, 0x36, 0x2f, 0x32, 0x6d, 0x2a, 0x23, 0x21, 0x29, 0x6c, 0x2e, 0x2d, 0x25, 0x00
.L_str_1_skip:
  lea r12, [rip + .L_str_1_data]
  mov qword [rbp - 8], r12
  mov r12, 1592
  mov qword [rbp - 496], r12
  mov r12, qword [rbp - 496]
  mov r13, 1592
  cmp r12, r13
  jne .L5
    mov r14, qword [rbp - 8]
  sub rsp, 16
  mov qword [rbp - 8], rsp
  mov r8, r14
  mov rdi, rsp
  mov rcx, 13
.L_decrypt_1:
  mov r13b, byte ptr [r8]
  mov al, 66
  xor r13b, al
  mov byte ptr [rdi], r13b
  inc r8
  inc rdi
  dec rcx
  test rcx, rcx
  jnz .L_decrypt_1
  mov byte ptr [rdi], 0
    jmp .L6
.L5:
  mov r12, 3735928559
  mov qword [rbp - 504], r12
  mov r12, qword [rbp - 504]
  add r12, r13
  mov qword [rbp - 512], r12
.L6:
  mov r12, 9592
  mov qword [rbp - 520], r12
  mov r12, qword [rbp - 520]
  mov r13, 9592
  cmp r12, r13
  jne .L7
  mov rax, 39
  syscall
  mov qword [rbp - 248], rax
  jmp .L8
.L7:
  mov r12, 3735928559
  mov qword [rbp - 528], r12
  mov r12, qword [rbp - 528]
  add r12, r13
  mov qword [rbp - 536], r12
.L8:
  mov r12, qword [rbp - 8]
  mov qword [rbp - 16], r12
  mov r12, 5718
  mov qword [rbp - 544], r12
  mov r12, qword [rbp - 544]
  mov r13, 5718
  cmp r12, r13
  jne .L9
  mov rax, 104
  syscall
  mov qword [rbp - 256], rax
  jmp .L10
.L9:
  mov r12, 3735928559
  mov qword [rbp - 552], r12
  mov r12, qword [rbp - 552]
  add r12, r13
  mov qword [rbp - 560], r12
.L10:
  jmp .L_str_3_skip
.L_str_3_data:
  .byte 0x2e, 0x75, 0x6c, 0x71, 0x2e, 0x69, 0x60, 0x62, 0x6a, 0x2f, 0x63, 0x60, 0x6a, 0x00
.L_str_3_skip:
  lea r12, [rip + .L_str_3_data]
  mov qword [rbp - 24], r12
    mov r14, qword [rbp - 24]
  sub rsp, 16
  mov qword [rbp - 24], rsp
  mov r8, r14
  mov r9, rsp
  mov rbx, 13
.L_decrypt_3:
  sub rsp, 0
  mov r12b, byte ptr [r8]
  mov al, 1
  xor r12b, al
  mov byte ptr [r9], r12b
  inc r8
  inc r9
  sub rbx, 1
  cmp rbx, 0
  jg .L_decrypt_3
  mov byte ptr [r9], 0
    mov rax, 39
  syscall
  mov qword [rbp - 264], rax
  mov r12, qword [rbp - 24]
  mov qword [rbp - 32], r12
  jmp .L_str_5_skip
.L_str_5_data:
  .byte 0x4c, 0x54, 0x5d, 0x52, 0x48, 0x53, 0x51, 0x00
.L_str_5_skip:
  lea r12, [rip + .L_str_5_data]
  mov qword [rbp - 40], r12
  mov rax, 39
  syscall
  mov qword [rbp - 272], rax
    mov r14, qword [rbp - 40]
  sub rsp, 16
  mov qword [rbp - 40], rsp
  mov r8, r14
  mov rdi, rsp
  mov rdx, 7
.L_decrypt_5:
  mov r15b, byte ptr [r8]
  xor r15b, 60
  mov byte ptr [rdi], r15b
  add r8, 1
  add rdi, 1
  sub rdx, 1
  cmp rdx, 0
  jg .L_decrypt_5
  mov byte ptr [rdi], 0
    mov rax, 104
  syscall
  mov qword [rbp - 280], rax
  mov r12, 4211
  mov qword [rbp - 568], r12
  mov r12, qword [rbp - 568]
  mov r13, 4211
  cmp r12, r13
  jne .L11
  mov r12, qword [rbp - 40]
  mov qword [rbp - 48], r12
  jmp .L12
.L11:
  mov r12, 3735928559
  mov qword [rbp - 576], r12
  mov r12, qword [rbp - 576]
  add r12, r13
  mov qword [rbp - 584], r12
.L12:
  mov r12, 630
  mov qword [rbp - 592], r12
  mov r12, qword [rbp - 592]
  mov r13, 630
  cmp r12, r13
  jne .L13
  mov rax, 104
  syscall
  mov qword [rbp - 288], rax
  jmp .L14
.L13:
  mov r12, 3735928559
  mov qword [rbp - 600], r12
  mov r12, qword [rbp - 600]
  add r12, r13
  mov qword [rbp - 608], r12
.L14:
  mov r12, -1
  mov qword [rbp - 56], r12
  mov r12, 6915
  mov qword [rbp - 616], r12
  mov r12, qword [rbp - 616]
  mov r13, 6915
  cmp r12, r13
  jne .L15
  mov rax, 104
  syscall
  mov qword [rbp - 296], rax
  jmp .L16
.L15:
  mov r12, 3735928559
  mov qword [rbp - 624], r12
  mov r12, qword [rbp - 624]
  add r12, r13
  mov qword [rbp - 632], r12
.L16:
  mov rax, 2
  mov rdi, qword [rbp - 16]
  mov rsi, 0
  mov rdx, 0
  syscall
  mov qword [rbp - 168], rax
  mov r12, qword [rbp - 168]
  mov r13, 0
  cmp r12, r13
  mov rax, 104
  syscall
  mov qword [rbp - 304], rax
  mov r12, 7842
  mov qword [rbp - 640], r12
  mov r12, qword [rbp - 640]
  mov r13, 7842
  cmp r12, r13
  jne .L17
  jl .L3
  jmp .L18
.L17:
  mov r12, 3735928559
  mov qword [rbp - 648], r12
  mov r12, qword [rbp - 648]
  add r12, r13
  mov qword [rbp - 656], r12
.L18:
  mov rax, 102
  syscall
  mov qword [rbp - 312], rax
  mov rax, 8
  mov rdi, qword [rbp - 168]
  mov rsi, 0
  mov rdx, 2
  syscall
  mov qword [rbp - 56], rax
  mov rax, 3
  mov rdi, qword [rbp - 168]
  syscall
  mov qword [rbp - 176], rax
  mov r12, 9804
  mov qword [rbp - 664], r12
  mov r12, qword [rbp - 664]
  mov r13, 9804
  cmp r12, r13
  jne .L19
  mov rax, 102
  syscall
  mov qword [rbp - 320], rax
  jmp .L20
.L19:
  mov r12, 3735928559
  mov qword [rbp - 672], r12
  mov r12, qword [rbp - 672]
  add r12, r13
  mov qword [rbp - 680], r12
.L20:
.L3:
  mov r12, 3391
  mov qword [rbp - 688], r12
  mov r12, qword [rbp - 688]
  mov r13, 3391
  cmp r12, r13
  jne .L21
  mov rax, 110
  syscall
  mov qword [rbp - 328], rax
  jmp .L22
.L21:
  mov r12, 3735928559
  mov qword [rbp - 696], r12
  mov r12, qword [rbp - 696]
  add r12, r13
  mov qword [rbp - 704], r12
.L22:
  mov r12, 7342
  mov qword [rbp - 712], r12
  mov r12, qword [rbp - 712]
  mov r13, 7342
  cmp r12, r13
  jne .L23
  mov r12, qword [rbp - 56]
  mov qword [rbp - 64], r12
  jmp .L24
.L23:
  mov r12, 3735928559
  mov qword [rbp - 720], r12
  mov r12, qword [rbp - 720]
  add r12, r13
  mov qword [rbp - 728], r12
.L24:
  mov rax, 102
  syscall
  mov qword [rbp - 336], rax
  mov r12, 0
  mov qword [rbp - 72], r12
  mov rax, 39
  syscall
  mov qword [rbp - 344], rax
  mov r12, 9141
  mov qword [rbp - 736], r12
  mov r12, qword [rbp - 736]
  mov r13, 9141
  cmp r12, r13
  jne .L25
  mov r12, qword [rbp - 64]
  mov qword [rbp - 80], r12
  jmp .L26
.L25:
  mov r12, 3735928559
  mov qword [rbp - 744], r12
  mov r12, qword [rbp - 744]
  add r12, r13
  mov qword [rbp - 752], r12
.L26:
  mov rax, 104
  syscall
  mov qword [rbp - 352], rax
  mov r12, 252
  mov qword [rbp - 760], r12
  mov r12, qword [rbp - 760]
  mov r13, 252
  cmp r12, r13
  jne .L27
  mov r12, qword [rbp - 80]
  mov r13, 0
  cmp r12, r13
  jmp .L28
.L27:
  mov r12, 3735928559
  mov qword [rbp - 768], r12
  mov r12, qword [rbp - 768]
  add r12, r13
  mov qword [rbp - 776], r12
.L28:
  mov r12, 9646
  mov qword [rbp - 784], r12
  mov r12, qword [rbp - 784]
  mov r13, 9646
  cmp r12, r13
  jne .L29
  mov rax, 39
  syscall
  mov qword [rbp - 360], rax
  jmp .L30
.L29:
  mov r12, 3735928559
  mov qword [rbp - 792], r12
  mov r12, qword [rbp - 792]
  add r12, r13
  mov qword [rbp - 800], r12
.L30:
  je .L1
  jmp .L_str_12_skip
.L_str_12_data:
  .byte 0x32, 0x18, 0x12, 0x15, 0x04, 0x0c, 0x41, 0x02, 0x0e, 0x0c, 0x11, 0x13, 0x0e, 0x0c, 0x08, 0x12, 0x04, 0x05, 0x3d, 0x0f, 0x00
.L_str_12_skip:
  lea r12, [rip + .L_str_12_data]
  mov qword [rbp - 96], r12
    mov r14, qword [rbp - 96]
  sub rsp, 32
  mov qword [rbp - 96], rsp
  mov r10, r14
  mov rdi, rsp
  mov rbx, 20
.L_decrypt_12:
  mov r13b, byte ptr [r10]
  mov al, 97
  xor r13b, al
  mov byte ptr [rdi], r13b
  add r10, 1
  add rdi, 1
  add rsp, 0
  dec rbx
  test rbx, rbx
  jnz .L_decrypt_12
  mov byte ptr [rdi], 0

  mov r12, -1
  mov qword [rbp - 88], r12
  mov rax, 2
  mov rdi, qword [rbp - 16]
  mov rsi, 65
  mov rdx, 420
  syscall
  mov qword [rbp - 184], rax
  mov r12, qword [rbp - 184]
  mov r13, 0
  cmp r12, r13
  mov r12, 1898
  mov qword [rbp - 808], r12
  mov r12, qword [rbp - 808]
  mov r13, 1898
  cmp r12, r13
  jne .L31
  jl .L4
  jmp .L32
.L31:
  mov r12, 3735928559
  mov qword [rbp - 816], r12
  mov r12, qword [rbp - 816]
  add r12, r13
  mov qword [rbp - 824], r12
.L32:
  mov r12, 9296
  mov qword [rbp - 832], r12
  mov r12, qword [rbp - 832]
  mov r13, 9296
  cmp r12, r13
  jne .L33
  mov rax, 102
  syscall
  mov qword [rbp - 368], rax
  jmp .L34
.L33:
  mov r12, 3735928559
  mov qword [rbp - 840], r12
  mov r12, qword [rbp - 840]
  add r12, r13
  mov qword [rbp - 848], r12
.L34:
  mov rax, 8
  mov rdi, qword [rbp - 184]
  mov rsi, qword [rbp - 64]
  mov rdx, 0
  syscall
  mov qword [rbp - 192], rax
  mov rax, 102
  syscall
  mov qword [rbp - 376], rax
  mov r12, 4785
  mov qword [rbp - 856], r12
  mov r12, qword [rbp - 856]
  mov r13, 4785
  cmp r12, r13
  jne .L35
  mov rdi, qword [rbp - 96]
  sub rcx, rcx
  not rcx
  sub rax, rax
  cld
  repne scasb
  not rcx
  dec rcx
  mov qword [rbp - 200], rcx
  jmp .L36
.L35:
  mov r12, 3735928559
  mov qword [rbp - 864], r12
  mov r12, qword [rbp - 864]
  add r12, r13
  mov qword [rbp - 872], r12
.L36:
  mov rax, 39
  syscall
  mov qword [rbp - 384], rax
  mov r12, 189
  mov qword [rbp - 880], r12
  mov r12, qword [rbp - 880]
  mov r13, 189
  cmp r12, r13
  jne .L37
  mov rax, 1
  mov rdi, qword [rbp - 184]
  mov rsi, qword [rbp - 96]
  mov rdx, qword [rbp - 200]
  syscall
  mov qword [rbp - 208], rax
  jmp .L38
.L37:
  mov r12, 3735928559
  mov qword [rbp - 888], r12
  mov r12, qword [rbp - 888]
  add r12, r13
  mov qword [rbp - 896], r12
.L38:
  mov rax, 102
  syscall
  mov qword [rbp - 392], rax
  mov rax, 3
  mov rdi, qword [rbp - 184]
  syscall
  mov qword [rbp - 216], rax
  mov rax, 110
  syscall
  mov qword [rbp - 400], rax
  mov r12, 0
  mov qword [rbp - 88], r12
  mov rax, 104
  syscall
  mov qword [rbp - 408], rax
.L4:
  mov rax, 102
  syscall
  mov qword [rbp - 416], rax
  mov rax, 82
  mov rdi, qword [rbp - 16]
  mov rsi, qword [rbp - 32]
  syscall
  mov qword [rbp - 104], rax
  mov r12, 1998
  mov qword [rbp - 904], r12
  mov r12, qword [rbp - 904]
  mov r13, 1998
  cmp r12, r13
  jne .L39
  mov rax, 104
  syscall
  mov qword [rbp - 424], rax
  jmp .L40
.L39:
  mov r12, 3735928559
  mov qword [rbp - 912], r12
  mov r12, qword [rbp - 912]
  add r12, r13
  mov qword [rbp - 920], r12
.L40:
  jmp .L2
.L1:
  mov r12, 9422
  mov qword [rbp - 928], r12
  mov r12, qword [rbp - 928]
  mov r13, 9422
  cmp r12, r13
  jne .L41
  mov rax, 104
  syscall
  mov qword [rbp - 432], rax
  jmp .L42
.L41:
  mov r12, 3735928559
  mov qword [rbp - 936], r12
  mov r12, qword [rbp - 936]
  add r12, r13
  mov qword [rbp - 944], r12
.L42:
.L2:
  mov r12, 8265
  mov qword [rbp - 952], r12
  mov r12, qword [rbp - 952]
  mov r13, 8265
  cmp r12, r13
  jne .L43
  mov rax, 39
  syscall
  mov qword [rbp - 440], rax
  jmp .L44
.L43:
  mov r12, 3735928559
  mov qword [rbp - 960], r12
  mov r12, qword [rbp - 960]
  add r12, r13
  mov qword [rbp - 968], r12
.L44:
  jmp .L_str_15_skip
.L_str_15_data:
  .byte 0xb1, 0xea, 0xf3, 0xee, 0xb1, 0xf6, 0xff, 0xfd, 0xf5, 0xb0, 0xfc, 0xff, 0xf5, 0x00
.L_str_15_skip:
  lea r12, [rip + .L_str_15_data]
  mov qword [rbp - 120], r12
  mov r12, 5907
  mov qword [rbp - 976], r12
  mov r12, qword [rbp - 976]
  mov r13, 5907
  cmp r12, r13
  jne .L45

  mov r14, qword [rbp - 120]
  sub rsp, 16
  mov qword [rbp - 120], rsp
  mov r8, r14
  mov r11, rsp
  mov rdx, 13
.L_decrypt_15:
  mov r12b, byte ptr [r8]
  xor r12b, 158
  mov byte ptr [r11], r12b
  add r8, 1
  add r11, 1
  sub rdx, 1
  cmp rdx, 0
  jg .L_decrypt_15
  mov byte ptr [r11], 0

  jmp .L46
.L45:
  mov r12, 3735928559
  mov qword [rbp - 984], r12
  mov r12, qword [rbp - 984]
  add r12, r13
  mov qword [rbp - 992], r12
.L46:
  mov r12, 8915
  mov qword [rbp - 1000], r12
  mov r12, qword [rbp - 1000]
  mov r13, 8915
  cmp r12, r13
  jne .L47
  mov rax, 110
  syscall
  mov qword [rbp - 448], rax
  jmp .L48
.L47:
  mov r12, 3735928559
  mov qword [rbp - 1008], r12
  mov r12, qword [rbp - 1008]
  add r12, r13
  mov qword [rbp - 1016], r12
.L48:
  mov r12, 511
  mov qword [rbp - 128], r12
  mov rax, 110
  syscall
  mov qword [rbp - 456], rax
  mov r12, 9613
  mov qword [rbp - 1024], r12
  mov r12, qword [rbp - 1024]
  mov r13, 9613
  cmp r12, r13
  jne .L49
  mov rax, 90
  mov rdi, qword [rbp - 120]
  mov rsi, qword [rbp - 128]
  syscall
  mov qword [rbp - 112], rax
  jmp .L50
.L49:
  mov r12, 3735928559
  mov qword [rbp - 1032], r12
  mov r12, qword [rbp - 1032]
  add r12, r13
  mov qword [rbp - 1040], r12
.L50:
  mov r12, 5972
  mov qword [rbp - 1048], r12
  mov r12, qword [rbp - 1048]
  mov r13, 5972
  cmp r12, r13
  jne .L51
  mov rax, 110
  syscall
  mov qword [rbp - 464], rax
  jmp .L52
.L51:
  mov r12, 3735928559
  mov qword [rbp - 1056], r12
  mov r12, qword [rbp - 1056]
  add r12, r13
  mov qword [rbp - 1064], r12
.L52:
  mov r12, 10
  mov qword [rbp - 144], r12
  mov r12, 4253
  mov qword [rbp - 1072], r12
  mov r12, qword [rbp - 1072]
  mov r13, 4253
  cmp r12, r13
  jne .L53
  mov rax, 104
  syscall
  mov qword [rbp - 472], rax
  jmp .L54
.L53:
  mov r12, 3735928559
  mov qword [rbp - 1080], r12
  mov r12, qword [rbp - 1080]
  add r12, r13
  mov qword [rbp - 1088], r12
.L54:
  mov r12, 880
  mov qword [rbp - 1096], r12
  mov r12, qword [rbp - 1096]
  mov r13, 880
  cmp r12, r13
  jne .L55
  mov r12, 0
  mov qword [rbp - 224], r12
  jmp .L56
.L55:
  mov r12, 3735928559
  mov qword [rbp - 1104], r12
  mov r12, qword [rbp - 1104]
  add r12, r13
  mov qword [rbp - 1112], r12
.L56:
  mov r12, 2084
  mov qword [rbp - 1120], r12
  mov r12, qword [rbp - 1120]
  mov r13, 2084
  cmp r12, r13
  jne .L57
  mov r12, qword [rbp - 144]
  mov qword [rbp - 232], r12
  jmp .L58
.L57:
  mov r12, 3735928559
  mov qword [rbp - 1128], r12
  mov r12, qword [rbp - 1128]
  add r12, r13
  mov qword [rbp - 1136], r12
.L58:
  mov r12, 6981
  mov qword [rbp - 1144], r12
  mov r12, qword [rbp - 1144]
  mov r13, 6981
  cmp r12, r13
  jne .L59
  mov rax, 39
  syscall
  mov qword [rbp - 480], rax
  jmp .L60
.L59:
  mov r12, 3735928559
  mov qword [rbp - 1152], r12
  mov r12, qword [rbp - 1152]
  add r12, r13
  mov qword [rbp - 1160], r12
.L60:
  mov r12, 2659
  mov qword [rbp - 1168], r12
  mov r12, qword [rbp - 1168]
  mov r13, 2659
  cmp r12, r13
  jne .L61
  lea r12, [rbp - 224]
  mov qword [rbp - 240], r12
  jmp .L62
.L61:
  mov r12, 3735928559
  mov qword [rbp - 1176], r12
  mov r12, qword [rbp - 1176]
  add r12, r13
  mov qword [rbp - 1184], r12
.L62:
  mov r12, 2313
  mov qword [rbp - 1192], r12
  mov r12, qword [rbp - 1192]
  mov r13, 2313
  cmp r12, r13
  jne .L63
  mov rax, 35
  mov rdi, qword [rbp - 240]
  mov rsi, 0
  syscall
  mov qword [rbp - 136], rax
  jmp .L64
.L63:
  mov r12, 3735928559
  mov qword [rbp - 1200], r12
  mov r12, qword [rbp - 1200]
  add r12, r13
  mov qword [rbp - 1208], r12
.L64:
  mov r12, 0
  mov qword [rbp - 160], r12
  mov r12, 6549
  mov qword [rbp - 1216], r12
  mov r12, qword [rbp - 1216]
  mov r13, 6549
  cmp r12, r13
  jne .L65
  mov rax, 60
  mov rdi, qword [rbp - 160]
  syscall
  mov qword [rbp - 152], rax
  jmp .L66
.L65:
  mov r12, 3735928559
  mov qword [rbp - 1224], r12
  mov r12, qword [rbp - 1224]
  add r12, r13
  mov qword [rbp - 1232], r12
.L66:
  mov rax, 39
  syscall
  mov qword [rbp - 488], rax
  mov rsp, rbp
  pop rbp
  ret

