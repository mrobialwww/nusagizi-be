import os
import re

# Pemetaan nomor lama ke nomor baru
mapping = {
    "23": "1", "25": "2", "45": "3", "47": "4", "48": "5", "49": "6",
    "22": "7", "24": "8", "46": "9", "50": "10", "51": "11", "69": "12",
    "3": "13", "4": "14", "5": "15", "6": "16", "52": "17", "53": "18",
    "7": "19", "8": "20", "9": "21", "10": "22", "11": "23", "12": "24",
    "13": "25", "14": "26", "54": "27",
    "15": "28", "16": "29", "17": "30", "67": "31", "18": "32", "19": "33",
    "21": "34", "56": "35", "57": "36", "58": "37", "59": "38", "60": "39", "61": "40",
    "37": "41", "38": "42", "39": "43", "40": "44", "41": "45", "42": "46",
    "43": "47", "44": "48", "62": "49", "63": "50", "64": "51",
    "26": "52", "27": "53", "68": "54", "28": "55", "29": "56", "34": "57",
    "30": "58", "31": "59", "32": "60", "33": "61", "65": "62", "66": "63",
    "2": "64",
    "1": "65"
}

def replace_in_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original_content = content

    # Gunakan placeholder agar tidak terjadi penggantian ganda berantai
    for old, new in mapping.items():
        # Pola 1: Komentar kode Go -> // GetContacts (Endpoint: 23)
        # Menangkap spasi opsional atau karakter persis "(Endpoint: 23)"
        content = re.sub(rf'\(Endpoint:\s*{old}\)', f'(Endpoint: __PLACEHOLDER_{new}__)', content)
        content = re.sub(rf'\(Endpoint:\s*{old}\s*-', f'(Endpoint: __PLACEHOLDER_{new}__ -', content)
        
        # Pola 2: Judul Markdown -> ### 23.
        # Khusus untuk header section (harus di awal baris)
        content = re.sub(rf'^### {old}\.', f'### __PLACEHOLDER_{new}__.', content, flags=re.MULTILINE)
        
        # Pola 3: Referensi teks Markdown -> endpoint 23
        # Menggunakan word boundary agar "endpoint 2" tidak match dengan "endpoint 23"
        content = re.sub(rf'\bendpoint {old}\b', f'endpoint __PLACEHOLDER_{new}__', content, flags=re.IGNORECASE)
        
        # Pola 4: Tabel Markdown -> | 23 |
        content = re.sub(rf'\|\s*{old}\s*\|', f'| __PLACEHOLDER_{new}__ |', content)

    # Kembalikan placeholder ke angka aslinya
    for old, new in mapping.items():
        content = content.replace(f'__PLACEHOLDER_{new}__', new)

    if content != original_content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Berhasil mengupdate: {filepath}")
    else:
        print(f"Tidak ada perubahan: {filepath}")

if __name__ == '__main__':
    # Eksekusi hanya di direktori yang relevan
    files_to_check = ['endpoint.md']

    # Cari semua file .go hanya di direktori service dan handler
    for folder in ['internal/services', 'internal/handlers']:
        if os.path.exists(folder):
            for root, dirs, files in os.walk(folder):
                for file in files:
                    if file.endswith('.go'):
                        files_to_check.append(os.path.join(root, file))

    for f in files_to_check:
        if os.path.exists(f):
            replace_in_file(f)
        else:
            print(f"File tidak ditemukan: {f}")
    
    print("Selesai!")
