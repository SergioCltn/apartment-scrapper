import contextlib
import logging
import sqlite3

import model
import service
import utils

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s",
)


def main():
    db_path = "../local/apartments.db"
    table = "apartments"
    with contextlib.closing(sqlite3.connect(db_path)) as conn:
        raw_apartments = utils.fetch_table(conn, table)

    apartments = [
        model.Apartment.from_raw_apartment(
            raw_apartment=raw_apartment,
        )
        for raw_apartment in raw_apartments
    ]

    count = service.get_address_counter(apartments=apartments)
    logging.info(count)


if __name__ == "__main__":
    main()
