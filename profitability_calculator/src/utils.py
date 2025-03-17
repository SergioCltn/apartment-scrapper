import sqlite3


def fetch_table(
    conn: sqlite3.Connection,
    table: str,
) -> list[dict]:
    conn.row_factory = sqlite3.Row  # Allows accessing columns by name
    cursor = conn.cursor()

    cursor.execute(f"SELECT * FROM {table};")
    rows = cursor.fetchall()

    return [dict(row) for row in rows]
