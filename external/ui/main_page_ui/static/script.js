document.addEventListener("DOMContentLoaded", () => {
    const editor = document.getElementById('editor');
    const highlightingContent = document.getElementById('highlighting-content');

    window.rawCompiledAsm = "";

    // Делаем функцию глобальной, чтобы её можно было вызывать из loadPsFile
    window.updateHighlighting = function() {
        let text = editor.value;

        if (text.length > 0 && text[text.length - 1] === '\n') {
            text += ' ';
        }

        highlightingContent.innerHTML = applySyntaxHighlighting(text);
    };

    editor.addEventListener('input', window.updateHighlighting);

    editor.addEventListener('scroll', () => {
        const highlightingLayer = document.getElementById('highlighting-layer');
        highlightingLayer.scrollTop = editor.scrollTop;
        highlightingLayer.scrollLeft = editor.scrollLeft;
    });

    function applySyntaxHighlighting(text) {
        text = text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');

        // 1. ИЩЕМ СТРОКИ И КОММЕНТАРИИ ЗА ОДИН ПРОХОД
        // ВАЖНО: класс теперь называется token-str, чтобы не конфликтовать со словом string!
        text = text.replace(/(\/\/.*|"(?:\\.|[^"\\])*")/g, function(match) {
            if (match.startsWith('//')) {
                return '<span class="token-comment">' + match + '</span>';
            } else {
                return '<span class="token-str">' + match + '</span>';
            }
        });

        // 2. РАЗБИВАЕМ ТЕКСТ ПО ЭТИМ СПАНАМ (чтобы не сломать их внутренности)
        const parts = text.split(/(<span class="token-comment">.*?<\/span>|<span class="token-str">.*?<\/span>)/g);

        for (let i = 0; i < parts.length; i++) {
            // Обрабатываем только то, что НЕ является строкой или комментарием
            if (!parts[i].startsWith('<span')) {
                let part = parts[i];

                // Ключевые слова
                part = part.replace(/\b(if|else|while|return|execute|fork|exit|sleep|copy|rename|write|delete|get_file_size|chmod|chown|setattr|useradd|passwd|connect)\b/g, '<span class="token-keyword">$1</span>');
                // Типы
                part = part.replace(/\b(int|string|qword|bool)\b/g, '<span class="token-type">$1</span>');
                // Числа (Hex и Dec)
                part = part.replace(/\b(0x[a-fA-F0-9]+|\d+)\b/g, '<span class="token-number">$1</span>');
                // Функции
                part = part.replace(/(\b[a-zA-Z_]\w*\b)(?=\s*\()/g, '<span class="token-function">$1</span>');

                parts[i] = part;
            }
        }
        return parts.join('');
    }

    // Запускаем подсветку при старте страницы
    window.updateHighlighting();
});

// --- ПОДСВЕТКА ASM ---
function applyAsmHighlighting(text) {
    if (!text) return "";
    text = text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    text = text.replace(/([;].*)/g, '<span class="token-comment">$1</span>');

    const parts = text.split(/(<span class="token-comment">.*?<\/span>)/g);
    for (let i = 0; i < parts.length; i++) {
        if (!parts[i].startsWith('<span')) {
            let part = parts[i];
            part = part.replace(/^([a-zA-Z_.0-9]+:)/gm, '<span class="token-keyword" style="color:#a5d6ff">$1</span>'); // Метки
            part = part.replace(/\b(mov|push|pop|call|ret|syscall|add|sub|xor|jmp|je|jne|cmp|lea|test)\b/gi, '<span class="token-keyword">$1</span>');
            part = part.replace(/\b(rax|rbx|rcx|rdx|rsi|rdi|rbp|rsp|r[8-9]|r1[0-5]|eax|ebx|ecx|edx|al|ah|bl|bh)\b/gi, '<span class="token-function">$1</span>');
            part = part.replace(/\b(0x[a-fA-F0-9]+|\d+)\b/g, '<span class="token-number">$1</span>');
            parts[i] = part;
        }
    }
    return parts.join('');
}

// --- ОСНОВНАЯ ФУНКЦИЯ КОМПИЛЯЦИИ ---
async function runCompilation() {
    const logsOutput = document.getElementById('logs-output');
    const compiledOutput = document.getElementById('compiled-output');

    logsOutput.textContent = "[INFO] Compiling...\n";
    compiledOutput.innerHTML = "";
    window.rawCompiledAsm = "";

    const payload = {
        code: document.getElementById('editor').value,
        optimizationLevel: document.getElementById('opt-level').value,
        debugMode: document.getElementById('debug-mode').checked,
        enableObfuscation: document.getElementById('enable-obfuscation').checked,
        enableSandboxNoise: document.getElementById('chk-sandbox').checked,
        enableStringCrypt: document.getElementById('chk-string').checked,
        enableOpaquePreds: document.getElementById('chk-opaque').checked,
        noiseFrequency: parseInt(document.getElementById('num-noise').value, 10),
        opaqueFrequency: parseInt(document.getElementById('num-opaque').value, 10),
        repeatObfuscator: parseInt(document.getElementById('obfs_qua').value, 10)
    };

    try {
        const response = await fetch('/api/compile', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

        const data = await response.json();
        logsOutput.textContent += data.logs;

        window.rawCompiledAsm = data.compiledCode;
        compiledOutput.innerHTML = applyAsmHighlighting(data.compiledCode);

    } catch (error) {
        logsOutput.textContent += `\n[ERROR] Build failed: ${error.message}\n`;
    }
}

// --- ФУНКЦИЯ СКАЧИВАНИЯ ФАЙЛА ---
function saveAsmFile() {
    if (!window.rawCompiledAsm) {
        alert("Нет скомпилированного кода для сохранения!");
        return;
    }

    const blob = new Blob([window.rawCompiledAsm], { type: "text/plain" });
    const url = URL.createObjectURL(blob);

    const a = document.createElement("a");
    a.href = url;
    a.download = "payload.s";
    document.body.appendChild(a);
    a.click();

    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

function openGraph(type) {
    if (type === 'pre') {
        window.open('/graph/clear', '_blank', 'width=1000,height=800');
    } else if (type === 'post') {
        window.open('/graph/obfuscated', '_blank', 'width=1000,height=800');
    }
}

// --- ФУНКЦИЯ ЗАГРУЗКИ ФАЙЛА .PS ---
function loadPsFile(event) {
    const file = event.target.files[0];
    if (!file) return;

    const fileName = file.name.toLowerCase();
    if (!fileName.endsWith('.ps') && !fileName.endsWith('.txt')) {
        alert("Пожалуйста, выберите файл с расширением .ps (или текстовый файл)");
        return;
    }

    const reader = new FileReader();

    reader.onload = function(e) {
        const contents = e.target.result;
        const editor = document.getElementById('editor');
        editor.value = contents;

        // Вызываем функцию обновления подсветки
        if (window.updateHighlighting) {
            window.updateHighlighting();
        } else {
            const event = new Event('input', { bubbles: true });
            editor.dispatchEvent(event);
        }

        const headerText = document.querySelector('.ide-container .panel-header span');
        if (headerText) headerText.textContent = file.name;

        const logsOutput = document.getElementById('logs-output');
        logsOutput.textContent += `[INFO] Файл ${file.name} успешно загружен.\n`;

        event.target.value = '';
    };

    reader.readAsText(file);
}