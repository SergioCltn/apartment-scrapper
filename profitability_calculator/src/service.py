import model


def get_address_counter(apartments: list[model.Apartment]) -> dict:
    keys_and_fvalues = (
        ("street", lambda apartment: apartment.address.street),
        ("number", lambda apartment: apartment.address.number),
        (
            "neighborhood",
            lambda apartment: apartment.address.neighborhood,
        ),
        ("district", lambda apartment: apartment.address.district),
        ("city", lambda apartment: apartment.address.city),
        ("region", lambda apartment: apartment.address.region),
        (
            "autonomous_community",
            lambda apartment: apartment.address.autonomous_community,
        ),
    )
    count = {}
    for apartment in apartments:
        for k, fv in keys_and_fvalues:
            if not count.get(k, False):
                count[k] = {}
            count[k][fv(apartment)] = count[k].get(fv(apartment), 0) + 1

    return count
