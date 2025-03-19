import re
import sqlite3


def get_digit(text: str) -> None | str:
    match = re.search(r"(\d+)", text)
    if match is not None:
        return match.group(1)
    return None


def fetch_table(
    conn: sqlite3.Connection,
    table: str,
) -> list[dict]:
    conn.row_factory = sqlite3.Row  # Allows accessing columns by name
    cursor = conn.cursor()

    cursor.execute(f"SELECT * FROM {table};")
    rows = cursor.fetchall()

    return [dict(row) for row in rows]
