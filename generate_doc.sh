#!/bin/bash

OUTPUT_FILE="DOC.md"

echo "📝 Menyusun dokumentasi ke $OUTPUT_FILE..."

# Tulis Header Utama
cat << 'EOF' > "$OUTPUT_FILE"
# Dokumentasi Struktur & Kode Project

Dokumen ini dihasilkan secara otomatis untuk memetakan seluruh struktur folder dan isi kode di dalam project ini.

## Struktur Project (Tree)

```text
EOF

# Jalankan perintah tree jika ada, jika tidak pakai find untuk simulasi tree sederhana
if command -v tree &> /dev/null; then
    tree -I "node_modules|.git|vendor|.next|dist" >> "$OUTPUT_FILE"
else
    find . -not -path '*/.*' -not -path './vendor*' -not -path './node_modules*' | sed -e 's/^[^\/]*\//⎹  /' -e 's/\/[^\/]*$/⎹__/' >> "$OUTPUT_FILE"
fi

echo '```' >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "## Isi Kode Berdasarkan File" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Mencari file kode yang valid
find . -type f \( -name "*.go" -o -name "go.mod" -o -name "*.js" -o -name "*.ts" -o -name "*.py" -o -name "*.php" -o -name "*.json" -o -name "*.html" -o -name "*.css" -o -name "*.esql" \) \
-not -path "*/.*" \
-not -path "./vendor/*" \
-not -path "./node_modules/*" \
-not -path "./dist/*" \
-not -path "./.next/*" | while read -r file; do
    
    # Hapus `./` di depan nama file agar rapi
    clean_path=$(echo "$file" | sed 's|^\./||')
    
    echo "### File: \`$clean_path\`" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    
    # Ambil ekstensi file
    ext="${file##*.}"
    
    # Tentukan syntax highlighting menggunakan kutip tunggal agar aman
    if [[ "$file" == *"go.mod" ]]; then
        echo '```text' >> "$OUTPUT_FILE"
    elif [[ "$ext" =~ ^(go|js|ts|py|php|json|html|css)$ ]]; then
        echo '```'"$ext" >> "$OUTPUT_FILE"
    else
        echo '```text' >> "$OUTPUT_FILE"
    fi
    
    # Masukkan isi file
    cat "$file" >> "$OUTPUT_FILE"
    
    echo "" >> "$OUTPUT_FILE"
    echo '```' >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    echo "---" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
done

echo "✅ Berhasil! File $OUTPUT_FILE telah dibuat."