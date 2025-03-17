import json
import logging
import re
from dataclasses import dataclass

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s",
)


@dataclass
class Address:
    street: str | None = None
    number: str | None = None
    neighborhood: str | None = None
    district: str | None = None
    city: str | None = None
    region: str | None = None
    autonomous_community: str | None = None

    @classmethod
    def from_location(cls, location: str) -> "Address":
        """Create an Address instance from a location string.

        Args:
            location (str): The location string to parse.

        Returns:
            Address: An Address instance with parsed data.

        """
        address = cls()
        parts = [p.strip() for p in location.split(",")]

        match = re.match(r"^(.*?)(?:[\s,-]+(\d+|s/n))?$", parts[0])
        if match:
            street = match.group(1).strip()
            address.street = (
                None if street.startswith("Barrio") else street
            )
            address.number = match.group(2) if match.group(2) else None

        # Handle addresses that might have number split to second field
        if not address.number and re.match(r"^-?\d+|s/n$", parts[1]):
            address.number = parts[1]

        for part in parts:
            if part.startswith("Distrito"):
                address.district = part
            elif part.startswith("Comarca"):
                address.region = part
            elif part.startswith("Barrio"):
                address.neighborhood = part

        address.city = parts[-3]
        address.region = parts[-2]
        address.autonomous_community = parts[-1]

        return address


@dataclass
class Details:
    n_rooms: int | None = None
    sqm_constructed: float | None = None
    sqm_usable: float | None = None
    bathrooms: int | None = None
    terrace: bool = False
    state: str | None = None
    orientation: str | None = None
    balcony: bool = False
    built_in: int | None = None
    heating: str | None = None
    reduced_mobility: str | None = None
    garage: str | None = None
    extra: str | None = None

    elevator: bool | None = None
    inside: bool | None = None
    floor: int | None = None

    kwh_sqm_year_emissions: float | None = None
    kg_co2_sqm_year_consumption: float | None = None

    tenant: str | None = None
    certificate: str | None = None

    storage_room: bool = False
    built_in_wardrobes: bool = False
    air_conditioning: bool = False
    semi_detached_house: bool = False
    green_areas: bool = False
    pool: bool = False
    new_building_development: bool = False

    @classmethod
    def from_raw_details(
        cls,
        raw_details: dict,
    ) -> "Details":
        details = cls()
        for k, value in raw_details.items():
            for val in value:
                v: str = val
                if " m² construidos, " in v and " m² útiles" in v:
                    numbers = [
                        int(num) for num in re.findall(r"(\d+)", v)
                    ]
                    if len(numbers) == 2:
                        (
                            details.sqm_constructed,
                            details.sqm_usable,
                        ) = numbers

                elif "m² construidos" in v:
                    details.sqm_constructed = float(
                        re.search(r"(\d+)", v).group(1),
                    )
                elif "m² útiles" in v:
                    details.sqm_usable = float(
                        re.search(r"(\d+)", v).group(1),
                    )
                elif "habitación" in v or "habitaciones" in v:
                    if "Sin " in v:
                        details.n_rooms = 0
                    else:
                        details.n_rooms = int(
                            re.search(r"(\d+)", v).group(1),
                        )
                elif "baño" in v:
                    details.bathrooms = int(
                        re.search(r"(\d+)", v).group(1),
                    )
                elif "Terraza y balcón" in v:
                    details.terrace = True
                    details.balcony = True
                elif "Terraza" in v:
                    details.terrace = True
                elif "Balcón" in v:
                    details.balcony = True
                elif "Segunda mano" in v:
                    details.state = v
                elif "Orientación" in v:
                    details.orientation = v.replace(
                        "Orientación ",
                        "",
                    )
                elif "Construido en " in v:
                    details.built_in = int(
                        re.search(r"(\d+)", v).group(1),
                    )
                elif "calefacción" in v.lower():
                    details.heating = v
                elif "movilidad reducida" in v:
                    details.reduced_mobility = v
                elif "garaje" in v:
                    details.garage = v
                elif "Sin ascensor" in v:
                    details.elevator = False
                elif "Con ascensor" in v:
                    details.elevator = True
                elif "planta " in v.lower() or "Bajo" in v:
                    if "Entreplanta " in v or "Bajo" in v:
                        details.floor = 0
                    else:
                        details.floor = int(
                            re.search(r"(\d+)(?=ª)", v).group(1),
                        )
                    if "exterior" in v:
                        details.inside = False
                    elif "interior" in v:
                        details.inside = True
                elif "interior" in v:
                    details.inside = True
                elif "exterior" in v:
                    details.inside = False
                elif "Emisiones: \n" in v:
                    details.kwh_sqm_year_emissions = float(
                        re.search(r"(\d+)", v).group(1),
                    )
                elif "Consumo: \n" in v:
                    details.kg_co2_sqm_year_consumption = float(
                        re.search(r"(\d+)", v).group(1),
                    )
                elif "Certificado" in k and "Consumo:" in v:
                    details.certificate = None
                elif (
                    "Consumo:" in v
                    or "Emisiones:" in v
                    or "No indicado" in v
                ):
                    pass
                elif (
                    "Inmueble exento" in v
                    or "En trámite" in v
                    or "No indicado" in v
                ):
                    details.certificate = v
                elif (
                    "Alquilada, con inquilinos" in v
                    or "Ocupada ilegalmente" in v
                ):
                    details.tenant = v
                elif "Trastero" in v:
                    details.storage_room = True
                elif "Aire acondicionado" in v:
                    details.air_conditioning = True
                elif "Armarios empotrados" in v:
                    details.built_in_wardrobes = True
                elif "Chalet adosado" in v:
                    details.semi_detached_house = True
                elif "Zonas verdes" in v:
                    details.green_areas = True
                elif "Piscina" in v:
                    details.pool = True
                elif "Promoción de obra nueva" in v:
                    details.new_building_development = True
                else:
                    logging.getLogger(__name__).warning(
                        "Value {%s: %s} from details could not be parsed, adding it to extra",
                        k,
                        v,
                    )
                    details.extra = v

        return details


@dataclass
class Apartment:
    id: int | None = None
    title: str | None = None
    property_price_euros: float | None = None
    price_per_sqm: float | None = None
    monthly_community_fees_euros: float | None = None
    location: str | None = None
    description: str | None = None

    address: Address | None = None

    details: Details | None = None

    @classmethod
    def from_raw_apartment(cls, raw_apartment: dict) -> "Apartment":
        """Create an Apartment instance from a raw apartment dictionary.

        Args:
            raw_apartment (dict): A dictionary containing raw apartment data.

        Returns:
            Apartment: An instance of the Apartment class populated with the provided data.

        """
        apartment_info = {
            "id": None,
            "title": None,
            "property_price_euros": None,
            "price_per_sqm": None,
            "monthly_community_fees_euros": None,
            "location": None,
            "description": None,
        }

        outer_keys_map = {
            "id": "id",
            "title": "title",
            "propertyPrice": "property_price_euros",
            "pricePerSqm": "price_per_sqm",
            "communityFees": "monthly_community_fees_euros",
            "description": "description",
            "location": "location",
        }

        address = None
        details = None

        for k, value in raw_apartment.items():
            if k == "location":
                address = Address.from_location(value)

            if k in outer_keys_map:
                apartment_info[outer_keys_map[k]] = value
            elif k == "details":
                raw_details = json.loads(value)
                details = Details.from_raw_details(
                    raw_details=raw_details,
                )

        apartment = cls(**apartment_info)
        apartment.address = address
        apartment.details = details

        return apartment
