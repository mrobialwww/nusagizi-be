### Contoh 1: Status "Sesuai Usia" (Skor KPSP 9-10)

```json
{
    "id": "b3b2c1a0-1234-5678-abcd-1234567890ab",
    "kpsp_score": 9,
    "next_check_date": "25-09-2026",
    "status": "Sesuai Usia",
    "created_at": "2026-08-25T10:00:00Z",
    "domains": [
        {
            "nerve_name": "Gross motor skills",
            "total_question": 2,
            "true_answer": 2
        },
        {
            "nerve_name": "Fine motor skills",
            "total_question": 3,
            "true_answer": 3
        },
        {
            "nerve_name": "Speech and language",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "nerve_name": "Socialization",
            "total_question": 2,
            "true_answer": 2
        }
    ]
}
```

### Contoh 2: Status "Perkembangan Meragukan" (Skor KPSP 7-8)

```json
{
    "id": "e7f8g9h0-9876-5432-fedc-0987654321fe",
    "kpsp_score": 7,
    "next_check_date": "10-09-2026",
    "status": "Perkembangan meragukan",
    "created_at": "2026-08-25T10:15:00Z",
    "domains": [
        {
            "nerve_name": "Gross motor skills",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "nerve_name": "Fine motor skills",
            "total_question": 2,
            "true_answer": 1
        },
        {
            "nerve_name": "Speech and language",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "nerve_name": "Socialization",
            "total_question": 2,
            "true_answer": 2
        }
    ]
}
```
