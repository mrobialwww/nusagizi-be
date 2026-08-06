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
            "developmental_domain": "gross_motor_skills",
            "total_question": 2,
            "true_answer": 2
        },
        {
            "developmental_domain": "fine_motor_skills",
            "total_question": 3,
            "true_answer": 3
        },
        {
            "developmental_domain": "speech_and_language",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "developmental_domain": "socialization",
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
            "developmental_domain": "gross_motor_skills",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "developmental_domain": "fine_motor_skills",
            "total_question": 2,
            "true_answer": 1
        },
        {
            "developmental_domain": "speech_and_language",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "developmental_domain": "socialization",
            "total_question": 2,
            "true_answer": 2
        }
    ]
}
```
